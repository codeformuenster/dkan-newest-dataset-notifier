package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/imroc/req"
)

type Datasets struct {
	Dataset []DatasetItem `json:"dataset"`
}

type DatasetItem struct {
	Modified    ISODate `json:"modified"`
	Issued      ISODate `json:"issued"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Identifier  string  `json:"identifier"`
}

type PackageResponse struct {
	Result []struct {
		URL string `json:"url"`
	} `json:"result"`
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

	datasets := []DatasetItem{}

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
func (d *Datasets) Compare(otherDatasets *Datasets) []DatasetItem {
	missing := []DatasetItem{}

	for _, dataset := range d.Dataset {
		hasDataset := false
	Inner:
		for _, otherDataset := range otherDatasets.Dataset {
			if dataset.Identifier == otherDataset.Identifier {
				hasDataset = true
				break Inner
			}
		}
		if hasDataset == false {
			missing = append(missing, dataset)
		}
	}
	return missing
}

// Compare compares a given Dataset with the current Dataset
// Only looks at Identifier
func (d *DatasetItem) Compare(otherDataset *DatasetItem) bool {
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

func (d *DatasetItem) ResolveURL() (string, error) {
	// "https://opendata.stadt-muenster.de/api/3/action/package_show?id=6e1bb0a6-fc86-4bcb-90ee-15f62fbcc82c"
	r, err := req.Get(fmt.Sprintf("https://opendata.stadt-muenster.de/api/3/action/package_show?id=%s", d.Identifier))
	if err != nil {
		return "", err
	}

	var foo PackageResponse
	err = r.ToJSON(&foo)
	if err != nil {
		return "", err
	}

	return foo.Result[0].URL, nil
}
