# go-webhook-receiver
GithubからのWebhookを受け取り、指定のコマンド(ここではnext.jsのビルド&デプロイ)を実行するスクリプト

## 構造
### ヘッダー検証
contextからGetHeaderを使ってヘッダーの値を取得できる
```go
c.GetHeader("X-GitHub-Event")
```

### 外部コマンド実行
```go
cmd := exec.Command("コマンド", "引数1", "引数2")
cmd.Dir = "作業ディレクトリ"
cmd.Run()
```