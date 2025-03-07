# コードの解説

## 構造体

```go
type Server struct {
  // Port is the port number to listen on.
  Port string
  // ImageDirPath is the path to the directory storing images.
  ImageDirPath string
}
```

- `type ~ struct {..}`で，`~`という名前の構造体を定義．
- `Port`や`ImageDirPath`はプロパティ名（Goでは**フィールド**名）

### 使用例

```go
func main() {
  // `Server` 構造体を作成（インスタンス化）
  s := Server{
    Port: "8080",
    ImageDirPath: "/var/images",
  }

  // 構造体のフィールドにアクセス
  fmt.Println("Server Port:", s.Port)  // Server Port: 8080
  fmt.Println("Image Directory:", s.ImageDirPath)  // Image Directory: /var/images
}
```

## `Run`メソッド

### メソッドと関数

```go
func (s Server) Run() int {..}
```

Goでは**構造体にメソッドを後から追加できる！**

- `(s Sever)`のように，レシーバー(ここでは`s`)を指定することで，構造体にメソッドを追加できる．

Goでは，**関数とメソッドは異なる概念**

- **関数**：
  - どの構造体にも属さない独立した処理
- **メソッド**：
  - 構造体に紐づいた関数

### レシーバー

オブジェクト指向言語における，`this`（**=自分自身のインスタンスを指す特別な変数**）のような役割をする．

#### レシーバー名

レシーバー名が異なっていても，**同じ構造体に追加**される！

つまり，それぞれのメソッドに対応したインスタンスが作られるわけではなく，レシーバーを指定した時点で，**全てのメソッドが同じ一つの構造体に追加**される．

```go
package main

import "fmt"

type Server struct {
  Port string
}

// `s` という名前のレシーバーを使う
func (s Server) Run() {
  fmt.Println("Running server on port:", s.Port)
}

// `s2` という名前のレシーバーを使う
func (s2 Server) Start() {
  fmt.Println("Starting the server on port:", s2.Port)
  s2.Run()  // `s2` を使って `Run()` を実行
}

func main() {
  s := Server{Port: "8080"}
  s.Run()
  s.Start()  // `s` の `Start()` を呼ぶ
}
```

→`Run()`，`Start()`メソッドの両方が`Server()`構造体に追加される．

ちなみに，`s.Run()`は`Run()`内で実行できるが，無限ループに注意する．

#### 値渡しと参照渡し

ポインタレシーバーも利用可能

```go
type Counter struct {
  Value int
}

// 値渡し（コピー）なので、元の `Counter` は変わらない
func (c Counter) Increment() {
  c.Value++
}

// ポインタレシーバーなら、元の `Counter` の値が変わる
func (c *Counter) IncrementPointer() {
  c.Value++
}

func main() {
  c := Counter{Value: 10}
  c.Increment()
  fmt.Println(c.Value)  // 10（変更されない）
  c.IncrementPointer()
  fmt.Println(c.Value)  // 11（変更される）
}
```

### ロガー

プログラムの動作を記録（ログ出力）するツール

```go
import (
  ..
  "log/slog"
  ..
)
```

ロガー用のライブラリ`"log/slog"`を`import`．

```go
  // set up logger
  // ここまではloggerは単なるローカル変数
  logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
  // この設定によって，宣言したloggerがグローバル変数に
  slog.SetDefault(logger)
  // STEP 4-6: set the log level to DEBUG
  slog.SetLogLoggerLevel(slog.LevelInfo)
```

#### STEP 4-6

**Goのログレベル**

- `slog.LevelDebug`
  - デバッグ用（詳細なログ）
- `slog.LevelInfo`
  - 一般的な情報ログ
- `slog.LevelWarn`
  - 警告ログ
- `slog.LevelError`
  - エラーログ

上に行けば行くほど詳細なログが出力され，下のレベルのログを全て包含する．

```go
slog.SetLogLoggerLevel(slog.LevelInfo)

slog.Debug("This is a debug log")  // ❌ 出力されない
slog.Info("This is an info log")   // ✅ 出力される
slog.Warn("This is a warning log") // ✅ 出力される
slog.Error("This is an error log") // ✅ 出力される
```

### ハンドラ関数

```go
mux.HandleFunc("メソッド /パス", ハンドラー関数)
```

Go の HTTP ルーター (http.ServeMux) に`メソッド /パス`というリクエストが来たときに，ハンドラ関数を実行する．

リクエストのURLパスに応じて実行するハンドラ関数を登録できる．

また，パス内で，`{..}`を使うと，動的な変数（=**プレースホルダ**）となり，ハンドラ関数内で変数として使うことができる．

```go
mux.HandleFunc("GET /items/{index}", h.GetItemByIndex)
```

ハンドラ関数内の`r`には，

- リクエストの URL (/items/2 など)
- リクエストのヘッダー
- クエリパラメータ (?key=value の部分)
- ボディ (GET では通常なし)
- リモートアドレス (クライアントの IP など)

などが入っており，このうち，プレースホルダは`r.PathValue("プレースホルダ名")`で取得できる．

また，ハンドラ関数から，別の関数にこの`r`を渡した場合でも，同様のリクエスト情報を渡すことができる．

