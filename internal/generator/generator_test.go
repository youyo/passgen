package generator

import (
	"errors"
	"strings"
	"testing"

	"github.com/youyo/passgen/internal/charset"
)

// Round 1: バリデーション系

func TestGenerate_LengthZero_ReturnsError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Length = 0
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error for length=0, got nil")
	}
	if !errors.Is(err, ErrLengthNotPositive) {
		t.Errorf("expected ErrLengthNotPositive, got %v", err)
	}
}

func TestGenerate_NegativeLength_ReturnsError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Length = -1
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error for length=-1, got nil")
	}
	if !errors.Is(err, ErrLengthNotPositive) {
		t.Errorf("expected ErrLengthNotPositive, got %v", err)
	}
}

func TestGenerate_RequiredExceedsLength_ReturnsError(t *testing.T) {
	cfg := Config{
		Length:  3,
		Lower:   1,
		Upper:   1,
		Digits:  1,
		Symbols: 1,
	}
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error when required > length, got nil")
	}
	if !errors.Is(err, ErrRequiredExceedsLength) {
		t.Errorf("expected ErrRequiredExceedsLength, got %v", err)
	}
}

// Round 2: 基本生成

func TestGenerate_DefaultConfig_Returns20Chars(t *testing.T) {
	cfg := DefaultConfig()
	pw, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("len(password) = %d, want 20", len(pw))
	}
}

func TestGenerate_Length100_ReturnsCorrectLength(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Length = 100
	pw, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 100 {
		t.Errorf("len(password) = %d, want 100", len(pw))
	}
}

func TestGenerate_Length4_AllCategories(t *testing.T) {
	cfg := Config{
		Length:  4,
		Lower:   1,
		Upper:   1,
		Digits:  1,
		Symbols: 1,
	}
	// 100回試行して毎回全カテゴリ含有を検証
	for i := 0; i < 100; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(pw) != 4 {
			t.Fatalf("len(password) = %d, want 4", len(pw))
		}
		assertContainsCategory(t, pw, charset.Lower, "Lower")
		assertContainsCategory(t, pw, charset.Upper, "Upper")
		assertContainsCategory(t, pw, charset.Digits, "Digits")
		assertContainsCategory(t, pw, charset.Symbols, "Symbols")
	}
}

// Round 3: カテゴリ保証（統計テスト）

func TestGenerate_DefaultConfig_ContainsAllCategories_Statistical(t *testing.T) {
	cfg := DefaultConfig()
	for i := 0; i < 1000; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		assertContainsCategory(t, pw, charset.Lower, "Lower")
		assertContainsCategory(t, pw, charset.Upper, "Upper")
		assertContainsCategory(t, pw, charset.Digits, "Digits")
		assertContainsCategory(t, pw, charset.Symbols, "Symbols")
	}
}

func TestGenerate_SymbolsMin3_Statistical(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Symbols = 3
	for i := 0; i < 100; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		count := countCharsIn(pw, charset.Symbols)
		if count < 3 {
			t.Errorf("iteration %d: symbol count = %d, want >= 3, password = %q", i, count, pw)
		}
	}
}

// Round 4: 一意性・文字セット検証

func TestGenerate_TwoConsecutive_AreDifferent(t *testing.T) {
	cfg := DefaultConfig()
	pw1, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pw2, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw1 == pw2 {
		t.Errorf("two consecutive passwords are identical: %q", pw1)
	}
}

func TestGenerate_NoAmbiguousCharacters(t *testing.T) {
	ambiguous := "lIO01"
	cfg := DefaultConfig()
	for i := 0; i < 1000; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		for _, c := range ambiguous {
			if strings.ContainsRune(pw, c) {
				t.Errorf("iteration %d: ambiguous char %q found in %q", i, c, pw)
			}
		}
	}
}

func TestGenerate_OnlyValidCharacters(t *testing.T) {
	all := charset.All()
	cfg := DefaultConfig()
	for i := 0; i < 100; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		for _, c := range pw {
			if !strings.ContainsRune(all, c) {
				t.Errorf("iteration %d: invalid char %q in %q", i, c, pw)
			}
		}
	}
}

// Round 5: secureRandomIndex 検証

