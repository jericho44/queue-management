package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_ = autoMigrate(db)

	return db, nil
}

func autoMigrate(db *sql.DB) error {
	queries := []string{
		`ALTER TABLE branches ADD COLUMN IF NOT EXISTS kiosk_enabled BOOLEAN NOT NULL DEFAULT TRUE`,
		`ALTER TABLE branches ADD COLUMN IF NOT EXISTS kiosk_mode VARCHAR(20) NOT NULL DEFAULT 'DUAL'`,
		`ALTER TABLE branches ADD COLUMN IF NOT EXISTS paper_size VARCHAR(10) NOT NULL DEFAULT '58mm'`,
		`ALTER TABLE branches ADD COLUMN IF NOT EXISTS receipt_header TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE branches ADD COLUMN IF NOT EXISTS receipt_footer TEXT NOT NULL DEFAULT 'Terima kasih atas kunjungan Anda. Harap menyimpan struk ini hingga dipanggil.'`,
		`ALTER TABLE branches ADD COLUMN IF NOT EXISTS auto_print BOOLEAN NOT NULL DEFAULT FALSE`,
		`CREATE TABLE IF NOT EXISTS invoices (
			id BIGSERIAL PRIMARY KEY,
			uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
			organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			invoice_number VARCHAR(100) NOT NULL UNIQUE,
			billing_period VARCHAR(7) NOT NULL,
			subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
			tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
			total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
			status VARCHAR(50) NOT NULL DEFAULT 'UNPAID',
			due_date TIMESTAMPTZ NOT NULL,
			paid_at TIMESTAMPTZ,
			snap_token VARCHAR(255),
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMPTZ
		)`,
		`CREATE TABLE IF NOT EXISTS invoice_items (
			id BIGSERIAL PRIMARY KEY,
			uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
			invoice_id BIGINT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
			description VARCHAR(255) NOT NULL,
			quantity INT NOT NULL DEFAULT 1,
			unit_price DECIMAL(12,2) NOT NULL DEFAULT 0,
			amount DECIMAL(12,2) NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS payments (
			id BIGSERIAL PRIMARY KEY,
			uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
			invoice_id BIGINT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
			transaction_id VARCHAR(100) UNIQUE,
			payment_type VARCHAR(50),
			gross_amount DECIMAL(12,2) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
			payload JSONB,
			paid_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		_, _ = db.Exec(q)
	}
	return nil
}

