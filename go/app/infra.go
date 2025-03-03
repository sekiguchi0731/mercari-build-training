package app

import (
	"context"
	"encoding/json"
	"errors"
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

// NewItemRepository creates a new itemRepository.
func NewItemRepository() ItemRepository {
	return &itemRepository{fileName: "./db/items.json"}
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	// STEP 4-2: add an implementation to store an item
	// 既存データの読み込み
	items, err := i.loadItems()
	if err != nil {
		return err
	}

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
	// defer make sure the file is closed after the function returns
	defer file.Close()

	var items []Item
	// make a new JSON decoder
	decoder := json.NewDecoder(file)
	// decode JSON from the file and store it in the items variable
	err = decoder.Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// saveItems saves items to the JSON file.
func (i *itemRepository) saveItems(items []Item) error {
	var file *os.File
	// check if the file exists
	if _, err := os.Stat(i.fileName); os.IsNotExist(err) {
		// if the file doesn't exist, create it
		file, err = os.Create(i.fileName)
		if err != nil {
			return err
		}
	} else {
		// if the file exists, open it and overwrite it
		// O_WRONLY: write only, O_TRUNC: truncate the file
		file, err = os.OpenFile(i.fileName, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
	}
	// defer make sure the file is closed after the function returns
	defer file.Close()

	// make a new JSON encoder
	encoder := json.NewEncoder(file)
	// encode items to JSON and write it to the file
	err := encoder.Encode(items)
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
