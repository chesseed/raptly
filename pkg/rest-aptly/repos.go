package aptly

import (
	"fmt"
)

type LocalRepo struct {
	Comment             string `json:"comment,omitempty"`
	DefaultComponent    string `json:"defaultComponent,omitempty"`
	DefaultDistribution string `json:"defaultDistribution,omitempty"`
	Name                string `json:"name,omitempty"`
}

// get list of repositories
func (c *Client) ReposList() ([]LocalRepo, error) {
	var repos []LocalRepo

	resp, err := c.client.R().
		SetResult(&repos).
		Get("api/repos")

	if err != nil {
		return repos, err
	} else if resp.IsSuccess() {
		return repos, nil
	}
	return repos, getError(resp)
}

type RepoCreateOptions struct {
	Comment             *string `json:",omitempty"`
	DefaultComponent    *string `json:",omitempty"`
	DefaultDistribution *string `json:",omitempty"`
	FromSnapshot        *string `json:",omitempty"`
}

// create new repository
func (c *Client) ReposCreate(name string, opts RepoCreateOptions) (LocalRepo, error) {
	var repo LocalRepo

	type CreatePayload struct {
		Name string
		RepoCreateOptions
	}

	resp, err := c.client.R().
		SetResult(&repo).
		SetBody(&CreatePayload{Name: name, RepoCreateOptions: opts}).
		Post("api/repos")

	if err != nil {
		return repo, err
	} else if resp.IsSuccess() {
		return repo, nil
	}
	return repo, getError(resp)
}

type RepoUpdateOptions struct {
	Comment             *string `json:",omitempty"`
	DefaultComponent    *string `json:",omitempty"`
	DefaultDistribution *string `json:",omitempty"`
	// new repository name
	Name *string `json:",omitempty"`
}

// edit/update existing repository
func (c *Client) ReposEdit(name string, opts RepoUpdateOptions) (LocalRepo, error) {
	var repo LocalRepo

	resp, err := c.client.R().
		SetResult(&repo).
		SetBody(opts).
		SetPathParam("name", name).
		Put("api/repos/{name}")

	if err != nil {
		return repo, err
	} else if resp.IsSuccess() {
		return repo, nil
	}
	return repo, getError(resp)
}

// get repository information
func (c *Client) ReposShow(name string) (LocalRepo, error) {
	var repo LocalRepo

	resp, err := c.client.R().
		SetResult(&repo).
		SetPathParam("name", name).
		Get("api/repos/{name}")

	if err != nil {
		return repo, err
	} else if resp.IsSuccess() {
		return repo, nil
	}
	return repo, getError(resp)
}

type RepoListPackageOption struct {
	Query         string
	WithDeps      bool
	NewestVersion bool // maximumVersion in API
}

func getListPackagesParams(conf *RepoListPackageOption) (map[string]string, error) {
	params := make(map[string]string)

	if conf != nil {
		if conf.Query == "" && conf.WithDeps {
			return nil, fmt.Errorf("withDeps requires a query")
		}
		if conf.Query != "" {
			params["q"] = conf.Query
		}
		if conf.WithDeps {
			params["withDeps"] = "1"
		}
		if conf.WithDeps {
			params["maximumVersion "] = "1"
		}
	}
	return params, nil
}

// Get list of packages keys
func (c *Client) ReposListPackages(name string, conf *RepoListPackageOption) ([]string, error) {

	params, err := getListPackagesParams(conf)
	if err != nil {
		return nil, err
	}

	var packages []string

	resp, err := c.client.R().
		SetPathParam("name", name).
		SetQueryParams(params).
		SetResult(&packages).
		Get("api/repos/{name}/packages")

	if err != nil {
		return packages, err
	} else if resp.IsSuccess() {
		return packages, nil
	}
	return packages, getError(resp)
}

