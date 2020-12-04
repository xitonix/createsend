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

func TestSmartEmailStatus_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		title    string
		expected SmartEmailStatus
		status   string
	}{
		{
			title:    "Unknown",
			status:   "Unknown",
			expected: UnknownSmartEmail,
		},
		{
			title:    "random string",
			status:   "random",
			expected: UnknownSmartEmail,
		},
		{
			title:    "active lowercase",
			status:   "active",
			expected: ActiveSmartEmail,
		},
		{
			title:    "active uppercase",
			status:   "ACTIVE",
			expected: ActiveSmartEmail,
		},
		{
			title:    "draft lowercase",
			status:   "draft",
			expected: DraftSmartEmail,
		},
		{
			title:    "draft uppercase",
			status:   "DRAFT",
			expected: DraftSmartEmail,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			var status SmartEmailStatus
			err := status.UnmarshalJSON([]byte(tC.status))
			if err != nil {
				t.Errorf("Expected error: nil, Actual: %q", err)
			}
			if tC.expected != status {
				t.Errorf("Expected %s, Actual: %s", tC.expected, status)
			}
		})
	}
}
