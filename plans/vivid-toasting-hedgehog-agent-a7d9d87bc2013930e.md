# passgen ロードマップ

> 仕様書: docs/specs/passgen_SPEC.md
> 方針: TDD必須 (Red -> Green -> Refactor)、マイルストーンは最小単位で分割

---

## ディレクトリ構成（提案）

```
passgen/
├── main.go                     # エントリポイント（cmd.Execute() を呼ぶだけ）
├── go.mod
├── go.sum
├── .gitignore
├── .goreleaser.yaml
├── LICENSE
├── README.md
├── Makefile
├── internal/
│   ├── charset/
│   │   ├── charset.go          # 文字セット定義・操作
│   │   └── charset_test.go
│   ├── generator/
│   │   ├── generator.go        # パスワード生成コアロジック
│   │   └── generator_test.go
│   ├── config/
│   │   ├── config.go           # 設定値の解決（CLI > 環境変数 > デフォルト）
│   │   └── config_test.go
│   ├── clipboard/
│   │   ├── clipboard.go        # クリップボード操作（pbcopy）
│   │   └── clipboard_test.go
│   └── cmd/
│       ├── root.go             # cobra ルートコマンド定義
│       ├── root_test.go
│       ├── completion.go       # zsh 補完サブコマンド
│       └── version.go          # バージョン情報（goreleaser ldflags）
├── docs/
│   └── specs/
│       └── passgen_SPEC.md
├── plans/
├── .github/
│   └── workflows/
│       ├── ci.yaml             # test + lint
│       └── release.yaml        # goreleaser
└── Formula/                    # or separate homebrew-tap repo
```

### 設計判断

**cobra を採用する理由:**
- `passgen completion zsh` サブコマンドが仕様に含まれており、cobra は completion を組み込みで提供する
- 標準 `flag` パッケージでは completion サブコマンドを自前実装する必要があり、工数が不釣り合い
- cobra は Go CLI のデファクトスタンダードであり、学習コスト・メンテナンスコストが低い

**internal/ を使う理由:**
- passgen はライブラリではなく CLI ツール。パッケージを外部公開する必要がない
- internal/ に置くことで API の公開範囲を明示的に制限し、将来の破壊的変更を自由に行える

**パッケージ分割の方針:**
- `charset`: 文字セットの定義と操作。純粋なデータ層で依存ゼロ
- `generator`: パスワード生成ロジック。charset に依存、crypto/rand を使用
- `config`: CLI フラグ・環境変数・デフォルト値の3層マージ。generator に渡す設定を構築
- `clipboard`: OS 依存の副作用を隔離。インターフェース経由でテスト可能に
- `cmd`: cobra コマンド定義。上記パッケージを組み合わせる薄いレイヤー

---

## マイルストーン一覧

| # | マイルストーン | 依存 | 主要成果物 |
|---|---|---|---|
| M01 | プロジェクト初期化 | なし | go.mod, main.go, .gitignore |
| M02 | 文字セット定義 | M01 | internal/charset/ |
| M03 | パスワード生成コアロジック | M02 | internal/generator/ |
| M04 | CLI基盤（cobra + length引数） | M01 | internal/cmd/root.go, main.go |
| M05 | カテゴリフラグ（--symbols等） | M03, M04 | internal/cmd/root.go 拡張 |
| M06 | --exclude フラグ | M02, M05 | charset 拡張, cmd 拡張 |
| M07 | 設定値の解決（環境変数） | M05 | internal/config/ |
| M08 | --no-copy / --no-print フラグ | M04 | cmd バリデーション拡張 |
| M09 | クリップボード連携 | M08 | internal/clipboard/ |
| M10 | エラーハンドリング統合 | M06, M09 | 横断的テスト追加 |
| M11 | シェル補完（zsh） | M04 | internal/cmd/completion.go |
| M12 | goreleaser 設定 | M04 | .goreleaser.yaml |
| M13 | CI/CD（GitHub Actions） | M12 | .github/workflows/ |
| M14 | Homebrew tap 設定 | M12 | goreleaser homebrew section |
| M15 | README・ドキュメント | M14 | README.md |

