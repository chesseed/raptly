package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/maxatome/tdhttpmock"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	type Body struct {
		Msg string
	}
	response := Body{
		Msg: "string",
	}

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/get",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, response)
		})

	reqJSON := client.get("get")
	resJSON, err := callAPIwithResult[Body](client, reqJSON)
	assert.NoError(t, err)
	assert.Equal(t, response, resJSON)
}

func TestPost(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	type Body struct {
		Msg string
	}
	response := Body{
		Msg: "string",
	}
	request := Body{
		Msg: "request body",
	}

	httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/post",
		tdhttpmock.JSONBody(td.JSON(`
	{
		"Msg": "request body"
	}
		`)),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, response)
		})

	reqJSON := client.post("post")
	reqJSON.SetBody(&request)
	resJSON, err := callAPIwithResult[Body](client, reqJSON)
	assert.NoError(t, err)
	assert.Equal(t, response, resJSON)
}

func TestPut(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	type Body struct {
		Msg string
	}
	response := Body{
		Msg: "string",
	}

	httpmock.RegisterResponder(http.MethodPut, "http://host.local/put",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, response)
		})

	reqJSON := client.put("put")
	resJSON, err := callAPIwithResult[Body](client, reqJSON)
	assert.NoError(t, err)
	assert.Equal(t, response, resJSON)
}

func TestDelete(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	type Body struct {
		Msg string
	}
	response := Body{
		Msg: "string",
	}

	httpmock.RegisterResponder(http.MethodDelete, "http://host.local/delete",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, response)
		})

	reqJSON := client.delete("delete")
	resJSON, err := callAPIwithResult[Body](client, reqJSON)
	assert.NoError(t, err)
	assert.Equal(t, response, resJSON)
}

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
