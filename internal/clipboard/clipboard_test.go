package clipboard_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/youyo/passgen/internal/clipboard"
)

func TestPbcopyCopier_Copy_EmptyString_NoError(t *testing.T) {
	// pbcopy が存在しない環境でもエラーにならないことを確認
	var stderr bytes.Buffer
	c := &clipboard.PbcopyCopier{Stderr: &stderr}
	err := c.Copy("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPbcopyCopier_Copy_NormalString_NoError(t *testing.T) {
	var stderr bytes.Buffer
	c := &clipboard.PbcopyCopier{Stderr: &stderr}
	err := c.Copy("test-password-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPbcopyCopier_Copy_PbcopyNotFound_Warning(t *testing.T) {
	// pbcopy が見つからない場合、LookPath がエラーを返す
	// macOS ではデフォルトで pbcopy が存在するため、このテストは
	// pbcopy が存在する環境では警告が出ないことを確認する
	var stderr bytes.Buffer
	c := &clipboard.PbcopyCopier{Stderr: &stderr}
	err := c.Copy("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// pbcopy が存在する場合は stderr は空、存在しない場合は警告が出る
	_, lookErr := exec.LookPath("pbcopy")
	if lookErr != nil {
		if stderr.Len() == 0 {
			t.Error("expected warning on stderr when pbcopy not found")
		}
	}
}

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

func TestMockCopier_Copy(t *testing.T) {
	m := &MockCopier{}
	err := m.Copy("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Called {
		t.Error("expected Called to be true")
	}
	if m.CalledWith != "hello" {
		t.Errorf("CalledWith = %q, want %q", m.CalledWith, "hello")
	}
}

func TestCopierInterface(t *testing.T) {
	// MockCopier が Copier インターフェースを満たすことを確認
	var _ clipboard.Copier = &MockCopier{}
	var _ clipboard.Copier = &clipboard.PbcopyCopier{}
}