---

## 各マイルストーン詳細

### M01: プロジェクト初期化

**ゴール:** ビルド・テスト実行可能な空のGoプロジェクトを構築する

**対象ファイル:**
- `go.mod`
- `main.go`（`func main()` のみ）
- `.gitignore`
- `Makefile`（`build`, `test`, `lint` ターゲット）

**作業内容:**
1. `go mod init github.com/youyo/passgen`
2. `main.go` に最小限の `func main()` を配置
3. Go 標準の `.gitignore` を作成
4. `Makefile` に `build`, `test`, `lint` ターゲットを定義
5. `go build ./...` と `go test ./...` が成功することを確認

**テストケース:**
- `go build ./...` がエラーなく完了する（ビルドスモークテスト）

**リスク:**
- 特になし。最も安全なマイルストーン

---

### M02: 文字セット定義

**ゴール:** 4カテゴリの文字セットを定義し、曖昧文字が確実に除外されていることをテストで保証する

**対象ファイル:**
- `internal/charset/charset.go`
- `internal/charset/charset_test.go`

**作業内容:**
1. 各カテゴリの文字列定数を定義
   - `Lower = "abcdefghijkmnopqrstuvwxyz"` （l 除外）
   - `Upper = "ABCDEFGHJKLMNPQRSTUVWXYZ"` （I, O 除外）
   - `Digits = "23456789"` （0, 1 除外）
   - `Symbols = "-_.~"`
2. `All()` 関数: 全カテゴリを結合して返す
3. `Exclude(base, excluded string) string` 関数: base から excluded の文字を除去
4. `Categories()` 関数: 個別カテゴリをスライスで返す

**TDDテストケース:**
- [ ] Lower に `l` が含まれないこと
- [ ] Upper に `I`, `O` が含まれないこと
- [ ] Digits に `0`, `1` が含まれないこと
- [ ] Symbols が `-_.~` の4文字であること
- [ ] All() が全カテゴリの結合であること
- [ ] All() の長さが 25 + 24 + 8 + 4 = 61 であること（仕様の文字数を数えて検証）
- [ ] Exclude("abcdef", "bd") が "acef" を返すこと
- [ ] Exclude で全文字を除外すると空文字列を返すこと

**依存:** M01

**リスク:**
- 仕様書の文字セットに typo がある可能性。テストで文字数を厳密に検証して検出する

---

### M03: パスワード生成コアロジック

**ゴール:** crypto/rand を使用し、各カテゴリ最低保証付きのパスワード生成関数を実装する

**対象ファイル:**
- `internal/generator/generator.go`
- `internal/generator/generator_test.go`

**作業内容:**
1. `Config` 構造体を定義:
   ```go
   type Config struct {
       Length  int
       Lower   int // 最低文字数
       Upper   int
       Digits  int
       Symbols int
       Charset charset.Set // Exclude 適用後の文字セット
   }
   ```
2. `Generate(cfg Config) (string, error)` を実装:
   - 各カテゴリから最低数を crypto/rand で生成
   - 残り (length - 必須合計) を全体文字セットから生成
   - Fisher-Yates シャッフル（crypto/rand ベース）
3. crypto/rand のラッパー関数: `secureRandomIndex(max int) (int, error)`

**TDDテストケース:**
- [ ] デフォルト設定（各1, length=20）で20文字のパスワードが生成される
- [ ] 生成パスワードに lower が最低1文字含まれる（1000回繰り返し統計テスト）
- [ ] 生成パスワードに upper が最低1文字含まれる
- [ ] 生成パスワードに digit が最低1文字含まれる
- [ ] 生成パスワードに symbol が最低1文字含まれる
- [ ] length=4, 各カテゴリ最低1 で、全4カテゴリが必ず含まれる
- [ ] 2回連続生成で異なるパスワードが返る（確率的に同一になる可能性は無視できるほど低い）
- [ ] length=100 で正しい長さのパスワードが生成される
- [ ] 生成パスワードに曖昧文字（l, I, O, 0, 1）が含まれない
- [ ] symbols 最低3を指定した場合、symbol が3文字以上含まれる
- [ ] エラーケース: length <= 0 で error を返す
- [ ] エラーケース: 必須合計 > length で error を返す

