package aptly

type AptlyAPIVersion struct {
	Version string
}
type Version = AptlyAPIVersion

func (c *Client) Version() (Version, error) {
	var version Version

	req := c.get("api/version").
		SetResult(&version)

	return version, c.send(req)
}

type StorageUsage struct {
	Free        uint64
	Total       uint64
	PercentFull float32
}

// StorageUsage get amount of memory used/free on disk
// since Aptly 1.6.0
func (c *Client) StorageUsage() (StorageUsage, error) {
	var storage StorageUsage

	req := c.get("api/storage").
		SetResult(&storage)

	return storage, c.send(req)
}
