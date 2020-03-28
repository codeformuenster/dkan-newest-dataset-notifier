package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/imroc/req"
)

const (
	defaultDataJSONURL = "https://opendata.stadt-muenster.de/data.json"
	defaultLocalPath   = "./data.json"
)

var (
	dataJSONURL string
	localPath   string
)

// How this works (at least in my head)
//
// 1. Download data.json
// 2. Compare with previous data.json

func main() {
	flag.StringVar(&dataJSONURL, "url", defaultDataJSONURL, "url of the remote json file containing dkan datasets")
	flag.StringVar(&localPath, "local_path", defaultLocalPath, "path to local json file for comparison")

	flag.Parse()

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
	missing := currDatasets.Compare(&prevDatasets)
	for _, m := range missing {
		url, err := m.ResolveURL()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s: %s %s \n", time.Time(m.Issued).Format("2006-01-02"), m.Title, url)
	}
}

func fetchDataset() ([]byte, error) {
	fmt.Printf("fetching dataset %s\n", dataJSONURL)
	r, err := req.Get(dataJSONURL)
	if err != nil {
		return nil, err
	}

	return r.ToBytes()
}

func loadPreviousDataset() ([]byte, error) {
	fmt.Printf("loading dataset from file %s\n", localPath)
	return ioutil.ReadFile(localPath)
}

func unmarshalDataset(datasetsBytes []byte) (Datasets, error) {
	d := Datasets{}
	err := json.Unmarshal(datasetsBytes, &d)
	if err != nil {
		return Datasets{}, err
	}
	return d, nil
}
