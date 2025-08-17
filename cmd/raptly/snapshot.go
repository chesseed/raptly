package main

import (
	"fmt"
	aptly "raptly/pkg/rest-aptly"
	"strings"
	"unicode/utf8"
)

type SnapshotCLI struct {
	List   snapshotListCmd   `kong:"cmd,help='get list of all created snapshots'"`
	Show   snapshotShowCmd   `kong:"cmd,help='display detailed information about snapshot'"`
	Create snapshotCreateCmd `kong:"cmd,help='create snapshot from local repository or mirror'"`
	Rename snapshotRenameCmd `kong:"cmd,help='changes name of the snapshot. Snapshot name should be unique'"`
	Diff   SnapshotDiffCmd   `kong:"cmd,help='displays difference in packages between two snapshots'"`
	Drop   snapshotDropCmd   `kong:"cmd,help='removes information about snapshot'"`
	Merge  snapshotMergeCmd  `kong:"cmd,help='merges several source snapshots into new destination snapshot'"`
}

type snapshotListCmd struct{}

func (c *snapshotListCmd) Run(ctx *Context) error {
	snaps, err := ctx.client.SnapshotList()
	if err != nil {
		return err
	}

	fmt.Println("List of snapshots:")
	for _, snap := range snaps {
		fmt.Printf(" * [%s] %s\n", snap.Name, snap.Description)
	}
	return nil
}

type snapshotShowCmd struct {
	Name         string `kong:"arg,help='snapshot name which has been given during snapshot creation'"`
	WithPackages bool   `kong:"name='with-packages',help='show detailed list of packages and versions stored in the mirror'"`
}

