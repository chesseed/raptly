// Package aptly provides a client to access the aptly REST API in go
package aptly

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPath(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		req := initRequest("GET", "test", nil)
		url, err := req.GetPath()
		assert.NoError(t, err)
		assert.Equal(t, "test", url)
	})

	t.Run("single replacement", func(t *testing.T) {
		req := initRequest("GET", "/{name}", nil)
		req.SetPathParam("name", "test")
		url, err := req.GetPath()
		assert.NoError(t, err)
		assert.Equal(t, "test", url)
	})

	t.Run("single replacement #2", func(t *testing.T) {
		req := initRequest("GET", "api/repos/{name}/snapshot", nil)
		req.SetPathParam("name", "test")
		url, err := req.GetPath()
		assert.NoError(t, err)
		assert.Equal(t, "api/repos/test/snapshot", url)
	})

	t.Run("map replacement", func(t *testing.T) {
		req := initRequest("GET", "{path1}/some/{path2}", nil)
		req.SetPathParams(map[string]string{
			"path1": "set",
			"path2": "path",
		})
		url, err := req.GetPath()
		assert.NoError(t, err)
		assert.Equal(t, "set/some/path", url)
	})

	t.Run("missing closing bracket", func(t *testing.T) {
		req := initRequest("GET", "{path1/some/{path2}", nil)
		req.SetPathParams(map[string]string{
			"path1": "set",
			"path2": "path",
		})
		_, err := req.GetPath()
		assert.EqualError(t, err, "missing closing bracket at '{path1/'")
	})

	t.Run("missing value", func(t *testing.T) {
		req := initRequest("GET", "{path}", nil)
		_, err := req.GetPath()
		assert.EqualError(t, err, "path parameter '{path}' not set")
	})
}
func TestGetURL(t *testing.T) {
	t.Run("no query params", func(t *testing.T) {
		req := initRequest("GET", "url", nil)
		url, err := req.GetURL("http://test")
		assert.NoError(t, err)
		assert.Equal(t, "http://test/url", url)
	})
	t.Run("base trailing slash", func(t *testing.T) {
		req := initRequest("GET", "url", nil)
		url, err := req.GetURL("http://test/")
		assert.NoError(t, err)
		assert.Equal(t, "http://test/url", url)
	})
	t.Run("base trailing slash + leading slash", func(t *testing.T) {
		req := initRequest("GET", "/url", nil)
		url, err := req.GetURL("http://test/")
		assert.NoError(t, err)
		assert.Equal(t, "http://test/url", url)
	})
	t.Run("with query params", func(t *testing.T) {
		req := initRequest("GET", "url", nil)
		req.SetQueryParams(map[string]string{
			"foo": "",
			"qux": "baz",
			"bar": "test",
		})
		url, err := req.GetURL("http://test/")
		assert.NoError(t, err)
		assert.Equal(t, "http://test/url?bar=test&foo=&qux=baz", url)
	})
}
