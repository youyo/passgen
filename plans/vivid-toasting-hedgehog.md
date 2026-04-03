# Plan: passgen ロードマップ作成

## Context

passgen は Go 言語で新規開発するパスワード生成 CLI ツール。仕様書 (`docs/specs/passgen_SPEC.md`) は完成済みだが、コードは一切存在しない（go.mod すらない）完全な新規プロジェクト。

ユーザー要件:
- **TDD 必須**（Red → Green → Refactor）
- **コード品質重視**
- **できる限り細かくマイルストーンを切る**

## 成果物

以下の 2 ファイルを `plans/` に作成する:

1. **`plans/passgen-roadmap.md`** — プロジェクト全体のロードマップ（14 マイルストーン）
2. **`plans/passgen-m01-project-init.md`** — M01 の詳細計画（即座に着手可能）

M02 以降の詳細計画は着手時に遅延生成する。

## アーキテクチャ決定

| # | 決定 | 理由 |
|---|------|------|
| 1 | **Kong 採用**（cobra ではなく） | 構造体タグベースの宣言的 API。env タグで環境変数サポートが組み込み。ユーザー指定 |
| 2 | **zsh completion は独自実装** | passgen のフラグは 7 個と少数。外部ライブラリ不要で zsh script を直接生成する方が依存が少なくシンプル |
| 3 | internal/ パッケージ | CLI ツールなので外部公開不要、将来の自由な変更を担保 |
| 4 | **4 パッケージ分割**（config 不要） | Kong の env タグが 3 層マージ（CLI > 環境変数 > デフォルト）を自動処理するため config パッケージは不要 |

## ディレクトリ構成

```
passgen/
├── main.go                     # エントリポイント（kong.Parse → Run）
├── go.mod / go.sum
├── .gitignore
├── .goreleaser.yaml
├── Makefile
├── README.md
├── internal/
│   ├── charset/                # 文字セット定義・操作（依存ゼロ）
│   │   ├── charset.go
│   │   └── charset_test.go
│   ├── generator/              # パスワード生成コアロジック（crypto/rand）
│   │   ├── generator.go
│   │   └── generator_test.go
│   ├── clipboard/              # pbcopy ラッパー（インターフェース経由）
│   │   ├── clipboard.go
│   │   └── clipboard_test.go
│   └── cli/                    # Kong CLI 定義・サブコマンド
│       ├── cli.go              # CLI 構造体（Kong タグ）
│       ├── cli_test.go
│       ├── completion.go       # zsh completion 独自実装
│       └── completion_test.go
├── docs/specs/
├── plans/
└── .github/workflows/
```

### Kong による CLI 構造体のイメージ

```go
type CLI struct {
    Length  int    `arg:"" optional:"" default:"20" help:"Password length" env:"PASSGEN_LENGTH"`
    Symbols int   `default:"1" help:"Minimum symbols" env:"PASSGEN_SYMBOLS"`
    Digits  int   `default:"1" help:"Minimum digits" env:"PASSGEN_DIGITS"`
    Upper   int   `default:"1" help:"Minimum uppercase" env:"PASSGEN_UPPER"`
    Lower   int   `default:"1" help:"Minimum lowercase" env:"PASSGEN_LOWER"`
    Exclude string `default:"" help:"Characters to exclude" env:"PASSGEN_EXCLUDE"`
    NoCopy  bool  `help:"Don't copy to clipboard"`
    NoPrint bool  `help:"Don't print to stdout"`

    Completion CompletionCmd `cmd:"" help:"Generate shell completion"`
}
```

**Kong のメリット:**
- `env:"PASSGEN_LENGTH"` で環境変数サポートが自動（config パッケージ不要）
- 優先順位 CLI > 環境変数 > デフォルト が Kong 側で自動解決
- 構造体タグのみで宣言的に定義、ボイラープレートが少ない

## 14 マイルストーン一覧

| # | マイルストーン | 依存 | 主要成果物 |
|---|---|---|---|
| M01 | プロジェクト初期化 | なし | go.mod, main.go, .gitignore, Makefile |
| M02 | 文字セット定義 | M01 | internal/charset/ |
| M03 | パスワード生成コアロジック | M02 | internal/generator/ |
| M04 | CLI 基盤（Kong + length 引数） | M01 | internal/cli/cli.go, main.go |
| M05 | カテゴリフラグ（--symbols 等） | M03, M04 | cli 構造体フラグ拡張 |
| M06 | --exclude フラグ | M02, M05 | charset + cli 拡張 |
| M07 | --no-copy / --no-print フラグ | M04 | cli バリデーション |
| M08 | クリップボード連携 | M07 | internal/clipboard/ |
| M09 | エラーハンドリング統合 | M06, M08 | 横断的テスト |
| M10 | シェル補完（zsh 独自実装） | M05 | cli/completion.go |
| M11 | goreleaser 設定 | M04 | .goreleaser.yaml |
| M12 | CI/CD（GitHub Actions） | M11 | .github/workflows/ |
| M13 | Homebrew tap 設定 | M11 | goreleaser brews section |
| M14 | README・ドキュメント | M13 | README.md |

### cobra → Kong 移行で変わった点

1. **M07 を M08 の前に移動**（旧 M07 環境変数サポートを削除）: Kong の env タグで自動解決されるため、環境変数専用マイルストーンが不要に
2. **パッケージ名を cmd → cli に変更**: Kong は cobra の cmd パターンと異なり、構造体定義が中心のため cli が適切
3. **マイルストーン総数が 15 → 14 に削減**: config パッケージ・環境変数マイルストーンが不要に

## 依存グラフ

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

クリティカルパス: **M01 → M02 → M03 → M05 → M06 → M09 → M14**

## 並行作業の可能性

- M02（charset）と M04（CLI 基盤）は M01 完了後に並行着手可能
- M07（no-copy/no-print）と M06（exclude）は M05 完了後に並行可能
- M10（zsh 補完）、M11（goreleaser）は M04/M05 完了後、他と並行可能
- M12（CI）と M13（Homebrew）は M11 完了後に並行可能

## 検証方法

- 各マイルストーン完了時に `go test ./...` が全パスすること
- M04 以降は `go build ./...` でバイナリが生成されること
- M09 で全エラーパスの統合テストが通ること
- M11 で `goreleaser check` がパスすること
- M12 で GitHub Actions が green になること
