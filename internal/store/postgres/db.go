package postgres

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"kbtu_go_6/internal/store"
	"log"
)

type DB struct {
	conn       *sqlx.DB
	goods      store.GoodsRepository
	categories store.CategoriesRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (D *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	D.conn = conn
	log.Println("DB has successful pinging")
	return nil
}

func (D *DB) Close() error {
	return D.conn.Close()
}
