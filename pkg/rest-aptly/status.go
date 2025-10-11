package aptly

type AptlyAPIVersion struct {
	Version string
}
type Version = AptlyAPIVersion

func (c *Client) Version() (Version, error) {
	req := c.get("api/version")
	return callAPIwithResult[Version](c, req)
}

type StorageUsage struct {
	Free        uint64
	Total       uint64
	PercentFull float32
}

// StorageUsage get amount of memory used/free on disk
// since Aptly 1.6.0
func (c *Client) StorageUsage() (StorageUsage, error) {
	req := c.get("api/storage")
	return callAPIwithResult[StorageUsage](c, req)
}
