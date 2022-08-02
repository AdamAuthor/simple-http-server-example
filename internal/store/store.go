package store

import (
	"context"
	"kbtu_go_6/internal/models"
)

type Store interface {
	Create(ctx context.Context, laptop *models.Laptop) error
	All(ctx context.Context) ([]*models.Laptop, error)
	ByID(ctx context.Context, id int) (*models.Laptop, error)
	Update(ctx context.Context, laptop *models.Laptop) error
	Delete(ctx context.Context, id int) error
}

// in future !!!
//type LaptopsRepository interface {
//	Create(ctx context.Context, laptop *models.Laptop) error
//	All(ctx context.Context) ([]*models.Laptop, error)
//	ByID(ctx context.Context, id int) (*models.Laptop, error)
//	Update(ctx context.Context, laptop *models.Laptop) error
//	Delete(ctx context.Context, id int) error
//}
