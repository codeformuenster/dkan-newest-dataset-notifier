package tooter

import (
	"github.com/codeformuenster/dkan-newest-dataset-notifier/externalservices"
)

type Tooter struct {
}

func NewTooter(twitterConfig externalservices.MastodonConfig) (Tooter, error) {
	return Tooter{}, nil
}

func (t *Tooter) SendToot(text string) error {
	return nil
}