func (c *snapshotShowCmd) Run(ctx *Context) error {
	snap, err := ctx.client.SnapshotShow(c.Name)
	if err != nil {
		return err
	}

	packages, err := ctx.client.SnapshotPackages(c.Name, aptly.ListPackagesOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", snap.Name)
	fmt.Printf("CreatedAt: %s\n", snap.CreatedAt)
	fmt.Printf("Description: %s\n", snap.Description)
	fmt.Printf("Number of packages: %v\n", len(packages))
	fmt.Print("Sources:\n")
	if snap.LocalRepos != nil {
		for _, lrepo := range snap.LocalRepos {
			fmt.Printf("  %s [%s]\n", lrepo.Name, snap.SourceKind)
		}
	}
	if snap.Snapshots != nil {
		for _, ssnap := range snap.Snapshots {
			fmt.Printf("  %s [%s]\n", ssnap.Name, snap.SourceKind)
		}
	}
	// if snap.RemoteRepos != nil {
	// 	for _, rrepo := range snap.RemoteRepos {
	// 		fmt.Printf("  %s [%s]\n", rrepo.Name, snap.SourceKind)
	// 	}
	// }

	if c.WithPackages {
		fmt.Print("Packages:\n")
		for _, pkg := range packages {
			fmt.Printf("  %s\n", pkg.Key)
		}
	}

	return nil
}

type snapshotDropCmd struct {
	Force bool   `kong:"help='drop snapshot even if it used as source in other snapshots'"`
	Name  string `kong:"arg,help='snapshot name which has been given during snapshot creation'"`
}

func (c *snapshotDropCmd) Run(ctx *Context) error {
	err := ctx.client.SnapshotDrop(c.Name, c.Force)
	if err != nil {
		return err
	}

	fmt.Printf("Snapshot `%s` has been dropped.\n", c.Name)
	return nil
}

type snapshotCreateCmd struct {
	Name struct {
		Name string `kong:"arg,help='name for the snapshot to be created'"`
		From struct {
			Repo struct {
				Repo *string `kong:"arg,help='local repository name to snapshot'"`
			} `kong:"cmd,help='create snapshot from current state of local package repository'"`
			Mirror struct {
				Mirror *string `kong:"arg,help='mirror name to snapshot'"`
			} `kong:"cmd,help='create snapshot from current state of remote mirror'"`
		} `kong:"cmd"`
		Empty struct{} `kong:"cmd,help='create empty snapshot'"`
	} `kong:"arg"`
}

func (c *snapshotCreateCmd) Run(ctx *Context) error {
	var snap aptly.Snapshot
	var err error
	if c.Name.From.Repo.Repo != nil {
		snap, err = ctx.client.SnapshotFromRepo(c.Name.Name, *c.Name.From.Repo.Repo, "")
	} else if c.Name.From.Mirror.Mirror != nil {
		snap, err = ctx.client.SnapshotFromMirror(c.Name.Name, *c.Name.From.Mirror.Mirror, "")
	} else {
		snap, err = ctx.client.SnapshotCreate(c.Name.Name, aptly.SnapshotCreateOptions{})
	}
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot '%s' successfully created.\n", snap.Name)
	return nil
}

type SnapshotDiffCmd struct {
	OnlyMatching bool   `kong:"name='only-matching',help='display diff only for package versions (don’t display missing packages)'"`
	Left         string `kong:"arg,help='snapshot name which is “on the left” during comparison'"`
	Right        string `kong:"arg,help='snapshot name which is “on the right” during comparison'"`
}

func padToWidth(str string, width int) string {
	runes := utf8.RuneCountInString(str)
	if runes < width {
		return str + strings.Repeat(" ", width-runes)
	}
	return str
}

func (c *SnapshotDiffCmd) Run(ctx *Context) error {
	diffs, err := ctx.client.SnapshotDiff(c.Left, c.Right, c.OnlyMatching)
	if err != nil {
		return err
	}

	if len(diffs) == 0 {
		fmt.Println("Snapshots are identical.")
		return nil
	}
	// minimum widths
	const (
		minArch    = 4  // len 'Arch'
		minPackage = 7  // len 'Package'
		minA       = 12 // len 'Version in A'
		minB       = 12 // len 'Version in B'
	)
	widthArch := minArch
	widthPackage := minPackage
	widthA := minA
	widthB := minB

	for _, pkgDiff := range diffs {
		if pkgDiff.Left != nil {
			widthArch = max(widthArch, utf8.RuneCountInString(pkgDiff.Left.Architecture))
			widthPackage = max(widthPackage, utf8.RuneCountInString(pkgDiff.Left.Package))
			widthA = max(widthA, utf8.RuneCountInString(pkgDiff.Left.Version))
		}
		if pkgDiff.Right != nil {
			widthArch = max(widthArch, utf8.RuneCountInString(pkgDiff.Right.Architecture))
			widthPackage = max(widthPackage, utf8.RuneCountInString(pkgDiff.Right.Package))
			widthB = max(widthB, utf8.RuneCountInString(pkgDiff.Right.Version))
		}
	}

	fmt.Printf("  %s | %s | %s | %s\n",
		padToWidth("Arch", widthArch),
		padToWidth("Package", widthPackage),
		padToWidth("Version in A", widthA),
		padToWidth("Version in B", widthB),
	)
	for _, pkgDiff := range diffs {

		// '!' for different version, '-' for missing, '+' for added
		indicator := "!"
		if pkgDiff.Left == nil {
			indicator = "+"
		} else if pkgDiff.Right == nil {
			indicator = "-"
		}
		arch := ""
		if pkgDiff.Left != nil {
			arch = pkgDiff.Left.Architecture
		} else {
			arch = pkgDiff.Right.Architecture
		}

		pkg := ""
		if pkgDiff.Left != nil {
			pkg = pkgDiff.Left.Package
		} else {
			pkg = pkgDiff.Right.Package
		}

		a := "-"
		if pkgDiff.Left != nil {
			a = pkgDiff.Left.Version
		}
		b := "-"
		if pkgDiff.Right != nil {
			b = pkgDiff.Right.Version
		}

		fmt.Printf("%s %s | %s | %s | %s\n",
			indicator,
			padToWidth(arch, widthArch),
			padToWidth(pkg, widthPackage),
			padToWidth(a, widthA),
			padToWidth(b, widthB),
		)
	}
	return nil
}

type snapshotRenameCmd struct {
	OldName string `kong:"arg,help='current snapshot name'"`
	NewName string `kong:"arg,help='new snapshot name'"`
}

func (c *snapshotRenameCmd) Run(ctx *Context) error {
	snap, err := ctx.client.SnapshotUpdate(c.OldName, aptly.SnapshotUpdateOptions{Name: c.NewName})
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot %s -> %s has been successfully renamed.\n", c.OldName, snap.Name)
	return nil
}

type snapshotMergeCmd struct {
	Destination string   `kong:"arg,help='name of the snapshot that would be created'"`
	Sources     []string `kong:"arg,help='list of snapshot names that would be merged together'"`
	Latest      bool     `kong:"name='latest',help='use only the latest version of each package'"`
	NoRemove    bool     `kong:"name='no-remove',help='don’t remove duplicate arch/name packages'"`
}

func (c *snapshotMergeCmd) Run(ctx *Context) error {
	snap, err := ctx.client.SnapshotMerge(c.Destination, c.Sources, aptly.SnapshotMergeOptions{Latest: c.Latest, NoRemove: c.NoRemove})
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot %s successfully created.\n", snap.Name)
	return nil
}
