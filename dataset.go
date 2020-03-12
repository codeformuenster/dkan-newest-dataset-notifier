package main

import (
	"strings"
	"time"
)

type Datasets struct {
	Dataset []Dataset `json:"dataset"`
}

type Dataset struct {
	Modified    ISODate `json:"modified"`
	Issued      ISODate `json:"issued"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Identifier  string  `json:"identifier"`
}

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
