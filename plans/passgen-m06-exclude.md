# M06: --exclude フラグ 詳細実装計画

## 概要
`--exclude` フラグを追加し、ユーザーが指定した文字をパスワード生成の文字セットから除外できるようにする。

## 設計判断

### 方法A（採用）: Config に Exclude フィールドを追加し、Generate 内部で適用
- **理由**: generator が自身の文字セット構築を完全に制御でき、空文字セット検証も内部で完結する
- CLI 層は単にフラグ値を Config に渡すだけ（責務分離が明確）
- テストも generator 単体で exclude の検証が可能

### 方法B（不採用）: CLI 側で Exclude 適用済み文字セットを Config に渡す
- generator の API が複雑化（カスタム文字セット受け取り）
- CLI 層にビジネスロジックが漏れる

## 変更ファイル一覧

### 1. `internal/generator/generator.go`
- Config に `Exclude string` フィールド追加
- Generate() 内で charset.Categories() の各カテゴリに charset.Exclude() 適用
- charset.All() にも Exclude 適用
- 除外後のカテゴリが空で min > 0 の場合のエラー追加
- 除外後の全文字セットが空の場合のエラー追加

### 2. `internal/generator/generator_test.go`
- TDD テストケース追加（後述）

### 3. `internal/cli/cli.go`
- CLI 構造体に `Exclude string` フィールド追加（`default:"" env:"PASSGEN_EXCLUDE"`）
- Run() で cfg.Exclude = c.Exclude を設定

### 4. `internal/cli/cli_test.go`
- --exclude フラグのパーステスト
- Run() 統合テスト追加

## エラー定義

```go
var ErrCategoryEmptyAfterExclude = errors.New("category charset is empty after exclusion")
var ErrAllCharsExcluded = errors.New("all characters excluded")
```

## 実装手順

### Step 1: generator エラー定義（Red）
- `ErrCategoryEmptyAfterExclude` と `ErrAllCharsExcluded` を定義
- テスト: Exclude で全文字除外 → ErrAllCharsExcluded 確認

### Step 2: Config.Exclude フィールド追加（Red → Green）
- Config に Exclude string 追加
- DefaultConfig() は Exclude: "" のまま
- Generate() 内で exclude ロジック実装:
  ```go
  // Exclude 適用
  cats := charset.Categories()
  for i := range cats {
      cats[i] = charset.Exclude(cats[i], cfg.Exclude)
  }
  all := charset.Exclude(charset.All(), cfg.Exclude)
  
  // 検証: min > 0 のカテゴリが空でないか
  catNames := []string{"lower", "upper", "digits", "symbols"}
  for i, min := range minimums {
      if min > 0 && len(cats[i]) == 0 {
          return "", fmt.Errorf("%w: %s has no characters after excluding %q", ErrCategoryEmptyAfterExclude, catNames[i], cfg.Exclude)
      }
  }
  if len(all) == 0 {
      return "", ErrAllCharsExcluded
  }
  ```

### Step 3: CLI フラグ追加（Red → Green）
- CLI 構造体に Exclude フィールド追加
- Run() で cfg.Exclude 設定

### Step 4: Refactor
- コード整理、不要な重複排除

## TDD テストケース

### generator_test.go

1. **TestGenerate_Exclude_abc_NoAbcInPassword**
   - Exclude: "abc", Length: 20
   - 100回試行、生成パスワードに a, b, c が含まれないことを確認

2. **TestGenerate_Exclude_AllLower_WithLower0_Success**
   - Exclude: charset.Lower 全体, Lower: 0, Upper: 1, Digits: 1, Symbols: 1
   - 正常に生成できることを確認

3. **TestGenerate_Exclude_AllLower_WithLower1_Error**
   - Exclude: charset.Lower 全体, Lower: 1
   - ErrCategoryEmptyAfterExclude を返すことを確認

4. **TestGenerate_Exclude_AllChars_AllMinZero_Error**
   - Exclude: charset.All() 全体, Lower: 0, Upper: 0, Digits: 0, Symbols: 0
   - ErrAllCharsExcluded を返すことを確認
   - 注意: min > 0 のカテゴリがあると ErrCategoryEmptyAfterExclude が先にトリガーされる

4b. **TestGenerate_Exclude_AllChars_WithMinimums_CategoryError**
   - Exclude: charset.All() 全体, Lower: 1, Upper: 1, Digits: 1, Symbols: 1
   - ErrCategoryEmptyAfterExclude を返すことを確認（最初の空カテゴリで発火）

5. **TestGenerate_Exclude_Empty_NoEffect**
   - Exclude: ""
   - 通常通り生成されることを確認（デフォルトと同じ）

6. **TestGenerate_Exclude_AlreadyExcludedChars_NoError**
   - Exclude: "lIO01"（曖昧文字＝既に除外済み）
   - 正常に生成されることを確認

7. **TestGenerate_Exclude_PartialCategory_Success**
   - Exclude: "abcdefghijkmnopqrstu" (lower の大部分を除外、vwxyz が残る)
   - Lower: 1 で正常に生成、パスワードに v,w,x,y,z のいずれかが含まれる

### cli_test.go

8. **TestCLI_ParseExcludeFlag**
   - `--exclude "abc"` パース → c.Exclude == "abc" 確認

9. **TestCLI_DefaultExcludeFlag**
   - デフォルト → c.Exclude == "" 確認

10. **TestCLI_Run_Exclude_abc_NoAbcInOutput**
    - Exclude: "abc", Length: 20 で Run()
    - 出力に a, b, c が含まれないことを確認

11. **TestCLI_Run_Exclude_AllChars_Error**
    - Exclude: charset.All() で Run()
    - エラーが返ることを確認

## リスク・注意点

- カテゴリ順序（Lower=0, Upper=1, Digits=2, Symbols=3）は charset.Categories() の実装に依存。テストで検証済みだが、変更があれば影響する
- Exclude 文字列にマルチバイト文字が含まれる場合も charset.Exclude() が rune ベースで処理するため安全
- Symbols が 4 文字のみなので、1 文字でも除外すると残り 3 文字になり、長いパスワードでの symbols 最低保証数が大きい場合に偏りが生じうる（機能的には問題なし）