**依存:** M02

**リスク:**
- crypto/rand はテストでモックしない方針とする（実際の乱数で統計的に検証）
- シャッフルの均一性を厳密にテストするのは難しいが、偏りの統計的検定は過剰なのでスキップ

---

### M04: CLI基盤（cobra + length引数）

**ゴール:** `passgen [length]` コマンドが動作し、パスワードを stdout に出力する最小限の CLI

**対象ファイル:**
- `internal/cmd/root.go`
- `internal/cmd/root_test.go`
- `main.go`（Execute() 呼び出しに変更）
- `go.sum`（cobra 依存追加）

**作業内容:**
1. `go get github.com/spf13/cobra`
2. cobra でルートコマンドを定義
3. 位置引数 `[length]` をパース（デフォルト 20）
4. generator.Generate() を呼び出して stdout に出力
5. `main.go` は `cmd.Execute()` を呼ぶだけ

**TDDテストケース:**
- [ ] 引数なしで20文字のパスワードが出力される
- [ ] `passgen 30` で30文字のパスワードが出力される
- [ ] `passgen 0` でエラーメッセージが出力される
- [ ] `passgen abc` でエラーメッセージが出力される
- [ ] `passgen -1` でエラーメッセージが出力される
- [ ] `passgen --help` でヘルプが表示される
- [ ] cobra の `cmd.SetOut()` / `cmd.SetErr()` でバッファに出力をキャプチャしてテスト

**依存:** M01（M03 がなくてもスタブで動作確認可能だが、M03 完了後に統合が自然）

**リスク:**
- cobra のバージョン選定。最新安定版を使用する

---

### M05: カテゴリフラグ（--symbols, --digits, --upper, --lower）

**ゴール:** 各カテゴリの最低文字数をフラグで指定可能にする

**対象ファイル:**
- `internal/cmd/root.go`（フラグ追加）
- `internal/cmd/root_test.go`

**作業内容:**
1. `--symbols int` (デフォルト 1), `--digits int`, `--upper int`, `--lower int` フラグを追加
2. フラグ値を `generator.Config` に渡す
3. 負の値のバリデーション

**TDDテストケース:**
- [ ] `--symbols 3` で symbol が3文字以上含まれるパスワードが生成される
- [ ] `--digits 0` で digit を含まないパスワードが生成される可能性がある（最低保証なし）
- [ ] `--lower 5 --upper 5 --digits 5 --symbols 5` (合計20) で length=20 のパスワードが生成される
- [ ] `--lower 10 --upper 10 --digits 10` (合計30) で length=20 のとき、エラーが返る
- [ ] `--symbols -1` でエラーが返る
- [ ] フラグ未指定時はデフォルト値 1 が使用される

**依存:** M03, M04

**リスク:**
- フラグ値が 0 の意味（「そのカテゴリは不要」）を明確に仕様化する必要あり

---

### M06: --exclude フラグ

**ゴール:** 指定文字を文字セットから除外する機能を追加する

**対象ファイル:**
- `internal/charset/charset.go`（Exclude 関数は M02 で実装済み）
- `internal/cmd/root.go`（フラグ追加）
- `internal/cmd/root_test.go`
- `internal/generator/generator_test.go`（統合テスト追加）

**作業内容:**
1. `--exclude string` フラグを追加
2. charset.Exclude() を使って文字セットからフラグ指定文字を除去
3. 除去後の文字セットが空になる場合のエラーハンドリング

