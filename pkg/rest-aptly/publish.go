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
	// use batch mode, required for versions older than 1.6.0, will be set automatically
	Batch bool `json:"Batch"`
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

func WithGpgKey(key string, passphrase string) PublishSigningOptions {
	return PublishSigningOptions{GpgKey: key, Passphrase: passphrase}
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

	req := c.get("api/publish").
		SetResult(&lists)

	return lists, c.send(req)
}

func (c *Client) PublishShow(distribution string, prefix string) (PublishedList, error) {
	var lists PublishedList

	req := c.get("api/publish/{prefix}/{name}").
		SetResult(&lists).
		SetPathParams(map[string]string{
			"name":   distribution,
			"prefix": escapePrefix(prefix),
		})

	return lists, c.send(req)
}

// PublishDrop deletes a published repo/snapshot
func (c *Client) PublishDrop(name string, prefix string, opts PublishDropOptions) error {

	params := make(map[string]string)

	if opts.Force {
		params["force"] = "1"
	}
	if opts.SkipCleanup {
		params["skipCleanup"] = "1"
	}

	req := c.delete("api/publish/{prefix}/{name}").
		SetPathParams(map[string]string{
			"name":   name,
			"prefix": escapePrefix(prefix),
		}).
		SetQueryParams(params)

	return c.send(req)
}

func (c *Client) PublishRepo(name string, prefix string, opts PublishOptions, sign PublishSigningOptions) (PublishedList, error) {

	reqBody := publishedRepoCreateParams{
		SourceKind: SourceLocalRepo,
		Sources: []SourceEntryRequest{
			{Name: name, Component: opts.Component},
		},
		Architectures: opts.Architectures,
		Distribution:  opts.Distribution,
		Signing:       sign,
	}
	// workaround for older aptly versions
	reqBody.Signing.Batch = true

	var list PublishedList
	req := c.post("api/publish/{prefix}").
		SetResult(&list).
		SetPathParams(map[string]string{
			"prefix": escapePrefix(prefix),
		}).
		SetBody(reqBody)

	return list, c.send(req)
}

func (c *Client) PublishSnapshot(name string, prefix string, opts PublishOptions, sign PublishSigningOptions) (PublishedList, error) {

	reqBody := publishedRepoCreateParams{
		SourceKind: SourceSnapshot,
		Sources: []SourceEntryRequest{
			{Name: name, Component: opts.Component},
		},
		Architectures: opts.Architectures,
		Distribution:  opts.Distribution,
		Signing:       sign,
	}
	// workaround for older aptly versions
	reqBody.Signing.Batch = true

	var list PublishedList
	req := c.post("api/publish/{prefix}").
		SetResult(&list).
		SetPathParams(map[string]string{
			"prefix": escapePrefix(prefix),
		}).
		SetBody(reqBody)

	return list, c.send(req)
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

// PublishUpdateOrSwitch updates published list to match repository
func (c *Client) PublishUpdateOrSwitch(prefix string, distribution string, opts PublishUpdateOptions) (PublishedList, error) {

	// workaround for older aptly versions
	opts.Signing.Batch = true

	var list PublishedList
	req := c.put("api/publish/{prefix}/{distribution}").
		SetPathParams(map[string]string{
			"prefix":       escapePrefix(prefix),
			"distribution": escapePrefix(distribution),
		}).
		SetBody(opts).
		SetResult(&list)

	return list, c.send(req)
}
