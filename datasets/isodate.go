package datasets

import (
	"fmt"
	"strings"
	"time"
)

type ISODate time.Time

const formatTemplate = "2006-01-02"

func (i *ISODate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	t, err := time.Parse(formatTemplate, s)
	if err != nil {
		return err
	}

	*i = ISODate(t)

	return nil
}

func (i *ISODate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", i)), nil
}

func (i ISODate) String() string {
	return time.Time(i).Format(formatTemplate)
}
