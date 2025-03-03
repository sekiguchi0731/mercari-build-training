package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	// STEP 5-1: uncomment this line
	// _ "github.com/mattn/go-sqlite3"
)

var errImageNotFound = errors.New("image not found")

type Item struct {
	ID   			int    `db:"id" json:"-"`
	Name 			string `db:"name" json:"name"`
	Category 	string `db:"category" json:"category"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// fileName is the path to the JSON file storing items.
	fileName string
}

// NewItemReposit creates a new itemRepository.
func NewItemRepository() ItemRepository {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(wd)
	}
	fmt.Println("here2")
	fmt.Println(wd)
	return &itemRepository{fileName: "./db/items.json"}
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-1: add an implementation to store an item
	// 既存データの読み込み
	fmt.Println("here3")
	items, err := i.loadItems()
	if err != nil {
		return err
	}
	fmt.Println("items:")
	fmt.Println(items)

	// 新規データの追加
	items = append(items, *item)

	// ファイルへの書き込み
	return i.saveItems(items)
}

// loadItems loads items from the JSON file.
func (i *itemRepository) loadItems() ([]Item, error) {
	file, err := os.Open(i.fileName)
	if err != nil {
		return nil, err
	}
	fmt.Println("file")
	fmt.Println(file)
	// defer make sure the file is closed after the function returns
	defer file.Close()

	var items []Item
	// make a new JSON decoder
	decoder := json.NewDecoder(file)
	// decode JSON from the file and store it in the items variable
	err = decoder.Decode(&items)
	if err != nil {
		fmt.Println("decode error")
		return nil, err
	}

	// debug print
	fmt.Println(items)

	return items, nil
}

// saveItems saves items to the JSON file.
func (i *itemRepository) saveItems(items []Item) error {
	file, err := os.Create(i.fileName)
	if err != nil {
		return err
	}
	// defer make sure the file is closed after the function returns
	defer file.Close()

	// make a new JSON encoder
	encoder := json.NewEncoder(file)
	// encode items to JSON and write it to the file
	err = encoder.Encode(items)
	if err != nil {
		return err
	}

	return nil
}

// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image

	return nil
}
