# データベース

## `SQLite3`のインストール

Macでインストール

```plaintext
brew install sqlite
```

インストールの確認：

```plaintext
sqlite3 --version
```

## Goで使うためのセットアップ

`go.mod`が存在するディレクトリに移動．

---

### Go Modulesとは

- プロジェクト単位でパッケージを管理する仕組み
- 外部のライブラリも適切に管理
  - 外部のライブラリは`go get`にて取得する．

---

```plaintext
go get github.com/mattn/go-sqlite3
```

`go.mod`にて，以下が追加される．

```mod
github.com/mattn/go-sqlite3 v1.14.24
```

go内での使用例：

```go
import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func main() {
  db, err := sql.Open("sqlite3", "test.db")
  if err != nil {
    panic(err)
  }
  defer db.Close()
}
```

## データベースにおける用語：「テーブル」・「スキーマ」

### テーブル（Table）とは

- **データを保存する「表（Excel のシートのようなもの）」**
- **「列（カラム）」と「行（レコード）」で構成される**
- データベース内に複数のテーブルを作成できる
- 各行（レコード）が 1 つのデータを表す

#### 例：`items` テーブル

| id | name   | category | image_name     |
|----|--------|---------|---------------|
| 1  | Jacket | Fashion | abc123.jpg    |
| 2  | T-Shirt| Fashion | def456.jpg    |
| 3  | Shoes  | Footwear | ghi789.jpg    |

- **`id`（主キー）** → 各アイテムを一意に識別
- **`name`** → 商品の名前
- **`category`** → 商品のカテゴリ
- **`image_name`** → 商品画像のファイル名

---

### **🔹 スキーマ（Schema）とは**

- **データベースの「設計図」や「ルール」**
- **どのようなテーブルがあり、どのカラムを持つかを定義する**
- **テーブルの構造や制約（データ型、NULL制約、キーなど）を決める**
- **「スキーマに沿ってテーブルを作成する」**

#### 例：`items` テーブルのスキーマ

```sql
CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- 商品ごとにユニークなID
    name TEXT NOT NULL,                   -- 商品の名前
    category TEXT NOT NULL,               -- 商品のカテゴリ
    image_name TEXT NOT NULL              -- 画像のパス
);
```

このスキーマがあることで、データベースは「どんなデータをどう保存するか」を理解できる

---

### 「テーブル」と「スキーマ」の関係

| **概念** | **役割** |
|----------|---------|
| **スキーマ** | **データベースの設計図・ルール**（どんなテーブルがあり、どんなカラムを持つか） |
| **テーブル** | **スキーマに基づいて作られるデータの入れ物**（実際のデータを保存） |

### 📌 実際の関係

1. **スキーマを定義**
   - `db/items.sql` にスキーマを記述
  
   ```sql
   CREATE TABLE items (
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       name TEXT NOT NULL,
       category TEXT NOT NULL,
       image_name TEXT NOT NULL
   );
   ```

2. **テーブルを作成**
   - `.read db/items.sql` を実行し、データベースに `items` テーブルを作成
3. **データを保存**
   - `INSERT INTO items (name, category, image_name) VALUES ('Jacket', 'Fashion', 'abc123.jpg');`

---

## ✅ まとめ

| **概念** | **説明** | **例** |
|---------|---------|---------|
| **スキーマ（Schema）** | データベースの設計図（どんなテーブルがあり，どんなデータ型か） | `CREATE TABLE` の定義 |
| **テーブル（Table）** | スキーマに基づいて作成されたデータの入れ物 | `items` テーブル |

🚀 **スキーマを定義し，それを元にテーブルを作成し，実際のデータを保存する！**

## `SQLite`におけるスキーマの定義方法

1. `SQLite`のインタラクティブモードで直接コマンドを実行する
2. `db/ファイル名.sql`にスキーマを定義し，インタラクティブモードでそれを実行（`.read`）する

通常は，二つ目の手法を取る．

理由：

- 再利用可能
- gitで変更履歴を追える
- 修正しやすい
- バグを発見しやすい

よって今回は，二つ目の手法について．

### スキーマを定義し，定義したスキーマに準ずるDBファイルを作成する

まず，`db/filename.sql`を作成する．

