package aptly

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetError(t *testing.T) {

	httpmock.Activate(t)
	client := resty.New()
	client.SetBaseURL("http://host.local")
	client.SetError(ApiError{})
	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient())

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/json/content-type",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(500, ApiError{Error: "json"})
		})
	httpmock.RegisterResponder(http.MethodGet, "http://host.local/plain/content-type",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(501, `{"error": "plain"}`), nil
		})
	httpmock.RegisterResponder(http.MethodGet, "http://host.local/invalid",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(502, ""), nil
		})

	resJson, err := client.R().Get("json/content-type")
	assert.NoError(t, err)
	assert.ErrorContains(t, getError(resJson), "json")
	resPlain, err := client.R().Get("plain/content-type")
	assert.NoError(t, err)
	assert.ErrorContains(t, getError(resPlain), "plain")
	resStatus, err := client.R().Get("invalid")
	assert.NoError(t, err)
	assert.ErrorContains(t, getError(resStatus), "unexpected response code 502")
}
