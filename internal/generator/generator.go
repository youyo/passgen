// Package generator はパスワード生成のコアロジックを提供する。
// crypto/rand を使用した暗号学的に安全な乱数生成を行う。
package generator

import (
	"errors"
	"fmt"

	"github.com/youyo/passgen/internal/charset"
)

// sentinel errors
var (
	// ErrLengthNotPositive は length が 0 以下の場合に返される。
	ErrLengthNotPositive = errors.New("length must be positive")
	// ErrRequiredExceedsLength は各カテゴリの最低文字数の合計が length を超える場合に返される。
	ErrRequiredExceedsLength = errors.New("required minimum characters exceeds length")
)

// Config はパスワード生成の設定を保持する。
type Config struct {
	Length  int // パスワード長（デフォルト: 20）
	Lower   int // 小文字の最低文字数（デフォルト: 1）
	Upper   int // 大文字の最低文字数（デフォルト: 1）
	Digits  int // 数字の最低文字数（デフォルト: 1）
	Symbols int // 記号の最低文字数（デフォルト: 1）
}

// DefaultConfig はデフォルト設定を返す。
func DefaultConfig() Config {
	return Config{
		Length:  20,
		Lower:   1,
		Upper:   1,
		Digits:  1,
		Symbols: 1,
	}
}

// Generate は Config に基づいてパスワードを生成する。
// 各カテゴリから最低文字数を保証し、残りを全文字セットからランダムに選択後、
// Fisher-Yates シャッフルで並びを均一にランダム化する。
func Generate(cfg Config) (string, error) {
	// バリデーション
	if cfg.Length <= 0 {
		return "", ErrLengthNotPositive
	}

	required := cfg.Lower + cfg.Upper + cfg.Digits + cfg.Symbols
	if required > cfg.Length {
		return "", fmt.Errorf("%w: need %d, got length %d", ErrRequiredExceedsLength, required, cfg.Length)
	}

	// 各カテゴリの文字セットと最低数をペアにする
	categories := charset.Categories()
	minimums := []int{cfg.Lower, cfg.Upper, cfg.Digits, cfg.Symbols}

	result := make([]byte, 0, cfg.Length)

	// 1. 各カテゴリから最低数をランダムに生成
	for i, min := range minimums {
		cat := categories[i]
		for j := 0; j < min; j++ {
			idx, err := secureRandomIndex(len(cat))
			if err != nil {
				return "", fmt.Errorf("generating category %d char: %w", i, err)
			}
			result = append(result, cat[idx])
		}
	}

	// 2. 残り(length - required)を全文字セットから生成
	all := charset.All()
	remaining := cfg.Length - required
	for i := 0; i < remaining; i++ {
		idx, err := secureRandomIndex(len(all))
		if err != nil {
			return "", fmt.Errorf("generating remaining char: %w", err)
		}
		result = append(result, all[idx])
	}

	// 3. Fisher-Yates シャッフル
	if err := shuffleBytes(result); err != nil {
		return "", fmt.Errorf("shuffling password: %w", err)
	}

	return string(result), nil
}
