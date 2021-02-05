package tweeter

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/externalservices"
)

type Tweeter struct {
	client *twitter.Client
}

func NewTweeter(twitterConfig externalservices.TwitterConfig) (Tweeter, error) {
	config := oauth1.NewConfig(twitterConfig.ConsumerKey, twitterConfig.ConsumerSecret)
	token := oauth1.NewToken(twitterConfig.AccessToken, twitterConfig.AccessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	httpClient.Timeout = 60 * time.Second

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:      twitter.Bool(true),
		IncludeEmail:    twitter.Bool(false),
		IncludeEntities: twitter.Bool(false),
	}
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return Tweeter{}, err
	}

	return Tweeter{
		client: client,
	}, nil
}

func (t *Tweeter) SendTweet(text string) error {
	_, _, err := t.client.Statuses.Update(text, nil)
	if err != nil {
		return err
	}
	return nil
}
