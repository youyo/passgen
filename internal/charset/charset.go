// Package charset は passgen で使用する文字セットを定義する。
// 曖昧文字（l, I, O, 0, 1）は全カテゴリから除外済み。
package charset

import "strings"

// 文字セット定数（曖昧文字除外済み）
const (
	Lower   = "abcdefghijkmnopqrstuvwxyz"  // 25文字（l 除外）
	Upper   = "ABCDEFGHJKLMNPQRSTUVWXYZ"   // 24文字（I, O 除外）
	Digits  = "23456789"                    // 8文字（0, 1 除外）
	Symbols = "-_.~"                        // 4文字
)

// All は全カテゴリを結合した文字列を返す。
func All() string {
	return Lower + Upper + Digits + Symbols
}

// Categories は各カテゴリを個別要素としたスライスを返す。
// 呼び出しごとに新しいスライスを返す。
func Categories() []string {
	return []string{Lower, Upper, Digits, Symbols}
}

// Exclude は base から excluded に含まれる文字を除去した新しい文字列を返す。
func Exclude(base, excluded string) string {
	if excluded == "" {
		return base
	}
	var b strings.Builder
	b.Grow(len(base))
	for _, r := range base {
		if !strings.ContainsRune(excluded, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
