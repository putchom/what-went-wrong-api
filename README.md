# What Went Wrong API

## 前提条件

- **Go**: 1.25+
- **Docker**: データベースの実行に必要

## セットアップ

1. **リポジトリのクローン**
   ```bash
   git clone https://github.com/putchom/what-went-wrong-api.git
   cd what-went-wrong-api
   ```

2. **依存関係のインストール**
   ```bash
   go mod download
   ```

3. **環境変数の設定**
   `.env.example`をコピーして`.env`ファイルを作成します:
   ```bash
   cp .env.example .env
   ```
   
   必要に応じて`.env`ファイルの値を編集してください:
   ```env
   POSTGRES_PASSWORD=password
   POSTGRES_USER=postgres
   POSTGRES_DB=postgres
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   AUTH0_DOMAIN=your-auth0-domain
   AUTH0_AUDIENCE=your-auth0-audience
   ```

4. **Swaggerのインストール (任意)**
   APIドキュメントを再生成する必要がある場合:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

## 認証 (Auth0)

このAPIの一部エンドポイントは保護されています。アクセスには有効なBearerトークンが必要です。

### トークンの取得方法
開発環境では、Auth0ダッシュボードからテスト用トークンを取得するか、Swagger UIの「Authorize」ボタンを使用してください。

### Swaggerでの認証
1. [Swagger UI](http://localhost:8080/swagger/index.html) にアクセスします。
2. 右上の **Authorize** ボタンをクリックします。
3. `Bearer <YOUR_TOKEN>` の形式でトークンを入力します（`Bearer ` プレフィックスを忘れずに）。
4. **Authorize** をクリックして閉じます。
5. これで保護されたエンドポイントを実行できます。

## アプリケーションの実行

### 1. データベースの起動
Docker Composeを使用してPostgreSQLデータベースを起動します。

```bash
docker compose up -d
```

- **ポート**: 5432
- **パスワード**: `password`

### 2. サーバーの実行
`go run`を使用してサーバーを直接実行できます:

```bash
go run ./cmd/server/main.go
```

サーバーは `http://localhost:8080` で起動します。

## APIドキュメント

サーバー起動中にSwagger UIを利用できます。

- **URL**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

コード修正後にSwaggerドキュメントを再生成するには:

```bash
swag init -g ./cmd/server/main.go -o cmd/docs
```

## プロジェクト構成

- `cmd/`: アプリケーションのエントリーポイント。
  - `server/`: APIサーバーのメインアプリケーション。
  - `docs/`: 生成されたSwaggerドキュメント。
- `internal/`: 外部から参照されないアプリケーション固有のコードやライブラリ。