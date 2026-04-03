package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	"github.com/youyo/passgen/internal/cli"
	"github.com/youyo/passgen/internal/clipboard"
)

// version, commit, date は goreleaser の ldflags で埋め込まれる。
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var c cli.CLI
	copier := &clipboard.PbcopyCopier{Stderr: os.Stderr}
	ctx := kong.Parse(&c,
		kong.Name("passgen"),
		kong.Description("シンプルかつ安全なパスワード生成 CLI"),
		kong.UsageOnError(),
		kong.Vars{"version": fmt.Sprintf("%s (%s, %s)", version, commit, date)},
		kong.BindTo(os.Stdout, (*io.Writer)(nil)),
		kong.BindTo(copier, (*clipboard.Copier)(nil)),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
