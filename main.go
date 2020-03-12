package main

import (
	"fmt"

	"github.com/imroc/req"
)

const url = "https://opendata.stadt-muenster.de/data.json"

// How this works (at least in my head)
//
// 1. Download data.json
// 2. Compare with previous

func main() {
	fmt.Printf("importing dataset %s\n", url)

	r, err := req.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

  // save a copy for later
  r.ToFile("data.json")

  r, err = req.Get("file:///home/gerald/code/dkan-newest-dataset-notifier/data.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	d := Datasets{}
	err = r.ToJSON(&d)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(d.Dataset[0].Title)
}

func downloadDataset() (string, error) {
	r, err := req.Get(url)
	if err != nil {
		return "", err
	}

  err = r.ToFile("data.json")
	if err != nil {
		return "", err
	}

  return "", nil
}
