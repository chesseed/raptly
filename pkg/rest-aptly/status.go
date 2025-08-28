package aptly

type AptlyAPIVersion struct {
	Version string
}
type Version = AptlyAPIVersion

func (c *Client) Version() (Version, error) {
	var version Version

	resp, err := c.client.R().
		SetResult(&version).
		Get("api/version")

	if err != nil {
		return version, err
	} else if resp.IsSuccess() {
		return version, nil
	}
	return version, getError(resp)
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

	resp, err := c.client.R().
		SetResult(&storage).
		Get("api/storage")

	if err != nil {
		return storage, err
	} else if resp.IsSuccess() {
		return storage, nil
	}
	return storage, getError(resp)
}
