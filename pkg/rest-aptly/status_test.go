package aptly

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	// mock to list out the articles
	httpmock.RegisterResponder("GET", "http://host.local/api/version", newRawJsonResponder(200, `{"Version": "1.5.0"}`))
	v, err := client.Version()
	assert.Nil(t, err)
	assert.Equal(t, v.Version, "1.5.0")
}

func TestStorageUsage(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	// mock to list out the articles
	httpmock.RegisterResponder("GET", "http://host.local/api/storage", newRawJsonResponder(200, `{"Total": 1000, "Free": 455, "PercentFull": 55.5}`))

	s, err := client.StorageUsage()
	assert.Nil(t, err)
	assert.Equal(t, s.Total, uint64(1000))
	assert.Equal(t, s.Free, uint64(455))
	assert.Equal(t, s.PercentFull, float32(55.5))
}
