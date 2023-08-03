package datasets

import (
	"encoding/json"
	"os"

	"github.com/imroc/req"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/s3"
)

type Datasets struct {
	Dataset []DatasetItem `json:"dataset"`
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

	// Field Issued does not do anything.. Its always the date of the request
	// This probably changed recently..
	// we're assuming the response is sorted so that newest is first.
	*d = Datasets{Dataset: datasets}

	return nil
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

func (d *Datasets) Reverse() {
	for i, j := 0, len(d.Dataset)-1; i < j; i, j = i+1, j-1 {
		d.Dataset[i], d.Dataset[j] = d.Dataset[j], d.Dataset[i]
	}
}

func FromURL(url string) (Datasets, error) {
	responseBytes, err := fetchDataset(url)
	if err != nil {
		return Datasets{}, err
	}

	return unmarshalDataset(responseBytes)
}

func FromPath(path string) (Datasets, error) {
	datasetBytes, err := loadDataset(path)
	if err != nil {
		return Datasets{}, err
	}

	return unmarshalDataset(datasetBytes)
}

func FromS3(s3 s3.S3) (Datasets, error) {
	s3Bytes, err := s3.FetchNewestFile()
	if err != nil {
		return Datasets{}, err
	}

	return unmarshalDataset(s3Bytes)
}

func (d *Datasets) Save(path string) error {
	datasetBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return os.WriteFile(path, datasetBytes, 0644)
}

func (d *Datasets) SaveToS3(path string, s3 s3.S3) error {
	datasetBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return s3.PutDataset(path, datasetBytes)
}

func fetchDataset(url string) ([]byte, error) {
	r, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

func loadDataset(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func unmarshalDataset(datasetsBytes []byte) (Datasets, error) {
	d := Datasets{}
	err := json.Unmarshal(datasetsBytes, &d)
	if err != nil {
		return Datasets{}, err
	}
	return d, nil
}
