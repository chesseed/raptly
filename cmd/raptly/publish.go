package main

import (
	"fmt"
	aptly "raptly/pkg/rest-aptly"
)

type PublishCLI struct {
	List     PublishListCmd     `kong:"cmd,help='Lists repositories that have been published with aptly publish snapshot and aptly publish repo.'"`
	Repo     PublishRepoCmd     `kong:"cmd,help='Publishes local repository directly, bypassing snapshot creation step.'"`
	Snapshot PublishSnapshotCmd `kong:"cmd,help='Publishes snapshot as repository to be consumed by apt.'"`
	Update   PublishUpdateCmd   `kong:"cmd,help='Re-publishes (updates) published local repository.'"`
	Switch   PublishSwitchCmd   `kong:"cmd,help='Switches in-place published repository with new snapshot contents.'"`
	Show     PublishShowCmd     `kong:"cmd,help='Shows detailed information of published repository. Since Aptly 1.6.0'"`
	Drop     PublishDropCmd     `kong:"cmd,help='Remove files belonging to published repository.'"`
}

func formatPublishedRepository(list *aptly.PublishedList) string {
	publishes := ""
	for i, src := range list.Sources {
		if i > 0 {
			publishes = publishes + ", "
		}
		publishes = publishes + fmt.Sprintf("%s: [%s]", src.Component, src.Name)
	}

	if list.SourceKind == "local" {
		return fmt.Sprintf("%s %v publishes local {%s}", list.Path, list.Architectures, publishes)
	}
	return fmt.Sprintf("%s %v publishes snaphot(s) {%s}", list.Path, list.Architectures, publishes)
}

type PublishListCmd struct{}

func (c *PublishListCmd) Run(ctx *Context) error {
	lists, err := ctx.client.PublishList()
	if err != nil {
		return err
	}
	fmt.Print("Published repositories:\n")
	for _, list := range lists {
		fmt.Printf(" * %s\n", formatPublishedRepository(&list))
	}
	return nil
}

type PublishShowCmd struct {
	Distribution string `kong:"arg"`
	Prefix       string `kong:"arg"`
}

func (c *PublishShowCmd) Run(ctx *Context) error {
	list, err := ctx.client.PublishShow(c.Distribution, c.Prefix)
	if err != nil {
		return err
	}
	fmt.Printf("Prefix: %s\n", list.Prefix)
	fmt.Printf("Distribution: %s\n", list.Distribution)
	fmt.Printf("Architectures: %v\n", list.Architectures)
	fmt.Print("Sources:\n")
	for _, src := range list.Sources {
		fmt.Printf("  %s: %s [%s]\n", src.Component, src.Name, list.SourceKind)
	}

	return nil
}

type PublishDropCmd struct {
	Distribution string `kong:"arg"`
	Prefix       string `kong:"arg"`
	ForceDrop    bool   `kong:"name='force-drop'"`
	SkipCleanup  bool   `kong:"name='skip-cleanup'"`
}

func (c *PublishDropCmd) Run(ctx *Context) error {
	opts := aptly.PublishDropOptions{
		Force:       c.ForceDrop,
		SkipCleanup: c.SkipCleanup,
	}

	err := ctx.client.PublishDrop(c.Distribution, c.Prefix, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Dropped %s/%s\n", c.Prefix, c.Distribution)

	return nil
}

type PublishRepoCmd struct {
	Name         string  `kong:"arg"`
	Prefix       string  `kong:"arg"`
	Distribution *string `kong:"help='distribution name to publish; guessed from local repository default distribution'"`
	Component    *string `kong:"help='component name to publish; it is taken from local repository default, otherwise it defaults to main'"`
}

// TODO signing options
func (c *PublishRepoCmd) Run(ctx *Context) error {
	opts := aptly.PublishOptions{
		Distribution: c.Distribution,
		Component:    c.Component,
	}

	list, err := ctx.client.PublishRepo(c.Name, c.Prefix, opts, aptly.WithoutSigning())
	if err != nil {
		return err
	}
	fmt.Printf("Published: %s\n", list.Path)

	return nil
}

type PublishSnapshotCmd struct {
	Name         string  `kong:"arg"`
	Prefix       string  `kong:"arg"`
	Distribution *string `kong:"help='distribution name to publish; guessed from local repository default distribution'"`
	Component    *string `kong:"help='component name to publish; it is taken from local repository default, otherwise it defaults to main'"`
}

// TODO signing options
func (c *PublishSnapshotCmd) Run(ctx *Context) error {
	opts := aptly.PublishOptions{
		Distribution: c.Distribution,
		Component:    c.Component,
	}

	list, err := ctx.client.PublishSnapshot(c.Name, c.Prefix, opts, aptly.WithoutSigning())
	if err != nil {
		return err
	}
	fmt.Printf("Published:  %s\n", list.Path)

	return nil
}

type PublishUpdateCmd struct {
	Distribution string `kong:"arg,help='distribution name of published repository'"`
	Prefix       string `kong:"arg"`
}

// TODO signing options
func (c *PublishUpdateCmd) Run(ctx *Context) error {
	opts := aptly.PublishUpdateOptions{
		Signing: aptly.WithoutSigning(),
	}

	list, err := ctx.client.PublishUpdateOrSwitch(c.Prefix, c.Distribution, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Publish for local repo %s %v publishes {%s: [%s]} has been successfully updated.\n", list.Path, list.Architectures, list.Sources[0].Component, list.Sources[0].Name)

	return nil
}

type PublishSwitchCmd struct {
	Distribution string `kong:"arg,help='distribution name of published repository'"`
	Prefix       string `kong:"arg"`
	Snapshot     string `kong:"arg"`
	Component    string `kong:""`
}

// TODO signing options
func (c *PublishSwitchCmd) Run(ctx *Context) error {
	opts := aptly.PublishUpdateOptions{
		Signing: aptly.WithoutSigning(),
		Snapshots: []aptly.SourceEntryRequest{
			{
				Name:      c.Snapshot,
				Component: &c.Component,
			},
		},
	}

	list, err := ctx.client.PublishUpdateOrSwitch(c.Prefix, c.Distribution, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Publish for snapshot %s %v publishes {%s: [%s]} has been successfully updated.\n", list.Path, list.Architectures, list.Sources[0].Component, list.Sources[0].Name)

	return nil
}
