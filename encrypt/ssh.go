// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func GenerateSshKeyPairs() (string, string, error) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return "", "", fmt.Errorf("encrypt: failed to generate random key (%v)", err)
	}

	pub, err := ssh.NewPublicKey(key.Public())
	if err != nil {
		return "", "", fmt.Errorf("encrypt: failed to create public key (%v)", err)
	}
	pubKeyStr := string(ssh.MarshalAuthorizedKey(pub))
	privKeyStr := marshalRSAPrivate(key)

	return pubKeyStr, privKeyStr, nil
}

func EncryptWithSshKey(plainText, publicKey []byte) (string, error) {
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to parse authorized key (%v)", err)
	}
	// To get back to an *rsa.PublicKey, we need to first upgrade to the
	// ssh.CryptoPublicKey interface
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)

	if len(plainText) <= 256 {
		// plainText is small enough to only use OAEP encryption; this will result in less bytes to transfer.
		encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plainText, nil)
		if err != nil {
			return "", fmt.Errorf("encrypt: failed to encrypt with OAEP method (%v)", err)
		}
		if len(encryptedBytes) != 256 {
			return "", errors.New("encrypt: invalid encrypted data length with OAEP")
		}
		return base64.StdEncoding.EncodeToString(encryptedBytes), nil
	}

	// otherwise, encrypt using AES256
	key, cipherText, err := encryptAES256(plainText)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, key, nil)
	if err != nil {
		return "", err
	}
	if len(encryptedBytes) != 256 {
		return "", errors.New("encrypt: invalid encrypted data length")
	}
	return base64.StdEncoding.EncodeToString(append(encryptedBytes, cipherText...)), nil
}

func DecryptWithSshKey(cipherText, privateKey string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	if len(data) < 256 {
		return nil, errors.New("encrypt: not enough data to decrypt")
	}

	block, _ := pem.Decode([]byte(privateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	aesData := data[256:]
	payload, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, data[:256], nil)
	if err != nil {
		return nil, fmt.Errorf("encrypt: failed to decrypt with OAEP (%v)", err)
	}

	if len(aesData) == 0 {
		return payload, nil
	}

	decryptedAESKey := payload
	decrypted, err := decryptAES(decryptedAESKey, aesData)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func marshalRSAPrivate(priv *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}))
}

// encryptAES256 returns a random passphrase and corresponding bytes encrypted with it
func encryptAES256(data []byte) ([]byte, []byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, nil, fmt.Errorf("encrypt: failed to generate random key (%v)", err)
	}

	n := len(data)
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, uint64(n)); err != nil {
		return nil, nil, fmt.Errorf("encrypt: failed to write binary data (%v)", err)
	}
	if _, err := buf.Write(data); err != nil {
		return nil, nil, fmt.Errorf("encrypt: failed to write buffer data (%v)", err)
	}

	paddingN := aes.BlockSize - (buf.Len() % aes.BlockSize)
	if paddingN > 0 {
		padding := make([]byte, paddingN)
		if _, err := rand.Read(padding); err != nil {
			return nil, nil, fmt.Errorf("encrypt: failed to generate random key with padding (%v)", err)
		}
		if _, err := buf.Write(padding); err != nil {
			return nil, nil, fmt.Errorf("encrypt: failed to write padding buffer (%v)", err)
		}
	}
	plaintext := buf.Bytes()

	sum := sha256.Sum256(plaintext)
	plaintext = append(sum[:], plaintext...)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("encrypt: failed to create cipher (%v)", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, fmt.Errorf("encrypt: failed to generate random key with iv (%v)", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plaintext)
	return key, cipherText, nil
}

func decryptAES(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("encrypt: failed to create cipher with key (%v)", err)
	}

	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("encrypt: cipherText too short to decrypt")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("encrypt: cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// TODO: works inplace when both args are the same
	mode.CryptBlocks(cipherText, cipherText)

	expectedSum := cipherText[:32]
	actualSum := sha256.Sum256(cipherText[32:])
	if !bytes.Equal(expectedSum, actualSum[:]) {
		return nil, fmt.Errorf("encrypt: sha256 mismatch %v vs %v", expectedSum, actualSum)
	}

	buf := bytes.NewReader(cipherText[32:])
	var n uint64
	if err = binary.Read(buf, binary.LittleEndian, &n); err != nil {
		return nil, fmt.Errorf("encrypt: failed to read binary data (%v)", err)
	}
	payload := make([]byte, n)
	if _, err = buf.Read(payload); err != nil {
		return nil, fmt.Errorf("encrypt: failed to read payload data (%v)", err)
	}

	return payload, nil
}
