# Roadmap: passgen

## Meta
| 項目 | 値 |
|------|---|
| ゴール | 日常利用に最適化されたシンプルかつ安全なパスワード生成 CLI ツール |
| 成功基準 | `brew install youyo/tap/passgen` でインストール → `passgen` でパスワード生成・クリップボードコピーが動作 |
| 制約 | Go 言語、Kong CLI フレームワーク、crypto/rand、configファイルなし（環境変数のみ） |
| 対象リポジトリ | /Users/youyo/src/github.com/youyo/passgen |
| 作成日 | 2026-04-02 |
| 最終更新 | 2026-04-02 16:50 |
| ステータス | 未着手 |

## Current Focus
- **マイルストーン**: M01 プロジェクト初期化
- **直近の完了**: ロードマップ作成
- **次のアクション**: M01 の実装を `/devflow:implement` で開始

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
- [ ] go mod init
- [ ] main.go（最小限の func main）
- [ ] .gitignore
- [ ] Makefile（build, test, lint）
- [ ] ビルド・テスト実行確認
- 📄 詳細: plans/passgen-m01-project-init.md

### M02: 文字セット定義
- [ ] 4 カテゴリの文字列定数（曖昧文字除外）
- [ ] All() / Exclude() / Categories() 関数
- [ ] TDD テスト（曖昧文字不在、文字数検証）
- 📄 詳細: 着手時に生成

### M03: パスワード生成コアロジック
- [ ] Config 構造体
- [ ] Generate() 関数（crypto/rand）
- [ ] 各カテゴリ最低保証 + シャッフル
- [ ] TDD テスト（統計的検証含む）
- 📄 詳細: 着手時に生成

### M04: CLI 基盤（Kong + length 引数）
- [ ] Kong CLI 構造体定義
- [ ] `passgen [length]` の位置引数パース
- [ ] generator.Generate() 呼び出し → stdout 出力
- [ ] main.go を kong.Parse → Run に変更
- 📄 詳細: 着手時に生成

### M05: カテゴリフラグ（--symbols, --digits, --upper, --lower）
- [ ] Kong 構造体にフラグ追加（各 default:"1"）
- [ ] フラグ値を generator.Config に渡す
- [ ] 負の値バリデーション
- [ ] TDD テスト
- 📄 詳細: 着手時に生成

### M06: --exclude フラグ
- [ ] --exclude string フラグ追加
- [ ] charset.Exclude() 適用
- [ ] 除外後の空集合エラー
- [ ] TDD テスト
- 📄 詳細: 着手時に生成

### M07: --no-copy / --no-print フラグ
- [ ] --no-copy, --no-print フラグ追加
- [ ] 同時指定禁止バリデーション（Validate() メソッド）
- [ ] --no-print 時の stdout 抑制
- [ ] TDD テスト
- 📄 詳細: 着手時に生成

### M08: クリップボード連携
- [ ] Copier インターフェース定義
- [ ] PbcopyCopier 実装（exec.Command）
- [ ] --no-copy 時のスキップ
- [ ] pbcopy 不在時の警告（エラーではない）
- [ ] TDD テスト（モック経由）
- 📄 詳細: 着手時に生成

### M09: エラーハンドリング統合
- [ ] 全エラー条件の統合テスト
- [ ] エラーメッセージフォーマット統一
- [ ] exit code 検証（正常: 0, エラー: 1）
- 📄 詳細: 着手時に生成

### M10: シェル補完（zsh 独自実装）
- [x] `passgen completion zsh` サブコマンド
- [x] zsh completion スクリプト生成
- [x] --short フラグ（eval 用ワンライナー）
- [x] TDD テスト
- 📄 詳細: plans/passgen-m10-completion.md

### M11: goreleaser 設定
- [ ] .goreleaser.yaml 作成
- [ ] ldflags バージョン埋め込み
- [ ] goreleaser check 検証
- [ ] スナップショットビルド確認
- 📄 詳細: 着手時に生成

### M12: CI/CD（GitHub Actions）
- [ ] ci.yaml（test + lint）
- [ ] release.yaml（goreleaser + GitHub App token）
- [ ] actionlint 検証
- 📄 詳細: 着手時に生成

### M13: Homebrew tap 設定
- [ ] goreleaser brews セクション追加
- [ ] tap リポジトリ設定
- [ ] goreleaser check 検証
- 📄 詳細: 着手時に生成

### M14: README・ドキュメント
- [ ] プロジェクト概要
- [ ] インストール方法（Homebrew, go install, GitHub Releases）
- [ ] 使用例（基本、フラグ、環境変数）
- [ ] zsh 補完設定方法
- 📄 詳細: 着手時に生成

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
