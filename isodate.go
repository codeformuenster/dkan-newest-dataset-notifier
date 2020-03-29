package main

import (
	"fmt"
	"strings"
	"time"
)

type ISODate time.Time

func (i *ISODate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*i = ISODate(t)

	return nil
}

func (i *ISODate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*i).Format("2006-01-02"))), nil
}
