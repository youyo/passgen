package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/youyo/passgen/internal/charset"
	"github.com/youyo/passgen/internal/cli"
	"github.com/youyo/passgen/internal/generator"
)

// parseCLI はテスト用にCLI引数をパースするヘルパー。
// パース成功時に *cli.CLI を返す。パースエラー時は error を返す。
func parseCLI(t *testing.T, args []string) (*cli.CLI, error) {
	t.Helper()
	var c cli.CLI
	parser, err := kong.New(&c)
	if err != nil {
		t.Fatalf("kong.New failed: %v", err)
	}
	_, err = parser.Parse(args)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func TestCLI_DefaultLength(t *testing.T) {
	c, err := parseCLI(t, []string{})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Length != 20 {
		t.Errorf("expected default length 20, got %d", c.Length)
	}
}

func TestCLI_CustomLength(t *testing.T) {
	c, err := parseCLI(t, []string{"30"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Length != 30 {
		t.Errorf("expected length 30, got %d", c.Length)
	}
}

func TestCLI_InvalidStringArg(t *testing.T) {
	_, err := parseCLI(t, []string{"abc"})
	if err == nil {
		t.Fatal("expected parse error for 'abc', got nil")
	}
}

func TestCLI_Run_DefaultLength(t *testing.T) {
	c := &cli.CLI{Length: 20}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_CustomLength(t *testing.T) {
	c := &cli.CLI{Length: 30}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 30 {
		t.Errorf("expected 30 chars, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_ZeroLength(t *testing.T) {
	c := &cli.CLI{Length: 0}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error for length 0, got nil")
	}
	if !strings.Contains(err.Error(), generator.ErrLengthNotPositive.Error()) {
		t.Errorf("expected ErrLengthNotPositive, got: %v", err)
	}
}

func TestCLI_Run_NegativeLength(t *testing.T) {
	c := &cli.CLI{Length: -1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error for length -1, got nil")
	}
	if !strings.Contains(err.Error(), generator.ErrLengthNotPositive.Error()) {
		t.Errorf("expected ErrLengthNotPositive, got: %v", err)
	}
}

// Round 1: カテゴリフラグパーステスト

func TestCLI_DefaultCategoryFlags(t *testing.T) {
	c, err := parseCLI(t, []string{})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Symbols != 1 {
		t.Errorf("Symbols = %d, want 1", c.Symbols)
	}
	if c.Digits != 1 {
		t.Errorf("Digits = %d, want 1", c.Digits)
	}
	if c.Upper != 1 {
		t.Errorf("Upper = %d, want 1", c.Upper)
	}
	if c.Lower != 1 {
		t.Errorf("Lower = %d, want 1", c.Lower)
	}
}

func TestCLI_ParseSymbolsFlag(t *testing.T) {
	c, err := parseCLI(t, []string{"--symbols", "3"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Symbols != 3 {
		t.Errorf("Symbols = %d, want 3", c.Symbols)
	}
}

func TestCLI_ParseAllCategoryFlags(t *testing.T) {
	c, err := parseCLI(t, []string{"--lower", "2", "--upper", "3", "--digits", "4", "--symbols", "5"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Lower != 2 {
		t.Errorf("Lower = %d, want 2", c.Lower)
	}
	if c.Upper != 3 {
		t.Errorf("Upper = %d, want 3", c.Upper)
	}
	if c.Digits != 4 {
		t.Errorf("Digits = %d, want 4", c.Digits)
	}
	if c.Symbols != 5 {
		t.Errorf("Symbols = %d, want 5", c.Symbols)
	}
}

// Round 2: Run() カテゴリフラグ統合テスト

func TestCLI_Run_Symbols3_ContainsAtLeast3Symbols(t *testing.T) {
	for i := 0; i < 50; i++ {
		c := &cli.CLI{Length: 20, Symbols: 3, Digits: 1, Upper: 1, Lower: 1}
		var buf bytes.Buffer
		err := c.Run(&buf)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		output := strings.TrimSpace(buf.String())
		count := 0
		for _, ch := range output {
			if strings.ContainsRune(charset.Symbols, ch) {
				count++
			}
		}
		if count < 3 {
			t.Errorf("iteration %d: symbol count = %d, want >= 3, password = %q", i, count, output)
		}
	}
}

func TestCLI_Run_Digits0_NoDigitGuarantee(t *testing.T) {
	c := &cli.CLI{Length: 20, Symbols: 1, Digits: 0, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_AllCategories5_Length20(t *testing.T) {
	c := &cli.CLI{Length: 20, Lower: 5, Upper: 5, Digits: 5, Symbols: 5}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
}

// Round 3: カテゴリフラグエラー系テスト

func TestCLI_Run_CategorySumExceedsLength(t *testing.T) {
	c := &cli.CLI{Length: 20, Lower: 10, Upper: 10, Digits: 10, Symbols: 1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error when category sum exceeds length, got nil")
	}
	if !strings.Contains(err.Error(), generator.ErrRequiredExceedsLength.Error()) {
		t.Errorf("expected ErrRequiredExceedsLength, got: %v", err)
	}
}

func TestCLI_Run_NegativeSymbols(t *testing.T) {
	c := &cli.CLI{Length: 20, Symbols: -1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error for negative symbols, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_NegativeDigits(t *testing.T) {
	c := &cli.CLI{Length: 20, Symbols: 1, Digits: -1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error for negative digits, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_NegativeUpper(t *testing.T) {
	c := &cli.CLI{Length: 20, Symbols: 1, Digits: 1, Upper: -1, Lower: 1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error for negative upper, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_NegativeLower(t *testing.T) {
	c := &cli.CLI{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: -1}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err == nil {
		t.Fatal("expected error for negative lower, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_OutputEndsWithNewline(t *testing.T) {
	c := &cli.CLI{Length: 20}
	var buf bytes.Buffer
	err := c.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("expected output to end with newline, got: %q", output)
	}
}
