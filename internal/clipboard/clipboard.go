// Package clipboard はクリップボードへのコピー機能を提供する。
package clipboard

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Copier はクリップボードへのコピー機能を抽象化するインターフェース。
type Copier interface {
	Copy(text string) error
}

// PbcopyCopier は macOS の pbcopy コマンドを使用した Copier 実装。
// pbcopy が不在または実行失敗の場合は stderr に警告を出力し nil を返す。
type PbcopyCopier struct {
	Stderr io.Writer
}

// Copy はテキストを pbcopy 経由でクリップボードにコピーする。
// pbcopy が存在しない場合やコピーに失敗した場合は警告を出力して nil を返す。
func (p *PbcopyCopier) Copy(text string) error {
	path, err := exec.LookPath("pbcopy")
	if err != nil {
		fmt.Fprintf(p.Stderr, "warning: pbcopy not found, clipboard copy skipped\n")
		return nil
	}

	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(p.Stderr, "warning: clipboard copy failed: %v\n", err)
		return nil
	}

	return nil
}
