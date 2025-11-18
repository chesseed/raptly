package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestPackagesSearch(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/packages",
		newRawJSONResponder(200, testPkgsSimple1.JSON))
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/packages",
		map[string]string{"q": "query"},
		newRawJSONResponder(200, testPkgsSimple2.JSON))
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/packages",
		map[string]string{"format": "details"},
		newRawJSONResponder(200, testPkgsDetailed.JSON))

	t.Run("simple without query", func(t *testing.T) {
		pkgs, err := client.PackagesSearch("", false)
		assert.NoError(t, err)
		assert.Equal(t, testPkgsSimple1.Pkgs, pkgs)
	})
	t.Run("simple with query", func(t *testing.T) {
		pkgs, err := client.PackagesSearch("query", false)
		assert.NoError(t, err)
		assert.Equal(t, testPkgsSimple2.Pkgs, pkgs)
	})

	t.Run("detailed", func(t *testing.T) {
		pkgs, err := client.PackagesSearch("", true)
		assert.NoError(t, err)
		assert.Equal(t, testPkgsDetailed.Pkgs, pkgs)
	})
}

func TestPackagesInfo(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/packages/Pamd64%20hello%203.0.0-2%2096e8a0deaf8fc95f", newRawJSONResponder(200, `
{
    "Architecture": "amd64",
    "Depends": "libc6 (>= 2.34)",
    "Description": " John's hello package\n John's package is written in C\n and prints a greeting.\n .\n It is awesome.\n",
    "Filename": "hello_3.0.0-2_amd64.deb",
    "FilesHash": "96e8a0deaf8fc95f",
    "Installed-Size": "23",
    "Key": "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
    "MD5sum": "be7cbf8cf38633a26b73c4511b2d597e",
    "Maintainer": "John Doe <john@doe.com>",
    "Package": "hello",
    "Priority": "optional",
    "SHA1": "3a4c46b150d3cbe8adb27c44b5b12cca3fd63668",
    "SHA256": "52417f0e39865af616b69514bb475a2b79d3c06b02d965236e3a1e66a035cc72",
    "SHA512": "a0fc5403436286c64a8e55a885d5ca1b0ac43407550ad19a012b9cecbcae14327a8d42975672cf7f6f957e2ae812dd2159862ae143beb7982bdf698a0109bade",
    "Section": "devel",
    "ShortKey": "Pamd64 hello 3.0.0-2",
    "Size": "2648",
    "Version": "3.0.0-2"
}
	`))

	pkg, err := client.PackagesInfo("Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f")
	assert.NoError(t, err)
	assert.Equal(t, Package{
		Architecture: "amd64",
		Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
		ShortKey:     "Pamd64 hello 3.0.0-2",
		FilesHash:    "96e8a0deaf8fc95f",
		Version:      "3.0.0-2",
		Package:      "hello",
	}, pkg)
}

func TestPackageFromKey(t *testing.T) {
	pkg, err := PackageFromKey("Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f")
	assert.NoError(t, err)
	assert.Equal(t, Package{
		Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
		Architecture: "amd64",
		FilesHash:    "96e8a0deaf8fc95f",
		Version:      "3.0.0-2",
		Package:      "hello",
	}, pkg)

	// with prefix
	pkg, err = PackageFromKey("xDPamd64 hello 3.0.0-2 96e8a0deaf8fc95f")
	assert.NoError(t, err)
	assert.Equal(t, Package{
		Key:          "xDPamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
		Architecture: "amd64",
		FilesHash:    "96e8a0deaf8fc95f",
		Version:      "3.0.0-2",
		Package:      "hello",
	}, pkg)

	// failure
	_, err = PackageFromKey("96e8a0deaf8fc95f")
	assert.Error(t, err)
}
