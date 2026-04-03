package cli_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/youyo/passgen/internal/charset"
	"github.com/youyo/passgen/internal/cli"
	"github.com/youyo/passgen/internal/clipboard"
	"github.com/youyo/passgen/internal/generator"
)

// MockCopier はテスト用の Copier 実装。
type MockCopier struct {
	Called     bool
	CalledWith string
	Err        error
}

func (m *MockCopier) Copy(text string) error {
	m.Called = true
	m.CalledWith = text
	return m.Err
}

// インターフェース適合性コンパイル時チェック
var _ clipboard.Copier = &MockCopier{}

// parseCLI はテスト用にCLI引数をパースするヘルパー。
// パース成功時に *cli.GenerateCmd を返す。パースエラー時は error を返す。
func parseCLI(t *testing.T, args []string) (*cli.GenerateCmd, error) {
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
	return &c.Generate, nil
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
	c := &cli.GenerateCmd{Length: 20}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_CustomLength(t *testing.T) {
	c := &cli.GenerateCmd{Length: 30}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 30 {
		t.Errorf("expected 30 chars, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_ZeroLength(t *testing.T) {
	c := &cli.GenerateCmd{Length: 0}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error for length 0, got nil")
	}
	if !strings.Contains(err.Error(), generator.ErrLengthNotPositive.Error()) {
		t.Errorf("expected ErrLengthNotPositive, got: %v", err)
	}
}

func TestCLI_Run_NegativeLength(t *testing.T) {
	c := &cli.GenerateCmd{Length: -1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
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
		c := &cli.GenerateCmd{Length: 20, Symbols: 3, Digits: 1, Upper: 1, Lower: 1}
		var buf bytes.Buffer
		mock := &MockCopier{}
		err := c.Run(&buf, mock)
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
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 0, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_AllCategories5_Length20(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Lower: 5, Upper: 5, Digits: 5, Symbols: 5}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
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
	c := &cli.GenerateCmd{Length: 20, Lower: 10, Upper: 10, Digits: 10, Symbols: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error when category sum exceeds length, got nil")
	}
	if !strings.Contains(err.Error(), generator.ErrRequiredExceedsLength.Error()) {
		t.Errorf("expected ErrRequiredExceedsLength, got: %v", err)
	}
}

func TestCLI_Run_NegativeSymbols(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: -1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error for negative symbols, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_NegativeDigits(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: -1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error for negative digits, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_NegativeUpper(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: -1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error for negative upper, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

func TestCLI_Run_NegativeLower(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: -1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error for negative lower, got nil")
	}
	if !strings.Contains(err.Error(), "must not be negative") {
		t.Errorf("expected 'must not be negative' error, got: %v", err)
	}
}

// Round 4: Exclude フラグテスト

func TestCLI_ParseExcludeFlag(t *testing.T) {
	c, err := parseCLI(t, []string{"--exclude", "abc"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Exclude != "abc" {
		t.Errorf("Exclude = %q, want %q", c.Exclude, "abc")
	}
}

func TestCLI_DefaultExcludeFlag(t *testing.T) {
	c, err := parseCLI(t, []string{})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.Exclude != "" {
		t.Errorf("Exclude = %q, want %q", c.Exclude, "")
	}
}

func TestCLI_Run_Exclude_abc_NoAbcInOutput(t *testing.T) {
	for i := 0; i < 50; i++ {
		c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: 1, Exclude: "abc"}
		var buf bytes.Buffer
		mock := &MockCopier{}
		err := c.Run(&buf, mock)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		output := strings.TrimSpace(buf.String())
		for _, ch := range "abc" {
			if strings.ContainsRune(output, ch) {
				t.Errorf("iteration %d: excluded char %q found in %q", i, ch, output)
			}
		}
	}
}

func TestCLI_Run_Exclude_AllChars_Error(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: 1, Exclude: charset.All()}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error when all chars excluded, got nil")
	}
}

func TestCLI_Run_OutputEndsWithNewline(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("expected output to end with newline, got: %q", output)
	}
}

// Round 5: --no-copy / --no-print フラグテスト

func TestCLI_ParseNoCopyFlag(t *testing.T) {
	c, err := parseCLI(t, []string{"--no-copy"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if !c.NoCopy {
		t.Error("NoCopy = false, want true")
	}
}

func TestCLI_ParseNoPrintFlag(t *testing.T) {
	c, err := parseCLI(t, []string{"--no-print"})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if !c.NoPrint {
		t.Error("NoPrint = false, want true")
	}
}

func TestCLI_DefaultNoCopyNoPrint(t *testing.T) {
	c, err := parseCLI(t, []string{})
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if c.NoCopy {
		t.Error("NoCopy = true, want false")
	}
	if c.NoPrint {
		t.Error("NoPrint = true, want false")
	}
}

func TestCLI_Validate_NoCopyAndNoPrint_Error(t *testing.T) {
	c := &cli.GenerateCmd{NoCopy: true, NoPrint: true}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected error when both --no-copy and --no-print, got nil")
	}
	if !strings.Contains(err.Error(), "--no-copy") || !strings.Contains(err.Error(), "--no-print") {
		t.Errorf("error message should mention both flags, got: %v", err)
	}
}

func TestCLI_Validate_NoCopyOnly_OK(t *testing.T) {
	c := &cli.GenerateCmd{NoCopy: true, NoPrint: false}
	err := c.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCLI_Validate_NoPrintOnly_OK(t *testing.T) {
	c := &cli.GenerateCmd{NoCopy: false, NoPrint: true}
	err := c.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCLI_Validate_NeitherFlag_OK(t *testing.T) {
	c := &cli.GenerateCmd{}
	err := c.Validate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCLI_Run_NoPrint_NoOutput(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, NoPrint: true, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output with --no-print, got: %q", buf.String())
	}
}

func TestCLI_Run_NoPrint_ExitCode0(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, NoPrint: true, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("expected no error (exit code 0) with --no-print, got: %v", err)
	}
}

func TestCLI_Run_NoCopy_StdoutOutput(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, NoCopy: true, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars output with --no-copy, got %d: %q", len(output), output)
	}
}

func TestCLI_Run_Default_StdoutOutput(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	if len(output) != 20 {
		t.Errorf("expected 20 chars output by default, got %d: %q", len(output), output)
	}
}

// Round 6: M08 クリップボード連携テスト

func TestCLI_Run_Default_CopierCalled(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.Called {
		t.Error("expected Copier.Copy() to be called by default")
	}
	// コピーされたテキストがstdout出力と一致すること
	output := strings.TrimSpace(buf.String())
	if mock.CalledWith != output {
		t.Errorf("Copier.Copy() received %q, want %q", mock.CalledWith, output)
	}
}

func TestCLI_Run_NoCopy_CopierNotCalled(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, NoCopy: true, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.Called {
		t.Error("expected Copier.Copy() NOT to be called with --no-copy")
	}
}

func TestCLI_Run_NoPrint_CopierCalled(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, NoPrint: true, NoCopy: false, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{}
	err := c.Run(&buf, mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.Called {
		t.Error("expected Copier.Copy() to be called with --no-print (but not --no-copy)")
	}
	// パスワードの長さが正しいこと
	if len(mock.CalledWith) != 20 {
		t.Errorf("Copier.Copy() received password of length %d, want 20", len(mock.CalledWith))
	}
}

func TestCLI_Run_CopierError_Propagated(t *testing.T) {
	c := &cli.GenerateCmd{Length: 20, Symbols: 1, Digits: 1, Upper: 1, Lower: 1}
	var buf bytes.Buffer
	mock := &MockCopier{Err: fmt.Errorf("clipboard error")}
	err := c.Run(&buf, mock)
	if err == nil {
		t.Fatal("expected error from Copier to propagate, got nil")
	}
	if !strings.Contains(err.Error(), "clipboard error") {
		t.Errorf("expected 'clipboard error', got: %v", err)
	}
}
