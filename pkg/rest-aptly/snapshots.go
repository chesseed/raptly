package aptly

// Snapshot is immutable state of repository: list of packages
type Snapshot struct {
	// Human-readable name
	Name string `json:"Name"`
	// Date of creation
	CreatedAt  string `json:"CreatedAt"`
	SourceKind string `json:"SourceKind"`
	// Sources
	Snapshots []*Snapshot `json:",omitempty"`
	//RemoteRepos []*RemoteRepo `json:",omitempty"`
	LocalRepos []*LocalRepo `json:",omitempty"`
	Packages   []string     `json:",omitempty"`

	// Description of how snapshot was created
	Description string

	Origin               string
	NotAutomatic         string
	ButAutomaticUpgrades string
}

func (c *Client) SnapshotList() ([]Snapshot, error) {
	var snaps []Snapshot

	resp, err := c.client.R().
		SetResult(&snaps).
		Get("api/snapshots")

	if err != nil {
		return snaps, err
	} else if resp.IsSuccess() {
		return snaps, nil
	}
	return snaps, getError(resp)
}

func (c *Client) SnapshotShow(name string) (Snapshot, error) {
	var snap Snapshot

	resp, err := c.client.R().
		SetResult(&snap).
		SetPathParam("name", name).
		Get("api/snapshots/{name}")

	if err != nil {
		return snap, err
	} else if resp.IsSuccess() {
		return snap, nil
	}
	return snap, getError(resp)
}

func (c *Client) SnapshotPackages(name string, opts ListPackagesOptions) ([]string, error) {
	params, err := opts.MakeParams()
	if err != nil {
		return nil, err
	}

	var pkgs []string

	resp, err := c.client.R().
		SetResult(&pkgs).
		SetPathParam("name", name).
		SetQueryParams(params).
		Get("api/snapshots/{name}/packages")

	if err != nil {
		return pkgs, err
	} else if resp.IsSuccess() {
		return pkgs, nil
	}
	return pkgs, getError(resp)
}

func (c *Client) SnapshotPackagesDetailed(name string, opts ListPackagesOptions) ([]Package, error) {
	params, err := opts.MakeParams()
	if err != nil {
		return nil, err
	}
	// return object instead of strings
	params["format"] = "details"

	var pkgs []Package

	resp, err := c.client.R().
		SetResult(&pkgs).
		SetPathParam("name", name).
		SetQueryParams(params).
		Get("api/snapshots/{name}/packages")

	if err != nil {
		return pkgs, err
	} else if resp.IsSuccess() {
		return pkgs, nil
	}
	return pkgs, getError(resp)
}

func (c *Client) SnapshotDrop(name string, force bool) error {
	params := make(map[string]string)
	if force {
		params["force"] = "1"
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		SetPathParam("name", name).
		Delete("api/snapshots/{name}")

	if err != nil {
		return err
	} else if resp.IsSuccess() {
		return nil
	}
	return getError(resp)
}

func (c *Client) SnapshotFromRepo(name string, repoName string, description string) (Snapshot, error) {
	var snap Snapshot

	type CreateParam struct {
		Name        string `json:"Name"`
		Description string `json:"Description,omitempty"`
	}

	params := CreateParam{
		Name: name, Description: description,
	}

	resp, err := c.client.R().
		SetResult(&snap).
		SetBody(params).
		SetPathParam("name", repoName).
		Post("api/repos/{name}/snapshots")

	if err != nil {
		return snap, err
	} else if resp.IsSuccess() {
		return snap, nil
	}
	return snap, getError(resp)
}

func (c *Client) SnapshotFromMirror(name string, mirror string, description string) (Snapshot, error) {
	var snap Snapshot

	type CreateParam struct {
		Name        string `json:"Name"`
		Description string `json:"Description,omitempty"`
	}

	params := CreateParam{
		Name: name, Description: description,
	}

	resp, err := c.client.R().
		SetResult(&snap).
		SetBody(params).
		SetPathParam("name", mirror).
		Post("api/mirrors/{name}/snapshots")

	if err != nil {
		return snap, err
	} else if resp.IsSuccess() {
		return snap, nil
	}
	return snap, getError(resp)
}

type DiffPkg struct {
	Architecture string `json:"architecture"`
	Name         string `json:"name"`
	Version      string `json:"version"`
}
type PackageDiff struct {
	Left  DiffPkg `json:"left"`
	Right DiffPkg `json:"right"`
}

func (c *Client) SnapshotDiff(left string, right string, onlyMatching bool) ([]PackageDiff, error) {
	var diff []PackageDiff

	params := make(map[string]string)

	if onlyMatching {
		params["onlyMatching"] = "1"
	}

	resp, err := c.client.R().
		SetResult(&diff).
		SetQueryParams(params).
		SetPathParam("left", left).
		SetPathParam("right", right).
		Get("api/snapshots/{left}/diff/{right}")

	if err != nil {
		return diff, err
	} else if resp.IsSuccess() {
		return diff, nil
	}
	return diff, getError(resp)
}

type SnapshotUpdateOptions struct {
	// new name for the snapshot
	Name *string `json:"Name,omitempty"`
	// new description for the snapshot
	Description *string `json:"Description,omitempty"`
}

func (c *Client) SnapshotUpdate(name string, opts SnapshotUpdateOptions) (Snapshot, error) {
	var snap Snapshot
	resp, err := c.client.R().
		SetResult(&snap).
		SetPathParam("name", name).
		SetBody(&opts).
		Put("api/snapshots/{name}")
	if err != nil {
		return snap, err
	} else if resp.IsSuccess() {
		return snap, nil
	}
	return snap, getError(resp)
}
