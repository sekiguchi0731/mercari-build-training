# Docker-compose

## フロントのビルド

フロントのディレクトリ下にも，`Dockerfile`を作成する．

### フロントの`Dockerfile`

```dockerfile
FROM node:20-alpine
```

- `node:20-alpine`ベースイメージを使用
- 軽量な`Node.js`環境を構築

```dockerfile
WORKDIR /app
```

- `/app`を作業ディレクトリに設定
- 以降のコマンドがこのディレクトリ内で実行される

```dockerfile
COPY package.json package-lock.json ./
```

- 依存関係 (`package.json`&`package-lock.json`) をコピー

```dockerfile
RUN npm ci
```

- `npm ci`で確実に`package-lock.json`に従ったインストールを行う

```dockerfile
COPY . .
```

- 全てのファイルを`/app`にコピー

```dockerfile
RUN addgroup -S mercari && adduser -S trainee -G mercari
```

- `mercari`グループと`trainee`ユーザを作成

```dockerfile
RUN chown -R trainee:mercari /app && chmod -R 755 /app
```

- `/app`の所有者を`trainee:mercari`に変更
- アクセス制御を適切に設定

```dockerfile
USER trainee
```

- `root`ではなく`trainee`ユーザーで実行

```dockerfile
EXPOSE 3000
```

- コンテナのポート`3000`を外部に開放

---

#### コンテナのポートとは

Dockerコンテナ内で動作しているアプリケーションがリッスン（待ち受け）するネットワークポートのこと

コンテナ内のアプリが**どのポートで通信を受け付けるか**を指定

`EXPOSE 3000`と記述すると，コンテナ内で3000番ポートが開放される．

- 「コンテナ内でポートを開放する」とは，**コンテナの内部でアプリケーションが特定のポートで通信を受け付けられるようにすること**

ただしこれは，**あくまでドキュメント的な宣言であり、ポートを外部に公開する設定ではない．**

#### 外部に開放とは

ポートフォワーディングとは，**ホストマシン（開発PCやサーバー）上のポートと，コンテナ内のポートを対応付けて通信できるようにすること**

Dockerでは`-p`オプションを使って，ホストのポートとコンテナのポートをマッピング

```sh
docker run -d -p 3000:3000 mercari-build-training/web:latest
```

- `-p 3000:3000`
  - ホスト（PC）の3000番ポート（前者）をコンテナの3000番ポート（後者）に接続する
  - ホストの`http://localhost:3000/`にアクセスすると，コンテナ内のWebアプリに接続できる

---

```dockerfile
CMD ["npm", "start", "--", "--host", "0.0.0.0"]
```

- `npm start`を実行し，外部からアクセスできるようにする

---

#### `0.0.0.0`の意味

Node.js（React, Express）や Flask などのフレームワークでは，**デフォルトで`localhost`しか受け付けない**設定になっていることが多い

そのため，コンテナ内でのアクセスを許可するために`0.0.0.0` にバインドする必要がある

→`0.0.0.0`にすると**コンテナ内のすべてのネットワークインターフェースで接続を受け付ける**ことができる

---

## フロントとバックの`Dockerfile`の違い

| 項目              | バックエンドの Dockerfile                     | フロントエンドの Dockerfile                     |
|------------------|---------------------------------|---------------------------------|
| ベースイメージ    | `golang:1.20-alpine`（例）       | `node:20-alpine`               |
| 作業ディレクトリ  | `WORKDIR /app`                 | `WORKDIR /app`                 |
| 依存関係のインストール | `go mod tidy` や `go build`   | `npm ci`（`package-lock.json` に基づく） |
| ソースコードのコピー | `COPY . .`                    | `COPY . .`                    |
| ユーザー作成      | `RUN adduser -D appuser`（例） | `RUN addgroup -S mercari && adduser -S trainee -G mercari` |
| ユーザー権限変更  | `USER appuser`（例）           | `USER trainee`                 |
| ポート公開        | `EXPOSE 9001`（例）            | `EXPOSE 3000`                  |
| 実行コマンド      | `CMD ["./server"]`（例）       | `CMD ["npm", "start", "--", "--host", "0.0.0.0"]` |
| 主な処理内容      | Go バックエンド API の実行      | React/Vite などのフロントエンドアプリの実行 |

## `docker-compose.yml`ファイル

`docker-compose.yml`は複数のコンテナ（フロントエンド、バックエンド、データベースなど）を**一つのコマンドで起動/停止**できるようにするためのファイル

フロントエンドとバックエンドを**一緒に管理するためのルートディレクトリ**に配置するのが最適

### 例文の解説

```yaml
version: "3.9"
```

- **Docker Compose のバージョン**を指定

```yaml
services:
```

- **コンテナとして起動するサービスの定義**
- ここでは `web`（アプリケーション）と `redis`（データベースキャッシュ）の**2つのサービス**を定義
  - `frontend`や`backend`も指定可能．

```yaml
  web:
```

- **`web` という名前のサービス** を定義
- ここでは，`web`が**アプリケーションのコンテナ**を意味

```yaml
    build: .
```

- **現在のディレクトリ (`.`) にある `Dockerfile` を使ってコンテナをビルド** する
  - `docker build .`と同じ意味
  - **`Dockerfile` を元に `web` のコンテナが作られる**

```yaml
    ports:
      - "8000:5000"
```

- **ホストの `8000` 番ポートを、コンテナ内の `5000` 番ポートにマッピング**
  - ホストの `http://localhost:8000/` にアクセスすると，**コンテナ内の `5000` 番ポートで動いているアプリにリクエストが届く**

```yaml
  redis:
```

- **`redis` という名前のサービス** を定義
  - **Redis（キャッシュ用データベース）を起動** するための設定

```yaml
    image: "redis:alpine"
```

- **`redis:alpine` という公式のRedisイメージを使ってコンテナを作成**
- `alpine`は軽量版のLinuxで，**より高速で小さい Docker イメージ**になる
  - **これは`docker run redis:alpine`と同じ意味

### 実行方法

1. `docker-compose`を使ってビルド

   ```sh
   % docker-compose build
   ```

2. コンテナを起動

   ```sh
   % docker-compose up -d
   ```

再起動（サーバダウン→ビルド→起動）

```sh
docker-compose down                          
docker-compose up --build -d
```

