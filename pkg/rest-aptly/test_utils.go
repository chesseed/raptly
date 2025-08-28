package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func newRawJSONResponse(status int, body string) *http.Response {
	resp := httpmock.NewStringResponse(status, body)
	resp.Header.Add("Content-Type", "application/json")
	return resp
}

func newRawJSONResponder(status int, body string) httpmock.Responder {
	return httpmock.ResponderFromResponse(newRawJSONResponse(status, body))
}

// activate httpmock and enable it in client
func clientForTest(t *testing.T, base string) *Client {
	httpmock.Activate(t)
	client := NewClient(base)
	// Get the underlying HTTP Client and set it to Mock
	httpmock.ActivateNonDefault(client.GetClient().GetClient())
	return client
}

// helper to assign pointer
func ptr[T any](v T) *T {
	return &v
}
