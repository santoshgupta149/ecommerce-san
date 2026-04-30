package order

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
)

// OrderRepository persists orders to MySQL (connection from .env via bootstrap → db.New).
type OrderRepository interface {
	GetByID(ctx context.Context, id int64) (*Order, error)
	Create(ctx context.Context, o *Order) (*Order, error)
}

type mysqlOrderRepository struct {
	db *sql.DB
}

// NewRepository returns a repository backed by the shared *sql.DB pool (MySQL).
func NewRepository(db *sql.DB) OrderRepository {
	return &mysqlOrderRepository{db: db}
}

func (r *mysqlOrderRepository) GetByID(ctx context.Context, id int64) (*Order, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, total, created_at
		 FROM orders
		 WHERE id = ?`,
		id,
	)

	var totalRaw sql.RawBytes
	o := &Order{}
	if err := row.Scan(&o.ID, &o.UserID, &totalRaw, &o.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	total, err := strconv.ParseFloat(string(totalRaw), 64)
	if err != nil {
		return nil, err
	}
	o.Total = total

	return o, nil
}

func (r *mysqlOrderRepository) Create(ctx context.Context, o *Order) (*Order, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO orders (user_id, total, created_at)
		 VALUES (?, ?, ?)`,
		o.UserID,
		o.Total,
		o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	o.ID = id

	return o, nil
}
