package aptly

type ServerVersion struct {
	Version string
}

func (c *Client) Version() (ServerVersion, error) {
	var version ServerVersion

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
