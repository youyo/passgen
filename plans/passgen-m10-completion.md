# M10: シェル補完（zsh 独自実装）詳細計画

## 概要
`passgen completion zsh` サブコマンドを追加し、zsh 補完スクリプトを stdout に出力する。
外部ライブラリ不使用の独自実装。

## 設計

### Kong サブコマンド構造

Kong でサブコマンドを追加するために CLI 構造体を変更する。
現在の CLI はトップレベルコマンドだが、Kong の `default:"1"` タグを使って
デフォルトコマンドとしてパスワード生成を維持しつつ、`completion` サブコマンドを追加する。

**方針**: Kong の構造体ネスト + `cmd:""` タグでサブコマンドを定義。

```go
// 新しいトップレベル構造体
type CLI struct {
    Generate   GenerateCmd   `cmd:"" default:"withargs" help:"パスワードを生成する"`
    Completion CompletionCmd `cmd:"" help:"シェル補完スクリプトを生成する"`
}

// 既存の CLI フィールドを GenerateCmd に移動
type GenerateCmd struct {
    Length  int    `arg:"" optional:"" default:"20" ...`
    Symbols int    `short:"s" default:"1" ...`
    // ... 既存フラグ全て
}

// 補完サブコマンド
type CompletionCmd struct {
    Zsh ZshCompletionCmd `cmd:"" help:"zsh 補完スクリプトを生成する"`
}

type ZshCompletionCmd struct {
    Short bool `help:"eval 用の短縮形式で出力する" default:"false"`
}
```

### zsh 補完スクリプト

`#compdef passgen` で始まる標準的な zsh 補完スクリプトを生成。
対象フラグ:
- `--symbols` / `-s`
- `--digits` / `-d`
- `--upper` / `-u`
- `--lower` / `-l`
- `--exclude` / `-e`
- `--no-copy`
- `--no-print`

位置引数: `[length]` (数値)

### --short フラグ

`--short` 指定時は `.zshrc` で `eval "$(passgen completion zsh --short)"` として使える
ワンライナー形式を出力。

## ファイル構成

1. `internal/cli/cli.go` - CLI 構造体をリファクタリング（GenerateCmd 分離）
2. `internal/cli/completion.go` - CompletionCmd, ZshCompletionCmd, zsh スクリプト生成
3. `internal/cli/completion_test.go` - TDD テスト
4. `main.go` - バインディング調整

## TDD テストケース

### Red Phase（失敗するテスト）
1. `TestZshCompletion_NonEmpty` - 出力が空でない
2. `TestZshCompletion_ContainsCompdef` - `#compdef passgen` を含む
3. `TestZshCompletion_ContainsAllFlags` - 全フラグ名が含まれる
4. `TestZshCompletion_ContainsShortFlags` - 短縮フラグ (-s, -d, -u, -l, -e) が含まれる
5. `TestZshCompletion_Short_OneLiner` - --short で短縮形式
6. `TestZshCompletion_ExitCode0` - エラーなし
7. `TestCompletionCmd_Parse` - `completion zsh` がパース可能
8. `TestCompletionCmd_Parse_Short` - `completion zsh --short` がパース可能
9. `TestGenerateCmd_DefaultBehavior` - 既存機能の回帰テスト

## 実装順序

1. テスト作成（Red）
2. CompletionCmd / ZshCompletionCmd 構造体定義
3. zsh 補完スクリプトテンプレート実装
4. CLI 構造体リファクタリング（GenerateCmd 分離）
5. main.go 調整
6. テスト通過確認（Green）
7. リファクタリング
