package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/jmoiron/sqlx"
	"github.com/pgx-contrib/pgxotel"
	"go.uber.org/fx"
	"log/slog"
	"platform/pkg/ctxlog"
	"platform/pkg/db/dblog"
	"platform/pkg/db/dbtx"
	"platform/pkg/db/tx"
	"time"
)

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`

	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	LogTracing      bool          `yaml:"log_tracing"`
}

func (c DbConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
}

func NewDB(cfg DbConfig, logger *slog.Logger, lc fx.Lifecycle) (dbtx.DBTX, *sql.DB, error) {
	pgConfig, err := pgx.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse postgres connection string: %w", err)
	}
	if cfg.LogTracing {
		pgConfig.Tracer = &tracelog.TraceLog{
			Logger:   dblog.NewLogger(),
			LogLevel: tracelog.LogLevelDebug,
		}
	} else {
		pgConfig.Tracer = pgxotel.NewTracer("pgx")
	}

	connStr := stdlib.RegisterConnConfig(pgConfig)

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open db connection: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := db.PingContext(ctxlog.WithLogger(ctx, logger))
			if err != nil {
				return fmt.Errorf("failed to ping db: %w", err)
			}
			return nil
		},
	})

	return tx.Wrap(db), db.DB, nil
}
