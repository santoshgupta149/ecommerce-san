package product

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, p *Product) (*Product, error)
	Update(ctx context.Context, p *Product) (*Product, error)
}

type mysqlProductRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) ProductRepository {
	return &mysqlProductRepository{db: db}
}

func (r *mysqlProductRepository) GetByID(ctx context.Context, id int64) (*Product, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, sku, description, category, brand, price, stock, image_url, is_active, created_at, updated_at
		 FROM products
		 WHERE id = ?`,
		id,
	)

	return scanProduct(row)
}

func (r *mysqlProductRepository) List(ctx context.Context) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, sku, description, category, brand, price, stock, image_url, is_active, created_at, updated_at
		 FROM products
		 ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *mysqlProductRepository) Create(ctx context.Context, p *Product) (*Product, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO products (name, sku, description, category, brand, price, stock, image_url, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name,
		p.SKU,
		p.Description,
		p.Category,
		p.Brand,
		p.Price,
		p.Stock,
		p.ImageURL,
		p.IsActive,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		if isDuplicateSKUError(err) {
			return nil, ErrDuplicateSKU
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (r *mysqlProductRepository) Update(ctx context.Context, p *Product) (*Product, error) {
	res, err := r.db.ExecContext(ctx,
		`UPDATE products
		 SET name = ?, sku = ?, description = ?, category = ?, brand = ?, price = ?, stock = ?, image_url = ?, is_active = ?, updated_at = ?
		 WHERE id = ?`,
		p.Name,
		p.SKU,
		p.Description,
		p.Category,
		p.Brand,
		p.Price,
		p.Stock,
		p.ImageURL,
		p.IsActive,
		p.UpdatedAt,
		p.ID,
	)
	if err != nil {
		if isDuplicateSKUError(err) {
			return nil, ErrDuplicateSKU
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	return r.GetByID(ctx, p.ID)
}

type productScanner interface {
	Scan(dest ...any) error
}

func scanProduct(scanner productScanner) (*Product, error) {
	var priceRaw sql.RawBytes
	p := &Product{}
	if err := scanner.Scan(
		&p.ID,
		&p.Name,
		&p.SKU,
		&p.Description,
		&p.Category,
		&p.Brand,
		&priceRaw,
		&p.Stock,
		&p.ImageURL,
		&p.IsActive,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	price, err := strconv.ParseFloat(string(priceRaw), 64)
	if err != nil {
		return nil, err
	}
	p.Price = price

	return p, nil
}

func isDuplicateSKUError(err error) bool {
	var me *mysql.MySQLError
	if !errors.As(err, &me) || me.Number != 1062 {
		return false
	}

	msg := strings.ToLower(me.Message)
	return strings.Contains(msg, "uk_products_sku") || (strings.Contains(msg, "for key") && strings.Contains(msg, "sku"))
}
