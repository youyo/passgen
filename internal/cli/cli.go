// Package cli は passgen の CLI 定義を提供する。
// Kong フレームワークの構造体タグで引数・フラグを宣言する。
package cli

import (
	"fmt"
	"io"

	"github.com/youyo/passgen/internal/generator"
)

// CLI は Kong CLI の構造体定義。
type CLI struct {
	Length int `arg:"" optional:"" default:"20" env:"PASSGEN_LENGTH" help:"パスワードの文字数（デフォルト: 20）"`
}

// Run はパスワードを生成して w に出力する。
// Kong の ctx.Run() から呼び出される。io.Writer はバインディング経由で注入。
func (c *CLI) Run(w io.Writer) error {
	cfg := generator.DefaultConfig()
	cfg.Length = c.Length

	password, err := generator.Generate(cfg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, password)
	return err
}
