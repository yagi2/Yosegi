# Yosegi 🌲

[![CI](https://github.com/yagi2/Yosegi/actions/workflows/ci.yml/badge.svg)](https://github.com/yagi2/Yosegi/actions/workflows/ci.yml)
[![Release](https://github.com/yagi2/Yosegi/actions/workflows/release.yml/badge.svg)](https://github.com/yagi2/Yosegi/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yagi2/yosegi)](https://goreportcard.com/report/github.com/yagi2/yosegi)
[![GitHub release](https://img.shields.io/github/release/yagi2/Yosegi.svg)](https://github.com/yagi2/Yosegi/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

美しいTUIインターフェースを備えたインタラクティブなgit worktree管理ツール

## 概要

Yosegiは、現代の「Vibe Coding」時代のために設計されたクロスプラットフォーム対応CLIツールで、git worktreeの直感的でビジュアルな管理を提供します。`tig`や`peco`のように、複数のgit worktreeを簡単に管理するための優れたビジュアルインターフェースを提供します。

## 機能

- 🎯 **インタラクティブUI**: Bubble TeaとLip Glossで構築された美しいターミナルインターフェース
- 🌲 **Worktree管理**: git worktreeをシームレスに作成、一覧表示、削除
- 🎨 **カスタマイズ可能なテーマ**: YAMLベースの色とUI設定
- ⚡ **キーボードナビゲーション**: Vimスタイルのナビゲーション（j/k）と矢印キー
- 🛡️ **安全機能**: 確認プロンプトと誤削除防止
- 🌍 **クロスプラットフォーム**: Windows、macOS、Linux完全対応
- 🐳 **Docker対応**: マルチアーキテクチャコンテナサポート
- 📦 **軽量**: シングルバイナリ、外部依存なし

## インストール

### Go Install（推奨）

Go 1.24以上がインストールされている場合：

```bash
go install github.com/yagi2/yosegi@latest
```

### プリビルドバイナリ

#### 自動インストール

**Linux/macOS:**
```bash
curl -L https://github.com/yagi2/Yosegi/releases/latest/download/yosegi_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv yosegi /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/yagi2/Yosegi/releases/latest/download/yosegi_Windows_x86_64.zip" -OutFile "yosegi.zip"
Expand-Archive -Path "yosegi.zip" -DestinationPath "."
Move-Item yosegi.exe C:\Windows\System32\
```

#### 手動ダウンロード

[リリースページ](https://github.com/yagi2/Yosegi/releases)から対応するプラットフォーム用のバイナリをダウンロード:

- **Linux**: `yosegi_Linux_x86_64.tar.gz` (AMD64), `yosegi_Linux_arm64.tar.gz` (ARM64)
- **macOS**: `yosegi_Darwin_x86_64.tar.gz` (Intel), `yosegi_Darwin_arm64.tar.gz` (Apple Silicon)
- **Windows**: `yosegi_Windows_x86_64.zip` (AMD64), `yosegi_Windows_arm64.zip` (ARM64)

### Docker

```bash
# 最新版を実行
docker run --rm -it ghcr.io/yagi2/yosegi:latest

# 特定のバージョン
docker run --rm -it ghcr.io/yagi2/yosegi:v1.0.0

# ローカルリポジトリをマウントして実行
docker run --rm -it -v $(pwd):/workspace -w /workspace ghcr.io/yagi2/yosegi:latest
```

### パッケージマネージャー

#### Homebrew（計画中）
```bash
# 近日対応予定
brew install yagi2/tap/yosegi
```

#### Chocolatey（計画中）
```powershell
# 近日対応予定
choco install yosegi
```

#### Scoop（計画中）
```powershell
# 近日対応予定
scoop install yosegi
```

### ソースからビルド

```bash
git clone https://github.com/yagi2/yosegi.git
cd yosegi
go build -o bin/yosegi .

# または開発用タスクランナーを使用
go install github.com/go-task/task/v3/cmd/task@latest
task build
```

## 使い方

### 基本コマンド

#### Worktreeの一覧表示
```bash
yosegi list     # または yosegi ls, yosegi l
```
現在のステータスインジケータ付きの全worktreeのインタラクティブリスト。

#### 新しいWorktreeの作成
```bash
yosegi new [branch]              # インタラクティブな作成
yosegi new feature-branch        # 指定したブランチで作成（ブランチが存在しない場合は自動作成）
yosegi new -b new-feature        # 明示的に新しいブランチとworktreeを作成
yosegi new -p ../feature feature # カスタムパスを指定
```

#### Worktreeの削除
```bash
yosegi remove   # または yosegi rm, yosegi delete
```
確認プロンプト付きの安全な削除。

### 設定

#### 設定の初期化
```bash
yosegi config init
```
`~/.config/yosegi/config.yaml`にデフォルト設定ファイルを作成します。

#### 現在の設定を表示
```bash
yosegi config show
```

### 設定ファイル

`~/.config/yosegi/config.yaml`の例：

```yaml
default_worktree_path: "../"
theme:
  primary: "#7C3AED"
  secondary: "#06B6D4" 
  success: "#10B981"
  warning: "#F59E0B"
  error: "#EF4444"
  muted: "#6B7280"
  text: "#F9FAFB"
git:
  auto_create_branch: true   # ブランチが存在しない場合、自動的に作成
  default_remote: "origin"
  exclude_patterns: []
ui:
  show_icons: true
  confirm_delete: true
  max_path_length: 50
aliases:
  ls: "list"
  rm: "remove"
```

## キーボードナビゲーション

- `↑/k`: 上に移動
- `↓/j`: 下に移動  
- `Enter`: 選択/実行
- `d`: 削除（削除モード時）
- `q`: 終了
- `Tab/Shift+Tab`: 入力フィールドのナビゲート

## 使用例

### 典型的なワークフロー

```bash
# 現在のworktreeを一覧表示（サブコマンド無しでも実行可能）
yosegi
# または
yosegi list

# 機能開発用の新しいworktreeを作成
yosegi new feature/user-auth

# 手動でディレクトリを移動
cd ../feature-user-auth

# 完了したらworktreeを削除
yosegi remove
```

### 高度な使い方

```bash
# カスタムパスと新しいブランチでworktreeを作成
yosegi new -b hotfix/urgent-fix -p ../hotfix

# 強制的にworktreeを削除（確認をスキップ）
yosegi remove --force

# インタラクティブ選択付きでworktreeパスを出力（シェルスクリプトで使用）
# TUIはstderrに表示され、選択結果はstdoutに出力される
yosegi list --print
# または
yosegi ls -p
```

### ディレクトリ移動の統合

Yosegiの`--print`フラグを使用することで、選択したworktreeに簡単に移動できます。このモードでは、TUIがstderrに表示され、選択結果がstdoutに出力されるため、コマンド置換と組み合わせて使用できます。

#### Bashの場合

```bash
# ~/.bashrcに追加
ycd() {
    local worktree=$(yosegi list --print)
    if [ -n "$worktree" ]; then
        cd "$worktree"
    fi
}

# より高度なバージョン（エラーハンドリング付き）
ycd() {
    local worktree
    worktree=$(yosegi list --print 2>/dev/null)
    if [ $? -eq 0 ] && [ -n "$worktree" ]; then
        cd "$worktree"
        echo "Changed to: $worktree"
    else
        echo "No worktree selected or error occurred"
    fi
}
```

#### Zshの場合

```zsh
# ~/.zshrcに追加
ycd() {
    local worktree=$(yosegi list --print)
    if [[ -n $worktree ]]; then
        cd $worktree
    fi
}
```

#### Fishの場合

```fish
# ~/.config/fish/functions/ycd.fishに保存
function ycd
    set worktree (yosegi list --print)
    if test -n "$worktree"
        cd $worktree
    end
end
```

#### ワンライナーでの使用

```bash
# コマンド置換を使用した直接移動（TUIで選択してから移動）
cd $(yosegi list --print)

# 短縮形
cd $(yosegi ls -p)

# 注: コマンド置換なしで`yosegi list`を使用すると、
# 最初の非currentワーキングツリーが自動的に選択されます
cd $(yosegi list)
```

## 開発

### 開発環境のセットアップ

```bash
git clone https://github.com/yagi2/yosegi.git
cd yosegi

# 依存関係のダウンロード
go mod download

# 開発用タスクランナーのインストール
go install github.com/go-task/task/v3/cmd/task@latest

# 利用可能なタスクを確認
task --list-all
```

### ビルド

```bash
# 開発用ビルド
go build -o bin/yosegi .

# リリース用ビルド（最適化済み）
task build-release

# クロスプラットフォームビルド
task build-all

# Taskを使用（推奨）
task build
```

### テスト

```bash
# すべてのテストを実行
go test ./...

# レースコンディション検出付き
go test -race ./...

# カバレッジ計測付き
go test -coverprofile=coverage.out ./...

# Taskを使用（推奨）
task test

# 短縮版テスト（開発時）
task test-short
```

### 品質チェック

```bash
# リンティング
go fmt ./...
go vet ./...

# セキュリティスキャン
gosec ./...

# Taskを使用（推奨）
task lint

# CI環境相当のチェック
task ci
```

### Docker開発

```bash
# Dockerイメージをビルド
docker build -t yosegi:dev .

# またはTaskを使用
task docker-build

# マルチアーキテクチャビルド
task docker-build-multi
```

### 利用可能なTaskコマンド

主要なタスクコマンド（`task --list-all`で全てを確認）:

- `task dev` - 開発モードで実行
- `task build` - 開発用ビルド
- `task build-release` - リリース用最適化ビルド
- `task test` - 全テスト実行（カバレッジ付き）
- `task test-short` - 短縮版テスト
- `task lint` - 全品質チェック実行
- `task clean` - ビルド成果物削除
- `task install` - GOPATH/binにインストール
- `task ci` - CI環境相当のチェック

## コントリビューション

コントリビューションを歓迎します！以下の手順に従ってください：

### コントリビューション手順

1. **リポジトリをフォーク**
2. **開発環境をセットアップ**
   ```bash
   git clone https://github.com/your-username/yosegi.git
   cd yosegi
   go mod download
   task dev
   ```
3. **フィーチャーブランチを作成**
   ```bash
   git checkout -b feature/amazing-feature
   ```
4. **変更を実装し、テストを追加**
   ```bash
   # 変更を実装
   # テストを実行
   task test
   # リンティングを実行
   task lint
   ```
5. **変更をコミット**
   ```bash
   git commit -m 'feat: add amazing feature'
   ```
6. **ブランチにプッシュ**
   ```bash
   git push origin feature/amazing-feature
   ```
7. **プルリクエストを作成**

### 開発ガイドライン

- **コミットメッセージ**: [Conventional Commits](https://conventionalcommits.org/)形式を使用
- **コードスタイル**: `gofmt`と`golangci-lint`に準拠
- **テスト**: 新機能には必ずテストを追加
- **ドキュメント**: READMEと関連ドキュメントを更新
- **Windows対応**: クロスプラットフォーム互換性を維持

### バグレポート・機能要求

- **バグレポート**: [Issues](https://github.com/yagi2/Yosegi/issues)で報告
- **機能要求**: Discussionsで提案
- **セキュリティ**: [SECURITY.md](SECURITY.md)に従って報告

## 動作要件

### 最小要件
- **Go**: 1.24以上（開発時）
- **Git**: 2.25以上（worktree機能サポート）
- **OS**: Windows 10+, macOS 10.15+, Linux (glibc 2.17+)

### 推奨環境
- **ターミナル**: True Color対応（24bit色）
- **フォント**: Nerd Font対応（アイコン表示用）
- **シェル**: Bash, Zsh, Fish, PowerShell

### 対応アーキテクチャ
- **x86_64** (AMD64)
- **ARM64** (Apple Silicon, ARM64)
- **32bit**: ARMv6, ARMv7（Linux）

## ライセンス

MIT License - 詳細は[LICENSE](LICENSE)ファイルを参照してください。

## 謝辞

Yosegiは以下の優れたプロジェクトとコミュニティの影響を受けて開発されました：

### インスピレーション
- **[tig](https://github.com/jonas/tig)**: 美しいGitインターフェースの先駆者
- **[peco](https://github.com/peco/peco)**: インタラクティブ選択の革命
- **[fzf](https://github.com/junegunn/fzf)**: 高速ファジーファインダー

### 技術スタック
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: エレガントなTUIフレームワーク
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)**: 美しいターミナルスタイリング
- **[Bubbles](https://github.com/charmbracelet/bubbles)**: 再利用可能なTUIコンポーネント
- **[Cobra](https://github.com/spf13/cobra)**: 強力なCLIライブラリ
- **[Go](https://golang.org/)**: シンプルで高性能な言語

### 開発ツール
- **[GoReleaser](https://goreleaser.com/)**: 自動化されたリリース管理
- **[GitHub Actions](https://github.com/features/actions)**: CI/CD パイプライン
- **[Task](https://taskfile.dev/)**: モダンなタスクランナー
- **[golangci-lint](https://golangci-lint.run/)**: 包括的なコード解析

### コミュニティ
すべてのコントリビューター、テスター、そしてフィードバックを提供してくださった皆様に感謝いたします。