**TDDテストケース:**
- [ ] `--exclude "abc"` で生成パスワードに a, b, c が含まれない
- [ ] `--exclude` で全 lower を除外しても、upper/digits/symbols のみで生成可能
- [ ] `--exclude` で全文字を除外するとエラーが返る
- [ ] `--lower 1 --exclude "abcdefghijkmnopqrstuvwxyz"` で lower カテゴリが空になりエラー
- [ ] `--exclude ""` （空文字列）は何も除外しない
- [ ] `--exclude` に曖昧文字（既に除外済み）を指定しても正常動作

**依存:** M02, M05

**リスク:**
- カテゴリ別に除外適用後、そのカテゴリの最低保証が満たせるか検証するロジックが必要

---

### M07: 設定値の解決（環境変数サポート）

**ゴール:** PASSGEN_* 環境変数をサポートし、CLI > 環境変数 > デフォルトの優先順位を実装する

**対象ファイル:**
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/cmd/root.go`（config パッケージを利用するよう変更）

**作業内容:**
1. `config.Resolve()` 関数: CLI フラグ値・環境変数・デフォルト値を3層マージ
2. cobra の `cmd.Flags().Changed("symbols")` で「CLI で明示指定されたか」を判定
3. 環境変数の読み取りとパース（整数変換、エラーハンドリング）

**TDDテストケース:**
- [ ] 環境変数もフラグもなし -> デフォルト値が使用される
- [ ] PASSGEN_LENGTH=30 設定時、引数なしで length=30 が使用される
- [ ] PASSGEN_LENGTH=30 かつ引数 `passgen 40` で length=40 が使用される（CLI 優先）
- [ ] PASSGEN_SYMBOLS=3 で symbols 最低3が使用される
- [ ] PASSGEN_SYMBOLS=3 かつ --symbols 5 で 5 が使用される（CLI 優先）
- [ ] PASSGEN_EXCLUDE="abc" が適用される
- [ ] PASSGEN_LENGTH=abc（不正値）でエラーが返る
- [ ] PASSGEN_LENGTH=-1 でエラーが返る
- [ ] 複数の環境変数を同時設定して正しく解決される
- [ ] t.Setenv() を使用してテスト内で環境変数を安全に設定

**依存:** M05

**リスク:**
- cobra の `Changed()` メソッドの挙動を正確に把握する必要あり
- テストで環境変数を操作するため、並列テストに注意（t.Setenv は Go 1.17+ で自動クリーンアップ）

---

### M08: --no-copy / --no-print フラグ

**ゴール:** 出力制御フラグを追加し、同時指定禁止のバリデーションを実装する

**対象ファイル:**
- `internal/cmd/root.go`（フラグ追加 + バリデーション）
- `internal/cmd/root_test.go`

**作業内容:**
1. `--no-copy` フラグ追加（デフォルト: false = clipboard コピーする）
2. `--no-print` フラグ追加（デフォルト: false = stdout 出力する）
3. 両方 true の場合にエラーを返す `PreRunE` バリデーション
4. `--no-print` 時は stdout への出力を抑制
5. clipboard 連携は M09 で実装するため、この段階では `--no-copy` は「フラグ受付」のみ

**TDDテストケース:**
- [ ] `--no-copy` のみ: 正常動作、stdout に出力される
- [ ] `--no-print` のみ: stdout に出力されない（exit code 0）
- [ ] `--no-copy --no-print` 同時指定: エラーメッセージが表示される
- [ ] フラグなし: デフォルト動作（stdout 出力あり）
- [ ] エラーメッセージが明確で、ユーザーに何が問題かわかる内容であること

**依存:** M04

**リスク:**
- `--no-print` 時にパスワードが一切出力されない（clipboard もまだない）状態が発生する。M09 完了までは `--no-print` 単独使用は実質的に無意味だが、フラグとバリデーションの実装は先行可能

---

### M09: クリップボード連携

**ゴール:** 生成したパスワードを macOS の pbcopy でクリップボードにコピーする

**対象ファイル:**
- `internal/clipboard/clipboard.go`
- `internal/clipboard/clipboard_test.go`
- `internal/cmd/root.go`（clipboard 呼び出し追加）

**作業内容:**
1. `Copier` インターフェースを定義: `Copy(text string) error`
2. `PbcopyCopier` 実装: `exec.Command("pbcopy")` に stdin でテキストを渡す
3. `--no-copy` フラグが true の場合はコピーをスキップ
4. pbcopy が存在しない場合のエラーハンドリング（警告を stderr に出力、処理は続行）

**TDDテストケース:**
- [ ] Copier インターフェースのモック実装で、cmd 統合テストが可能
- [ ] PbcopyCopier.Copy() が pbcopy コマンドを正しく呼び出す（CI では pbcopy がない可能性があるため、統合テストは build tag で分離）
- [ ] `--no-copy` 時に Copier.Copy() が呼ばれない
- [ ] pbcopy が存在しない環境でエラーではなく警告が出る
- [ ] Copy に空文字列を渡してもエラーにならない

**依存:** M08

**リスク:**
- CI 環境（Linux）に pbcopy がない。`//go:build darwin` のビルドタグか、コマンド存在チェックで対応
- テストで実際のクリップボードを汚染しないよう、インターフェースによるモックが必須

