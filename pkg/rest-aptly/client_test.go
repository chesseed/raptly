package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	type Body struct {
		Msg string
	}
	response := Body{
		Msg: "string",
	}

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/get",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, response)
		})

	resJSON, err := client.get("get").SetResult(&Body{}).Send()
	assert.NoError(t, err)
	assert.Equal(t, &response, resJSON.Result().(*Body))
}

func TestGetError(t *testing.T) {

	httpmock.Activate(t)
	client := NewClient("http://host.local")

	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/json/content-type",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(500, APIError{ErrorMsg: ptr("json")})
		})
	httpmock.RegisterResponder(http.MethodGet, "http://host.local/plain/content-type",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(501, `{"error": "plain"}`), nil
		})
	httpmock.RegisterResponder(http.MethodGet, "http://host.local/invalid",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(502, ""), nil
		})

	resJSON, err := client.get("json/content-type").Send()
	assert.NoError(t, err)
	assert.ErrorContains(t, getError(resJSON), "json")
	resPlain, err := client.get("plain/content-type").Send()
	assert.NoError(t, err)
	assert.ErrorContains(t, getError(resPlain), "plain")
	resStatus, err := client.get("invalid").Send()
	assert.NoError(t, err)
	assert.ErrorContains(t, getError(resStatus), "unexpected response code 502")
}
