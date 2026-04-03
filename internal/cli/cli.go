// Package cli は passgen の CLI 定義を提供する。
// Kong フレームワークの構造体タグで引数・フラグを宣言する。
package cli

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	"github.com/youyo/passgen/internal/clipboard"
	"github.com/youyo/passgen/internal/generator"
)

// CLI は Kong CLI のトップレベル構造体定義。
type CLI struct {
	Generate   GenerateCmd      `cmd:"" default:"withargs" help:"パスワードを生成する"`
	Completion CompletionCmd    `cmd:"" help:"シェル補完スクリプトを生成する"`
	Version    kong.VersionFlag `name:"version" help:"バージョン情報を表示する"`
}

// GenerateCmd はパスワード生成コマンド。デフォルトコマンドとして動作する。
type GenerateCmd struct {
	Length  int    `arg:"" optional:"" default:"20" env:"PASSGEN_LENGTH" help:"パスワードの文字数（デフォルト: 20）"`
	Symbols int    `short:"s" default:"1" env:"PASSGEN_SYMBOLS" help:"記号の最低文字数（デフォルト: 1）"`
	Digits  int    `short:"d" default:"1" env:"PASSGEN_DIGITS" help:"数字の最低文字数（デフォルト: 1）"`
	Upper   int    `short:"u" default:"1" env:"PASSGEN_UPPER" help:"大文字の最低文字数（デフォルト: 1）"`
	Lower   int    `short:"l" default:"1" env:"PASSGEN_LOWER" help:"小文字の最低文字数（デフォルト: 1）"`
	Exclude string `short:"e" default:"" env:"PASSGEN_EXCLUDE" help:"除外する文字"`
	NoCopy  bool   `name:"no-copy" default:"false" help:"クリップボードへのコピーを無効化"`
	NoPrint bool   `name:"no-print" default:"false" help:"stdout への出力を無効化"`
}

// Validate は Kong のバリデーション hook。パース後、Run() の前に呼び出される。
func (c *GenerateCmd) Validate() error {
	if c.NoCopy && c.NoPrint {
		return fmt.Errorf("--no-copy and --no-print cannot be used together")
	}
	return nil
}

// Run はパスワードを生成して w に出力する。
// Kong の ctx.Run() から呼び出される。io.Writer と clipboard.Copier はバインディング経由で注入。
func (c *GenerateCmd) Run(w io.Writer, copier clipboard.Copier) error {
	if c.Symbols < 0 || c.Digits < 0 || c.Upper < 0 || c.Lower < 0 {
		return fmt.Errorf("category minimum values must not be negative")
	}

	cfg := generator.DefaultConfig()
	cfg.Length = c.Length
	cfg.Symbols = c.Symbols
	cfg.Digits = c.Digits
	cfg.Upper = c.Upper
	cfg.Lower = c.Lower
	cfg.Exclude = c.Exclude

	password, err := generator.Generate(cfg)
	if err != nil {
		return err
	}

	if !c.NoPrint {
		_, err = fmt.Fprintln(w, password)
		if err != nil {
			return err
		}
	}

	if !c.NoCopy {
		if err := copier.Copy(password); err != nil {
			return err
		}
	}

	return nil
}
