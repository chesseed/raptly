package aptly

import "fmt"

// Snapshot is immutable state of repository: list of packages
type Snapshot struct {
	// Human-readable name
	Name string `json:"Name"`
	// Date of creation
	CreatedAt  string `json:"CreatedAt"`
	SourceKind string `json:"SourceKind"`
	// Sources
	Snapshots []Snapshot `json:",omitempty"`
	//RemoteRepos []*RemoteRepo `json:",omitempty"`
	LocalRepos []LocalRepo `json:",omitempty"`
	Packages   []string    `json:",omitempty"`

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

func (c *Client) SnapshotPackages(name string, opts ListPackagesOptions) ([]Package, error) {
	params, err := opts.MakeParams()
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetPathParam("name", name).
		SetQueryParams(params).
		Get("api/snapshots/{name}/packages")

	if err != nil {
		return nil, err
	}
	return responseToPackages(resp, opts.Detailed)
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

type SnapshotCreateOptions struct {
	// Description for the snapshot
	Description string
	// List of source snapshots
	PackageRefs []string
	// List of package refs
	SourceSnapshots []string
}

func (c *Client) SnapshotCreate(name string, opts SnapshotCreateOptions) (Snapshot, error) {
	var snap Snapshot

	type createParam struct {
		Name            string   `json:"Name"`
		Description     string   `json:"Description,omitempty"`
		PackageRefs     []string `json:"PackageRefs,omitempty"`
		SourceSnapshots []string `json:"SourceSnapshots,omitempty"`
	}

	resp, err := c.client.R().
		SetResult(&snap).
		SetBody(createParam{
			Name: name, Description: opts.Description, PackageRefs: opts.PackageRefs, SourceSnapshots: opts.SourceSnapshots,
		}).
		Post("api/snapshots")

	if err != nil {
		return snap, err
	} else if resp.IsSuccess() {
		return snap, nil
	}
	return snap, getError(resp)
}

type PackageDiff struct {
	Left  *Package
	Right *Package
}

func (c *Client) SnapshotDiff(left string, right string, onlyMatching bool) ([]PackageDiff, error) {

	params := make(map[string]string)
	if onlyMatching {
		params["onlyMatching"] = "1"
	}

	type pkgDiff struct {
		Left  *string
		Right *string
	}
	var diffs []pkgDiff

	resp, err := c.client.R().
		SetResult(&diffs).
		SetQueryParams(params).
		SetPathParam("left", left).
		SetPathParam("right", right).
		Get("api/snapshots/{left}/diff/{right}")

	if err != nil {
		return nil, err
	} else if resp.IsSuccess() {
		diff := make([]PackageDiff, 0, len(diffs))

		for _, d := range diffs {
			var left, right *Package

			if d.Left != nil {
				leftPkg, err := PackageFromKey(*d.Left)
				if err != nil {
					return nil, err
				}
				left = &leftPkg
			}
			if d.Right != nil {
				rightPkg, err := PackageFromKey(*d.Right)
				if err != nil {
					return nil, err
				}
				right = &rightPkg
			}
			diff = append(diff, PackageDiff{Left: left, Right: right})
		}
		return diff, nil
	}
	return nil, getError(resp)
}

type SnapshotUpdateOptions struct {
	// new name for the snapshot
	Name string `json:"Name,omitempty"`
	// new description for the snapshot
	Description string `json:"Description,omitempty"`
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

type SnapshotMergeOptions struct {
	// use only the latest version of each package
	Latest bool
	// don’t remove duplicate arch/name packages
	NoRemove bool
}

// SnapshotMerge create snapshot by merging many into a single one
//
// since aptly 1.6.0
func (c *Client) SnapshotMerge(destination string, sources []string, opts SnapshotMergeOptions) (Snapshot, error) {
	var snap Snapshot
	type mergeRequest struct {
		Sources []string
	}
	// check for simple errors before hitting the server
	if len(sources) == 0 {
		return snap, fmt.Errorf("minimum one source snapshot is required")
	}
	if opts.Latest && opts.NoRemove {
		return snap, fmt.Errorf("minimum one source snapshot is required")
	}

	params := make(map[string]string)
	if opts.Latest {
		params["latest"] = "1"
	}
	if opts.NoRemove {
		params["no-remove"] = "1"
	}

	resp, err := c.client.R().
		SetResult(&snap).
		SetPathParam("name", destination).
		SetQueryParams(params).
		SetBody(&mergeRequest{Sources: sources}).
		Post("api/snapshots/{name}/merge")
	if err != nil {
		return snap, err
	} else if resp.IsSuccess() {
		return snap, nil
	}
	return snap, getError(resp)
}
