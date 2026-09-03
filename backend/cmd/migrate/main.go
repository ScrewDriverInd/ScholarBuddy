package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type migration struct {
	version  int64
	upPath   string
	downPath string
}

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down" && os.Args[1] != "reset") {
		fail("usage: go run ./cmd/migrate [up|down|reset]")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fail("DATABASE_URL is required")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fail("connect to database: %v", err)
	}
	defer db.Close()
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	err = db.Ping(pingCtx)
	cancel()
	if err != nil {
		fail("ping database: %v", err)
	}

	var runErr error
	switch os.Args[1] {
	case "up":
		runErr = migrateUp(ctx, db)
	case "down":
		runErr = migrateDown(ctx, db)
	case "reset":
		runErr = executeFile(ctx, db, "scripts/reset_database.sql")
	}
	if runErr != nil {
		fail("migration failed: %v", runErr)
	}
	fmt.Printf("migration %s completed\n", os.Args[1])
}

func migrateUp(ctx context.Context, db *pgxpool.Pool) error {
	if _, err := db.Exec(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())"); err != nil {
		return err
	}
	migrations, err := discoverMigrations()
	if err != nil {
		return err
	}
	for _, item := range migrations {
		var applied bool
		if err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)", item.version).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := executeMigration(ctx, db, item.upPath, item.version, true); err != nil {
			return err
		}
		fmt.Printf("applied %s\n", filepath.Base(item.upPath))
	}
	return nil
}

func migrateDown(ctx context.Context, db *pgxpool.Pool) error {
	migrations, err := discoverMigrations()
	if err != nil {
		return err
	}
	var version int64
	if err := db.QueryRow(ctx, "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version); err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no applied migrations")
		}
		return err
	}
	for _, item := range migrations {
		if item.version == version {
			if item.downPath == "" {
				return fmt.Errorf("no down migration for version %d", version)
			}
			if err := executeMigration(ctx, db, item.downPath, item.version, false); err != nil {
				return err
			}
			fmt.Printf("rolled back %s\n", filepath.Base(item.downPath))
			return nil
		}
	}
	return fmt.Errorf("migration version %d is not present locally", version)
}

func discoverMigrations() ([]migration, error) {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}
	byVersion := map[int64]*migration{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}
		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}
		item := byVersion[version]
		if item == nil {
			item = &migration{version: version}
			byVersion[version] = item
		}
		path := filepath.Join("migrations", entry.Name())
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			item.upPath = path
		}
		if strings.HasSuffix(entry.Name(), ".down.sql") {
			item.downPath = path
		}
	}
	result := make([]migration, 0, len(byVersion))
	for _, item := range byVersion {
		if item.upPath == "" {
			return nil, fmt.Errorf("migration %d has no up file", item.version)
		}
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].version < result[j].version })
	return result, nil
}

func executeMigration(ctx context.Context, db *pgxpool.Pool, path string, version int64, up bool) error {
	sql, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	txCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	tx, err := db.BeginTx(txCtx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(txCtx)
	if _, err := tx.Exec(txCtx, string(sql)); err != nil {
		return err
	}
	if up {
		_, err = tx.Exec(txCtx, "INSERT INTO schema_migrations (version) VALUES ($1)", version)
	} else {
		_, err = tx.Exec(txCtx, "DELETE FROM schema_migrations WHERE version=$1", version)
	}
	if err != nil {
		return err
	}
	return tx.Commit(txCtx)
}

func executeFile(ctx context.Context, db *pgxpool.Pool, path string) error {
	sql, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	resetCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	tx, err := db.BeginTx(resetCtx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(resetCtx)
	if _, err := tx.Exec(resetCtx, string(sql)); err != nil {
		return err
	}
	return tx.Commit(resetCtx)
}

func fail(format string, args ...any) { fmt.Fprintf(os.Stderr, format+"\n", args...); os.Exit(1) }
