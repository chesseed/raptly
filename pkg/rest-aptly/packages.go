package aptly

import "encoding/json"

type Package struct {
	Key       string
	ShortKey  string
	FilesHash string
	Extras    map[string]string
}

func (p *Package) UnmarshalJSON(bs []byte) (err error) {

	m := make(map[string]string)

	if err = json.Unmarshal(bs, &m); err == nil {
		p.Key = m["Key"]
		p.ShortKey = m["ShortKey"]
		p.FilesHash = m["FilesHash"]
		delete(m, "Key")
		delete(m, "ShortKey")
		delete(m, "FilesHash")
		p.Extras = m
	}

	return err
}

// Get list of packages keys
//
// since Aptly 1.6.0
func (c *Client) PackagesSearch(query string) ([]string, error) {

	params := make(map[string]string)
	if query != "" {
		params["q"] = query
	}

	var packages []string

	resp, err := c.client.R().
		SetQueryParams(params).
		SetResult(&packages).
		Get("api/packages")

	if err != nil {
		return packages, err
	} else if resp.IsSuccess() {
		return packages, nil
	}
	return packages, getError(resp)
}

// Get list of packages with detailed information
//
// since Aptly 1.6.0
func (c *Client) PackagesSearchDetailed(query string) ([]Package, error) {

	params := make(map[string]string)
	if query != "" {
		params["q"] = query
	}
	params["format"] = "details"

	var packages []Package

	resp, err := c.client.R().
		SetQueryParams(params).
		SetResult(&packages).
		Get("api/packages")

	if err != nil {
		return packages, err
	} else if resp.IsSuccess() {
		return packages, nil
	}
	return packages, getError(resp)
}
