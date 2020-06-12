package createsend

import (
	"fmt"

	"github.com/xitonix/createsend/mock"
)

type bodyMock struct {
	Value               string `json:"value"`
	forceMarshalFailure bool
}

func newBodyMock(forceMarshalFailure bool) *bodyMock {
	return &bodyMock{
		Value:               "value",
		forceMarshalFailure: forceMarshalFailure,
	}
}

func (b *bodyMock) MarshalJSON() ([]byte, error) {
	if b.forceMarshalFailure {
		return nil, mock.ErrDeliberate
	}
	return []byte(fmt.Sprintf(`{"value": "%s"}`, b.Value)), nil
}
