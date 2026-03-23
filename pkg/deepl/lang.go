// opyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package deepl

import "strings"

// langNames maps uppercase DeepL language codes to Japanese display names.
var langNames = map[string]string{
	"EN": "英語",
	"ES": "スペイン語",
	"FR": "フランス語",
	"DE": "ドイツ語",
	"IT": "イタリア語",
	"PT": "ポルトガル語",
	"RU": "ロシア語",
	"KO": "韓国語",
	"ZH": "中国語",
	"AR": "アラビア語",
	"HI": "ヒンディー語",
	"TH": "タイ語",
	"VI": "ベトナム語",
	"ID": "インドネシア語",
	"TR": "トルコ語",
	"NL": "オランダ語",
	"PL": "ポーランド語",
	"SV": "スウェーデン語",
	"DA": "デンマーク語",
	"FI": "フィンランド語",
	"NO": "ノルウェー語",
	"UK": "ウクライナ語",
	"JA": "日本語",
}

// LangName returns the Japanese display name for a DeepL language code.
// It accepts codes in any case (e.g. "en", "EN", "En").
func LangName(code string) string {
	if name, ok := langNames[strings.ToUpper(code)]; ok {
		return name
	}
	return code
}
