package db

import (
	"database/sql"
	_ "embed"

	appkit "github.com/TrueBlocks/trueblocks-art/packages/appkit/v2"
)

//go:embed schema.sql
var schema string

type DB struct {
	conn *sql.DB
}

type Item struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Open opens (creating if necessary) the database and applies the embedded
// schema. The schema must stay idempotent so reopening an existing database
// is safe.
func Open(path string) (*DB, error) {
	conn, err := appkit.OpenSQLite(path, &appkit.SQLiteOpts{Schema: schema})
	if err != nil {
		return nil, err
	}
	return &DB{conn: conn}, nil
}

func (d *DB) Close() error { return d.conn.Close() }

func (d *DB) ListItems() ([]Item, error) {
	rows, err := d.conn.Query("SELECT id, name FROM items ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := []Item{}
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Name); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (d *DB) AddItem(name string) error {
	_, err := d.conn.Exec("INSERT INTO items (name) VALUES (?)", name)
	return err
}
