# M11: goreleaser 設定 - 実装詳細計画

## 概要
goreleaser v2 形式の設定ファイルを作成し、バージョン情報のビルド時埋め込みと --version フラグ対応を実装する。

## 変更ファイル一覧

### 1. `.goreleaser.yaml`（新規作成）
- v2 形式（`version: 2`）
- builds: darwin/arm64（メイン）、linux/amd64
- ldflags: `-s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}`
- archives: tar.gz（Linux）、zip（Darwin）
- checksum / changelog 設定

### 2. `main.go`（変更）
- `version`, `commit`, `date` パッケージ変数を追加（ldflags で上書き）
- Kong の `kong.Vars` でバージョン情報を渡す

### 3. `internal/cli/cli.go`（変更）
- CLI 構造体に `Version kong.VersionFlag` フィールドを追加

### 4. `Makefile`（変更）
- build ターゲットに ldflags を追加

## 検証方法
1. `goreleaser check` がエラーなし
2. `goreleaser build --snapshot --clean` でバイナリ生成
3. ビルドされたバイナリで `--version` フラグが動作

## TDD
- `--version` フラグのテストは Kong のパース挙動に依存するため、統合テストレベルで確認
