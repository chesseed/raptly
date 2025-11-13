package main

import (
	"fmt"
	"os"
	aptly "raptly/pkg/rest-aptly"
)

type publishCLI struct {
	List     publishListCmd     `kong:"cmd,help='Lists repositories that have been published with aptly publish snapshot and aptly publish repo.'"`
	Repo     publishRepoCmd     `kong:"cmd,help='Publishes local repository directly, bypassing snapshot creation step.'"`
	Snapshot publishSnapshotCmd `kong:"cmd,help='Publishes snapshot as repository to be consumed by apt.'"`
	Update   publishUpdateCmd   `kong:"cmd,help='Re-publishes (updates) published local repository.'"`
	Switch   publishSwitchCmd   `kong:"cmd,help='Switches in-place published repository with new snapshot contents.'"`
	Show     publishShowCmd     `kong:"cmd,help='Shows detailed information of published repository. Since Aptly 1.6.0'"`
	Drop     publishDropCmd     `kong:"cmd,help='Remove files belonging to published repository.'"`
}

func formatPublishedRepository(list *aptly.PublishedList) string {
	publishes := ""
	for i, src := range list.Sources {
		if i > 0 {
			publishes += ", "
		}
		publishes += fmt.Sprintf("%s: [%s]", src.Component, src.Name)
	}

	if list.SourceKind == "local" {
		return fmt.Sprintf("%s %v publishes local {%s}", list.Path, list.Architectures, publishes)
	}
	return fmt.Sprintf("%s %v publishes snaphot(s) {%s}", list.Path, list.Architectures, publishes)
}

type publishListCmd struct{}

func (c *publishListCmd) Run(ctx *Context) error {
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

type publishShowCmd struct {
	Distribution string `kong:"arg"`
	Prefix       string `kong:"arg"`
}

func (c *publishShowCmd) Run(ctx *Context) error {
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

type publishDropCmd struct {
	Distribution string `kong:"arg"`
	Prefix       string `kong:"arg"`
	ForceDrop    bool   `kong:"name='force-drop'"`
	SkipCleanup  bool   `kong:"name='skip-cleanup'"`
}

func (c *publishDropCmd) Run(ctx *Context) error {
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

// common signing options
type signingCommands struct {
	Skip    bool   `kong:"name='skip-signing',help='do not sign Release files with GPG'"`
	GpgKey  string `kong:"name='gpg-key',help='GPG key ID to use when signing the release, if not specified default key is used'"`
	Keyring string `kong:"help='GPG keyring to use (instead of default)'"`

	Passphrase     string `kong:"name='passphrase',help='GPG passphrase to unlock private key (possibly insecure)'"`
	PassFile       string `kong:"name='passphrase-file',help='GPG passphrase file to unlock private key, on the local machine NOT the server'"`
	RemotePassFile string `kong:"name='remote-passphrase-file',help='GPG passphrase file to unlock private key, on the server'"`
}

func (cmd *signingCommands) MakeSigningOptions() (aptly.PublishSigningOptions, error) {
	if cmd.Skip {
		return aptly.WithoutSigning(), nil
	}
	opts := aptly.PublishSigningOptions{
		Skip:           false,
		GpgKey:         cmd.GpgKey,
		Passphrase:     cmd.Passphrase,
		PassphraseFile: cmd.RemotePassFile,
	}
	if cmd.PassFile != "" {
		b, err := os.ReadFile(cmd.PassFile)
		if err != nil {
			return opts, err
		}
		opts.Passphrase = string(b)
	}

	return opts, nil
}

type publishRepoCmd struct {
	Name         string          `kong:"arg"`
	Prefix       string          `kong:"arg"`
	Distribution *string         `kong:"help='distribution name to publish; guessed from local repository default distribution'"`
	Component    *string         `kong:"help='component name to publish; it is taken from local repository default, otherwise it defaults to main'"`
	Signing      signingCommands `kong:"embed"` // shared
}

// TODO signing options
func (c *publishRepoCmd) Run(ctx *Context) error {
	opts := aptly.PublishOptions{
		Distribution: c.Distribution,
		Component:    c.Component,
	}

	signing, err := c.Signing.MakeSigningOptions()
	if err != nil {
		return err
	}

	list, err := ctx.client.PublishRepo(c.Name, c.Prefix, opts, signing)
	if err != nil {
		return err
	}
	fmt.Printf("Published: %s\n", list.Path)

	return nil
}

type publishSnapshotCmd struct {
	Name         string          `kong:"arg"`
	Prefix       string          `kong:"arg"`
	Distribution *string         `kong:"help='distribution name to publish; guessed from local repository default distribution'"`
	Component    *string         `kong:"help='component name to publish; it is taken from local repository default, otherwise it defaults to main'"`
	Signing      signingCommands `kong:"embed"` // shared
}

// TODO signing options
func (c *publishSnapshotCmd) Run(ctx *Context) error {
	signing, err := c.Signing.MakeSigningOptions()
	if err != nil {
		return err
	}
	opts := aptly.PublishOptions{
		Distribution: c.Distribution,
		Component:    c.Component,
	}

	list, err := ctx.client.PublishSnapshot(c.Name, c.Prefix, opts, signing)
	if err != nil {
		return err
	}
	fmt.Printf("Published:  %s\n", list.Path)

	return nil
}

type publishUpdateCmd struct {
	Distribution string          `kong:"arg,help='distribution name of published repository'"`
	Prefix       string          `kong:"arg"`
	Signing      signingCommands `kong:"embed"` // shared
}

// TODO signing options
func (c *publishUpdateCmd) Run(ctx *Context) error {
	signing, err := c.Signing.MakeSigningOptions()
	if err != nil {
		return err
	}
	opts := aptly.PublishUpdateOptions{
		Signing: signing,
	}

	list, err := ctx.client.PublishUpdateOrSwitch(c.Prefix, c.Distribution, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Publish for local repo %s %v publishes {%s: [%s]} has been successfully updated.\n", list.Path, list.Architectures, list.Sources[0].Component, list.Sources[0].Name)

	return nil
}

type publishSwitchCmd struct {
	Distribution string          `kong:"arg,help='distribution name of published repository'"`
	Prefix       string          `kong:"arg"`
	Snapshot     string          `kong:"arg"`
	Component    string          `kong:""`
	Signing      signingCommands `kong:"embed"` // shared
}

// TODO signing options
func (c *publishSwitchCmd) Run(ctx *Context) error {
	signing, err := c.Signing.MakeSigningOptions()
	if err != nil {
		return err
	}
	opts := aptly.PublishUpdateOptions{
		Signing: signing,
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
