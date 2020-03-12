package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/imroc/req"
)

const url = "https://opendata.stadt-muenster.de/data.json"

// How this works (at least in my head)
//
// 1. Download data.json
// 2. Compare with previous data.json

func main() {
	previousBytes, err := loadPreviousDataset()
	if err != nil {
		fmt.Println(err)
		return
	}
	prevDatasets, err := unmarshalDataset(previousBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	responseBytes, err := fetchDataset()
	if err != nil {
		fmt.Println(err)
		return
	}
	currDatasets, err := unmarshalDataset(responseBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(prevDatasets.Size())
	fmt.Println(currDatasets.Size())
}

func fetchDataset() ([]byte, error) {
	fmt.Printf("fetching dataset %s\n", url)
	r, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

func loadPreviousDataset() ([]byte, error) {
	fmt.Printf("loading dataset from file %s\n", url)
	return ioutil.ReadFile("/home/gerald/code/dkan-newest-dataset-notifier/data.json")
}

func unmarshalDataset(datasetsBytes []byte) (Datasets, error) {
	d := Datasets{}
	err := json.Unmarshal(datasetsBytes, &d)
	if err != nil {
		return Datasets{}, err
	}
	return d, nil
}