---

### M10: エラーハンドリング統合

**ゴール:** 全エラーパスを横断的に検証し、ユーザーに適切なエラーメッセージを提供する

**対象ファイル:**
- `internal/cmd/root_test.go`（統合テスト追加）
- 各パッケージのエラーメッセージ見直し

**作業内容:**
1. 全エラー条件の統合テストを作成
2. エラーメッセージのフォーマット統一（stderr 出力、exit code 1）
3. エッジケースの網羅

**TDDテストケース:**
- [ ] `passgen 0` -> エラー "length must be positive"
- [ ] `passgen 3` (デフォルト最低4) -> エラー "length 3 is too short: minimum required characters is 4"
- [ ] `--exclude` で全文字除外 -> エラー "no characters available after exclusion"
- [ ] `--lower 1 --exclude <全lower文字>` -> エラー "lower charset is empty after exclusion but lower minimum is 1"
- [ ] `--no-copy --no-print` -> エラー "--no-copy and --no-print cannot be used together"
- [ ] PASSGEN_LENGTH=abc -> エラー "invalid PASSGEN_LENGTH value"
- [ ] 全エラーが stderr に出力され、exit code が 1 であること
- [ ] 正常終了時の exit code が 0 であること

**依存:** M06, M09

**リスク:**
- 既存テストとの重複。このマイルストーンは「統合テスト」として、個別パッケージテストでは検証しにくいエンドツーエンドのシナリオに集中する

---

### M11: シェル補完（zsh）

**ゴール:** `passgen completion zsh` コマンドで zsh 補完スクリプトを出力する

**対象ファイル:**
- `internal/cmd/completion.go`
- `internal/cmd/root.go`（サブコマンド登録）

**作業内容:**
1. cobra 組み込みの `GenZshCompletion()` を利用
2. `passgen completion zsh` サブコマンドを追加
3. `--short` フラグで短縮形式の出力（eval 用ワンライナー対応、仕様書に記載あり）

**TDDテストケース:**
- [ ] `passgen completion zsh` が非空の出力を返す
- [ ] 出力に `#compdef passgen` が含まれる（zsh 補完スクリプトの慣例）
- [ ] `passgen completion bash` 等の未サポートシェルでエラーまたは適切なメッセージ
- [ ] exit code 0 で正常終了

**依存:** M04

**リスク:**
- cobra の GenZshCompletion は基本機能のみ。カスタム補完（例: length の候補）が必要な場合は追加作業

---

### M12: goreleaser 設定

**ゴール:** goreleaser でクロスプラットフォームバイナリをビルド・リリースできる設定を作成する

**対象ファイル:**
- `.goreleaser.yaml`
- `internal/cmd/version.go`（ldflags でバージョン埋め込み）

**作業内容:**
1. `.goreleaser.yaml` を作成
   - builds: darwin/arm64（macOS Apple Silicon メイン、他は任意）
   - ldflags: `-s -w -X main.version={{.Version}}`
   - archives 設定
