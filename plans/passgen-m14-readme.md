# M14: README・ドキュメント 実装詳細計画

## 概要
passgen リポジトリの README.md と LICENSE（MIT）ファイルを作成する。
OSSとして公開するため、英語で記述する。

## 作成ファイル

### 1. README.md
以下のセクションを順番に記述:

#### 1.1 プロジェクト概要
- 1行説明: "Simple and secure password generator CLI tool"
- 設計思想を1-2行で補足

#### 1.2 特徴（Features）
- Secure by default（crypto/rand）
- Simple CLI（minimal flags）
- URL-safe characters（`-_.~`）
- Ambiguous characters excluded（`0OoIl1` 等）
- Clipboard integration（auto-copy）

#### 1.3 インストール方法（Installation）
3つの方法を記載:
1. **Homebrew**: `brew install youyo/tap/passgen`
2. **go install**: `go install github.com/youyo/passgen@latest`
3. **GitHub Releases**: バイナリダウンロードリンク

#### 1.4 使用例（Usage）
- 基本: `passgen` → 20文字パスワード生成
- 長さ指定: `passgen 32`
- フラグ例: `passgen --symbols 2 --digits 3`
- --no-copy / --no-print の例
- --exclude の例

#### 1.5 文字セットの説明（Character Set）
テーブル形式で4カテゴリを表示:
| Category | Characters | Count |
|----------|-----------|-------|
| Lower    | abcdefghijkmnopqrstuvwxyz | 23 |
| Upper    | ABCDEFGHJKLMNPQRSTUVWXYZ | 23 |
| Digits   | 23456789 | 8 |
| Symbols  | -_.~ | 4 |

曖昧文字除外の理由を簡潔に説明

#### 1.6 フラグ一覧表（Flags）
cli.go の GenerateCmd 構造体から抽出:
| Flag | Short | Default | Env | Description |
|------|-------|---------|-----|-------------|
| (length) | - | 20 | PASSGEN_LENGTH | Password length |
| --symbols | -s | 1 | PASSGEN_SYMBOLS | Minimum symbol count |
| --digits | -d | 1 | PASSGEN_DIGITS | Minimum digit count |
| --upper | -u | 1 | PASSGEN_UPPER | Minimum uppercase count |
| --lower | -l | 1 | PASSGEN_LOWER | Minimum lowercase count |
| --exclude | -e | "" | PASSGEN_EXCLUDE | Characters to exclude |
| --no-copy | - | false | - | Disable clipboard copy |
| --no-print | - | false | - | Disable stdout output |

#### 1.7 環境変数一覧表（Environment Variables）
優先順位: CLI flag > Environment variable > Default

| Variable | Description | Default |
|----------|-------------|---------|
| PASSGEN_LENGTH | Password length | 20 |
| PASSGEN_SYMBOLS | Minimum symbols | 1 |
| PASSGEN_DIGITS | Minimum digits | 1 |
| PASSGEN_UPPER | Minimum uppercase | 1 |
| PASSGEN_LOWER | Minimum lowercase | 1 |
| PASSGEN_EXCLUDE | Characters to exclude | "" |

#### 1.8 シェル補完（Shell Completion）
zsh のみ対応。設定方法:
```bash
# .zshrc に追加
eval "$(passgen completion zsh --short)"
```

#### 1.9 ライセンス（License）
MIT License — LICENSE ファイルへのリンク

### 2. LICENSE
MIT License テンプレート。Copyright (c) 2026 youyo

## 検証
- README.md が markdown として正しくレンダリングされること
- すべてのフラグ・環境変数が漏れなく記載されていること
- インストール方法が .goreleaser.yaml と整合していること

## コミット
- メッセージ: 「docs: README.md と LICENSE を追加」
