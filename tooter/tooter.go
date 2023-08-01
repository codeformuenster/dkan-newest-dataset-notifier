package tooter

import (
	"context"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/externalservices"

	"github.com/mattn/go-mastodon"
)

type Tooter struct {
	client  *mastodon.Client
	context context.Context
}

func NewTooter(mastodonConfig externalservices.MastodonConfig) (Tooter, error) {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       mastodonConfig.Server,
		ClientID:     mastodonConfig.ClientID,
		ClientSecret: mastodonConfig.ClientSecret,
		AccessToken:  mastodonConfig.AccessToken,
	})

	ctx := context.Background()
	err := c.Authenticate(ctx, mastodonConfig.Email, mastodonConfig.Password)
	if err != nil {
		return Tooter{}, err
	}

	return Tooter{c, ctx}, nil
}

func (t *Tooter) SendToot(text string) error {
	toot := mastodon.Toot{Status: text}

	_, err := t.client.PostStatus(t.context, &toot)
	if err != nil {
		return err
	}
	return nil
}
