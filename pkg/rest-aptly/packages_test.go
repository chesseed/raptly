package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestPackagesSearch(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("GET", "http://host.local/api/packages",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"pkg1", "pkg2", "pkg3"})
		})
	httpmock.RegisterResponderWithQuery("GET", "http://host.local/api/packages",
		map[string]string{"q": "query"},
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"queried1", "queried2", "queried3"})
		})

	t.Run("without query", func(t *testing.T) {
		pkgs, err := client.PackagesSearch("")
		assert.Nil(t, err)
		assert.Equal(t, []string{"pkg1", "pkg2", "pkg3"}, pkgs)

	})
	t.Run("with query", func(t *testing.T) {
		pkgs, err := client.PackagesSearch("query")
		assert.Nil(t, err)
		assert.Equal(t, []string{"queried1", "queried2", "queried3"}, pkgs)

	})
}

// helper to assign pointer
func ptr[T any](v T) *T {
	return &v
}

func TestPackagesSearchDetailed(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponderWithQuery("GET", "http://host.local/api/packages",
		map[string]string{"format": "details"},
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `
[
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
    },
    {
        "Architecture": "amd64",
        "Auto-Built-Package": "debug-symbols",
        "Build-Ids": "7a50c209d451f1dd8c2103771fc96c2142ee059c",
        "Depends": "hello (= 3.0.0-2)",
        "Description": " debug symbols for hello\n",
        "Filename": "hello-dbgsym_3.0.0-2_amd64.deb",
        "FilesHash": "185cc47ca86a934c",
        "Installed-Size": "16",
        "Key": "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
        "MD5sum": "1464a3c2ad70765dbc349fc4a4b6eb2a",
        "Maintainer": "John Doe <john@doe.com>",
        "Package": "hello-dbgsym",
        "Priority": "optional",
        "SHA1": "3183e2c73091e5fa992e64b8ed392a59d7442a6a",
        "SHA256": "21dc7e8f5fafcf4683c233e715860fbf38328b376f3aba8b20a70ab2843b18a8",
        "SHA512": "a01b4d7559683cf5ca752842659acd719f17fe33ece94b773ca8aa3ee9c66085899e44050b0ebf79d9a8593548d9dd6f929a55e86bc4ea6b72ec52b2a43ef9bb",
        "Section": "debug",
        "ShortKey": "Pamd64 hello-dbgsym 3.0.0-2",
        "Size": "2628",
        "Source": "hello",
        "Version": "3.0.0-2"
    }
]
			`)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		})

	pkgs, err := client.PackagesSearchDetailed("")
	assert.Nil(t, err)
	assert.Equal(t, []Package{
		{
			Architecture: "amd64",
			Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
			ShortKey:     "Pamd64 hello 3.0.0-2",
			FilesHash:    "96e8a0deaf8fc95f",
			Version:      "3.0.0-2",
			Package:      "hello",
		},
		{
			Architecture: "amd64",
			Key:          "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
			ShortKey:     "Pamd64 hello-dbgsym 3.0.0-2",
			FilesHash:    "185cc47ca86a934c",
			Version:      "3.0.0-2",
			Package:      "hello-dbgsym",
			Source:       ptr("hello"),
		},
	}, pkgs)
}

func TestPackagesInfo(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("GET", "http://host.local/api/packages/Pamd64%20hello%203.0.0-2%2096e8a0deaf8fc95f",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `
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
			`)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		})

	pkg, err := client.PackagesInfo("Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f")
	assert.Nil(t, err)
	assert.Equal(t, Package{
		Architecture: "amd64",
		Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
		ShortKey:     "Pamd64 hello 3.0.0-2",
		FilesHash:    "96e8a0deaf8fc95f",
		Version:      "3.0.0-2",
	}, pkg)
}
