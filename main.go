package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/datasets"
	"github.com/codeformuenster/dkan-newest-dataset-notifier/tweeter"
)

const (
	defaultDataJSONURL = "https://opendata.stadt-muenster.de/data.json"
	defaultLocalPath   = ""
)

var (
	dataJSONURL, localPath                                 string
	consumerKey, consumerSecret, accessToken, accessSecret string
	enableTweeter                                          bool
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
	flag.BoolVar(&enableTweeter, "enable-tweeter", false, "enable the creation of tweets")

	flag.StringVar(&dataJSONURL, "url", defaultDataJSONURL, "url of the remote json file containing dkan datasets")
	flag.StringVar(&localPath, "local-path", defaultLocalPath, "path to local json file for comparison")

	flag.Parse()

	t, tweeterAvailable := setupTweeter()

	if tweeterAvailable == false {
		fmt.Println("disabling tweeter, no tweets will be created")
	}

	if localPath == "" {
		localPath = makeDataPath(time.Now().Add(-24 * time.Hour))
		fmt.Printf("empty local-path flag, assuming path %s\n", localPath)
	}

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
		text := fmt.Sprintf("\"%s\" ist nun als Datensatz im Open-Data-Portal Münster verfügbar %s #OpenData #Münster", m.Title, url)
		textLen := len(text)
		if textLen > 280 {
			fmt.Print("!!!! ")
		}
		fmt.Printf("%d %s\n", textLen, text)

		if tweeterAvailable == true {
			err = t.SendTweet(text)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	err = currDatasets.Save(makeDataPath(time.Now()))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func setupTweeter() (tweeter.Tweeter, bool) {
	if enableTweeter == false ||
		consumerKey == "" ||
		consumerSecret == "" ||
		accessToken == "" ||
		accessSecret == "" {
		return tweeter.Tweeter{}, false
	}

	t, err := tweeter.NewTweeter(consumerKey, consumerSecret, accessToken, accessSecret)
	if err != nil {
		fmt.Println(err)
		return t, false
	}

	return t, true
}

func makeDataPath(date time.Time) string {
	filename := fmt.Sprintf("data-%s.json", date.Format("2006-01-02"))
	dir, err := os.Getwd()
	if err != nil {
		return "./" + filename
	}

	return path.Join(dir, filename)
}
