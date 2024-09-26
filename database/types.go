package database

import (
	"database/sql"
	"time"
)

func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullTime returns a sql.NullTime with the given time.Time. If the time is
// the zero value, the NullTime is invalid.
func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: t != time.Time{},
	}
}
