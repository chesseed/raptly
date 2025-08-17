package aptly

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
)

// used in SnapshotPackages(Detailed) and RepoPackages(Detailed)
type ListPackagesOptions struct {
	// package query, see https://www.aptly.info/doc/feature/query/
	Query string
	// include dependencies when evaluating package query
	WithDeps bool
	// return the highest version for each package name
	MaximumVersion bool
	// include all information
	Detailed bool
}

func (opts *ListPackagesOptions) MakeParams() (map[string]string, error) {
	params := make(map[string]string)

	if opts.Query == "" && opts.WithDeps {
		return nil, fmt.Errorf("withDeps requires a query")
	}
	if opts.Query != "" {
		params["q"] = opts.Query
	}
	if opts.WithDeps {
		params["withDeps"] = "1"
	}
	if opts.WithDeps {
		params["maximumVersion"] = "1"
	}
	if opts.Detailed {
		params["format"] = "details"
	}
	return params, nil
}

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
	Provides *[]string

	Source *string

	//Extras map[string]string
}

var packageRegex = regexp.MustCompile(`^(\S*)P(\S+)\s(\S+)\s(\S+)\s(\S+)$`)

func PackageFromKey(key string) (Package, error) {
	matched := packageRegex.FindStringSubmatch(key)

	if matched == nil {
		return Package{}, fmt.Errorf("could not match '%s'", key)
	}
	// ignore prefix matched[1] for now
	return Package{Key: key, Architecture: matched[2], Package: matched[3], Version: matched[4], FilesHash: matched[5]}, nil
}

func responseToPackages(resp *resty.Response, detailed bool) ([]Package, error) {
	if resp.IsError() {
		return nil, getError(resp)
	}

	var packages []Package
	var err error
	if detailed {
		err = json.Unmarshal(resp.Body(), &packages)
	} else {
		var keys []string
		err = json.Unmarshal(resp.Body(), &keys)
		if err != nil {
			return nil, err
		}

		packages = make([]Package, 0, len(keys))
		for _, key := range keys {
			p, err := PackageFromKey(key)
			if err != nil {
				return nil, err
			}
			packages = append(packages, p)
		}
	}

	if err != nil {
		return packages, err
	}
	return packages, nil
}

// Get list of packages keys
//
// since Aptly 1.6.0
func (c *Client) PackagesSearch(query string, detailed bool) ([]Package, error) {

	params := make(map[string]string)
	if query != "" {
		params["q"] = query
	}
	if detailed {
		params["format"] = "details"
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get("api/packages")

	if err != nil {
		return nil, err
	}
	return responseToPackages(resp, detailed)
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
