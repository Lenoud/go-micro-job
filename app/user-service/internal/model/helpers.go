package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// NullString nullable column wrapper
type NullString struct {
	sql.NullString
}

func (ns *NullString) Scan(value interface{}) error {
	return ns.NullString.Scan(value)
}

func (ns NullString) Value() string {
	if ns.Valid {
		return ns.NullString.String
	}
	return ""
}

// NullTime nullable time wrapper
type NullTime sql.NullTime

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullTime) Scan(value interface{}) error {
	return (*sql.NullTime)(nt).Scan(value)
}

func (nt NullTime) Value() (driver.Value, error) {
	return sql.NullTime(nt).Value()
}

func normalizeIDs(ids []string) []string {
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" {
			result = append(result, id)
		}
	}
	return result
}

func buildInClauseArgs(ids []string) (string, []interface{}) {
	normalized := normalizeIDs(ids)
	placeholders := make([]string, len(normalized))
	args := make([]interface{}, len(normalized))
	for i, id := range normalized {
		placeholders[i] = "?"
		args[i] = id
	}
	return strings.Join(placeholders, ","), args
}
