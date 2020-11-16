package campaigns

import (
	"encoding/json"
	"strings"
	"time"
)

type CmTime time.Time

func (j *CmTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = CmTime(t)
	return nil
}

func (j CmTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (t CmTime) Equal(u CmTime) bool {
	return time.Time(t).Equal(time.Time(u))
}
