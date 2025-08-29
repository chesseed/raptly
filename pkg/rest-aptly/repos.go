package aptly

type LocalRepo struct {
	Comment             string `json:"comment,omitempty"`
	DefaultComponent    string `json:"defaultComponent,omitempty"`
	DefaultDistribution string `json:"defaultDistribution,omitempty"`
	Name                string `json:"name,omitempty"`
}

// ReposList get the list of local repositories
func (c *Client) ReposList() ([]LocalRepo, error) {
	var repos []LocalRepo

	req := c.get("api/repos").
		SetResult(&repos)

	return repos, c.send(req)
}

type RepoCreateOptions struct {
	Comment             string `json:",omitempty"`
	DefaultComponent    string `json:",omitempty"`
	DefaultDistribution string `json:",omitempty"`
	FromSnapshot        string `json:",omitempty"`
}

// ReposCreate create new local repository
func (c *Client) ReposCreate(name string, opts RepoCreateOptions) (LocalRepo, error) {
	var repo LocalRepo

	type CreatePayload struct {
		Name string
		RepoCreateOptions
	}

	req := c.post("api/repos").
		SetResult(&repo).
		SetBody(&CreatePayload{Name: name, RepoCreateOptions: opts})

	return repo, c.send(req)
}

type RepoUpdateOptions struct {
	Comment             string `json:",omitempty"`
	DefaultComponent    string `json:",omitempty"`
	DefaultDistribution string `json:",omitempty"`
	// new repository name
	Name string `json:",omitempty"`
}

// ReposEdit edit/update existing local repository
func (c *Client) ReposEdit(name string, opts RepoUpdateOptions) (LocalRepo, error) {
	var repo LocalRepo

	req := c.put("api/repos/{name}").
		SetPathParam("name", name).
		SetResult(&repo).
		SetBody(opts)

	return repo, c.send(req)
}

// ReposShow get repository information
func (c *Client) ReposShow(name string) (LocalRepo, error) {
	var repo LocalRepo

	req := c.get("api/repos/{name}").
		SetPathParam("name", name).
		SetResult(&repo)

	return repo, c.send(req)
}

// ReposListPackages get list of packages
func (c *Client) ReposListPackages(name string, opts ListPackagesOptions) ([]Package, error) {

	params, err := opts.MakeParams()
	if err != nil {
		return nil, err
	}

	req := c.get("api/repos/{name}/packages").
		SetPathParam("name", name).
		SetQueryParams(params)

	return sendPackagesRequest(req, opts.Detailed)
}

// ReposDrop delete the local repository
func (c *Client) ReposDrop(name string, force bool) error {

	params := make(map[string]string)
	if force {
		params["force"] = "1"
	}

	req := c.delete("api/repos/{name}").
		SetPathParam("name", name).
		SetQueryParams(params)

	return c.send(req)
}

type RepoAddResult struct {
	FailedFiles []string `json:"FailedFiles"`
	// Report      map[string]string
}

type RepoAddOptions struct {
	// when adding package that conflicts with existing package, remove existing package
	ForceReplace bool
	// remove files that have been imported successfully into repository
	NoRemove bool
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

	req := c.post("api/repos/{name}/file/{dir}/{file}").
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetPathParam("file", filename).
		SetQueryParams(params).
		SetResult(&result)

	return result, c.send(req)
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

	req := c.post("api/repos/{name}/file/{dir}").
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetQueryParams(params).
		SetResult(&result)

	return result, c.send(req)
}

type RepoIncludeOptions struct {
	// when adding package that conflicts with existing package, remove existing package
	ForceReplace bool
	// remove files that have been imported successfully into repository
	NoRemove bool
	// accept unsigned .changes files
	AcceptUnsigned bool
	// disable verification of .changes file signature
	IgnoreSignature bool
}

// ReposIncludeFile include previously uploaded changes to repository
//
// Note: does not check files, it's the caller's responsibility to ensure the file is a valid changes file
func (c *Client) ReposIncludeFile(repo string, directory string, filename string, opts RepoIncludeOptions) (RepoAddResult, error) {

	params := make(map[string]string)
	if opts.NoRemove {
		params["noRemoveFiles"] = "1"
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

	req := c.post("api/repos/{name}/include/{dir}/{file}").
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetPathParam("file", filename).
		SetQueryParams(params).
		SetResult(&result)

	return result, c.send(req)
}

// ReposIncludeDirectory include previously uploaded directory to repository
func (c *Client) ReposIncludeDirectory(repo string, directory string, opts RepoIncludeOptions) (RepoAddResult, error) {

	params := make(map[string]string)
	if opts.NoRemove {
		params["noRemoveFiles"] = "1"
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

	req := c.post("api/repos/{name}/include/{dir}").
		SetPathParam("name", repo).
		SetPathParam("dir", directory).
		SetQueryParams(params).
		SetResult(&result)

	return result, c.send(req)
}

type pkgRefList struct {
	PackageRefs []string
}

func (c *Client) ReposRemovePackages(repo string, keys []string) (LocalRepo, error) {
	refs := pkgRefList{PackageRefs: keys}

	var result LocalRepo

	req := c.delete("api/repos/{name}/packages").
		SetPathParam("name", repo).
		SetBody(&refs).
		SetResult(&result)

	return result, c.send(req)
}
