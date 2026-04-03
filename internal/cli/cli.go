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
	Length  int `arg:"" optional:"" default:"20" env:"PASSGEN_LENGTH" help:"パスワードの文字数（デフォルト: 20）"`
	Symbols int `short:"s" default:"1" env:"PASSGEN_SYMBOLS" help:"記号の最低文字数（デフォルト: 1）"`
	Digits  int `short:"d" default:"1" env:"PASSGEN_DIGITS" help:"数字の最低文字数（デフォルト: 1）"`
	Upper   int `short:"u" default:"1" env:"PASSGEN_UPPER" help:"大文字の最低文字数（デフォルト: 1）"`
	Lower   int `short:"l" default:"1" env:"PASSGEN_LOWER" help:"小文字の最低文字数（デフォルト: 1）"`
}

// Run はパスワードを生成して w に出力する。
// Kong の ctx.Run() から呼び出される。io.Writer はバインディング経由で注入。
func (c *CLI) Run(w io.Writer) error {
	if c.Symbols < 0 || c.Digits < 0 || c.Upper < 0 || c.Lower < 0 {
		return fmt.Errorf("category minimum values must not be negative")
	}

	cfg := generator.DefaultConfig()
	cfg.Length = c.Length
	cfg.Symbols = c.Symbols
	cfg.Digits = c.Digits
	cfg.Upper = c.Upper
	cfg.Lower = c.Lower

	password, err := generator.Generate(cfg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, password)
	return err
}
