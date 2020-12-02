package transactional

import (
	"encoding/json"
	"strings"
)

// SmartEmailStatus represents a smart email status.
type SmartEmailStatus uint8

const (
	// UnknownSmartEmail unknown status.
	UnknownSmartEmail SmartEmailStatus = iota
	// ActiveSmartEmail active smart email.
	ActiveSmartEmail
	// DraftSmartEmail draft smart email.
	DraftSmartEmail
)

const (
	unknownSmartEmailStr = `unknown`
	activeSmartEmailStr  = `active`
	draftSmartEmailStr   = `draft`
)

var (
	smartEmailStatusToValue = map[string]SmartEmailStatus{
		activeSmartEmailStr: ActiveSmartEmail,
		draftSmartEmailStr:  DraftSmartEmail,
	}

	smartEmailStatusFromValue = map[SmartEmailStatus]string{
		ActiveSmartEmail: activeSmartEmailStr,
		DraftSmartEmail:  draftSmartEmailStr,
	}
)

// MarshalJSON marshal the object into json bytes.
func (r SmartEmailStatus) MarshalJSON() ([]byte, error) {
	typeStr, ok := smartEmailStatusFromValue[r]
	if !ok {
		return json.Marshal(unknownSmartEmailStr)
	}
	return json.Marshal(typeStr)
}

// UnmarshalJSON unmarshal json bytes back to object.
func (r *SmartEmailStatus) UnmarshalJSON(b []byte) error {
	value := strings.ToLower(strings.Trim(string(b), "\""))
	rt, ok := smartEmailStatusToValue[value]
	if !ok {
		rt = UnknownSmartEmail
	}
	*r = rt
	return nil
}

// String Stringer implementation
func (r *SmartEmailStatus) String() string {
	return smartEmailStatusFromValue[*r]
}
