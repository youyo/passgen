package cli

import (
	"fmt"
	"io"
)

// CompletionCmd はシェル補完スクリプト生成のサブコマンド。
type CompletionCmd struct {
	Zsh ZshCompletionCmd `cmd:"" help:"zsh 補完スクリプトを生成する"`
}

// ZshCompletionCmd は zsh 補完スクリプトを生成するサブコマンド。
type ZshCompletionCmd struct {
	Short bool `help:"eval 用の短縮形式で出力する" default:"false"`
}

// Run は zsh 補完スクリプトを w に出力する。
func (c *ZshCompletionCmd) Run(w io.Writer) error {
	if c.Short {
		return c.runShort(w)
	}
	return c.runFull(w)
}

// runFull は #compdef passgen で始まる完全な zsh 補完スクリプトを出力する。
func (c *ZshCompletionCmd) runFull(w io.Writer) error {
	_, err := fmt.Fprint(w, zshCompletionScript)
	return err
}

// runShort は eval 用のワンライナー形式を出力する。
func (c *ZshCompletionCmd) runShort(w io.Writer) error {
	_, err := fmt.Fprint(w, zshCompletionShort)
	return err
}

const zshCompletionScript = `#compdef passgen

_passgen() {
    local -a flags
    flags=(
        '(-s --symbols)'{-s,--symbols}'[記号の最低文字数（デフォルト: 1）]:number:'
        '(-d --digits)'{-d,--digits}'[数字の最低文字数（デフォルト: 1）]:number:'
        '(-u --upper)'{-u,--upper}'[大文字の最低文字数（デフォルト: 1）]:number:'
        '(-l --lower)'{-l,--lower}'[小文字の最低文字数（デフォルト: 1）]:number:'
        '(-e --exclude)'{-e,--exclude}'[除外する文字]:string:'
        '--no-copy[クリップボードへのコピーを無効化]'
        '--no-print[stdout への出力を無効化]'
        '--help[ヘルプを表示する]'
    )

    _arguments -s \
        "${flags[@]}" \
        '1::length:'
}
`

const zshCompletionShort = `_passgen() { local -a flags; flags=('(-s --symbols)'{-s,--symbols}'[記号の最低文字数]:number:' '(-d --digits)'{-d,--digits}'[数字の最低文字数]:number:' '(-u --upper)'{-u,--upper}'[大文字の最低文字数]:number:' '(-l --lower)'{-l,--lower}'[小文字の最低文字数]:number:' '(-e --exclude)'{-e,--exclude}'[除外する文字]:string:' '--no-copy[クリップボードへのコピーを無効化]' '--no-print[stdout への出力を無効化]' '--help[ヘルプを表示する]'); _arguments -s "${flags[@]}" '1::length:'; }; compdef _passgen passgen
`
