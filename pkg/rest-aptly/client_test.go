package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
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

	req := client.get("get")
	res, err := client.Send(req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestPost(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	type Body struct {
		Msg string
	}
	response := Body{
		Msg: "string",
	}

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/post",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, response)
		})

	req := client.post("post")
	res, err := client.Send(req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, 200, res.StatusCode)
	}
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

	req := client.put("put")
	res, err := client.Send(req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, 200, res.StatusCode)
	}
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

	req := client.delete("delete")
	res, err := client.Send(req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, 200, res.StatusCode)
	}
}

func TestCheckResponseForError(t *testing.T) {
	// ok
	resNoError := httpmock.NewStringResponse(200, "")
	noErr := checkResponseForError(resNoError)
	assert.NoError(t, noErr)

	// no json error struct
	resOnlyCode := httpmock.NewStringResponse(400, "")
	errOnlyCode := checkResponseForError(resOnlyCode)
	assert.EqualError(t, errOnlyCode, "unexpected status code 400")

	// no json provides "{}"
	resWrongJson := httpmock.NewStringResponse(520, "{}")
	errWrongJson := checkResponseForError(resWrongJson)
	assert.EqualError(t, errWrongJson, "unexpected status code 520")

	// json error object
	resError := httpmock.NewStringResponse(400, "{\"error\": \"test\"}")
	errJsonErr := checkResponseForError(resError)
	assert.EqualError(t, errJsonErr, "test")
}