2. `version` 変数を main.go または cmd パッケージに定義
3. `goreleaser check` で設定の妥当性を検証

**TDDテストケース:**
- [ ] `goreleaser check` がエラーなく完了する
- [ ] `goreleaser build --snapshot --clean` でバイナリが生成される
- [ ] 生成バイナリに `--version` フラグでバージョンが表示される（ldflags 検証はスナップショットビルドで確認）

**依存:** M04

**リスク:**
- goreleaser v2 と v1 で設定フォーマットが異なる。v2 を使用する前提で進める

---

### M13: CI/CD（GitHub Actions）

**ゴール:** PR 時にテスト・lint を実行し、tag push 時に goreleaser でリリースする

**対象ファイル:**
- `.github/workflows/ci.yaml`
- `.github/workflows/release.yaml`

**作業内容:**
1. `ci.yaml`: push/PR トリガー、`go test ./...`、`golangci-lint run`
2. `release.yaml`: tag push (`v*`) トリガー、goreleaser action
3. GitHub App token を使用したリリース（仕様書記載）

**TDDテストケース:**
- [ ] CI ワークフローが構文的に正しいこと（`actionlint` で検証）
- [ ] CI でテストが全て pass すること（初回 push で確認）

**依存:** M12

**リスク:**
- GitHub App token の設定は手動作業が必要。ワークフロー内での token 取得方法を明確にする
- golangci-lint のバージョン固定

---

### M14: Homebrew tap 設定

**ゴール:** `brew install youyo/tap/passgen` でインストール可能にする

**対象ファイル:**
- `.goreleaser.yaml`（brews セクション追加）
- 別リポジトリ `homebrew-tap` が必要な場合はその旨を記載

**作業内容:**
1. goreleaser の `brews` セクションに Homebrew tap 設定を追加
2. tap リポジトリ（`youyo/homebrew-tap`）の作成（別途）
3. GitHub App token で tap リポジトリへの Formula push を自動化

**TDDテストケース:**
- [ ] `goreleaser check` が brews 設定を含めてエラーなし
- [ ] リリース後に tap リポジトリに Formula が生成される（手動確認）

**依存:** M12

**リスク:**
- tap リポジトリへの書き込み権限。GitHub App token のスコープ設定が必要

---

### M15: README・ドキュメント

**ゴール:** ユーザーがインストール・使用方法を理解できる README を作成する

**対象ファイル:**
- `README.md`

**作業内容:**
1. プロジェクト概要
2. インストール方法（Homebrew, GitHub Releases, go install）
3. 使用例（基本、フラグ、環境変数）
4. 文字セットの説明
5. zsh 補完の設定方法
6. ライセンス

**テストケース:**
- [ ] README 内のコマンド例が実際に動作する（手動確認）
- [ ] リンクが全て有効である

**依存:** M14（全機能完成後）

**リスク:**
- 特になし

---

## 実装順序の依存グラフ

```
M01 ──┬── M02 ── M03 ──┐
      │                 │
      ├── M04 ──────────┼── M05 ── M06 ──┐
      │    │            │    │            │
      │    │            │    └── M07      │
      │    │            │                 │
      │    ├── M08 ── M09 ────────────────┼── M10
      │    │                              │
      │    ├── M11                        │
      │    │                              │
      │    └── M12 ──┬── M13             │
      │              └── M14             │
      │                                   │
      └───────────────────────────────────┴── M15
```

## クリティカルパス

**M01 -> M02 -> M03 -> M05 -> M06 -> M10 -> M15**

これが最長の依存チェーンであり、プロジェクト全体のボトルネックとなる。M04（CLI基盤）は M02/M03 と並行して着手可能だが、M05 で合流する。

## 並行作業の可能性

- M02（charset）と M04（CLI基盤）は M01 完了後に並行着手可能
- M11（zsh補完）、M12（goreleaser）は M04 完了後、他のフラグ実装と並行可能
- M13（CI）と M14（Homebrew）は M12 完了後に並行可能
- M07（環境変数）と M08（no-copy/no-print）は M05 完了後に並行可能
