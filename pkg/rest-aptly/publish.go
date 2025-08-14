package aptly

import (
	"strings"
)

type SourceEntry struct {
	Component string `json:"Component,omitempty"`
	Name      string `json:"Name,omitempty"`
}

type SourceEntryRequest struct {
	Component *string `json:"Component,omitempty"`
	Name      string  `json:"Name,omitempty"`
}

const (
	SourceLocalRepo  = "local"
	SourceSnapshot   = "snapshot"
	SourceRemoteRepo = "repo"
)

type PublishedList struct {
	//AcquireByHash bool `json:"acquireByHash,omitempty"`
	// which architectures are published
	Architectures []string
	Distribution  string
	Label         string
	Origin        string
	Prefix        string
	Path          string
	// local or snapshot
	SourceKind string
	Sources    []SourceEntry
}

type PublishDropOptions struct {
	Force       bool `json:",omitempty"`
	SkipCleanup bool `json:",omitempty"`
}

func escapePrefix(prefix string) string {
	if prefix == "." {
		return ":."
	}
	s1 := strings.Replace(prefix, "_", "__", -1)
	s2 := strings.Replace(s1, "/", "_", -1)

	return s2
}

type PublishOptions struct {
	Architectures []string
	Distribution  *string
	Component     *string
}

type PublishSigningOptions struct {
	// skip signing
	Skip bool `json:"Skip"`
	// GPG key ID to use when signing the release, if not specified default key is used
	GpgKey string `json:"GpgKey,omitempty"`
	// GPG keyring to use (instead of default)
	Keyring string `json:"Keyring,omitempty"`
	// GPG secret keyring to use (instead of default) Note: depreciated with gpg2
	SecretKeyring string `json:"SecretKeyring,omitempty"`
	// GPG passphrase to unlock private key (possibly insecure)
	Passphrase string `json:"Passphrase,omitempty"`
	// GPG passphrase file to unlock private key (possibly insecure)
	PassphraseFile string `json:"PassphraseFile,omitempty"`
}

func WithoutSigning() PublishSigningOptions {
	return PublishSigningOptions{Skip: true}
}

type publishedRepoCreateParams struct {
	// 'local' for local repositories and 'snapshot' for snapshots
	SourceKind string               `json:"SourceKind"`
	Sources    []SourceEntryRequest `json:"Sources"`
	// Distribution name, if missing Aptly would try to guess from sources
	Distribution *string `json:"Distribution,omitempty"`
	// // Value of Label: field in published repository stanza
	// Label optional.Option[string] `json:"Label"`
	// // Value of Origin: field in published repository stanza
	// //Origin optional.Option[string] `json:"Origin"`
	// // when publishing, overwrite files in pool/ directory without notice
	// ForceOverwrite bool `json:"ForceOverwrite"`
	// Override list of published architectures
	Architectures []string `json:"Architectures,omitempty"`
	// GPG options
	Signing PublishSigningOptions `json:"Signing"`
	// // Setting to yes indicates to the package manager to not install or upgrade packages from the repository without user consent
	// NotAutomatic string `json:"NotAutomatic"`
	// // setting to yes excludes upgrades from the NotAutomic setting
	// ButAutomaticUpgrades string `json:"ButAutomaticUpgrades"`
	// // Don't generate contents indexes
	// SkipContents optional.Option[bool] `json:"SkipContents"`
	// // Don't remove unreferenced files in prefix/component
	// SkipCleanup optional.Option[bool] `json:"SkipCleanup"`
	// // Skip bz2 compression for index files
	// SkipBz2 optional.Option[bool] `json:"SkipBz2"`
	// // Provide index files by hash
	// AcquireByHash optional.Option[bool] `json:"AcquireByHash"`
	// // Enable multiple packages with the same filename in different distributions
	// MultiDist optional.Option[bool] `json:"MultiDist"`
}

func (c *Client) PublishList() ([]PublishedList, error) {
	var lists []PublishedList

	resp, err := c.client.R().
		SetResult(&lists).
		Get("api/publish")

	if err != nil {
		return lists, err
	} else if resp.IsSuccess() {
		return lists, nil
	}
	return lists, getError(resp)
}

