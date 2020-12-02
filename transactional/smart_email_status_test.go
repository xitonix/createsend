package transactional

import (
	"fmt"
	"testing"
)

func TestSmartEmailStatus_MarshalJSON(t *testing.T) {
	testCases := []struct {
		title    string
		status   SmartEmailStatus
		expected string
	}{
		{
			title:    "Unknown",
			expected: fmt.Sprintf("%q", unknownSmartEmailStr),
		},
		{
			title:    "Active",
			status:   ActiveSmartEmail,
			expected: fmt.Sprintf("%q", activeSmartEmailStr),
		},
		{
			title:    "Draft",
			status:   DraftSmartEmail,
			expected: fmt.Sprintf("%q", draftSmartEmailStr),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			marshalled, _ := tC.status.MarshalJSON()
			if tC.expected != string(marshalled) {
				t.Errorf("Expected %s, Actual: %s", tC.expected, string(marshalled))
			}
		})
	}
}
