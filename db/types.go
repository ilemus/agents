package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONTags is a custom type representing an array of tags stored as a JSON array in SQLite/Turso.
type JSONTags []string

// Value implements the driver.Valuer interface.
func (j JSONTags) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONTags) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("failed to scan JSONTags: unsupported type %T", value)
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, j)
}

// Float32Vector is a custom type representing a vector of float32 stored as a JSON array in SQLite/Turso.
type Float32Vector []float32

// Value converts the slice into a JSON-marshaled string/bytes for the database
func (v Float32Vector) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	// Marshal the slice to a JSON byte array
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal float32 vector: %w", err)
	}
	// Return it as a string (or []byte, which drivers handle cleanly)
	return string(bytes), nil
}

// Scan converts the database value (JSON string/bytes) back into a Go slice
func (v *Float32Vector) Scan(value interface{}) error {
	if value == nil {
		*v = nil
		return nil
	}

	var bytes []byte
	switch data := value.(type) {
	case []byte:
		bytes = data
	case string:
		bytes = []byte(data)
	default:
		return fmt.Errorf("unsupported data type for Float32Vector: %T", value)
	}

	return json.Unmarshal(bytes, v)
}
