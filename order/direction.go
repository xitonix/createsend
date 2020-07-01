package order

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Direction order direction.
type Direction int8

const (
	// ASC ascending.
	ASC Direction = iota
	// DESC descending.
	DESC
)

var (
	directionToString = map[Direction]string{
		ASC:  "asc",
		DESC: "desc",
	}

	stringToDirection = map[string]Direction{
		"asc":  ASC,
		"desc": DESC,
	}
)

// UnmarshalJSON parses the json bytes into a Direction value.
func (d *Direction) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return fmt.Errorf("order direction should be a string, got %s", bytes)
	}
	direction, ok := stringToDirection[strings.ToLower(value)]
	if !ok {
		return fmt.Errorf("invalid order direction %q", value)
	}
	*d = direction
	return nil
}

// String Stringer implementation
func (d Direction) String() string {
	return directionToString[d]
}
