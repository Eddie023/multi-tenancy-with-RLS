package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, uri string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("DB parse config: %w", err)
	}

	// set the tenant id into this connection's setting
	cfg.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		tenantID, err := GetTenantID(ctx)
		if err != nil {
			slog.Error("unable to get tenantID from ctx", "err", err)
			return false
		}

		slog.Info("db pool: before acquire", "current tenantID", tenantID)

		_, err = conn.Exec(ctx, "SELECT set_tenant($1)", tenantID)
		if err != nil {
			slog.Error("unable to set tenantID", "err", err)
			return false
		}

		return true
	}

	// set the setting to be empty before this connection is released to the pool
	cfg.AfterRelease = func(conn *pgx.Conn) bool {
		_, err = conn.Exec(ctx, "SELECT set_tenant($1)", "0")
		if err != nil {
			slog.Error("unable to release tenantID", "err", err)
			return false
		}

		slog.Info("db pool: after release", "resetting tenantID", 0)

		return true
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("DB connect: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("DB ping: %w", err)
	}

	return &DB{pool: pool}, nil
}

func (d *DB) Close() {
	d.pool.Close()
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	TenantID    int     `json:"tenantId"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// GetProducts returns all products for specific tenant based on tenantID from context.
func (d *DB) GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := d.pool.Query(ctx, `SELECT id, tenant_id, name, description, price FROM product;`)
	if err != nil {
		slog.Error("query failed", "err", err)
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product

		err = rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.Price)
		if err != nil {
			slog.Error("query failed", "err", err)
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
