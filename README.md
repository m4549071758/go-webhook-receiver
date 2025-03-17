# go-webhook-receiver
GithubからのWebhookを受け取り、指定のコマンド(ここではnext.jsのビルド&デプロイ)を実行するスクリプト

## 構造
### ヘッダー検証
contextからGetHeaderを使ってヘッダーの値を取得できる
```go
c.GetHeader("X-GitHub-Event")
```