func (c *Client) PublishShow(distribution string, prefix string) (PublishedList, error) {
	var lists PublishedList

	resp, err := c.client.R().
		SetResult(&lists).
		SetPathParams(map[string]string{
			"name":   distribution,
			"prefix": escapePrefix(prefix),
		}).
		Get("api/publish/{prefix}/{name}")

	if err != nil {
		return lists, err
	} else if resp.IsSuccess() {
		return lists, nil
	}
	return lists, getError(resp)
}

// Drop a published repo/snapshot
func (c *Client) PublishDrop(name string, prefix string, opts PublishDropOptions) error {
	var lists PublishedList

	params := make(map[string]string)

	if opts.Force {
		params["force"] = "1"
	}
	if opts.SkipCleanup {
		params["skipCleanup"] = "1"
	}

	resp, err := c.client.R().
		SetResult(&lists).
		SetPathParams(map[string]string{
			"name":   name,
			"prefix": escapePrefix(prefix),
		}).
		SetQueryParams(params).
		Delete("api/publish/{prefix}/{name}")

	if err != nil {
		return err
	} else if resp.IsSuccess() {
		return nil
	}
	return getError(resp)
}

func (c *Client) PublishRepo(name string, prefix string, opts PublishOptions, sign PublishSigningOptions) (PublishedList, error) {

	req := publishedRepoCreateParams{
		SourceKind: SourceLocalRepo,
		Sources: []SourceEntryRequest{
			{Name: name, Component: opts.Component},
		},
		Architectures: opts.Architectures,
		Distribution:  opts.Distribution,
		Signing:       sign,
	}

	var list PublishedList
	resp, err := c.client.R().
		SetResult(&list).
		SetPathParams(map[string]string{
			"prefix": escapePrefix(prefix),
		}).
		SetBody(req).
		Post("api/publish/{prefix}")

	if err != nil {
		return list, err
	} else if resp.IsSuccess() {
		return list, nil
	}
	return list, getError(resp)
}

func (c *Client) PublishSnapshot(name string, prefix string, opts PublishOptions, sign PublishSigningOptions) (PublishedList, error) {

	req := publishedRepoCreateParams{
		SourceKind: SourceSnapshot,
		Sources: []SourceEntryRequest{
			{Name: name, Component: opts.Component},
		},
		Architectures: opts.Architectures,
		Distribution:  opts.Distribution,
		Signing:       sign,
	}

	var list PublishedList
	resp, err := c.client.R().
		SetResult(&list).
		SetPathParams(map[string]string{
			"prefix": escapePrefix(prefix),
		}).
		SetBody(req).
		Post("api/publish/{prefix}")

	if err != nil {
		return list, err
	} else if resp.IsSuccess() {
		return list, nil
	}
	return list, getError(resp)
}

type PublishUpdateOptions struct {
	// when publishing, overwrite files in pool/ directory without notice
	ForceOverwrite bool `json:"ForceOverwrite"`
	// GPG options
	Signing PublishSigningOptions `json:"Signing"`
	// Don't generate contents indexes
	SkipContents *bool `json:"SkipContents,omitempty"`
	// Skip bz2 compression for index files
	SkipBz2 *bool `json:"SkipBz2,omitempty"`
	// Don't remove unreferenced files in prefix/component
	SkipCleanup *bool `json:"SkipCleanup,omitempty"`
	// only when updating published snapshots, list of objects 'Component/Name'
	Snapshots []SourceEntryRequest `json:"Snapshots,omitempty"`
	// Provide index files by hash
	AcquireByHash *bool `json:"AcquireByHash,omitempty"`
	// Enable multiple packages with the same filename in different distributions
	MultiDist *bool `json:"MultiDist,omitempty"`
}

// Update published list to match repository
func (c *Client) PublishUpdateOrSwitch(prefix string, distribution string, opts PublishUpdateOptions) (PublishedList, error) {

	var list PublishedList
	resp, err := c.client.R().
		SetResult(&list).
		SetPathParams(map[string]string{
			"prefix":       escapePrefix(prefix),
			"distribution": escapePrefix(distribution),
		}).
		SetBody(opts).
		Put("api/publish/{prefix}/{distribution}")

	if err != nil {
		return list, err
	} else if resp.IsSuccess() {
		return list, nil
	}
	return list, getError(resp)
}
