package order

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SuppressionListField client suppression list order field.
type SuppressionListField int8

const (
	// BySuppressedEmailAddress email address.
	BySuppressedEmailAddress SuppressionListField = iota
	// BySuppressionDate suppression date.
	BySuppressionDate
)

var (
	suppressionFiledToString = map[SuppressionListField]string{
		BySuppressedEmailAddress: "email",
		BySuppressionDate:        "date",
	}

	stringToSuppressionFiled = map[string]SuppressionListField{
		"email": BySuppressedEmailAddress,
		"date":  BySuppressionDate,
	}
)

// UnmarshalJSON parses the json bytes into a SuppressionListField value.
func (f *SuppressionListField) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return fmt.Errorf("order-by filed should be a string, got %s", bytes)
	}
	field, ok := stringToSuppressionFiled[strings.ToLower(value)]
	if !ok {
		return fmt.Errorf("invalid order-by field %q", value)
	}
	*f = field
	return nil
}

// String Stringer implementation
func (f SuppressionListField) String() string {
	return suppressionFiledToString[f]
}
