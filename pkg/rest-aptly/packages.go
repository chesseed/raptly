package aptly

type Package struct {
	Key       string
	ShortKey  string
	FilesHash string
	//
	Version      string
	Architecture string
	// package name
	Package string
	// List of virtual packages this package provides
	Provides []string

	Source *string

	//Extras map[string]string
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

// Get package by key
func (c *Client) PackagesInfo(key string) (Package, error) {

	var pkg Package

	resp, err := c.client.R().
		SetPathParam("key", key).
		SetResult(&pkg).
		Get("api/packages/{key}")

	if err != nil {
		return pkg, err
	} else if resp.IsSuccess() {
		return pkg, nil
	}
	return pkg, getError(resp)
}
