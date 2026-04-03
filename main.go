package main

import (
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/youyo/passgen/internal/cli"
)

func main() {
	var c cli.CLI
	ctx := kong.Parse(&c,
		kong.Name("passgen"),
		kong.Description("シンプルかつ安全なパスワード生成 CLI"),
		kong.UsageOnError(),
		kong.BindTo(os.Stdout, (*io.Writer)(nil)),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