```go
func (s *Handlers) GetItemByIndex(w http.ResponseWriter, r *http.Request) {
  // 1️⃣ `{index}` の値を取得
  indexStr := r.PathValue("index")

  // 2️⃣ 文字列を整数に変換
  index, err := strconv.Atoi(indexStr)
  if err != nil {
    http.Error(w, "Invalid index", http.StatusBadRequest)
    return
  }

  // 3️⃣ `items.json` からデータ取得
  items, err := s.itemRepo.GetAll()
  if err != nil {
    http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
    return
  }

  // 4️⃣ インデックスが範囲外でないかチェック
  if index < 0 || index >= len(items) {
    http.Error(w, "Index out of range", http.StatusNotFound)
    return
  }

  // 5️⃣ `index` 番目のアイテムを取得し、JSON で返す
  item := items[index]

  w.Header().Set("Content-Type", "application/json")
  err = json.NewEncoder(w).Encode(item)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}
```




## `AddItem`メソッド

### 入力

```go
func (s *Handlers) AddItem(w http.ResponseWriter, r *http.Request) {..}
```

- `w`：HTTP レスポンスを書き込むためのオブジェクト
- `r`：クライアントから送られたリクエストデータ（POST /items の内容）

具体例：

```plaintext
% curl \
  -X POST \
  --url 'http://localhost:9001/items' \
  -d 'name=jacket'
```

これが`r`

```
{"message":"item received: jacket"}
```

これが`w`

### 処理部分1

```go
ctx := r.Context()
```

- HTTPリクエストの「コンテキスト（context.Context）」を取得する
- キャンセル処理やタイムアウト処理に使う（未使用なら削除しても動く）
- Go の context は「リクエストごとの状態管理」に使われる

```go
req, err := parseAddItemRequest(r)
if err != nil {
  http.Error(w, err.Error(), http.StatusBadRequest)
  return
}
```

- リクエスト (r) を解析して、AddItemRequest 型のデータに変換する
- 失敗した場合は err にエラー情報が入る

#### `parseAddItemRequest`関数

```go
func parseAddItemRequest(r *http.Request) (*AddItemRequest, error) {..}
```

入力：`r`（型は`*http.Request`）

出力：型が`(*AddItemRequest, error)`

- `error`型はGoの組み込み型

```go
  req := &AddItemRequest{
    Name: r.FormValue("name"),
    // STEP 4-2: add a category field
  }
```

構造体`AddItemRequest`のインスタンスを作成し，そのアドレスを代入

##### `:=`

**宣言**と**代入**を同時に行う演算子

Goでは通常，

- 宣言：
  - `var 変数名 型`
- 初期化：
  - `変数名 = 初期値`

を行う必要があるが，これらを**同時に行い**，かつ，**動的に型推論**してくれる．

##### `&`

Cと同様，**アドレス演算子**

ここでは，構造体のアドレスを返す．

##### `FormValue`メソッド

`curl -d`で送られた**フォームデータ**を取得

### 処理部分2

```go
  item := &Item{
    Name: req.Name,
    // STEP 4-2: add a category field
    // STEP 4-4: add an image field
  }
```

上記同様，インスタンスを作成しアドレス代入

`Item`構造体は，`infra.go`ファイル内に以下のように記載．

```go
type Item struct {
  ID   int    `db:"id" json:"-"`
  Name string `db:"name" json:"name"`
}
```

- **Goでは，同じディレクトリにあるファイルは「同じパッケージ」として扱われる**
- つまり、**infra.go にある Item は，同じディレクトリ内の他のファイルから import なしでそのまま使える**

#### 構造体タグ

Goでは，構造体のフィールド（=プロパティ）に，データベースやJSONとの対応関係を指定できる．

```go
`db:"id" json:"-"`
```

なお，`"-"`はJSONファイルに記述しないことを意味する．

### 処理部分3

```go
message := fmt.Sprintf("item received: %s", item.Name)
slog.Info(message)
err = s.itemRepo.Insert(ctx, item)
if err != nil {
  slog.Error("failed to store item: ", "error", err)
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}
```

#### `fmt`パッケージ

Go の標準ライブラリで，文字列フォーマットや出力を行うパッケージ．

`Sprintf`メソッドは，`printf`のように文字列をフォーマットして返す．

#### `Handlers`構造体

`s`は`Handlers`構造体のアドレスレシーバーなので，`Handlers`構造体を見る．

```go
type Handlers struct {
  // imgDirPath is the path to the directory storing images.
  imgDirPath string
  itemRepo   ItemRepository
}
```

#### `ItemRepository`インタフェース・`itemRepository`構造体

```go
type ItemRepository interface {
  Insert(ctx context.Context, item *Item) error
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
  // fileName is the path to the JSON file storing items.
  fileName string
}
```

##### 大文字と小文字

Goでは，**大文字で始まる名前は「公開（エクスポート）」、小文字で始まる名前は「非公開」**

- このファイルでは，`ItemRepository`インタフェースは`app`パッケージ以外でもimportすれば使える
- `itemRepository`構造体は他のパッケージや他の構造体からは見えない．

