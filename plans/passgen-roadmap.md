# Roadmap: passgen

## Meta
| 項目 | 値 |
|------|---|
| ゴール | 日常利用に最適化されたシンプルかつ安全なパスワード生成 CLI ツール |
| 成功基準 | `brew install youyo/tap/passgen` でインストール → `passgen` でパスワード生成・クリップボードコピーが動作 |
| 制約 | Go 言語、Kong CLI フレームワーク、crypto/rand、configファイルなし（環境変数のみ） |
| 対象リポジトリ | /Users/youyo/src/github.com/youyo/passgen |
| 作成日 | 2026-04-02 |
| 最終更新 | 2026-04-03 10:00 |
| ステータス | 完了 |

## Current Focus
- **マイルストーン**: 全14マイルストーン完了
- **直近の完了**: M14 README・ドキュメント
- **次のアクション**: GitHub リポジトリへ push → タグ付け → リリース

## Architecture Decisions
| # | 決定 | 理由 | 日付 |
|---|------|------|------|
| 1 | Kong 採用（cobra ではなく） | 構造体タグベースの宣言的 API。env タグで環境変数の 3 層マージが自動。ユーザー指定 | 2026-04-02 |
| 2 | zsh completion は独自実装 | フラグ 7 個と少数。外部ライブラリ不要で zsh script を直接生成する方が依存が少なくシンプル | 2026-04-02 |
| 3 | internal/ パッケージ | CLI ツールなので外部公開不要、将来の自由な変更を担保 | 2026-04-02 |
| 4 | 4 パッケージ分割（charset / generator / clipboard / cli） | Kong の env タグが 3 層マージを自動処理するため config パッケージ不要 | 2026-04-02 |

## Directory Structure
```
passgen/
├── main.go
├── go.mod / go.sum
├── .gitignore
├── .goreleaser.yaml
├── Makefile
├── internal/
│   ├── charset/       # 文字セット定義・操作（依存ゼロ）
│   ├── generator/     # パスワード生成コアロジック（crypto/rand）
│   ├── clipboard/     # pbcopy ラッパー（インターフェース経由）
│   └── cli/           # Kong CLI 定義・サブコマンド・completion
├── docs/specs/
├── plans/
└── .github/workflows/
```

## Progress

### M01: プロジェクト初期化
- [x] go mod init
- [x] main.go（最小限の func main）
- [x] .gitignore
- [x] Makefile（build, test, lint）
- [x] ビルド・テスト実行確認
- 📄 詳細: plans/passgen-m01-project-init.md

### M02: 文字セット定義
- [x] 4 カテゴリの文字列定数（曖昧文字除外）
- [x] All() / Exclude() / Categories() 関数
- [x] TDD テスト（曖昧文字不在、文字数検証）
- 📄 詳細: plans/passgen-m02-charset.md

### M03: パスワード生成コアロジック
- [x] Config 構造体
- [x] Generate() 関数（crypto/rand）
- [x] 各カテゴリ最低保証 + シャッフル
- [x] TDD テスト（統計的検証含む）
- 📄 詳細: plans/passgen-m03-generator.md

### M04: CLI 基盤（Kong + length 引数）
- [x] Kong CLI 構造体定義
- [x] `passgen [length]` の位置引数パース
- [x] generator.Generate() 呼び出し → stdout 出力
- [x] main.go を kong.Parse → Run に変更
- 📄 詳細: plans/passgen-m04-cli-base.md

### M05: カテゴリフラグ（--symbols, --digits, --upper, --lower）
- [x] Kong 構造体にフラグ追加（各 default:"1"）
- [x] フラグ値を generator.Config に渡す
- [x] 負の値バリデーション
- [x] TDD テスト
- 📄 詳細: plans/passgen-m05-category-flags.md

### M06: --exclude フラグ
- [x] --exclude string フラグ追加
- [x] charset.Exclude() 適用
- [x] 除外後の空集合エラー
- [x] TDD テスト
- 📄 詳細: plans/passgen-m06-exclude.md

### M07: --no-copy / --no-print フラグ
- [x] --no-copy, --no-print フラグ追加
- [x] 同時指定禁止バリデーション（Validate() メソッド）
- [x] --no-print 時の stdout 抑制
- [x] TDD テスト
- 📄 詳細: plans/passgen-m07-output-flags.md

### M08: クリップボード連携
- [x] Copier インターフェース定義
- [x] PbcopyCopier 実装（exec.Command）
- [x] --no-copy 時のスキップ
- [x] pbcopy 不在時の警告（エラーではない）
- [x] TDD テスト（モック経由）
- 📄 詳細: plans/passgen-m08-clipboard.md

### M09: エラーハンドリング統合
- [x] 全エラー条件の統合テスト
- [x] エラーメッセージフォーマット統一
- [x] exit code 検証（正常: 0, エラー: 1）
- 📄 詳細: plans/passgen-m09-error-handling.md

### M10: シェル補完（zsh 独自実装）
- [x] `passgen completion zsh` サブコマンド
- [x] zsh completion スクリプト生成
- [x] --short フラグ（eval 用ワンライナー）
- [x] TDD テスト
- 📄 詳細: plans/passgen-m10-completion.md

### M11: goreleaser 設定
- [x] .goreleaser.yaml 作成
- [x] ldflags バージョン埋め込み
- [x] goreleaser check 検証
- [x] スナップショットビルド確認
- 📄 詳細: plans/passgen-m11-goreleaser.md

### M12: CI/CD（GitHub Actions）
- [x] ci.yaml（test + lint）
- [x] release.yaml（goreleaser + GitHub App token）
- [x] actionlint 検証（YAML 構文検証パス、actionlint バイナリはインストール不可のため YAML パーサーで代替）
- 📄 詳細: plans/passgen-m12-cicd.md

### M13: Homebrew tap 設定
- [x] goreleaser homebrew_casks セクション追加（brews は v2.10+ で非推奨）
- [x] tap リポジトリ設定（youyo/homebrew-tap）
- [x] goreleaser check 検証
- 📄 詳細: plans/passgen-m13-homebrew.md

### M14: README・ドキュメント
- [x] プロジェクト概要
- [x] インストール方法（Homebrew, go install, GitHub Releases）
- [x] 使用例（基本、フラグ、環境変数）
- [x] zsh 補完設定方法
- 📄 詳細: plans/passgen-m14-readme.md

## Dependency Graph
```
M01 ──┬── M02 ── M03 ──┐
      │                 │
      ├── M04 ──────────┼── M05 ──┬── M06 ──┐
      │    │            │         │          │
      │    ├── M07 ── M08         └── M10    │
      │    │                                 │
      │    └── M11 ──┬── M12                │
      │              └── M13                 │
      │                                      │
      └──────────────────────────────────────┴── M09 ── M14
```

**クリティカルパス:** M01 → M02 → M03 → M05 → M06 → M09 → M14

## Blockers
なし

## Changelog
| 日時 | 種別 | 内容 |
|------|------|------|
| 2026-04-02 16:50 | 作成 | ロードマップ初版作成。Kong 採用、14 マイルストーン、TDD 必須 |
| 2026-04-03 10:00 | 完了 | 全14マイルストーン完了。devflow:cycle で自律実行 |
