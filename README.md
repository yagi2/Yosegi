# CLI Vibe Go Template
GitHub Codespaces + Go 1.24 + Claude Codeで動く、CLIツール開発用テンプレートです。
起動直後から `main.go` を編集してすぐ実行できます。

## 特徴

- **GitHub Codespaces**: クラウド開発環境ですぐに開始
- **Go 1.24**: 最新のGo言語環境
- **Claude Code**: AI支援コーディング環境を内蔵
- **VS Code拡張**: Go開発に最適化された拡張機能セット

## 使い方

### 基本的なGo開発

```bash
go run main.go arg1 arg2
```

### Claude Codeを使用したAI支援開発

```bash
# 対話的なコーディング支援
claude

# 一回だけの質問
claude -p "Go でCLIツールを作る方法を教えて"

# 前回の会話を継続
claude -c
```

## セットアップ

### 1. 前提条件
- [Claude MAX](https://claude.ai/) サブスクリプション契約
- API キーは不要（サブスクリプションを直接使用）

### 2. Codespace作成・認証
1. このリポジトリで "Code" > "Codespaces" > "Create codespace on main"
2. 環境構築を待つ（数分程度、Claude Codeが自動インストールされます）
3. ターミナルで認証: `claude /login`

## 詳細なセットアップ手順

### Claude Codeの認証

Codespaceが起動したら、Claude MAXサブスクリプションで認証します：

```bash
claude /login
```

これにより以下が実行されます：
1. ブラウザウィンドウが開くか、URLが提供されます
2. Claude.aiでの認証画面にリダイレクトされます
3. 既存のClaude MAXサブスクリプションが使用されます
4. 認証情報がCodespace内にローカル保存されます

### インストール確認

Claude Codeが正常に動作しているかテストします：

```bash
claude --version
```

### 認証に関する注意点

- **API Key不要**: Claude MAXサブスクリプションを直接使用
- **一回設定**: 認証はCodespace内で永続化されます
- **アカウント切り替え**: `claude /login` で別アカウントに切り替え可能

## トラブルシューティング

### 認証関連の問題
- ログインに失敗する場合: `claude /login` を再実行
- [Claude.ai](https://claude.ai/) でClaude MAXサブスクリプションが有効か確認
- ブラウザアクセスが可能な環境か確認

### インストール関連の問題
- Node.jsが利用可能か確認: `node --version`
- 手動でClaude Codeをインストール: `npm install -g @anthropic-ai/claude-code`

### 権限関連の問題
- 必要に応じてsudoで実行: `sudo npm install -g @anthropic-ai/claude-code`

### ブラウザアクセスの問題
- ヘッドレス環境の場合、Claude CodeがURLを提供します
- 提供されたURLをClaude.aiにログイン済みのブラウザで開いてください

## 活用のヒント

- Go開発支援にClaude Codeを活用
- CLIツール実装パターンのアドバイスを求める
- コードレビューや改善提案を受ける
- Goアプリケーションのデバッグと最適化

詳細情報: [Claude Code ドキュメント](https://docs.anthropic.com/en/docs/claude-code)