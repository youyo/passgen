package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/youyo/passgen/internal/cli"
)

// === TDD Red Phase: zsh 補完テスト ===

func TestZshCompletion_NonEmpty(t *testing.T) {
	cmd := &cli.ZshCompletionCmd{}
	var buf bytes.Buffer
	err := cmd.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output from zsh completion")
	}
}

func TestZshCompletion_ContainsCompdef(t *testing.T) {
	cmd := &cli.ZshCompletionCmd{}
	var buf bytes.Buffer
	err := cmd.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "#compdef passgen") {
		t.Errorf("output should contain '#compdef passgen', got:\n%s", output)
	}
}

func TestZshCompletion_ContainsAllFlags(t *testing.T) {
	cmd := &cli.ZshCompletionCmd{}
	var buf bytes.Buffer
	err := cmd.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()

	flags := []string{
		"--symbols",
		"--digits",
		"--upper",
		"--lower",
		"--exclude",
		"--no-copy",
		"--no-print",
	}
	for _, flag := range flags {
		if !strings.Contains(output, flag) {
			t.Errorf("output should contain %q", flag)
		}
	}
}

func TestZshCompletion_ContainsShortFlags(t *testing.T) {
	cmd := &cli.ZshCompletionCmd{}
	var buf bytes.Buffer
	err := cmd.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()

	shortFlags := []string{"-s", "-d", "-u", "-l", "-e"}
	for _, flag := range shortFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("output should contain short flag %q", flag)
		}
	}
}

func TestZshCompletion_Short_OneLiner(t *testing.T) {
	cmd := &cli.ZshCompletionCmd{Short: true}
	var buf bytes.Buffer
	err := cmd.Run(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	// --short 形式では #compdef を含まない（eval 用ワンライナー）
	if strings.Contains(output, "#compdef passgen") {
		t.Error("--short output should not contain '#compdef passgen'")
	}
	// compdef 関数呼び出しを含むべき
	if !strings.Contains(output, "compdef") {
		t.Error("--short output should contain 'compdef' function call")
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output from zsh completion --short")
	}
}

func TestZshCompletion_ExitCode0(t *testing.T) {
	cmd := &cli.ZshCompletionCmd{}
	var buf bytes.Buffer
	err := cmd.Run(&buf)
	if err != nil {
		t.Errorf("expected no error (exit code 0), got: %v", err)
	}
}
