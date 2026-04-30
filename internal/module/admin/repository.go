package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// AdminRepository is the persistence boundary for admins (implemented by mysqlAdminRepository).
type AdminRepository interface {
	CreateAdmin(ctx context.Context, admin *Admin) error
	FindByEmail(ctx context.Context, email string) (*Admin, error)
	FindByMobile(ctx context.Context, mobile string) (*Admin, error)
}

type mysqlAdminRepository struct {
	db *sql.DB
}

// NewRepository returns a MySQL-backed AdminRepository.
func NewRepository(db *sql.DB) AdminRepository {
	return &mysqlAdminRepository{db: db}
}

func (r *mysqlAdminRepository) FindByEmail(ctx context.Context, email string) (*Admin, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, password, mobile, role, created_at, updated_at
		 FROM admins WHERE email = ? LIMIT 1`,
		email,
	)
	return scanAdmin(row)
}

func (r *mysqlAdminRepository) FindByMobile(ctx context.Context, mobile string) (*Admin, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, first_name, last_name, email, password, mobile, role, created_at, updated_at
		 FROM admins WHERE mobile = ? LIMIT 1`,
		mobile,
	)
	return scanAdmin(row)
}

func (r *mysqlAdminRepository) CreateAdmin(ctx context.Context, admin *Admin) error {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO admins (first_name, last_name, email, password, mobile, role, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		admin.FirstName,
		admin.LastName,
		admin.Email,
		admin.Password,
		admin.Mobile,
		admin.Role,
		admin.CreatedAt,
		admin.UpdatedAt,
	)
	if err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return err
		}
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	admin.ID = id
	return nil
}

func scanAdmin(row *sql.Row) (*Admin, error) {
	var a Admin
	if err := row.Scan(
		&a.ID,
		&a.FirstName,
		&a.LastName,
		&a.Email,
		&a.Password,
		&a.Mobile,
		&a.Role,
		&a.CreatedAt,
		&a.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}
