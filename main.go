package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/datasets"
	"github.com/codeformuenster/dkan-newest-dataset-notifier/tweeter"
)

const (
	defaultDataJSONURL = "https://opendata.stadt-muenster.de/data.json"
	defaultLocalPath   = "./data.json"
)

var (
	dataJSONURL                                            string
	localPath                                              string
	consumerKey, consumerSecret, accessToken, accessSecret string
)

// How this works (at least in my head)
//
// - load local (previous data.json)
// - Download data.json
// - Compare with previous data.json

func main() {
	flag.StringVar(&consumerKey, "consumer-key", "", "twitter oauth1 consumerKey")
	flag.StringVar(&consumerSecret, "consumer-secret", "", "twitter oauth1 consumerSecret")
	flag.StringVar(&accessToken, "access-token", "", "twitter oauth1 accessToken")
	flag.StringVar(&accessSecret, "access-secret", "", "twitter oauth1 accessSecret")

	flag.StringVar(&dataJSONURL, "url", defaultDataJSONURL, "url of the remote json file containing dkan datasets")
	flag.StringVar(&localPath, "local-path", defaultLocalPath, "path to local json file for comparison")

	flag.Parse()

	t, tweeterAvailable := setupTweeter()

	prevDatasets, err := datasets.FromPath(localPath)
	if err != nil {
		fmt.Println(err)
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
		text := fmt.Sprintf("Datensatz \"%s\" nun verfügbar im Open-Data-Portal Münster %s #OpenData #Münster", m.Title, url)
		fmt.Printf("%d %s\n", len(text), text)

		if tweeterAvailable == true {
			err = t.SendTweet(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	err = currDatasets.Save(fmt.Sprintf("./data-%s.json", time.Now().Format("2006-01-02")))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func setupTweeter() (tweeter.Tweeter, bool) {
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return tweeter.Tweeter{}, false
	}

	t, err := tweeter.NewTweeter(consumerKey, consumerSecret, accessToken, accessSecret)
	if err != nil {
		fmt.Println(err)
		return t, false
	}

	return t, true
}
