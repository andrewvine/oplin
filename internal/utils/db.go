package utils

import (
	"database/sql"
	"encoding/json"
	"github.com/tabbed/pqtype"
	"time"
)

// ToPQRawMessageType converts a json.RawMessage to a pqtype.NullRawMessage
func ToPQRawMessageType(r json.RawMessage) pqtype.NullRawMessage {
	if r == nil || len(r) == 0 {
		return pqtype.NullRawMessage{Valid: false}
	}
	return pqtype.NullRawMessage{Valid: true, RawMessage: r}
}

// NullString returns a sql.NullString with the given string
func NullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

// NullInt64 returns a sql.NullInt64 with the given int64
func NullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

// NullTime returns a sql.NullTime with the given time.Time
func NullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// NowUTC returns the current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}

// NowUTCAsNullTime returns the current time in UTC as a sql.NullTime
func NowUTCAsNullTime() sql.NullTime {
	return sql.NullTime{Time: NowUTC(), Valid: true}
}

// IsNoRowsError returns true if the error is a sql.ErrNoRows
func IsNoRowsError(err error) bool {
	return err == sql.ErrNoRows
}

// MergeFacets merges two json.RawMessages into one
func MergeFacets(a, b []byte) ([]byte, error) {
	var mapA map[string]interface{}
	if len(a) > 0 {
		if err := json.Unmarshal(a, &mapA); err != nil {
			return nil, err
		}
	}
	var mapB map[string]interface{}
	if len(b) > 0 {
		if err := json.Unmarshal(b, &mapB); err != nil {
			return nil, err
		}
	}
	mapC := make(map[string]interface{})
	for k, v := range mapA {
		mapC[k] = v
	}
	for k, v := range mapB {
		mapC[k] = v
	}
	return json.Marshal(mapC)
}
