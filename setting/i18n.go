// Copyright 2023 The GitBundle Inc. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

// defaultI18nLangNames must be a slice, we need the order
var defaultI18nLangNames = []string{
	"en-US", "English",
	"zh-CN", "简体中文",
	// "zh-HK", "繁體中文（香港）",
	// "zh-TW", "繁體中文（台灣）",
	// "de-DE", "Deutsch",
	// "fr-FR", "Français",
	// "nl-NL", "Nederlands",
	// "lv-LV", "Latviešu",
	// "ru-RU", "Русский",
	// "uk-UA", "Українська",
	// "ja-JP", "日本語",
	// "es-ES", "Español",
	// "pt-BR", "Português do Brasil",
	// "pt-PT", "Português de Portugal",
	// "pl-PL", "Polski",
	// "bg-BG", "Български",
	// "it-IT", "Italiano",
	// "fi-FI", "Suomi",
	// "tr-TR", "Türkçe",
	// "cs-CZ", "Čeština",
	// "sr-SP", "Српски",
	// "sv-SE", "Svenska",
	// "ko-KR", "한국어",
	// "el-GR", "Ελληνικά",
	// "fa-IR", "فارسی",
	// "hu-HU", "Magyar nyelv",
	// "id-ID", "Bahasa Indonesia",
	// "ml-IN", "മലയാളം",
}

func defaultI18nLangs() (res []string) {
	for i := 0; i < len(defaultI18nLangNames); i += 2 {
		res = append(res, defaultI18nLangNames[i])
	}
	return
}

func defaultI18nNames() (res []string) {
	for i := 0; i < len(defaultI18nLangNames); i += 2 {
		res = append(res, defaultI18nLangNames[i+1])
	}
	return
}