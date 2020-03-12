package main

import (
	"encoding/json"
	"sort"
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

func (d *Datasets) UnmarshalJSON(b []byte) error {
	jsonMsg := map[string]json.RawMessage{}

	err := json.Unmarshal(b, &jsonMsg)
	if err != nil {
		return err
	}

	if jsonMsg["dataset"] == nil {
		*d = Datasets{}
		return nil
	}

	datasets := []Dataset{}

	err = json.Unmarshal(jsonMsg["dataset"], &datasets)
	if err != nil {
		return err
	}

	// sort datasets by Issued field (Date)
	// The response looks like it is already sorted, but to be sure
	// sort it again
	sort.Slice(datasets, func(i, j int) bool {
		return time.Time(datasets[i].Issued).After(time.Time(datasets[j].Issued))
	})

	*d = Datasets{Dataset: datasets}

	return nil
}

func (d *Datasets) Size() int {
	return len(d.Dataset)
}

// Compare compares the given Datasets with the current Datasets
// func (d *Datasets) Compare(otherDatasets *Datasets) []*Dataset {
// 	// primitively compare the first element
// 	for _, dataset := range d.Dataset {

// 	}
// }

// Compare compares a given Dataset with the current Dataset
// Only looks at Identifier
func (d *Dataset) Compare(otherDataset *Dataset) bool {
	if d.Identifier != otherDataset.Identifier {
		return false
	}
	if d.Title != otherDataset.Title {
		return false
	}
	if time.Time(d.Issued).Equal(time.Time(otherDataset.Issued)) == false {
		return false
	}

	return true
}
