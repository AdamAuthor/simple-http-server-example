package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"kbtu_go_6/internal/models"
	"kbtu_go_6/internal/store"
)

func (D *DB) Goods() store.GoodsRepository {
	if D.goods == nil {
		D.goods = NewGoodsRepository(D.conn)
	}

	return D.goods
}

type GoodsRepository struct {
	conn *sqlx.DB
}

func NewGoodsRepository(conn *sqlx.DB) store.GoodsRepository {
	return &GoodsRepository{conn: conn}
}

func (g GoodsRepository) Create(ctx context.Context, good *models.Good) error {
	_, err := g.conn.Exec("INSERT INTO goods(name, category_id) VALUES ($1, $2)", good.Name, good.CategoryID)
	if err != nil {
		return err
	}
	return nil
}

func (g GoodsRepository) All(ctx context.Context) ([]*models.Good, error) {
	goods := make([]*models.Good, 0)
	if err := g.conn.Select(&goods, "SELECT * FROM goods"); err != nil {
		return nil, err
	}
	return goods, nil
}

func (g GoodsRepository) ByID(ctx context.Context, id int) (*models.Good, error) {
	good := new(models.Good)
	if err := g.conn.Get(good, "SELECT id, name, category_id FROM goods WHERE id=$1", id); err != nil {
		return nil, err
	}
	return good, nil
}

func (g GoodsRepository) Update(ctx context.Context, good *models.Good) error {
	_, err := g.conn.Exec("UPDATE goods SET name = $1 WHERE id = $2", good.Name, good.ID)
	if err != nil {
		return err
	}
	return nil
}

func (g GoodsRepository) Delete(ctx context.Context, id int) error {
	_, err := g.conn.Exec("DELETE FROM goods WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
