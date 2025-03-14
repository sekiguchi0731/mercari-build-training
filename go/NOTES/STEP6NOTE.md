# テスト

## コード解説

コードの内容：

```go
func TestSayHello(t *testing.T) {
    cases := map[string]struct{
        name string
        want string
    }{
        "Alice": {
            name: "Alice",
            want: "Hello, Alice!"
        }
        "empty": {
            name: "",
            want: "Hello!"
        }
    }

    for name, tt := range cases {
        t.Run(name, func(t *testing.T) {
            got := sayHello(tt.name)

            // 期待する返り値と実際に得た値が同じか確認
            if tt.want != got {
                // 期待する返り値と実際に得た値が異なる場合は、エラーを表示
                t.Errorf("unexpected result of sayHello: want=%v, got=%v", tt.want, got)
            }
        })
    }
}
```

### `map`

Pythonにおける`dict`と似ている．

基本的な`map`の宣言

```go
myMap := map[string]int{
    "apple":  3,
    "banana": 5,
    "orange": 2,
}
```

### `for i, j := range cases`

Pythonにおける`for key, value in dict.items()`と同じ．

### Goのテスト用関数`t.Run()`

Pythonにおける`subTest(name, lambda: ...)`

### `t.Error()`と`t.Fatal()`の違い

| **関数** | **動作** | **使いどころ** |
|---------|---------|-------------|
| `t.Error()` | **エラーを記録するが、テストは続行** | **複数のアサーションをチェックしたい場合** |
| `t.Fatal()` | **エラーを記録し、テストを即座に終了** | **エラーが発生した時点でテストを中断すべき場合** |

## テストの実施

特定の関数のテスト：

```plaintext
go test -run TestParseAddItemRequest
```

### `t.Pallarel()`とは

Go のテストで**並列実行**を可能にする関数です

#### 動作

`t.Parallel()`を呼び出すと，そのテスト関数は他のテストと並行して実行される．

同じテストファイル内の他の`t.Parallel()`を含むテストも同時に実行される．

デフォルトでは、Go の`go test`はテストを**直列（順番）**に実行するが，`t.Parallel()`を使うと**並列で実行**される．

#### メリット

テストの実行時間が短縮される（特に多数のテストがある場合）．

副作用のないテストを並列化できる（独立したテストケースを最適化）．

#### 注意点

共有データ（グローバル変数など）を変更するテストでは使用を避ける！

---

### `mock_infra.go`のメモ

上記のテストを実行する際に，`server_test.go`がビルドできないエラーが発生．

流れとしては，

- `server_test.go`内の`mockIR`が型エラー
  - 元々のinterfaceにシグネチャを追加したため，mock内の型と整合性がつかずエラー
- 型エラーのため，`server_test.go`がビルドできず，テストも実行されない
  - `go test -run ..`は，ファイル全体をビルドしてから関数をテスト

また，`mock_infra.go`は自動作成されたため，編集できない．

よって，以下を実行し，モックファイルを作り直す．

```plaintext
go run go.uber.org/mock/mockgen -source=infra.go -package=app -destination=mock_infra.go
```

---

### テストの出力

1. そのまま比較する．

   - 比較対象同士が単純なものであるときに使える．
     - 文字列や整数など．

  ```go
  if res != want {
    t.Errorf("...")
  }
  ```

1. `cmp.Diff`を用いる

   - 比較対象同士が複雑なものであるときに，フィールドごとの差分を見やすく出力できる
     - JSONやmap構造のときなど．

  ```go
  expected := `{"message": "Hello, world!"}`
  actual := `{"message": "Hello world"}`
  fmt.Println(cmp.Diff(expected, actual))
  ```

  ```plaintext
  -{"message": "Hello, world!"}
  +{"message": "Hello world"}
  ```

---

#### `cmp.Diff`の比較結果が`map[string]string`と（JSON）`string`で違う場合の修正方法

例：

