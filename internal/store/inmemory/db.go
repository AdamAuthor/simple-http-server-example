package inmemory

import (
	"context"
	"fmt"
	"kbtu_go_6/internal/models"
	"kbtu_go_6/internal/store"
	"sync"
)

// DB saves information about laptops
type DB struct {
	data map[int]*models.Laptop
	mu   *sync.RWMutex
}

// NewDB is the function for creating basic Database
func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Laptop),
		mu:   new(sync.RWMutex),
	}
}

// Create for creating new element in DB
func (db *DB) Create(ctx context.Context, laptop *models.Laptop) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[laptop.ID] = laptop

	return nil
}

// All is used for reading all elements in DB
func (db *DB) All(ctx context.Context) ([]*models.Laptop, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	laptops := make([]*models.Laptop, 0, len(db.data))
	for _, laptop := range db.data {
		laptops = append(laptops, laptop)
	}

	return laptops, nil
}

// ByID is used for reading elements by id in DB
func (db *DB) ByID(ctx context.Context, id int) (*models.Laptop, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	laptop, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No laptop eith id: %d", id)
	}

	return laptop, nil
}

// Update is used for updating elements in DB
func (db *DB) Update(ctx context.Context, laptop *models.Laptop) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[laptop.ID] = laptop
	return nil
}

// Delete is used for deleting elements by id in DB
func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, id)
	return nil
}
