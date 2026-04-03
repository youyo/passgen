# M07: --no-copy / --no-print フラグ 実装計画

## 概要
パスワード出力の制御フラグを CLI に追加する。
- `--no-copy`: clipboard コピーを抑制（M08 で clipboard 実装するまではフラグ受付のみ）
- `--no-print`: stdout 出力を抑制
- 両方同時指定はエラー

## 要件
1. `--no-copy` bool フラグ追加（デフォルト: false = clipboard コピーする）
2. `--no-print` bool フラグ追加（デフォルト: false = stdout 出力する）
3. 両方 true の場合にエラー（Kong の `Validate()` hook で実装）
4. `--no-print` 時は stdout への出力を抑制（Run 内で条件分岐）
5. `--no-copy` は clipboard 実装（M08）まではフラグ受付のみ（パース成功、動作影響なし）

## 設計

### CLI 構造体変更
```go
type CLI struct {
    // 既存フィールド...
    NoCopy  bool `name:"no-copy"  default:"false" help:"クリップボードへのコピーを無効化"`
    NoPrint bool `name:"no-print" default:"false" help:"stdout への出力を無効化"`
}
```

### Validate() メソッド（新規）
Kong の `Validate()` hook を利用。パース後、`Run()` の前に自動呼出される。
```go
func (c *CLI) Validate() error {
    if c.NoCopy && c.NoPrint {
        return fmt.Errorf("--no-copy and --no-print cannot be used together")
    }
    return nil
}
```

### Run() メソッド変更
```go
func (c *CLI) Run(w io.Writer) error {
    // ...生成ロジック（変更なし）...
    
    if !c.NoPrint {
        _, err = fmt.Fprintln(w, password)
        if err != nil {
            return err
        }
    }
    // --no-copy は M08 で clipboard 実装時に条件分岐を追加
    return nil
}
```

## TDD テストケース

### Red フェーズ（先に書く失敗テスト）

#### 1. パーステスト
- `TestCLI_ParseNoCopyFlag`: `--no-copy` → NoCopy=true
- `TestCLI_ParseNoPrintFlag`: `--no-print` → NoPrint=true
- `TestCLI_DefaultNoCopyNoPrint`: フラグなし → 両方 false

#### 2. Validate テスト
- `TestCLI_Validate_NoCopyAndNoPrint_Error`: 両方 true → エラー
- `TestCLI_Validate_NoCopyOnly_OK`: NoCopy のみ → nil
- `TestCLI_Validate_NoPrintOnly_OK`: NoPrint のみ → nil
- `TestCLI_Validate_NeitherFlag_OK`: 両方 false → nil

#### 3. Run テスト
- `TestCLI_Run_NoPrint_NoOutput`: NoPrint=true → buf が空
- `TestCLI_Run_NoPrint_ExitCode0`: NoPrint=true → err=nil（正常終了）
- `TestCLI_Run_NoCopy_StdoutOutput`: NoCopy=true → stdout に出力あり（既存動作）
- `TestCLI_Run_Default_StdoutOutput`: フラグなし → stdout に出力あり

#### 4. エラーメッセージテスト
- `TestCLI_Validate_ErrorMessage_Clear`: エラーメッセージに `--no-copy` と `--no-print` が含まれる

## 実装手順

1. **Red**: テストファイルに上記テストケースを追加 → 全て失敗確認
2. **Green**: CLI 構造体にフィールド追加 → Validate() 実装 → Run() 条件分岐 → テスト通過
3. **Refactor**: 既存の負の値バリデーションを Validate() に移動するか検討（スコープ外なら見送り）

## 影響範囲
- `internal/cli/cli.go`: 構造体フィールド追加 + Validate() + Run() 修正
- `internal/cli/cli_test.go`: テスト追加
- `main.go`: 変更なし（Kong が Validate() を自動呼出し）

## リスク
- Kong の `Validate()` hook が期待通りに動作するか → Kong ドキュメントで確認済み、構造体に `Validate() error` メソッドがあればパース後に呼ばれる
- 既存テストへの影響 → NoCopy/NoPrint はデフォルト false なので既存テストは影響なし
