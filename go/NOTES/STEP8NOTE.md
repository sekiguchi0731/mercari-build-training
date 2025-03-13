# CI

## CI（Continuous Integration: 継続的インテグレーション）とは

**CI（継続的インテグレーション）** とは、  
ソフトウェア開発において、**コードの変更を頻繁にリポジトリへ統合（マージ）し、 自動でビルドやテストを実行するプロセス** のことです。

### 目的

- **開発中のコードの品質を維持**
- **バグを早期に発見・修正**
- **チーム開発をスムーズにする**

CI の主な流れ：

1. **開発者がコードを Git リポジトリにプッシュ（push）**
2. **CI ツールが自動でビルドを実行**
3. **自動テスト（ユニットテスト・統合テストなど）を実行**
4. **問題がなければ、コードが統合（merge）される**

### 構成

#### 1️⃣ ソースコード管理（SCM）

- GitHub, GitLab, Bitbucket などのリポジトリ管理ツール
- CI は、リポジトリへのプッシュやプルリクエストをトリガーに動作する

#### 2️⃣ CI ツール

- **GitHub Actions**
- **GitLab CI/CD**
- **CircleCI**
- **Jenkins**
- **Travis CI**  
などがあり、リポジトリの変更を検知してビルドやテストを実行する。

#### 3️⃣ ビルド & テストの自動化

- `Docker` などを使って環境を統一
- `pytest`（Python）や `Jest`（JavaScript）などのテストツールで自動テストを実行

### メリット

✅ **バグの早期発見** → 小さな変更ごとにテストするため、不具合を素早く特定できる  
✅ **品質の向上** → コードの変更ごとにテストが実行されるため、不具合の混入を防げる  
✅ **開発スピードの向上** → 自動化により、手作業のミスを減らし、デプロイまでの時間を短縮  
✅ **チーム開発がスムーズ** → 誰かのコードが原因で他の開発者の環境が壊れるのを防ぐ  

### 実際の動作（GitHub Actions の例）

例えば、GitHub Actions を使って `Python` プロジェクトの CI を実装する場合、`.github/workflows/ci.yml` に以下のような設定を記述します。

```yaml
name: CI Pipeline

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: リポジトリをチェックアウト
        uses: actions/checkout@v2

      - name: Python をセットアップ
        uses: actions/setup-python@v4
        with:
          python-version: "3.12"

      - name: 依存関係をインストール
        run: pip install -r requirements.txt

      - name: テストを実行
        run: pytest
```

1. **GitHub にコードが push されたら、この CI ワークフローが起動**
2. **`Python 3.12` の環境をセットアップ**
3. **`pip install -r requirements.txt` で依存関係をインストール**
4. **`pytest` でテストを実行**
5. **テストが失敗した場合、プルリクエストのマージがブロックされる**

---

## CI と CD（Continuous Deployment）の違い

| **項目** | **CI（継続的インテグレーション）** | **CD（継続的デリバリー / デプロイ）** |
|----------|--------------------------------|--------------------------------|
| **目的** | コードの統合・自動テスト | 自動デプロイ（本番環境への適用） |
| **タイミング** | プルリクエストやマージ時 | CI が成功した後、本番環境へ |
| **主なツール** | GitHub Actions, GitLab CI, CircleCI | AWS CodeDeploy, ArgoCD, Kubernetes |

✅ **CI（継続的インテグレーション）** ＝ **コード変更ごとに自動でビルド・テストを行い、品質を保つ仕組み**  
✅ **ツール（GitHub Actions, Jenkins など）を使って、リポジトリの変更を監視し、自動でテストを実行**  
✅ **バグを早期に発見し、開発スピードを向上させるために導入される**  
✅ **CI + CD を組み合わせると、デプロイまで自動化できる**

## Step 6との比較

**`go test -run TestHelloHandler` のように、手動でユニットテストや統合テストを実行していた作業を、CI によって自動化する** ということ．

Step6では，開発者は以下のように **手動でテスト** を行っていた：

```sh
go test -run TestHelloHandler
```

**CIを導入すると**...

- `go test` の実行を **自動化**
- すべてのテストファイルを **自動で検出**
- 変更があれば **即座にテストを実行**
- もし失敗したら、**通知を送る（GitHub Actions なら PR に赤バツがつく）**

