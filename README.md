# Yosegi 🌲

美しいTUIインターフェースを備えたインタラクティブなgit worktree管理ツール

## 概要

Yosegiは、現代の「Vibe Coding」時代のために設計されたCLIツールで、git worktreeの直感的でビジュアルな管理を提供します。`tig`や`peco`のように、複数のgit worktreeを簡単に管理するための優れたビジュアルインターフェースを提供します。

## 機能

- 🎯 **インタラクティブUI**: Bubble TeaとLip Glossで構築された美しいターミナルインターフェース
- 🌲 **Worktree管理**: git worktreeをシームレスに作成、一覧表示、削除
- 🎨 **カスタマイズ可能なテーマ**: YAMLベースの色とUI設定
- ⚡ **キーボードナビゲーション**: Vimスタイルのナビゲーション（j/k）と矢印キー
- 🛡️ **安全機能**: 確認プロンプトと誤削除防止

## インストール

### Go Install（推奨）

```bash
go install github.com/yagi2/yosegi@latest
```

### プリビルドバイナリ

[リリースページ](https://github.com/yagi2/Yosegi/releases)から対応するプラットフォーム用のバイナリをダウンロード:

```bash
# Linux/macOS (自動取得)
curl -L https://github.com/yagi2/Yosegi/releases/latest/download/yosegi_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv yosegi /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/yagi2/Yosegi/releases/latest/download/yosegi_Windows_x86_64.zip" -OutFile "yosegi.zip"
Expand-Archive -Path "yosegi.zip" -DestinationPath "."
```

### ソースからビルド

```bash
git clone https://github.com/yagi2/yosegi.git
cd yosegi
go build -o bin/yosegi .
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

### ビルド
```bash
# 通常のビルド
go build -o bin/yosegi .

# またはTaskを使用
task build
```

### テスト
```bash
# すべてのテストを実行
go test ./...

# またはTaskを使用
task test
```

### リンティング
```bash
# 手動でリンティング
go fmt ./...
go vet ./...

# またはTaskを使用
task lint
```

## コントリビューション

1. リポジトリをフォーク
2. フィーチャーブランチを作成（`git checkout -b feature/amazing-feature`）
3. 変更をコミット（`git commit -m 'Add amazing feature'`）
4. ブランチにプッシュ（`git push origin feature/amazing-feature`）
5. プルリクエストを作成

## 動作要件

- Go 1.24以上
- worktree機能をサポートするGit
- カラー表示対応のターミナル

## ライセンス

MIT License - 詳細は[LICENSE](LICENSE)ファイルを参照してください。

## 謝辞

- `tig`や`peco`などの優れたビジュアルインターフェースにインスパイアされました
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)と[Lip Gloss](https://github.com/charmbracelet/lipgloss)で構築
- CLIフレームワークに[Cobra](https://github.com/spf13/cobra)を使用