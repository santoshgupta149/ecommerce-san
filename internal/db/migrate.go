package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

// InitSchema creates the minimal tables needed for the current learning scaffold.
// In a real production system, use migrations (e.g., golang-migrate, goose).
func (d *DB) InitSchema(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			mobile VARCHAR(20) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			UNIQUE KEY uk_users_email (email),
			UNIQUE KEY uk_users_mobile (mobile)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL,
			total DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			INDEX idx_orders_user_id (user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS admins (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			mobile VARCHAR(20) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'admin',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			UNIQUE KEY uk_admins_email (email),
			UNIQUE KEY uk_admins_mobile (mobile)
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			sku VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			category VARCHAR(255) NOT NULL,
			brand VARCHAR(255) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			stock BIGINT NOT NULL DEFAULT 0,
			image_url TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			UNIQUE KEY uk_products_sku (sku)
		)`,
	}

	for _, s := range stmts {
		if _, err := d.SQL.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("init schema: %w", err)
		}
	}

	// Legacy DBs: add `mobile` if missing (no UNIQUE here — next step adds the index).
	if _, err := d.SQL.ExecContext(ctx,
		`ALTER TABLE users ADD COLUMN mobile VARCHAR(20) NOT NULL DEFAULT '' AFTER email`,
	); err != nil {
		var me *mysql.MySQLError
		if !errors.As(err, &me) || me.Number != 1060 {
			return fmt.Errorf("migrate users.mobile column: %w", err)
		}
	}

	if err := d.ensureUsersMobileUnique(ctx); err != nil {
		return err
	}
	if err := d.ensureProductsSchema(ctx); err != nil {
		return err
	}
	return nil
}

const addUsersMobileUniqueSQL = `ALTER TABLE users ADD UNIQUE KEY uk_users_mobile (mobile)`
const addProductsSKUUniqueSQL = `ALTER TABLE products ADD UNIQUE KEY uk_products_sku (sku)`

// ensureUsersMobileUnique adds uk_users_mobile. If duplicates block the index (1062), we remove
// duplicate rows keeping the smallest id per mobile, then retry once.
// If you have FKs to users, review deletions or fix data manually instead.
func (d *DB) ensureUsersMobileUnique(ctx context.Context) error {
	_, err := d.SQL.ExecContext(ctx, addUsersMobileUniqueSQL)
	if err == nil {
		return nil
	}

	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return fmt.Errorf("migrate users.mobile unique: %w", err)
	}
	if me.Number == 1061 {
		return nil
	}
	if me.Number != 1062 {
		return fmt.Errorf("migrate users.mobile unique: %w", err)
	}

	log.Printf("migration: duplicate mobile values block UNIQUE; deleting duplicate user rows (keeping smallest id per mobile), then retrying index")

	res, err := d.SQL.ExecContext(ctx, `
		DELETE u1 FROM users u1
		INNER JOIN users u2 ON u1.mobile = u2.mobile AND u1.id > u2.id
	`)
	if err != nil {
		return fmt.Errorf("dedupe users.mobile: %w", err)
	}
	n, _ := res.RowsAffected()
	if n > 0 {
		log.Printf("migration: removed %d duplicate user row(s) by mobile", n)
	}

	_, err = d.SQL.ExecContext(ctx, addUsersMobileUniqueSQL)
	if err != nil {
		var me2 *mysql.MySQLError
		if errors.As(err, &me2) && me2.Number == 1061 {
			return nil
		}
		return fmt.Errorf("migrate users.mobile unique (after dedupe): %w", err)
	}
	return nil
}

func (d *DB) ensureProductsSchema(ctx context.Context) error {
	stmts := []string{
		`ALTER TABLE products ADD COLUMN sku VARCHAR(100) NOT NULL AFTER name`,
		`ALTER TABLE products ADD COLUMN description TEXT NOT NULL AFTER sku`,
		`ALTER TABLE products ADD COLUMN category VARCHAR(255) NOT NULL AFTER description`,
		`ALTER TABLE products ADD COLUMN brand VARCHAR(255) NOT NULL AFTER category`,
		`ALTER TABLE products ADD COLUMN price DECIMAL(10,2) NOT NULL AFTER brand`,
		`ALTER TABLE products ADD COLUMN stock BIGINT NOT NULL DEFAULT 0 AFTER price`,
		`ALTER TABLE products ADD COLUMN image_url TEXT NOT NULL AFTER stock`,
		`ALTER TABLE products ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE AFTER image_url`,
		`ALTER TABLE products ADD COLUMN created_at TIMESTAMP NOT NULL AFTER is_active`,
		`ALTER TABLE products ADD COLUMN updated_at TIMESTAMP NOT NULL AFTER created_at`,
	}

	for _, stmt := range stmts {
		if _, err := d.SQL.ExecContext(ctx, stmt); err != nil {
			var me *mysql.MySQLError
			if !errors.As(err, &me) || me.Number != 1060 {
				return fmt.Errorf("migrate products schema: %w", err)
			}
		}
	}

	_, err := d.SQL.ExecContext(ctx, addProductsSKUUniqueSQL)
	if err == nil {
		return nil
	}

	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return fmt.Errorf("migrate products.sku unique: %w", err)
	}
	if me.Number == 1061 {
		return nil
	}
	return fmt.Errorf("migrate products.sku unique: %w", err)
}
