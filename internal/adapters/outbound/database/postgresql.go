package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAdapter struct {
	pool *pgxpool.Pool
}

func NowConnection(cfg DBConfig) (*PostgresAdapter, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config : %w", err)
	}

	configurePool(poolCfg, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection pool : %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres database")
	}

	adapter := &PostgresAdapter{
		pool: pool,
	}

	return adapter, nil
}

func configurePool(poolCfg *pgxpool.Config, cfg DBConfig) {
	maxConns := cfg.MaxConns
	if maxConns <= 0 {
		maxConns = 25
	}

	minConns := cfg.MinConns
	if minConns < 0 {
		minConns = 5
	}
	if minConns > maxConns {
		minConns = maxConns
	}

	poolCfg.MaxConns = int32(maxConns)
	poolCfg.MinConns = int32(minConns)
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = 1 * time.Minute
}

func (a *PostgresAdapter) Stats() *pgxpool.Stat { return a.pool.Stat() }
func (a *PostgresAdapter) Exec(ctx context.Context, expr string, intrfc ...interface{}) (pgconn.CommandTag, error) {
	return a.Exec(ctx, expr, intrfc)
}

func (a *PostgresAdapter) Query(ctx context.Context, expr string, intrfc ...interface{}) (pgx.Rows, error) {
	return a.Query(ctx, expr, intrfc)
}

func (a *PostgresAdapter) QueryROW(ctx context.Context, expr string, intrfc ...interface{}) pgx.Rows {
	return a.QueryROW(ctx, expr, intrfc)
}