func TestSecureRandomIndex_ReturnsWithinRange(t *testing.T) {
	for i := 0; i < 10000; i++ {
		idx, err := secureRandomIndex(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if idx < 0 || idx >= 10 {
			t.Errorf("index %d out of range [0, 10)", idx)
		}
	}
}

func TestSecureRandomIndex_Distribution(t *testing.T) {
	const (
		max      = 6
		trials   = 60000
		expected = trials / max // 10000
	)
	counts := make([]int, max)
	for i := 0; i < trials; i++ {
		idx, err := secureRandomIndex(max)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		counts[idx]++
	}
	tolerance := float64(expected) * 0.05 // ±5%
	for i, c := range counts {
		diff := float64(c) - float64(expected)
		if diff < -tolerance || diff > tolerance {
			t.Errorf("value %d: count=%d, expected=%d (±%.0f)", i, c, expected, tolerance)
		}
	}
}

func TestSecureRandomIndex_MaxOne_ReturnsZero(t *testing.T) {
	for i := 0; i < 100; i++ {
		idx, err := secureRandomIndex(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if idx != 0 {
			t.Errorf("secureRandomIndex(1) = %d, want 0", idx)
		}
	}
}

func TestSecureRandomIndex_MaxZero_ReturnsError(t *testing.T) {
	_, err := secureRandomIndex(0)
	if err == nil {
		t.Fatal("expected error for max=0, got nil")
	}
}

// Round 6: DefaultConfig 検証

func TestDefaultConfig_Values(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Length != 20 {
		t.Errorf("Length = %d, want 20", cfg.Length)
	}
	if cfg.Lower != 1 {
		t.Errorf("Lower = %d, want 1", cfg.Lower)
	}
	if cfg.Upper != 1 {
		t.Errorf("Upper = %d, want 1", cfg.Upper)
	}
	if cfg.Digits != 1 {
		t.Errorf("Digits = %d, want 1", cfg.Digits)
	}
	if cfg.Symbols != 1 {
		t.Errorf("Symbols = %d, want 1", cfg.Symbols)
	}
}

// Round 7: Exclude テスト

func TestGenerate_Exclude_abc_NoAbcInPassword(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Exclude = "abc"
	for i := 0; i < 100; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		for _, ch := range "abc" {
			if strings.ContainsRune(pw, ch) {
				t.Errorf("iteration %d: excluded char %q found in %q", i, ch, pw)
			}
		}
	}
}

func TestGenerate_Exclude_AllLower_WithLower0_Success(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Exclude = charset.Lower
	cfg.Lower = 0
	pw, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("len(password) = %d, want 20", len(pw))
	}
	// パスワードに小文字が含まれないことを確認
	for _, ch := range pw {
		if strings.ContainsRune(charset.Lower, ch) {
			t.Errorf("excluded lower char %q found in %q", ch, pw)
		}
	}
}

func TestGenerate_Exclude_AllLower_WithLower1_Error(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Exclude = charset.Lower
	cfg.Lower = 1
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error when lower category is empty after exclusion, got nil")
	}
	if !errors.Is(err, ErrCategoryEmptyAfterExclude) {
		t.Errorf("expected ErrCategoryEmptyAfterExclude, got %v", err)
	}
}

func TestGenerate_Exclude_AllChars_AllMinZero_Error(t *testing.T) {
	cfg := Config{
		Length:  20,
		Lower:   0,
		Upper:   0,
		Digits:  0,
		Symbols: 0,
		Exclude: charset.All(),
	}
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error when all chars excluded, got nil")
	}
	if !errors.Is(err, ErrAllCharsExcluded) {
		t.Errorf("expected ErrAllCharsExcluded, got %v", err)
	}
}

func TestGenerate_Exclude_AllChars_WithMinimums_CategoryError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Exclude = charset.All()
	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error when all chars excluded with minimums, got nil")
	}
	if !errors.Is(err, ErrCategoryEmptyAfterExclude) {
		t.Errorf("expected ErrCategoryEmptyAfterExclude, got %v", err)
	}
}

func TestGenerate_Exclude_Empty_NoEffect(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Exclude = ""
	pw, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("len(password) = %d, want 20", len(pw))
	}
}

func TestGenerate_Exclude_AlreadyExcludedChars_NoError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Exclude = "lIO01" // 曖昧文字は既に除外済み
	pw, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 20 {
		t.Errorf("len(password) = %d, want 20", len(pw))
	}
}

func TestGenerate_Exclude_PartialLower_Success(t *testing.T) {
	// lower の大部分を除外、vwxyz が残る
	cfg := DefaultConfig()
	cfg.Exclude = "abcdefghijkmnopqrstu"
	cfg.Lower = 1
	for i := 0; i < 50; i++ {
		pw, err := Generate(cfg)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		// 除外文字が含まれないことを確認
		for _, ch := range cfg.Exclude {
			if strings.ContainsRune(pw, ch) {
				t.Errorf("iteration %d: excluded char %q found in %q", i, ch, pw)
			}
		}
		// v, w, x, y, z のいずれかが含まれることを確認（lower min=1）
		hasLower := false
		for _, ch := range pw {
			if strings.ContainsRune("vwxyz", ch) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			t.Errorf("iteration %d: no remaining lower char (vwxyz) found in %q", i, pw)
		}
	}
}

// ヘルパー関数

func assertContainsCategory(t *testing.T, pw, category, name string) {
	t.Helper()
	for _, c := range pw {
		if strings.ContainsRune(category, c) {
			return
		}
	}
	t.Errorf("password %q does not contain any %s character", pw, name)
}

func countCharsIn(pw, set string) int {
	count := 0
	for _, c := range pw {
		if strings.ContainsRune(set, c) {
			count++
		}
	}
	return count
}
