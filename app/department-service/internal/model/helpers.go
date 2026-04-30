package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

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
