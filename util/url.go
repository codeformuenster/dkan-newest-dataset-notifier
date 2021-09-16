package util

import (
	"net/url"
)

func MakeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
