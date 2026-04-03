# M05: カテゴリフラグ（--symbols, --digits, --upper, --lower）

## 概要
Kong CLI 構造体に 4 つのカテゴリフラグを追加し、各文字カテゴリの最低文字数を制御可能にする。

## 前提条件
- M04 完了: CLI 構造体に Length フィールド、Run(io.Writer) メソッド存在
- generator.Config に Lower/Upper/Digits/Symbols フィールド存在
- generator.DefaultConfig() でデフォルト値（各1）取得可能

## 変更ファイル

### 1. internal/cli/cli.go
Kong 構造体に 4 フラグ追加:

```go
type CLI struct {
    Length  int `arg:"" optional:"" default:"20" env:"PASSGEN_LENGTH" help:"パスワードの文字数（デフォルト: 20）"`
    Symbols int `short:"s" default:"1" env:"PASSGEN_SYMBOLS" help:"記号の最低文字数（デフォルト: 1）"`
    Digits  int `short:"d" default:"1" env:"PASSGEN_DIGITS"  help:"数字の最低文字数（デフォルト: 1）"`
    Upper   int `short:"u" default:"1" env:"PASSGEN_UPPER"   help:"大文字の最低文字数（デフォルト: 1）"`
    Lower   int `short:"l" default:"1" env:"PASSGEN_LOWER"   help:"小文字の最低文字数（デフォルト: 1）"`
}
```

Run() メソッド内でマッピング + バリデーション:

```go
func (c *CLI) Run(w io.Writer) error {
    // 負の値バリデーション
    if c.Symbols < 0 || c.Digits < 0 || c.Upper < 0 || c.Lower < 0 {
        return fmt.Errorf("category minimum values must not be negative")
    }

    cfg := generator.DefaultConfig()
    cfg.Length  = c.Length
    cfg.Symbols = c.Symbols
    cfg.Digits  = c.Digits
    cfg.Upper   = c.Upper
    cfg.Lower   = c.Lower

    password, err := generator.Generate(cfg)
    if err != nil {
        return err
    }

    _, err = fmt.Fprintln(w, password)
    return err
}
```

### 2. internal/cli/cli_test.go
追加テストケース:

## TDD テストケース

### Round 1: パーステスト
1. **TestCLI_DefaultCategoryFlags** - フラグ未指定時はデフォルト1
   - `parseCLI(t, []string{})` → Symbols=1, Digits=1, Upper=1, Lower=1
2. **TestCLI_ParseSymbolsFlag** - `--symbols 3` でパース成功
   - `parseCLI(t, []string{"--symbols", "3"})` → Symbols=3
3. **TestCLI_ParseAllCategoryFlags** - 全フラグ同時指定
   - `parseCLI(t, []string{"--lower", "2", "--upper", "3", "--digits", "4", "--symbols", "5"})` → 各値確認

### Round 2: Run() 統合テスト
4. **TestCLI_Run_Symbols3_ContainsAtLeast3Symbols** - --symbols 3 で symbol 3文字以上（50回ループ）
   - CLI{Length:20, Symbols:3, Digits:1, Upper:1, Lower:1} → Run × 50 → 毎回出力中に記号3文字以上
5. **TestCLI_Run_Digits0_NoDigitGuarantee** - --digits 0 で digit 最低保証なし
   - CLI{Length:20, Symbols:1, Digits:0, Upper:1, Lower:1} → Run → エラーなし（digit不在も許容）
6. **TestCLI_Run_AllCategories5_Length20** - 合計20 = length で成功
   - CLI{Length:20, Lower:5, Upper:5, Digits:5, Symbols:5} → Run → 正常、20文字

### Round 3: エラー系テスト
7. **TestCLI_Run_CategorySumExceedsLength** - 合計 > length でエラー
   - CLI{Length:20, Lower:10, Upper:10, Digits:10, Symbols:1} → Run → ErrRequiredExceedsLength
8. **TestCLI_Run_NegativeSymbols** - --symbols -1 でエラー
   - CLI{Length:20, Symbols:-1, Digits:1, Upper:1, Lower:1} → Run → "must not be negative"
9. **TestCLI_Run_NegativeDigits** - --digits -1 でエラー
10. **TestCLI_Run_NegativeUpper** - --upper -1 でエラー
11. **TestCLI_Run_NegativeLower** - --lower -1 でエラー

## エラーメッセージ
- 負の値: `"category minimum values must not be negative"`
- 合計超過: generator.ErrRequiredExceedsLength（既存）

## 環境変数サポート
Kong の `env` タグにより自動:
- `PASSGEN_SYMBOLS` → Symbols
- `PASSGEN_DIGITS` → Digits
- `PASSGEN_UPPER` → Upper
- `PASSGEN_LOWER` → Lower

## 実装順序
1. cli_test.go に Round 1 テスト追加 → Red
2. cli.go に構造体フィールド追加 → Green
3. cli_test.go に Round 2 テスト追加 → Red
4. cli.go の Run() にマッピング追加 → Green
5. cli_test.go に Round 3 テスト追加 → Red
6. cli.go の Run() にバリデーション追加 → Green
7. Refactor: エラーメッセージを sentinel error 化するか検討
8. go test ./... で全 green 確認

## 変更しないファイル
- internal/generator/generator.go（変更不要、既に Config に全フィールドあり）
- internal/charset/（変更不要）
- main.go（変更不要）
