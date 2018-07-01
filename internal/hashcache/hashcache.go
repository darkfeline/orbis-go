// Package hashcache implements a cache for file content hash values.
package hashcache

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// HashCache caches the hash value of files.
type HashCache struct {
	db *sql.DB
}

// New creates a new HashCache.
func New(p string) (*HashCache, error) {
	d, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, err
	}
	if err := setupTable(d); err != nil {
		d.Close()
		return nil, err
	}
	return &HashCache{d}, nil
}

func (c *HashCache) Close() error {
	return c.db.Close()
}

type noRow struct{}

func (n noRow) Error() string {
	return "no such row"
}

func IsNoRow(e error) bool {
	_, ok := e.(noRow)
	return ok
}

// Get retrieves the cached Hash value for a file.
func (c *HashCache) Get(path string, i os.FileInfo) (string, error) {
	q := `SELECT hexdigest FROM sha256_cache
WHERE path=? AND mtime=? AND size=?`
	rows, err := c.db.Query(q, path, i.ModTime(), i.Size())
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", err
		}
		return "", noRow{}
	}
	var s string
	err = rows.Scan(&s)
	if err != nil {
		panic(err)
	}
	return s, nil
}

// Set sets the cached Hash value for a file.
func (c *HashCache) Set(path string, i os.FileInfo, hash string) error {
	q := `INSERT OR REPLACE INTO sha256_cache
(path, mtime, size, hexdigest) VALUES (?, ?, ?, ?)`
	_, err := c.db.Exec(q, path, i.ModTime(), i.Size(), hash)
	return err
}

func setupTable(d *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS sha256_cache (
path TEXT NOT NULL,
mtime INT NOT NULL,
size INT NOT NULL,
hexdigest TEXT NOT NULL,
CONSTRAINT path_u UNIQUE (path))`
	_, err := d.Exec(q)
	return err
}
