package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4/stdlib"

	"OtusGolang/23-sql/internal/repository"
)

var _ repository.BaseRepo = (*Repo)(nil)

type Repo struct {
	db *sql.DB
}

func (r *Repo) Connect(ctx context.Context, dsn string) (err error) {
	r.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return r.db.PingContext(ctx)
}

///go:embed migrations/*.sql
//var embedMigrations embed.FS

func (r *Repo) Migrate(ctx context.Context, migrate string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(r.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) GetBooks(ctx context.Context) ([]repository.Book, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT title, created_at, updated_at, meta FROM books
	`)
	if err != nil {
		return nil, fmt.Errorf("cannot select: %w", err)
	}
	defer rows.Close()

	var books []repository.Book

	for rows.Next() {
		var b repository.Book

		var updatedAt sql.NullTime

		if err := rows.Scan(
			&b.Title,
			&b.CreatedAt,
			&updatedAt,
			&b.Meta,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}

		if updatedAt.Valid {
			b.UpdatedAt = updatedAt.Time
		}
		books = append(books, b)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return books, rows.Err()
}