```plaintext
unexpected response body (-want +got): 
      any(
    -       map[string]string{"message": "Hello, world!"},
    +       string("{\"message\":\"Hello, world!\"}\n"),
      )
```

JSON文字列を，`map[string]string`型にする．

```go
var gotBody map[string]string
if err := json.NewDecoder(res.Body).Decode(&gotBody); err != nil {
  t.Fatalf("failed to decode response body: %v", err)
}
```

## モック

### 概要

モック（Mock）とは，テストのために作成された**偽物のオブジェクト**のこと  
実際の処理を行わず，決められたデータや動作を返すことで，単体テストを簡単にする．  
テスト対象（関数や構造体）が依存している外部のシステムを置き換える役割を持つ．

### 対象と目的

#### 対象

Mockする対象は，`interface`．

つまり，`interface`を対象にMockを生成するので，ファイル内に`interface`がない場合，Mockのメリットがなくなる．

また，**テストしたい関数自体をMockするのではない**！

→Mockは，**外部依存な関数について，外部の値を使わず，返り値を指定するもの**だから．

つまり，テストしたい関数内部に，外部依存な関数が存在するときに，（本来のテスト対象は外部依存な部分ではないので）余計なコストがかかってしまう．  
これを防ぐために，テストしたい関数に含まれる**外部依存な関数をMockする**．

したがって，外部依存な関数は`interface`のシグネチャに含めるべき．

#### 目的

**外部依存を排除して，テストを簡単にする**こと．

例えば，データベースアクセスを含む関数をテストする場合，実際のデータベースを用意したり，セットアップや実行時間が長くなったりして，難しい．

そこで，Mockを使って，`interface`に含まれる関数の，外部依存を排除した部分のみをテストする．

したがって，データベースから発生するエラーなどはモックでは取得できないため，**Integration Test**を実行する．

### 生成方法

```plaintext
go run go.uber.org/mock/mockgen -source=infra.go -package=app -destination=mock_infra.go
```

1. `mockgen -source=<ファイル名>`を実行すると，`<ファイル名>`にある`interface`を自動的に解析し，それを満たす`Mock<Interface名>`を生成する．
2. もし`interface`が複数あれば，それぞれに対応する`Mock<Interface名>`も複数生成される.

### 実際のテスト実装

```go
func TestAddItem(t *testing.T) {
  t.Parallel()

  type wants struct {
    code int
  }
  cases := map[string]struct {
    args     map[string]string
    injector func(m *MockItemRepository)
    wants
  }{
    "ok: correctly inserted": {
      args: map[string]string{
        "name":     "used iPhone 16e",
        "category": "phone",
      },
      injector: func(m *MockItemRepository) {
        // STEP 6-3: define mock expectation
        // succeeded to insert
        
      },
      wants: wants{
        code: http.StatusOK,
      },
    },
    "ng: failed to insert": {
      args: map[string]string{
        "name":     "used iPhone 16e",
        "category": "phone",
      },
      injector: func(m *MockItemRepository) {
        // STEP 6-3: define mock expectation
        // failed to insert
      },
      wants: wants{
        code: http.StatusInternalServerError,
      },
    },
  }
...
}
```

この関数の目的は．`AddItem`をテストすること！

しかし，`AddItem`に外部依存な`Insert`が含まれているため，この部分をMockする．

#### `injector`の役割

Mockに`EXPECT().Return()`を設定する関数．

```go
m.EXPECT().Mockしたい関数(引数).Return(Mock関数から期待した返り値)
```

このようにすることで，もしMockした`Insert`がエラーなしなのに，`AddItem`が`http.StatusInternalServerError`を返したら，`AddItem`内で，`Insert`以外の部分でロジックなどのエラーがあることをテストできる！

## エンドツーエンドテスト