##### インタフェースと構造体

- インタフェース
  - メソッドの定義（シグネチャ）の集合のみ記述できる
  - つまり，インターフェースには**メソッド名と引数・戻り値**しか定義できない
    - メソッドの中身の実装はできない！
  - 外部に提供する機能だけを定義すべき．
    - 外部のコードから非公開メソッドを呼べるようになってしまうため．
- 構造体
  - メソッドに加え，データも実装可能
  - **構造体がインタフェース内のメソッドを実装する**
    - 具体的には，構造体をレシーブ？して，そのメソッドとして定義する．

```go
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
}

func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item

	return nil
}
```

流れとしては，

1. インターフェースにメソッドA（シグネチャ）を定義
2. 非公開構造体を定義
3. その非公開構造体に、メソッドA をレシーバー付きで定義
4. その非公開構造体にメソッドA が加わる
5. メソッドA を持つ構造体は、インターフェースに属するので「インターフェースを満たす」

**「インタフェースを満たす」**=**インタフェース内のメソッドを全て実装している**

##### コンストラクタ関数

```go
func NewItemRepository() ItemRepository {
  return &itemRepository{fileName: "items.json"}
}
```

なぜ作る必要があるか？

- 直接`itemRepository`を作ると，`itemRepository`が`ItemRepository`インタフェースを満たすからといって，返る型は`itemRepository`型
  - 例えば，以下は`*itemRepository`型（非公開型）
  
  ```go
  repo := &itemRepository{fileName: "items.json"}
  ```

  - これを`ItemRepository`型（公開型）にするには，明示的にキャストしなければならない．
- 一方で，コンストラクタ関数は戻り値の型を`ItemRepository`型と宣言しているため，作られるインスタンスは`ItemRepository`型

### 処理部分4

```go
err = s.itemRepo.Insert(ctx, item)
```

#### `Insert`メソッド

実装すべき部分．

流れとしては，

1. これまでのJSONファイルの中身を取得
   1. なければ新規作成
2. その中身にitemを加える
3. JSONファイルをセーブする．

#### スライスとは

動的にサイズを変更できる配列．

##### スライスの追加

```go
slice = append(slice, 4)
```

##### スライスの取得

```go
fmt.PrintLn(slice[3:5])
```

##### `decoder.Decode()`

`decoder.Decode(ある型のポインタ変数)`とすると，`decoder`の中身をその型に変換しつつ，ポインタ変数に変換後のものを代入してくれる！

返り値はエラー

## `GetItem`メソッド

### `itemRepository`構造体と`Handlers`構造体

- `itemRepository`構造体に属するもの
  - `Insert`（公開）
  - `loadItem`（非公開）
  - `saveItems`（非公開）
- `Handlers`構造体に属するメソッド
  - `imgDirPath`（非公開）
  - `itemRepo`（非公開）
    - `itemRepository`型なので，橋渡し！
  - `Hello`（公開）
  - `AddItem`（公開）
  - `GetItem`（公開）

## SHA-256でハッシュ化とは

データを SHA-256（Secure Hash Algorithm 256-bit） というアルゴリズムを使って，固定長（256ビット）のハッシュ値に変換すること．

### Goの`crypto/sha256`パッケージの`sha256.Sum256`

- `[]byte`型のデータをハッシュ化し，32bitのバイト配列`[32]byte`を返す

#### `[]byte`型

`string`と一意

```go
s := "hello"
b := []byte(s) // 文字列 → バイトスライス
fmt.Println(b) // [104 101 108 108 111]
```

`string`に戻せる

```go
b := []byte{104, 101, 108, 108, 111}
s := string(b) // バイトスライス → 文字列
fmt.Println(s) // "hello"
```

ファイルの読み書きはバイト型

```go
func main() {
  // ファイルを読み込む（[]byte 型）
  data, err := os.ReadFile("example.txt")
  if err != nil {
    panic(err)
  }

  // `data` は `[]byte`
  fmt.Println(string(data)) // ファイルの中身を文字列として出力
}
```

`string`は変更できないが，`[]byte`は変更できる

```go
s := "hello"
// s[0] = 'H'  // コンパイルエラー
```

```go
b := []byte("hello")
b[0] = 'H'  // 変更可能 ✅
fmt.Println(string(b)) // "Hello"
```

## `-d`と`-F`の違い

|  | **`-d` (`application/x-www-form-urlencoded`)** | **`-F` (`multipart/form-data`)** |
|----|--------------------------------|-----------------------------|
| **使い方** | `-d "key=value"` | `-F "key=value"` |
| **データの種類** | **文字列や数値のみ** | **文字列 + ファイルも送信可能** |
| **Go 側の取得方法（文字列）** | `r.FormValue("key")` | `r.FormValue("key")` |
| **Go 側の取得方法（ファイル）** | ❌ `FormValue` では取得不可 | ✅ `r.FormFile("key")` を使う |
| **用途** | テキストデータの送信 | ファイルアップロード |

- `-F`であっても，データが文字列や数値なら，`r.FormValue`で取得可能．
- ただし，今回は

> 文字列で保存
  とあるように，`FormValue`として扱う