package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/youyo/passgen/internal/charset"
)

// testBinary はテスト用にビルドされた passgen バイナリのパス。
var testBinary string

// TestMain はテスト用バイナリを一度だけビルドし、全統合テストで再利用する。
func TestMain(m *testing.M) {
	// テスト用バイナリをビルド
	tmpDir := os.TempDir()
	testBinary = filepath.Join(tmpDir, "passgen_integration_test")
	cmd := exec.Command("go", "build", "-o", testBinary, "../../.")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build test binary: " + err.Error())
	}
	os.Exit(m.Run())
}

// runPassgen はテスト用バイナリを実行し、stdout, stderr, exit code を返す。
func runPassgen(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(testBinary, args...)
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run passgen: %v", err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

// === 正常系 ===

func TestIntegration_Default_ExitCode0(t *testing.T) {
	stdout, stderr, exitCode := runPassgen(t, "--no-copy")
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", exitCode, stderr)
	}
	output := strings.TrimSpace(stdout)
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
	// 正常時にstderrにエラーメッセージがないこと（pbcopy warning は許容）
	if strings.Contains(stderr, "error:") {
		t.Errorf("unexpected error in stderr: %s", stderr)
	}
}

func TestIntegration_CustomLength_ExitCode0(t *testing.T) {
	stdout, stderr, exitCode := runPassgen(t, "--no-copy", "10")
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", exitCode, stderr)
	}
	output := strings.TrimSpace(stdout)
	if len(output) != 10 {
		t.Errorf("expected 10 chars, got %d: %q", len(output), output)
	}
}

func TestIntegration_NoPrint_ExitCode0_NoStdout(t *testing.T) {
	stdout, stderr, exitCode := runPassgen(t, "--no-print")
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", exitCode, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Errorf("expected empty stdout with --no-print, got: %q", stdout)
	}
}

func TestIntegration_NoCopy_ExitCode0(t *testing.T) {
	stdout, _, exitCode := runPassgen(t, "--no-copy")
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0", exitCode)
	}
	output := strings.TrimSpace(stdout)
	if len(output) != 20 {
		t.Errorf("expected 20 chars, got %d: %q", len(output), output)
	}
}

// === エラー系: generator sentinel errors ===

func TestIntegration_Length0_Error(t *testing.T) {
	_, stderr, exitCode := runPassgen(t, "--no-copy", "0")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for length 0")
	}
	if !strings.Contains(stderr, "length must be positive") {
		t.Errorf("stderr should contain 'length must be positive', got: %s", stderr)
	}
}

func TestIntegration_Length3_RequiredExceedsLength(t *testing.T) {
	_, stderr, exitCode := runPassgen(t, "--no-copy", "3")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for length 3 with default minimums")
	}
	if !strings.Contains(stderr, "required minimum characters exceeds length") {
		t.Errorf("stderr should contain 'required minimum characters exceeds length', got: %s", stderr)
	}
}

func TestIntegration_ExcludeAllChars_Error(t *testing.T) {
	allChars := charset.All()
	_, stderr, exitCode := runPassgen(t, "--no-copy", "--exclude", allChars)
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code when all chars excluded")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr should contain error message, got: %s", stderr)
	}
}

func TestIntegration_ExcludeAllLower_CategoryEmpty(t *testing.T) {
	_, stderr, exitCode := runPassgen(t, "--no-copy", "--lower", "1", "--exclude", charset.Lower)
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code when lower category empty after exclude")
	}
	if !strings.Contains(stderr, "category charset is empty after exclusion") {
		t.Errorf("stderr should contain 'category charset is empty after exclusion', got: %s", stderr)
	}
}

// === エラー系: Validate() errors ===

func TestIntegration_NoCopyAndNoPrint_Error(t *testing.T) {
	_, stderr, exitCode := runPassgen(t, "--no-copy", "--no-print")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for --no-copy --no-print")
	}
	if !strings.Contains(stderr, "--no-copy") || !strings.Contains(stderr, "--no-print") {
		t.Errorf("stderr should mention both flags, got: %s", stderr)
	}
}

// === エラー系: パースエラー ===

func TestIntegration_InvalidArg_Error(t *testing.T) {
	_, stderr, exitCode := runPassgen(t, "abc")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for invalid arg 'abc'")
	}
	if !strings.Contains(stderr, "error") {
		t.Errorf("stderr should contain error message, got: %s", stderr)
	}
}

// === stderr出力検証 ===

func TestIntegration_ErrorsGoToStderr(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string // stderr に含まれるべき文字列
	}{
		{
			name: "length 0",
			args: []string{"--no-copy", "0"},
			want: "length must be positive",
		},
		{
			name: "required exceeds length",
			args: []string{"--no-copy", "3"},
			want: "required minimum characters exceeds length",
		},
		{
			name: "no-copy and no-print",
			args: []string{"--no-copy", "--no-print"},
			want: "--no-copy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runPassgen(t, tt.args...)
			if exitCode == 0 {
				t.Fatal("expected non-zero exit code")
			}
			if !strings.Contains(stderr, tt.want) {
				t.Errorf("stderr should contain %q, got: %s", tt.want, stderr)
			}
			// エラー時にstdoutにパスワードが出力されないこと
			// （Kongのusageがstdoutに出る場合があるが、パスワード文字列ではない）
			trimmed := strings.TrimSpace(stdout)
			if trimmed != "" && !strings.Contains(stdout, "Usage:") && !strings.Contains(stdout, "Flags:") {
				t.Errorf("unexpected stdout on error: %q", stdout)
			}
		})
	}
}

// === 正常系のstdout/stderr分離検証 ===

func TestIntegration_Success_NoErrorInStderr(t *testing.T) {
	_, stderr, exitCode := runPassgen(t, "--no-copy")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if strings.Contains(stderr, "error:") {
		t.Errorf("stderr should not contain 'error:' on success, got: %s", stderr)
	}
}

// === M10: completion zsh 統合テスト ===

func TestIntegration_CompletionZsh_ExitCode0(t *testing.T) {
	stdout, stderr, exitCode := runPassgen(t, "completion", "zsh")
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "#compdef passgen") {
		t.Errorf("stdout should contain '#compdef passgen', got:\n%s", stdout)
	}
}

func TestIntegration_CompletionZsh_ContainsAllFlags(t *testing.T) {
	stdout, _, exitCode := runPassgen(t, "completion", "zsh")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	flags := []string{"--symbols", "--digits", "--upper", "--lower", "--exclude", "--no-copy", "--no-print"}
	for _, flag := range flags {
		if !strings.Contains(stdout, flag) {
			t.Errorf("stdout should contain %q", flag)
		}
	}
}

func TestIntegration_CompletionZsh_ContainsCompdef(t *testing.T) {
	stdout, _, exitCode := runPassgen(t, "completion", "zsh")
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout, "compdef _passgen passgen") {
		t.Error("output should contain 'compdef _passgen passgen' for eval usage")
	}
}
