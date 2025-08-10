package aptly

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	// mock to list out the articles
	httpmock.RegisterResponder("GET", "http://host.local/api/version",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, json.RawMessage(`{"Version": "1.5.0"}`))
		})
	v, err := client.Version()
	assert.Nil(t, err)
	assert.Equal(t, v.Version, "1.5.0")
}

func TestStorageUsage(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	// mock to list out the articles
	httpmock.RegisterResponder("GET", "http://host.local/api/storage",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, json.RawMessage(`{"Total": 1000, "Free": 455, "PercentFull": 55.5}`))
		})
	s, err := client.StorageUsage()
	assert.Nil(t, err)
	assert.Equal(t, s.Total, uint64(1000))
	assert.Equal(t, s.Free, uint64(455))
	assert.Equal(t, s.PercentFull, float32(55.5))
}
