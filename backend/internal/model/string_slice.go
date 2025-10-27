package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Custom implementation of the []string serializer

type StringSlice []string

// Value implements the driver.Valuer interface.
// This defines how the slice is stored in the database.
// Due to commas being dangerous no element may include a comma
func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "", nil
	}

	for _, v := range s {
		if strings.Contains(v, ",") {
			return "", fmt.Errorf("unsafe string, %s", s)
		}
	}

	return strings.Join(s, ","), nil
}

// Scan implements the sql.Scanner interface.
// This defines how the database value is converted back into go.
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		b, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("failed to scan StringSlice, %v", value)
		}

		str = string(b)
	}

	if str == "" {
		*s = []string{}
	} else {
		*s = strings.Split(str, ",")
	}

	return nil
}
