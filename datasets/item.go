package datasets

import (
	"fmt"

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

func (d *DatasetItem) ResolveURL() (string, error) {
	r, err := req.Get(fmt.Sprintf("https://opendata.stadt-muenster.de/api/3/action/package_show?id=%s", d.Identifier))
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
