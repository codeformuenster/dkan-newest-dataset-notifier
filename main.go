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
	"github.com/codeformuenster/dkan-newest-dataset-notifier/tooter"
	"github.com/codeformuenster/dkan-newest-dataset-notifier/util"
)

const defaultDKANInstance = "https://opendata.stadt-muenster.de"

var (
	dkanInstanceURL, localPath, externalServicesConfigPath string
	enableTweeter, allowEmpty                              bool
)

// How this works (at least in my head)
//
// - load local (previous data.json)
// - Download data.json
// - Compare with previous data.json

func main() {
	flag.BoolVar(&enableTweeter, "enable-tooter", false, "enable the creation of toots")
	flag.BoolVar(&allowEmpty, "allow-empty", false, "allow empty previous dataset, for initialization")

	flag.StringVar(&dkanInstanceURL, "url", defaultDKANInstance, "base url of the dkan instance (https://...)")
	flag.StringVar(&localPath, "local-path", "", "path to local json file for comparison")

	flag.StringVar(&externalServicesConfigPath, "config-path", "", "path to local json external services configuration")

	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var err error

	// validate + constructr dkan urls
	datasetsURL, err := util.MakeURL(fmt.Sprintf("%s/%s", dkanInstanceURL, "data.json"))
	if err != nil {
		log.Printf("Could not create valid datasets URL")
		log.Panicln(err)
	}

	cfg, err := externalservices.FromFile(externalServicesConfigPath)
	if err != nil {
		log.Println(err)
		log.Println("Disabling external services")
	}

	s3Instance, s3Available := setupS3(cfg.S3Config)

	t, tooterAvailable, err := setupMastodon(cfg.MastodonConfig)
	if err != nil {
		log.Panicln(err)
	}
	if !tooterAvailable {
		log.Println("disabling tooter, no toots will be created")
	}

	var prevDatasets datasets.Datasets

	if s3Available {
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
		log.Println("Reading previous datasets failed, assuming empty")
	}

	if !allowEmpty && len(prevDatasets.Dataset) == 0 {
		log.Println("Empty previous datasets not allowed")
		return
	}

	currDatasets, err := datasets.FromURL(datasetsURL)
	if err != nil {
		log.Panicln(err)
	}

	missing := currDatasets.Compare(&prevDatasets)
	for _, m := range missing {
		tootText, err := m.ToTootText(dkanInstanceURL)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("%d %s\n", len(tootText), tootText)

		if tooterAvailable {
			err = t.SendToot(tootText)
			if err != nil {
				log.Println(err)
				continue
			}
			time.Sleep(30 * time.Second)
		}
	}

	{
		var err error

		if s3Available {
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
}

func setupMastodon(cfg externalservices.MastodonConfig) (tooter.Tooter, bool, error) {
	if !enableTweeter || !cfg.Validate() {
		return tooter.Tooter{}, false, nil
	}

	t, err := tooter.NewTooter(cfg)
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