### 👩‍💻 手動で行っていたこと

開発者はこれまで：

1. 必要があればモックを作成（Mocking）
2. **ユニットテスト** を実行（関数単位のテスト）

   ```sh
   go test -run TestHelloHandler
   ```

3. **統合テスト** を実行（異なるモジュールが正しく動くか）

   ```sh
   go test -run TestDatabaseIntegration
   ```

4. **エンドツーエンド（E2E）テスト** を実行（全体の動作確認）

   ```sh
   go test -run TestAPIEndpoints
   ```

5. テストが全部通ったら、GitHub にプッシュ

---

### **🤖 CI（自動化）の場合**

開発者が `git push` した瞬間に、**CI が以下を自動で実行** します：

1. **テストファイルを検出**
   - `*_test.go` のようなファイルを CI ツールが自動で見つける
2. **モックを含めたユニットテストを実行**
   - `go test ./...` で全てのテストを実行
3. **統合テストを実行**
   - 例えば、データベースとの接続をテスト
4. **E2Eテスト（エンドツーエンドテスト）を実行**
   - API のエンドポイントを実際に叩いてテスト
5. **テスト結果を通知**
   - もし失敗したら、GitHub や Slack に通知
   - GitHub の PR に「テスト通過 ✅」「テスト失敗 ❌」を表示

---

### CI で Go のテストを自動化する設定

例えば，**GitHub Actions を使って Go のテストを自動化** する場合，  
リポジトリの `.github/workflows/go-ci.yml` に以下を追加

```yaml
name: Go CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - name: リポジトリをチェックアウト
        uses: actions/checkout@v2

      - name: Go をセットアップ
        uses: actions/setup-go@v4
        with:
          go-version: "1.21"

      - name: 依存関係をインストール
        run: go mod tidy

      - name: ユニットテストを実行
        run: go test -v ./...

      - name: 統合テストを実行
        run: go test -tags=integration ./...
```

CI を導入すると **「テストを実行し忘れる」ことがなくなる** ので，チーム開発でも安心してコードを統合できるようになる．

## `yml`コードの解説

```yml
name: build
```

- ワークフロー名（GitHub Actions に表示される）

```yml
run-name: ${{ github.actor }} is building ${{ github.ref_name }} 🚀
```

- 実行時の表示名（GitHub Actions の実行履歴に表示）
  - `github.actor`, `github.ref_name`はGitHub Actionsが提供する変数
  - `github.actor`はpushしたユーザ名，`github.ref_name`はpush先のブランチ名

```yml
on: [push]
```

- ワークフローのトリガーを設定
  - この場合，`push`が実行されるとトリガー

```yml
env:
 REGISTRY: ghcr.io
 IMAGE_NAME: ${{ github.repository }}
```

- 環境変数を定義
  - `REGISTRY`：Dockerイメージのpush先を指定．
    - `ghcr.io`はGitHub Container Registoryの略で，GitHubのレジストリを使うなら，ここにpushする
    - そのほかにも，`docker.io`などもある
  - `IMAGE_NAME`：Dockerイメージ名
    - `github.repository`なので，リポジトリ名

---

### `jobs`とは

ジョブとは，**GitHub Actions のワークフロー内で実行される個々のタスクのまとまり**

1つのワークフローには**複数のジョブを定義でき**，並行実行（並列）や順番に実行（依存関係）も可能

並列実行の例：

```yml
jobs:
  job1:  # ジョブ1の名前
    runs-on: ubuntu-latest  # 実行環境（OS）
    steps:
      - name: ジョブ1のステップ1
        run: echo "Hello from Job 1"
      - name: ジョブ1のステップ2
        run: echo "Job 1 finished!"

  job2:  # ジョブ2の名前
    runs-on: ubuntu-latest  # 実行環境（OS）
    steps:
      - name: ジョブ2のステップ1
        run: echo "Hello from Job 2"
```

この例では，`job1`と`job2`が並行して実行される
→`job1`が終わるのを待たずに`job2`も実行される

依存関係の例：

依存関係（順序）を作るには，`needs`を使う．

```yml
jobs:
  job1:
    runs-on: ubuntu-latest
    steps:
      - name: ジョブ1のステップ1
        run: echo "Hello from Job 1"

  job2:
    needs: job1  # `job1` が終わってから `job2` を実行
    runs-on: ubuntu-latest
    steps:
      - name: ジョブ2のステップ1
        run: echo "Hello from Job 2"
```

