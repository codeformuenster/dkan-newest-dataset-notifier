package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/datasets"
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
// - load local (previous data.json)
// - Download data.json
// - Compare with previous data.json

func main() {
	flag.StringVar(&dataJSONURL, "url", defaultDataJSONURL, "url of the remote json file containing dkan datasets")
	flag.StringVar(&localPath, "local_path", defaultLocalPath, "path to local json file for comparison")

	flag.Parse()
	prevDatasets, err := datasets.FromPath(localPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	currDatasets, err := datasets.FromURL(dataJSONURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	missing := currDatasets.Compare(&prevDatasets)
	for _, m := range missing {
		url, err := m.ResolveURL()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s: %s %s \n", m.Issued, m.Title, url)
	}

	err = currDatasets.Save(fmt.Sprintf("./data-%s.json", time.Now().Format("2006-01-02")))
	if err != nil {
		fmt.Println(err)
		return
	}
}
