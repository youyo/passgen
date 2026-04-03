package generator

import (
	"crypto/rand"
	"errors"
)

// secureRandomIndex は [0, max) の範囲で暗号学的に安全な乱数インデックスを返す。
// モジュラスバイアスを回避するため rejection sampling を使用する。
func secureRandomIndex(max int) (int, error) {
	if max <= 0 {
		return 0, errors.New("max must be positive")
	}
	if max == 1 {
		return 0, nil
	}

	// rejection sampling でモジュラスバイアスを回避
	// threshold は max の倍数で 256 以下の最大値
	threshold := 256 - (256 % max)
	var buf [1]byte
	for {
		if _, err := rand.Read(buf[:]); err != nil {
			return 0, err
		}
		r := int(buf[0])
		if r < threshold {
			return r % max, nil
		}
	}
}

// shuffleBytes は Fisher-Yates アルゴリズムでバイト列をインプレースシャッフルする。
// crypto/rand ベースの secureRandomIndex を使用するため均一なシャッフルが保証される。
func shuffleBytes(b []byte) error {
	for i := len(b) - 1; i > 0; i-- {
		j, err := secureRandomIndex(i + 1)
		if err != nil {
			return err
		}
		b[i], b[j] = b[j], b[i]
	}
	return nil
}
