package app

import (
	"context"
	"errors"
	"fmt"
	// STEP 5-1: uncomment this line
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

var errImageNotFound = errors.New("image not found")

type Item struct {
	ID   			int    `db:"id" json:"-"`
	Name 			string `db:"name" json:"name"`
	Category 	string `db:"category" json:"category"`
	ImageName string `db:"image_name" json:"image_name"`
}

// Please run `go generate ./...` to generate the mock implementation
// ItemRepository is an interface to manage items.
//
//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -package=${GOPACKAGE} -destination=./mock_$GOFILE
type ItemRepository interface {
	Insert(ctx context.Context, item *Item) error
	GetItems() ([]Item, error)
}

// itemRepository is an implementation of ItemRepository
type itemRepository struct {
	// db is a database connection
	db *sql.DB
}

// NewItemRepository connects db and creates a new itemRepository.
func NewItemRepository(dbPath string) (ItemRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	// if the table does not exist, create it
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			image_name TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return &itemRepository{db: db}, nil
}

// Insert inserts an item into the repository.
func (i *itemRepository) Insert(ctx context.Context, item *Item) error {
	_, err := i.db.Exec(
		"INSERT INTO items (name, category, image_name) VALUES (?, ?, ?)",
		item.Name, item.Category, item.ImageName,
	)
	if err != nil {
		return fmt.Errorf("failed to insert an item: %w", err)
	}
	return nil
}

// GetItems returns all items from the repository.
func (i *itemRepository) GetItems() ([]Item, error) {
	rows, err := i.db.Query("SELECT id, name, category, image_name FROM items")
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	// defer make sure rows.Close() is called after the function returns
	defer rows.Close()

	var items []Item
	// iterate over the rows
	for rows.Next() {
		var item Item
		err:= rows.Scan(&item.ID, &item.Name, &item.Category, &item.ImageName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}


// StoreImage stores an image and returns an error if any.
// This package doesn't have a related interface for simplicity.
func StoreImage(fileName string, image []byte) error {
	// STEP 4-4: add an implementation to store an image

	return nil
}
