# M09: エラーハンドリング統合 - 実装計画

## 目的
全エラー条件をCLIレベルでエンドツーエンドに検証する統合テストを作成する。
個別パッケージテスト（generator_test, cli_test）では検証しにくい、
Kong パース → Validate → Run → エラー出力 → exit code の一連のフローを統合的にテストする。

## 現状分析

### 既存エラーハンドリング
1. **generator パッケージ**: 4つのsentinel error定義済み
   - `ErrLengthNotPositive`: length <= 0
   - `ErrRequiredExceedsLength`: カテゴリ最低数合計 > length
   - `ErrCategoryEmptyAfterExclude`: exclude後にカテゴリ空
   - `ErrAllCharsExcluded`: exclude後に全文字空
2. **cli.Validate()**: `--no-copy` + `--no-print` 同時指定エラー
3. **cli.Run()**: 負の値バリデーション（`must not be negative`）
4. **main.go**: `ctx.FatalIfErrorf(err)` でエラー時にstderrに出力+exit(1)

### 不足している検証
- CLIバイナリレベルでの exit code 検証
- stderr出力の検証（エラーメッセージがstderrに出力されること）
- Kong パースエラーの統合検証（不正引数）

## 実装方針

### テストファイル
`internal/cli/cli_integration_test.go` に統合テストを追加する。

### テスト手法
`os/exec` で `go run` またはビルド済みバイナリを実行し、stdout/stderr/exit code を検証する。
テスト用にヘルパー関数 `runPassgen(t, args...)` を作成し、(stdout, stderr, exitCode) を返す。

### テストケース（TDD Red → Green → Refactor）

#### 正常系
1. `passgen` → exit code 0, stdout にパスワード出力
2. `passgen 10` → exit code 0, stdout に10文字出力
3. `passgen --no-copy` → exit code 0
4. `passgen --no-print` → exit code 0, stdout 空

#### エラー系
5. `passgen 0` → exit code 1, stderr に "length must be positive"
6. `passgen -1` → Kong パースエラー（`-1`はフラグとして解釈される可能性）→ exit code 1
7. `passgen 3` → exit code 1, stderr に "required minimum characters exceeds length"
8. `passgen --exclude <全文字>` → exit code 1, stderr にエラーメッセージ
9. `passgen --lower 1 --exclude <全lower文字>` → exit code 1, stderr に "category charset is empty"
10. `passgen --no-copy --no-print` → exit code 1, stderr に "--no-copy" と "--no-print"
11. `passgen abc` → exit code 1（パースエラー）

#### stderr検証
12. 正常時にstderrにエラーメッセージがないこと
13. エラー時にstderrにエラーメッセージがあること

## 実装ステップ

### Step 1: テストヘルパー作成
```go
func buildPassgen(t *testing.T) string    // テスト用バイナリビルド
func runPassgen(t *testing.T, args ...string) (stdout, stderr string, exitCode int)
```

### Step 2: Red - 失敗テスト作成
全テストケースを作成（まだ通る前提だが、統合テストなので既存コードで通るはず）

### Step 3: Green - 必要に応じてcli.goのエラーメッセージフォーマット調整

### Step 4: Refactor - テーブルドリブンテストに整理

## ファイル変更一覧
| ファイル | 変更内容 |
|---------|---------|
| `internal/cli/cli_integration_test.go` | 新規: 統合テスト |
| `internal/cli/cli.go` | 変更なし（既存実装で十分） |

## リスク
- `go run` のテスト速度が遅い → `TestMain` でバイナリをビルドして再利用
- CI環境でのpbcopy不在 → `--no-copy` フラグで回避
- `-1` 等の負数がKongのフラグパーサーに干渉 → テストで挙動を確認

## 完了条件
- [ ] 全テストケースがGreen
- [ ] `go test ./...` が全パス
- [ ] エラー時は必ずexit code 1
- [ ] エラーメッセージは必ずstderrに出力
