# M08: クリップボード連携 - 実装詳細計画

## 概要
パスワード生成後にクリップボードへ自動コピーする機能を実装する。
macOS の `pbcopy` コマンドを使用し、DI パターンでテスト可能にする。

## アーキテクチャ

### パッケージ構成
```
internal/clipboard/
├── clipboard.go       # Copier インターフェース + PbcopyCopier 実装
└── clipboard_test.go  # ユニットテスト
```

### Copier インターフェース
```go
// Copier はクリップボードへのコピー機能を抽象化するインターフェース。
type Copier interface {
    Copy(text string) error
}
```

### PbcopyCopier 実装
```go
// PbcopyCopier は macOS の pbcopy コマンドを使用した Copier 実装。
type PbcopyCopier struct{}

func (p *PbcopyCopier) Copy(text string) error {
    // exec.LookPath("pbcopy") で存在チェック
    // 不在の場合: stderr に警告を出力し nil を返す（エラーではない）
    // 存在する場合: exec.Command("pbcopy") の stdin に text を書き込み実行
}
```

### NopCopier（--no-copy 用ではなく、DI のフォールバック）
不要。`--no-copy` の場合は CLI.Run() 内で Copier.Copy() を呼ばないだけ。

## 変更対象ファイル

### 1. `internal/clipboard/clipboard.go` (新規)
- `Copier` インターフェース定義
- `PbcopyCopier` 構造体 + `Copy()` メソッド
- `pbcopy` 不在時: `fmt.Fprintf(os.Stderr, "warning: pbcopy not found, clipboard copy skipped\n")` → `return nil`
- `pbcopy` 実行失敗時: `fmt.Fprintf(os.Stderr, "warning: clipboard copy failed: %v\n", err)` → `return nil`
- **重要**: clipboard コピー失敗はすべて警告扱い。パスワード生成自体は成功とする。

### 2. `internal/clipboard/clipboard_test.go` (新規)
- PbcopyCopier の存在チェックロジックのテスト
- 空文字列コピーでエラーにならないことの確認

### 3. `internal/cli/cli.go` (変更)
- `Run()` のシグネチャ: `func (c *CLI) Run(w io.Writer, copier clipboard.Copier) error`
- Kong の BindTo で Copier を注入（io.Writer と同じパターン）
- `!c.NoCopy` の場合に `copier.Copy(password)` を呼び出し

### 4. `internal/cli/cli_test.go` (変更)
- MockCopier 追加
- Run() 呼び出しに MockCopier を渡す
- 新規テストケース追加
- 既存テストの Run() 呼び出しを更新（MockCopier 追加）

### 5. `main.go` (変更)
- `clipboard.PbcopyCopier{}` を生成
- `kong.BindTo(&copier, (*clipboard.Copier)(nil))` で DI

## TDD テストケース

### Red → Green サイクル

#### Cycle 1: Copier インターフェース + MockCopier
- テスト: MockCopier が Copy() を呼べること
- 実装: Copier インターフェース定義

#### Cycle 2: PbcopyCopier.Copy() - 空文字列
- テスト: `PbcopyCopier.Copy("")` がエラーを返さない
- 実装: PbcopyCopier.Copy() の基本実装

#### Cycle 3: CLI.Run() で Copier.Copy() が呼ばれる
- テスト: デフォルト（NoCopy=false）で MockCopier.Copy() が1回呼ばれ、生成されたパスワードが渡される
- 実装: cli.Run() に copier パラメータ追加 + Copy() 呼び出し

#### Cycle 4: --no-copy 時に Copier.Copy() が呼ばれない
- テスト: NoCopy=true で MockCopier.Copy() が呼ばれない
- 実装: `if !c.NoCopy { copier.Copy(password) }` 条件分岐

#### Cycle 5: --no-print 時でも Copier.Copy() が呼ばれる
- テスト: NoPrint=true, NoCopy=false で MockCopier.Copy() が呼ばれる
- 実装: （既に Cycle 3-4 で実装済みのはず）

#### Cycle 6: 既存テストの修正
- 既存の Run() テストすべてに MockCopier を追加して通す

## MockCopier 設計

```go
type MockCopier struct {
    Called    bool
    CalledWith string
    Err      error // テストで任意のエラーを返すために使用
}

func (m *MockCopier) Copy(text string) error {
    m.Called = true
    m.CalledWith = text
    return m.Err
}
```

## DI パターン（main.go）

```go
func main() {
    var c cli.CLI
    copier := &clipboard.PbcopyCopier{}
    ctx := kong.Parse(&c,
        kong.Name("passgen"),
        kong.Description("シンプルかつ安全なパスワード生成 CLI"),
        kong.UsageOnError(),
        kong.BindTo(os.Stdout, (*io.Writer)(nil)),
        kong.BindTo(copier, (*clipboard.Copier)(nil)),
    )
    err := ctx.Run()
    ctx.FatalIfErrorf(err)
}
```

## エラーハンドリング方針
- `pbcopy` 不在 → stderr 警告 + `return nil`
- `pbcopy` 実行失敗 → stderr 警告 + `return nil`
- クリップボードはあくまで利便機能。失敗してもパスワード生成は成功とする。
- **CLI.Run() では Copier.Copy() のエラーは無視しない** → Copier 実装側で警告処理して nil を返す設計

## 実装順序
1. `internal/clipboard/clipboard.go` - インターフェース + 実装
2. `internal/clipboard/clipboard_test.go` - ユニットテスト
3. `internal/cli/cli.go` - Run() 変更
4. `internal/cli/cli_test.go` - MockCopier + テスト更新
5. `main.go` - DI 設定
6. 全テスト実行で green 確認
