package main

import (
	"flag"
	"fmt"
	"log"
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

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	t, tweeterAvailable, err := setupTweeter()
	if err != nil {
		log.Println(err)
	}
	if tweeterAvailable == false {
		log.Println("disabling tweeter, no tweets will be created")
	}

	if localPath == "" {
		localPath = makeDataPath(time.Now().Add(-24 * time.Hour))
		log.Printf("empty local-path flag, assuming path %s\n", localPath)
	}

	prevDatasets, err := datasets.FromPath(localPath)
	if err != nil {
		log.Println(err)
	}

	currDatasets, err := datasets.FromURL(dataJSONURL)
	if err != nil {
		log.Panicln(err)
	}

	missing := currDatasets.Compare(&prevDatasets)
	for _, m := range missing {
		tweetText, err := m.ToTweetText()
		if err != nil {
			log.Panicln(err)
		}

		log.Printf("%d %s\n", len(tweetText), tweetText)

		if tweeterAvailable == true {
			err = t.SendTweet(tweetText)
			if err != nil {
				log.Panicln(err)
			}
		}
	}

	err = currDatasets.Save(makeDataPath(time.Now()))
	if err != nil {
		log.Panicln(err)
	}
}

func setupTweeter() (tweeter.Tweeter, bool, error) {
	if enableTweeter == false ||
		consumerKey == "" ||
		consumerSecret == "" ||
		accessToken == "" ||
		accessSecret == "" {
		return tweeter.Tweeter{}, false, nil
	}

	t, err := tweeter.NewTweeter(consumerKey, consumerSecret, accessToken, accessSecret)
	if err != nil {
		return t, false, err
	}

	return t, true, nil
}

func makeDataPath(date time.Time) string {
	filename := fmt.Sprintf("data-%s.json", date.Format("2006-01-02"))
	dir, err := os.Getwd()
	if err != nil {
		return "./" + filename
	}

	return path.Join(dir, filename)
}
