package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/datasets"
	"github.com/codeformuenster/dkan-newest-dataset-notifier/externalservices"
	"github.com/codeformuenster/dkan-newest-dataset-notifier/s3"
	"github.com/codeformuenster/dkan-newest-dataset-notifier/tweeter"
)

const defaultDataJSONURL = "https://opendata.stadt-muenster.de/data.json"

var (
	dataJSONURL, localPath, externalServicesConfigPath string
	enableTweeter                                      bool
)

// How this works (at least in my head)
//
// - load local (previous data.json)
// - Download data.json
// - Compare with previous data.json

func main() {
	flag.BoolVar(&enableTweeter, "enable-twitter", false, "enable the creation of tweets")

	flag.StringVar(&dataJSONURL, "url", defaultDataJSONURL, "url of the remote json file containing dkan datasets")
	flag.StringVar(&localPath, "local-path", "", "path to local json file for comparison")

	flag.StringVar(&externalServicesConfigPath, "config-path", "", "path to local json external services configuration")

	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var err error

	cfg, err := externalservices.FromFile(externalServicesConfigPath)
	if err != nil {
		log.Println(err)
		log.Println("Disabling external services")
	}

	s3Instance, s3Available := setupS3(cfg.S3Config)

	t, tweeterAvailable, err := setupTweeter(cfg.TwitterConfig)
	if err != nil {
		log.Println(err)
	}
	if tweeterAvailable == false {
		log.Println("disabling tweeter, no tweets will be created")
	}

	var prevDatasets datasets.Datasets

	if s3Available == true {
		prevDatasets, err = datasets.FromS3(s3Instance)
	} else {
		if localPath == "" {
			localPath = makeDataPath(time.Now().Add(-24 * time.Hour))
			log.Printf("empty local-path flag, assuming path %s\n", localPath)
		}

		prevDatasets, err = datasets.FromPath(localPath)
	}
	// handle error of prev dataset fetch
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
			log.Println(err)
			continue
		}

		log.Printf("%d %s\n", len(tweetText), tweetText)

		if tweeterAvailable == true {
			err = t.SendTweet(tweetText)
			if err != nil {
				log.Println(err)
				continue
			}
			time.Sleep(30 * time.Second)
		}
	}

	if s3Available == true {
		if len(missing) != 0 {
			err = currDatasets.SaveToS3(fmt.Sprintf("data-%s.json", time.Now().Format("2006-01-02")), s3Instance)
		}
	} else {
		err = currDatasets.Save(makeDataPath(time.Now()))

	}
	if err != nil {
		log.Panicln(err)
	}
}

func setupTweeter(cfg externalservices.TwitterConfig) (tweeter.Tweeter, bool, error) {
	if enableTweeter == false ||
		cfg.Validate() == false {
		return tweeter.Tweeter{}, false, nil
	}

	t, err := tweeter.NewTweeter(cfg)
	if err != nil {
		return t, false, err
	}

	return t, true, nil
}

func setupS3(cfg externalservices.S3Config) (s3.S3, bool) {
	return s3.NewS3(cfg), cfg.Validate()
}

func makeDataPath(date time.Time) string {
	filename := fmt.Sprintf("data-%s.json", date.Format("2006-01-02"))
	dir, err := os.Getwd()
	if err != nil {
		return "./" + filename
	}

	return path.Join(dir, filename)
}
