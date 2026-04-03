# M12: CI/CD（GitHub Actions）詳細計画

## 概要

GitHub Actions で CI（テスト・lint）と Release（goreleaser）の 2 ワークフローを構築する。

## ファイル構成

```
.github/
└── workflows/
    ├── ci.yaml        # CI: テスト + lint
    └── release.yaml   # Release: goreleaser + GitHub App token
```

## 1. ci.yaml

### トリガー

```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
```

### ジョブ: test

| ステップ | 内容 |
|----------|------|
| actions/checkout@v4 | コードチェックアウト |
| actions/setup-go@v5 | Go セットアップ（go-version-file: go.mod） |
| go test -v -race ./... | テスト実行 |
| go vet ./... | 静的解析 |

### 設計判断

- Go バージョンは `go-version-file: go.mod` で自動取得（ハードコードしない）
- golangci-lint は含めない（現時点では go vet のみ）
- マトリクスビルドは不要（単一 Go バージョン）
- キャッシュは actions/setup-go が自動管理

## 2. release.yaml

### トリガー

```yaml
on:
  push:
    tags:
      - "v*"
```

### ジョブ: release

| ステップ | 内容 |
|----------|------|
| actions/create-github-app-token@v1 | GitHub App token 生成 |
| actions/checkout@v4 | コードチェックアウト（fetch-depth: 0） |
| actions/setup-go@v5 | Go セットアップ（go-version-file: go.mod） |
| goreleaser/goreleaser-action@v6 | goreleaser でリリース |

### GitHub App Token

```yaml
- uses: actions/create-github-app-token@v1
  id: app-token
  with:
    app-id: ${{ secrets.APP_ID }}
    private-key: ${{ secrets.APP_PRIVATE_KEY }}
```

- `secrets.APP_ID` と `secrets.APP_PRIVATE_KEY` はリポジトリ設定で手動登録
- 生成されたトークンを checkout と goreleaser の両方で使用
- `fetch-depth: 0` は goreleaser の changelog 生成に必要

### goreleaser 設定

```yaml
- uses: goreleaser/goreleaser-action@v6
  with:
    distribution: goreleaser
    version: "~> v2"
    args: release --clean
  env:
    GITHUB_TOKEN: ${{ steps.app-token.outputs.token }}
```

## 3. 検証

- [ ] actionlint で両ワークフローを検証
- [ ] YAML 構文の正当性確認

## 4. 手動作業（ワークフロー外）

- GitHub App の作成と設定
- `APP_ID` シークレットの登録
- `APP_PRIVATE_KEY` シークレットの登録

## TDD 適用

このマイルストーンは YAML 設定ファイルのため、TDD サイクルは適用外。
代わりに actionlint による静的検証を実施する。

## 完了条件

- [x] .github/workflows/ci.yaml 作成
- [x] .github/workflows/release.yaml 作成
- [x] actionlint 検証パス
- [ ] 実際の push/PR で CI が動作（マージ後に確認）
- [ ] タグ push でリリースが動作（手動確認）
