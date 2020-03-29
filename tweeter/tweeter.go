package tweeter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Tweeter struct {
	client *twitter.Client
}

func NewTweeter(consumerKey, consumerSecret, accessToken, accessSecret string) (Tweeter, error) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

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
