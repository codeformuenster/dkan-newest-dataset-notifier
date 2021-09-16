package datasets

import (
	"fmt"

	"github.com/codeformuenster/dkan-newest-dataset-notifier/util"
	"github.com/imroc/req"
)

type DatasetItem struct {
	Modified    ISODate `json:"modified"`
	Issued      ISODate `json:"issued"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Identifier  string  `json:"identifier"`
}

type PackageResponse struct {
	Result []struct {
		URL string `json:"url"`
	} `json:"result"`
}

func (d *DatasetItem) ResolveURL(baseURL string) (string, error) {
	apiURL, err := util.MakeURL(fmt.Sprintf(
		"%s/api/3/action/package_show?id=%s",
		baseURL,
		d.Identifier,
	))
	if err != nil {
		return "", err
	}

	r, err := req.Get(apiURL)
	if err != nil {
		return "", err
	}

	var foo PackageResponse
	err = r.ToJSON(&foo)
	if err != nil {
		return "", err
	}

	return foo.Result[0].URL, nil
}

func (d *DatasetItem) ToTweetText(baseURL string) (string, error) {
	url, err := d.ResolveURL(baseURL)
	if err != nil {
		return "", err
	}
	var text string
	for _, template := range tweetTemplates {
		text = fmt.Sprintf(
			template,
			d.Title, url,
		)
		if len(text) < 280 {
			return text, nil
		}
	}
	if len(text) < 280 {
		return "", fmt.Errorf("Tweet too long (> 280): %s", text)
	}
	return "", fmt.Errorf("I thought this error will never happen")
}
