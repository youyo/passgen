package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
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