// Get the list of packages with full information
func (c *Client) ReposListPackagesDetailed(name string, conf *RepoListPackageOption) ([]Package, error) {

	params, err := getListPackagesParams(conf)
	if err != nil {
		return nil, err
	}
	// return object instead of strings
	params["format"] = "details"

	var packages []Package

	resp, err := c.client.R().
		SetPathParam("name", name).
		SetQueryParams(params).
		SetResult(&packages).
		Get("api/repos/{name}/packages")

	if err != nil {
		return packages, err
	} else if resp.IsSuccess() {
		return packages, nil
	}
	return packages, getError(resp)
}

// Remove the repository
func (c *Client) ReposDrop(name string, force bool) error {

	params := make(map[string]string)
	if force {
		params["force"] = "1"
	}

	resp, err := c.client.R().
		SetPathParam("name", name).
		SetQueryParams(params).
		Delete("api/repos/{name}")

	if err != nil {
		return err
	} else if resp.IsSuccess() {
		return nil
	}
	return getError(resp)
}

type RepoAddResult struct {
	FailedFiles []string `json:"FailedFiles"`
	// Report      map[string]string
}

type RepoAddOptions struct {
	ForceReplace bool
	NoRemove     bool
}

func (c *Client) ReposAddFile(repo string, directory string, filename string, opts RepoAddOptions) (RepoAddResult, error) {

	params := make(map[string]string)
	if opts.NoRemove {
		params["noRemove"] = "1"
	}
	if opts.ForceReplace {
		params["forceReplace"] = "1"
	}

	var result RepoAddResult

	resp, err := c.client.R().
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetPathParam("file", filename).
		SetQueryParams(params).
		SetResult(&result).
		Post("api/repos/{name}/file/{dir}/{file}")

	if err != nil {
		return result, err
	} else if resp.IsSuccess() {
		return result, nil
	}
	return result, getError(resp)
}

func (c *Client) ReposAddDirectory(repo string, directory string, opts RepoAddOptions) (RepoAddResult, error) {

	params := make(map[string]string)
	if opts.NoRemove {
		params["noRemove"] = "1"
	}
	if opts.ForceReplace {
		params["forceReplace"] = "1"
	}

	var result RepoAddResult

	resp, err := c.client.R().
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetQueryParams(params).
		SetResult(&result).
		Post("api/repos/{name}/file/{dir}")

	if err != nil {
		return result, err
	} else if resp.IsSuccess() {
		return result, nil
	}
	return result, getError(resp)
}

type RepoIncludeOptions struct {
	ForceReplace    bool
	NoRemove        bool
	AcceptUnsigned  bool
	IgnoreSignature bool
}

func (c *Client) ReposIncludeFile(repo string, directory string, filename string, opts RepoIncludeOptions) (RepoAddResult, error) {

	params := make(map[string]string)
	if opts.NoRemove {
		params["noRemove"] = "1"
	}
	if opts.ForceReplace {
		params["forceReplace"] = "1"
	}
	if opts.AcceptUnsigned {
		params["acceptUnsigned"] = "1"
	}
	if opts.IgnoreSignature {
		params["ignoreSignature"] = "1"
	}

	var result RepoAddResult

	resp, err := c.client.R().
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetPathParam("file", filename).
		SetQueryParams(params).
		SetResult(&result).
		Post("api/repos/{name}/include/{dir}/{file}")

	if err != nil {
		return result, err
	} else if resp.IsSuccess() {
		return result, nil
	}
	return result, getError(resp)
}

func (c *Client) ReposIncludeDirectory(repo string, directory string, opts RepoIncludeOptions) (RepoAddResult, error) {

	params := make(map[string]string)
	if opts.NoRemove {
		params["noRemove"] = "1"
	}
	if opts.ForceReplace {
		params["forceReplace"] = "1"
	}
	if opts.AcceptUnsigned {
		params["acceptUnsigned"] = "1"
	}
	if opts.IgnoreSignature {
		params["ignoreSignature"] = "1"
	}

	var result RepoAddResult

	resp, err := c.client.R().
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetQueryParams(params).
		SetResult(&result).
		Post("api/repos/{name}/include/{dir}")

	if err != nil {
		return result, err
	} else if resp.IsSuccess() {
		return result, nil
	}
	return result, getError(resp)
}