| **テストの種類** | **範囲** | **目的** | **特徴** | **例** |
|--------------|-------|-------|-------|-------|
| **ユニットテスト（Unit Test）** | 関数・メソッド単位 | **小さい単位での動作確認** | - 個々の関数が正しく動くかを検証<\br>- 外部依存を持たないようにする | `AddItem()` の処理が正しく動くか |
| **結合テスト（Integration Test）** | 複数のコンポーネント間 | **API や DB とのやりとりの確認** | - モジュール間の連携が正しく動くか検証<\br>- DBやAPIを実際に呼び出す | API 経由でDBにデータが登録されるか |
| **E2E テスト（End-to-End Test）** | システム全体 | **ユーザー視点の動作検証** | - フロントエンド・バックエンド・DBの連携をチェック<\br>- 実際のユーザー操作をシミュレーション | ユーザーがボタンを押したらデータがDBに反映されるか |

- 最初の "End" → ユーザーがアクションを起こすポイント（例: フロントエンドでボタンを押す）
- 最後の "End" → アクションの結果が反映されるポイント（例: データがDBに保存され、フロントに正しく表示される）

**システムの「入力」と「出力」の流れを，最初から最後まで検証するのが**E2Eテスト！

### 仮DBセットアップ関数

#### 目的

- テスト用のデータベースを一時的に作成し，テストが終わったら削除する．
- データベースの初期化（テーブル作成など）を行う
- エラーが発生した場合には適切にリソースを解放する．

#### 引数と返り値

```go
func setupDB(t *testing.T) (db *sql.DB, closers []func(), e error) {..}
```

| **パラメータ** | **説明** |
|---------|---------|
| `t *testing.T` | Goのテスト用の`t`オブジェクト．エラーハンドリング用 |

| **返り値** | **説明** |
|---------|---------|
| db *sql.DB` | セットアップされた SQLite データベースのインスタンス |
| `closers []func()` | 後でクリーンアップするための関数のリスト（ファイル削除、DB クローズなど）|
| `e error` | エラーが発生した場合に返される |

#### 内部処理

```go
  t.Helper()
```

- 補助関数であることをGoに伝える．

```go
  defer func() {
    if e != nil {
      for _, c := range closers {
        c()
      }
    }
  }()
```

- `setupDB`が終了した後に，もし`e`に値があれば，`closers`内の関数を全て実行する．

```go
  // create a temporary file for e2e testing
  f, err := os.CreateTemp(".", "*.sqlite3")
  if err != nil {
    return nil, nil, err
  }
  closers = append(closers, func() {
    f.Close()
    os.Remove(f.Name())
  })
```

- `os.CreateTemp(".", "*.sqlite3")` で一時的なSQLite データベースファイルを作成
- クリーンアップ用に`closers`に`f.Close()`と`os.Remove(f.Name())`を追加（テストが終わったらファイルを削除）

```go
  // set up tables
  db, err = sql.Open("sqlite3", f.Name())
  if err != nil {
    return nil, nil, err
  }
  closers = append(closers, func() {
    db.Close()
  })
```

- `sql.Open("sqlite3", f.Name())`で SQLite データベースに接続
- `closers`に`db.Close()`を追加（テストが終わったらDBを閉じる）

### E2Eテスト関数

```go
  t.Cleanup(func() {
    for _, c := range closers {
      c()
    }
  })
```

- `t.Cleanup()`に関数を登録すると，テストの最後に必ず実行される！
- 登録された`closers`内の関数を順番に呼び出して，リソースを適切に解放する！

#### `defer`ではなく`t.Cleanup()`を使う理由

- `defer`は関数を抜けるときに 即座に実行される
- `t.Cleanup()`は，テスト関数`(TestXxx)`の最後にまとめて実行される
- テストがパニック (`t.Fatalf()`) で終了しても，`t.Cleanup()`は確実に実行される！
  - `defer`は`TestExample()`を抜けるとすぐに実行されるので，`file.Name()`を後のテストコードで使おうとするとすでに削除されてしまっている可能性がある

