package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// UserRepository persists users to MySQL (connection from .env via bootstrap → db.New).
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
}

var ErrDuplicateEmail = errors.New("email already exists")
var ErrDuplicateMobile = errors.New("mobile number already exists")

type mysqlUserRepository struct {
	db *sql.DB
}

// NewRepository returns a repository backed by the shared *sql.DB pool (MySQL).
func NewRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, mobile, created_at
		 FROM users
		 WHERE id = ?`,
		id,
	)

	u := &User{}
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Mobile, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *mysqlUserRepository) Create(ctx context.Context, u *User) (*User, error) {
	// Check first so we don’t hit INSERT + 1062. (No row is committed on duplicate INSERT, but
	// InnoDB can still advance AUTO_INCREMENT, which looks like “skipped” ids.)
	if err := r.ensureEmailAndMobileUnique(ctx, u); err != nil {
		return nil, err
	}

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, email, mobile, created_at)
		 VALUES (?, ?, ?, ?)`,
		u.Name,
		u.Email,
		u.Mobile,
		u.CreatedAt,
	)
	if err != nil {
		var me *mysql.MySQLError
		// Race: two requests can pass the SELECTs then one INSERT wins.
		if errors.As(err, &me) && me.Number == 1062 {
			msg := strings.ToLower(me.Message)
			if strings.Contains(msg, "uk_users_mobile") {
				return nil, ErrDuplicateMobile
			}
			if strings.Contains(msg, "uk_users_email") || (strings.Contains(msg, "for key") && strings.Contains(msg, "email")) {
				return nil, ErrDuplicateEmail
			}
			return nil, err
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (r *mysqlUserRepository) ensureEmailAndMobileUnique(ctx context.Context, u *User) error {
	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM users WHERE email = ? LIMIT 1`, u.Email).Scan(&one)
	if err == nil {
		return ErrDuplicateEmail
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = r.db.QueryRowContext(ctx, `SELECT 1 FROM users WHERE mobile = ? LIMIT 1`, u.Mobile).Scan(&one)
	if err == nil {
		return ErrDuplicateMobile
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}
