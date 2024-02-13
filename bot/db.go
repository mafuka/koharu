package bot

import (
	"database/sql"

	"golang.org/x/xerrors"
)

// DB includes a database instance.
type DB struct {
	db *sql.DB
}

// NewDB connects to the sqlite3 database located at path, or creates one if it does not exist.
func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, xerrors.Errorf("%w", err)
	}

	return &DB{db}, nil
}

func (d *DB) Close() error {
	if err := d.db.Close(); err != nil {
		return xerrors.Errorf("Error closing database: %w", err)
	}
	return nil
}

func (d *DB) Init() error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	if err := createConfigTable(tx); err != nil {
		return err
	}

	return nil
}

func createConfigTable(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS config (
		id INTEGER PRIMARY KEY,
		key TEXT,
		value TEXT
	)`)
	return err
}