``plaintext
touch ./db/filename.sql
```

まず，`SQLite`のインタラクティブモードに入る

```plaintext
sqlite3 データベースファイル名
```

その後，スキーマを以下の形式で定義する．

```sql
CREATE TABLE テーブル名 (
    {各列の指定方式}
);
```

以下も可能

```sql
CREATE TABLE IF NOT EXISTS テーブル名 (
    {各列の指定方式}
);
```

各列の指定方式は以下

```plaintext
列名 データ型 [制約] [追加のルール]
```

- []：Optional

| **要素** | **説明** | **例** |
|---------|---------|--------|
| **列名** | カラムの名前 | `id`, `name`, `category` |
| **データ型** | 保存するデータの種類 | `INTEGER`, `TEXT`, `REAL` |
| **制約（Constraints）** | `NOT NULL`, `PRIMARY KEY`, `UNIQUE` など | `NOT NULL`, `DEFAULT` |
| **追加のルール** | `AUTOINCREMENT`, `DEFAULT` など | `AUTOINCREMENT`, `DEFAULT 'unknown'` |

### DBファイルを実行

まず，sqlite3のインタラクティブモードに入る

```plaintext
sqlite3 ファイル名.sqlite3
```

次に，スキーマを定義したDBファイルを実行する．

```plaintext
sqlite> .read items.sql
```

これで，ファイル名.sqlite3のスキーマが設定される．

---

#### `.db`ファイルと`.sqlite3`ファイルの違い

| **用途** | **`.db` を使う場合** | **`.sqlite3` を使う場合** |
|---------|------------------|------------------|
| **一般的なデータベース** | **他のデータベースシステム（MySQL, PostgreSQL など）と統一するために `.db` を使うことが多い** | **SQLite であることを明確にするために `.sqlite3` を使う** |
| **プロジェクトのデフォルト** | **シンプルに `database.db` などにする場合** | **複数の DB システムがあるとき `database.sqlite3` で区別する** |
| **外部ツールとの互換性** | **汎用的なデータベース名として `.db` を使うと、様々なツールで開きやすい** | **SQLite であることを示したい場合は `.sqlite3` を使う** |

`.sqlite2` \(\in\) `.db`

### `.sql`ファイルと`.sqlite3`ファイルの違い

| **項目** | **`.sqlite3`（SQLite データベース）** | **`.sql`（SQL スクリプト）** |
|---------|---------------------------|----------------------|
| **役割** | **データベース本体**（データ & スキーマ） | **スキーマ & データを作成するためのスクリプト** |
| **データの保存** | **データを格納する（バイナリ形式）** | **データを含む SQL 文を書くことはできるが、データベースそのものではない** |
| **実行方法** | `sqlite3 mercari.sqlite3` で開く | `sqlite3 mercari.sqlite3 < db/items.sql` で適用 |
| **データの有無** | **データが入っている（データベース）** | **SQL コマンドを持つ（命令書）** |

`.sql`は`.sqlite3`の設計図にすぎず，`.sqlite3`こそがデータベースそのもの.

---

## `SQLite`と`Go`を繋げる

```go
import(
  _ "github.com/mattn/go-sqlite3"
  "database/sql"
  )
```

### `Insert`関数

JSONでは**ロード→変更→保存**の3ステップが必要だが，DBならSQLの`INSERT`を実行するだけですむ．

### `GetItems`関数

1. DBファイルから該当列を`SELECT`する
2. `SELECT`した行列について，一行ずつ`rows.Scan`し，DBデータを構造体にマッピングする．
  
   - `rows.Scan(&item.ID, &item.Name, &item.Category, &item.ImageName)`は，`rows`の`i`行目の0列目を`item.ID`に格納
   - なお，`Scan()`の引数はポインタであり，アドレスが渡されると，そのメモリアドレスにデータが格納される．

3. `items []Item`に追加


## デバッグ用メモ

```plaintext
% curl \
  -X GET \
  --url 'http://localhost:9001/'
```

```plaintext
% curl \
  -X POST \
  --url 'http://localhost:9001/items' \
  -F 'name=jacket' \
  -F 'category=fashion' \
  -F 'image=./images/default.jpg'
```

```plaintext
% curl \                           
  -X GET \ 
  --url 'http://localhost:9001/items/11'
```

```plaintext
% curl \
  -X GET \
  --url 'http://localhost:9001/images/30de3839c2cc5185d8fb630cf03fa5f72b05898376eaa039a28256d90a8fb178.jpg'
```