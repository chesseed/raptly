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

	req := c.get("api/snapshots").
		SetResult(&snaps)

	return snaps, c.send(req)
}

func (c *Client) SnapshotShow(name string) (Snapshot, error) {
	var snap Snapshot

	req := c.get("api/snapshots/{name}").
		SetResult(&snap).
		SetPathParam("name", name)

	return snap, c.send(req)
}

func (c *Client) SnapshotPackages(name string, opts ListPackagesOptions) ([]Package, error) {
	params, err := opts.MakeParams()
	if err != nil {
		return nil, err
	}

	req := c.get("api/snapshots/{name}/packages").
		SetPathParam("name", name).
		SetQueryParams(params)

	return sendPackagesRequest(req, opts.Detailed)
}

func (c *Client) SnapshotDrop(name string, force bool) error {
	params := make(map[string]string)
	if force {
		params["force"] = "1"
	}

	req := c.delete("api/snapshots/{name}").
		SetQueryParams(params).
		SetPathParam("name", name)

	return c.send(req)
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

	req := c.post("api/repos/{name}/snapshots").
		SetPathParam("name", repoName).
		SetResult(&snap).
		SetBody(params)

	return snap, c.send(req)
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

	req := c.post("api/mirrors/{name}/snapshots").
		SetPathParam("name", mirror).
		SetResult(&snap).
		SetBody(params)

	return snap, c.send(req)
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

	req := c.post("api/snapshots").
		SetResult(&snap).
		SetBody(createParam{
			Name: name, Description: opts.Description, PackageRefs: opts.PackageRefs, SourceSnapshots: opts.SourceSnapshots,
		})

	return snap, c.send(req)
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

	req := c.get("api/snapshots/{left}/diff/{right}").
		SetPathParam("left", left).
		SetPathParam("right", right).
		SetResult(&diffs).
		SetQueryParams(params)

	err := c.send(req)

	if err == nil {
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
	} else {
		return nil, err
	}
}

type SnapshotUpdateOptions struct {
	// new name for the snapshot
	Name string `json:"Name,omitempty"`
	// new description for the snapshot
	Description string `json:"Description,omitempty"`
}

func (c *Client) SnapshotUpdate(name string, opts SnapshotUpdateOptions) (Snapshot, error) {
	var snap Snapshot
	req := c.put("api/snapshots/{name}").
		SetPathParam("name", name).
		SetResult(&snap).
		SetBody(&opts)

	return snap, c.send(req)
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

	req := c.post("api/snapshots/{name}/merge").
		SetPathParam("name", destination).
		SetResult(&snap).
		SetQueryParams(params).
		SetBody(&mergeRequest{Sources: sources})

	return snap, c.send(req)
}
