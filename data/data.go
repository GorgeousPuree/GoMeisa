package data

import (
	"Gomeisa"
	"database/sql"
)

type Creator interface {
	Create() (int, error)
}

// No way to see an Insert error from QueryRow().
// The problem is if QueryRow() with an INSERT...RETURNING fails,
// then the error is a rather unhelpful no rows in result set.
// Details: https://github.com/lib/pq/issues/77
type DbRowFromInsert struct {
	rows *sql.Rows
}

// Using sql.Row.Scan() as a template for this method, but modifying the errors.
func (row *DbRowFromInsert) Scan(dest ...interface{}) error {
	defer row.rows.Close()
	row.rows.Next()

	// There may be no rows, but Scan anyway to get the errors...
	// If there are no rows because of a db constraint error this is when those errors will be returned.
	err := row.rows.Scan(dest...)
	if err != nil {
		return err
	}

	return nil
}

// Usage:
func InsertReturning(query string, args ...interface{}) (*DbRowFromInsert, error) {
	rows, err := Gomeisa.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return &DbRowFromInsert{rows: rows}, nil
}