1️⃣ `job1`が実行される
2️⃣ `job1`が完了すると`job2`が実行される

このように`needs`を使えば，ジョブの実行順序を制御できる！

---

```yml
jobs:
  build:
    runs-on: ubuntu-latest
```

- ジョブの定義
  - ジョブは最新のUbuntu環境で実行される

```yml
    permissions:
     contents: read
     packages: write
```

- ジョブの権限設定
  - ジョブに対して何の権限を与えるかの定義
  - `contents: read`：ジョブに対して，コンテンツ（コード）を読み取る権限を付与
  - `packages: write`：パッケージとは，作成されるDockerイメージのこと．コンテナイメージをghcrにpushするための権限

### ステップ１：リポジトリをチェックアウト

```yml
    steps:
    # Checkout repository
    - name: Checkout
      uses: actions/checkout@v3
```

チェックアウトとは，`git checkout main`と同じで，**Gitのリポジトリから特定のブランチやコミットのコードを取得する操作**のこと．

GitHub Actionsでは，GitHubの仮想環境（この場合はUbuntu）でコードを実行する．

リポジトリのコードを GitHub の仮想環境（ワークスペース）にダウンロードする必要があるため，`checkout`をまず行う．

- `name: ..`：ステップ名（ログに表示される）
- `uses: actions/checkout@v3`：
  - `actions/checkout@v3`は，GitHub Actionsによって提供される公式アクション（関数のようなもの．）

### ステップ２：コンテナレジストリ（GitHub Container Registry）にログイン

```yml
    - name: Log in to the Container registry
      uses: docker/login-action@f054a8b539a109f9f41c372932f1ae047eff08c9
      with:
       registry: ${{ env.REGISTRY }}
       username: ${{ github.actor }}
       password: ${{ secrets.GITHUB_TOKEN }}
```

- `uses: docker/login-action@f054a8b539a109f9f41c372932f1ae047eff08c9`
  - Docker社が提供する**公式アクション（ログイン）**
    - `@`以降は，使用する公式アクションのバージョンを指定する．
      - `@v2`：`v2`タグがついた最新の安定版を使用
      - `@f0..`：Docker社が開発したログインアクションの，特定のコミットバージョンを使用（推奨）
      - `@main`：最新の`main`ブランチの安定版を使用
  - `password: ${{ secrets.GITHUB_TOKEN }}`：
    - CI/DI用に，一時的にGitHubが自動で発行するトークン
    - プルリクエストや`docker push`の認証に使用できる
    - ワークフローの実行中のみ有効で，実行が終わると無効化される（セキュリティ的に安全）

### ステップ３：Dockerのメタデータ（タグ・ラベル）を抽出

```yml
    - name: Extract metadata (tags, labels) for Docker
      id: meta
      uses: docker/metadata-action@98669ae865ea3cffbcbaa878cf57c20bbf1c6c38
      with:
       images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
```

- `id: meta`
  - ステップにつけるタグ（識別子）のようなもの．ステップの出力を別のステップから参照できる．
  - `${{ steps.meta.outputs.〇〇 }}`という形で．
  - スコープは，同じ`jobs`内．
- `uses: docker/metadata-action@9..`
  - 以下のような 出力（outputs）を自動生成する：
    - `outputs`のキー：`tags`
      - Dockerイメージのタグ
    - `outputs`のキー：`labels`
      - Dockerイメージのラベル
  - `docker/metadata-action`が`images:`をもとに，適切な`tags`や`labels`を自動生成する．

### ステップ４：変数の確認

```yml
    - name: Check variables
      run: |
        echo 'Current path:'
        pwd
        echo 'Tag: ${{ steps.meta.outputs.tags }}'
        echo 'Label: ${{ steps.meta.labels.tags }}'
```

### ステップ5：Dockerイメージのビルド&プッシュ

```yml
    - name: Build and push Docker image
      uses: docker/build-push-action@ad44023a93711e3deb337508980b4b5e9bcdc5dc
      with:
        context: <go>
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
```

- `context`
  - `Dockerfile`が存在する，ビルドするディレクトリ
- `push`
  - `true`にすると，GitHub Container Registryにプッシュ
