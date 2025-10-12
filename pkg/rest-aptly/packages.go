package aptly

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// ListPackagesOptions is used in SnapshotPackages(Detailed) and RepoPackages(Detailed)
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
	if opts.MaximumVersion {
		params["maximumVersion"] = "1"
	}
	if opts.Detailed {
		params["format"] = "details"
	}
	return params, nil
}

type Package struct {
	// always available

	// unique package identifier
	Key          string
	FilesHash    string
	Version      string
	Architecture string

	// only available when details format used

	ShortKey string
	// package name
	Package string
	// List of virtual packages this package provides
	Provides *[]string

	Source *string
	//Extras map[string]string
}

var packageRegex = regexp.MustCompile(`^(\S*)P(\S+)\s(\S+)\s(\S+)\s(\S+)$`)

// PackageFromKey convert aptly key to Package
func PackageFromKey(key string) (Package, error) {
	matched := packageRegex.FindStringSubmatch(key)

	if matched == nil {
		return Package{}, fmt.Errorf("could not match '%s'", key)
	}
	// ignore prefix matched[1] for now
	return Package{Key: key, Architecture: matched[2], Package: matched[3], Version: matched[4], FilesHash: matched[5]}, nil
}

func sendPackagesRequest(c *Client, req *request, detailed bool) ([]Package, error) {
	resp, err := callAPIWithResponse(c, req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var packages []Package
	if detailed {
		err = json.NewDecoder(resp.Body).Decode(&packages)
	} else {
		var keys []string
		err = json.NewDecoder(resp.Body).Decode(&keys)
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

// PackagesSearch returns list of packages
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

	req := c.get("api/packages").
		SetQueryParams(params)

	return sendPackagesRequest(c, req, detailed)
}

// PackagesInfo returns the package by key
func (c *Client) PackagesInfo(key string) (Package, error) {
	req := c.get("api/packages/{key}").
		SetPathParam("key", key)
	return callAPIwithResult[Package](c, req)
}
