# M13: Homebrew tap 設定

## 概要
goreleaser の brews セクションを追加し、タグ push 時に youyo/homebrew-tap リポジトリへ Formula を自動 push する。

## 前提条件
- M11: .goreleaser.yaml 作成済み（goreleaser v2 形式）
- M12: release.yaml で GitHub App token（RELEASE_TOKEN）取得済み
- tap リポジトリ youyo/homebrew-tap は手動作成（本タスク対象外）

## 変更対象ファイル

### 1. .goreleaser.yaml
homebrew_casks セクションを追加（goreleaser v2.10+ では brews は非推奨、homebrew_casks が後継）:

```yaml
homebrew_casks:
  - ids:
      - passgen
    binaries:
      - passgen
    repository:
      owner: youyo
      name: homebrew-tap
      token: "{{ .Env.GITHUB_TOKEN }}"
    directory: Casks
    homepage: "https://github.com/youyo/passgen"
    description: "Simple and secure password generator CLI tool"
    license: "MIT"
```

設定ポイント:
- `ids`: archives の id と一致させる（passgen）
- `binaries`: インストールするバイナリ名
- `repository.token`: GitHub App token を環境変数経由で渡す
- `directory`: Casks ディレクトリに配置
- `homepage`, `description`, `license`: Cask のメタデータ

### 2. release.yaml（変更不要）
既に GitHub App token が GITHUB_TOKEN として goreleaser に渡されている。
brews の repository.token に `{{ .Env.GITHUB_TOKEN }}` を指定することで、
同じトークンで tap リポジトリへの push が可能。

## 検証
- `goreleaser check` がエラーなしで通ること

## インストール確認（リリース後）
```bash
brew install youyo/tap/passgen
```

## リスク
- GitHub App に youyo/homebrew-tap リポジトリへの write 権限が必要
- tap リポジトリが存在しない場合、リリース時に Formula push が失敗する（リリース自体は成功）
