package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// func TestGet(t *testing.T) {
// 	httpmock.Activate(t)
// 	client := NewClient("http://host.local")

// 	// Get the underlying HTTP Client and set it to Mock
// 	httpmock.ActivateNonDefault(client.GetClient().GetClient())

// 	type Body struct {
// 		Msg string
// 	}
// 	response := Body{
// 		Msg: "string",
// 	}

// 	httpmock.RegisterResponder(http.MethodGet, "http://host.local/get",
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewJsonResponse(200, response)
// 		})

// 	resJSON, err := client.get("get").SetResult(&Body{}).Send()
// 	assert.NoError(t, err)
// 	assert.Equal(t, &response, resJSON.Result().(*Body))
// }

func TestGetError(t *testing.T) {

	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/json/content-type",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(500, APIError{ErrorMsg: ptr("json")})
		})
	httpmock.RegisterResponder(http.MethodGet, "http://host.local/plain/content-type",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(501, `{"error": "plain"}`), nil
		})
	httpmock.RegisterResponder(http.MethodGet, "http://host.local/invalid",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(502, ""), nil
		})
	// with content-type header
	reqJSON := client.get("json/content-type")
	errJSON := callAPI(client, reqJSON)
	assert.ErrorContains(t, errJSON, "json")
	// without content-type header
	reqPlain := client.get("plain/content-type")
	errPlain := callAPI(client, reqPlain)
	assert.ErrorContains(t, errPlain, "plain")
	// non aptly error
	reqInvalid := client.get("invalid")
	errInvalid := callAPI(client, reqInvalid)
	assert.ErrorContains(t, errInvalid, "unexpected response code 502")
}
