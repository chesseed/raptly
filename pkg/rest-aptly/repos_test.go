package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/maxatome/tdhttpmock"
	"github.com/stretchr/testify/assert"
)

func TestReposList(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/repos",
		newRawJsonResponder(200, `
[
    {
        "Name": "secondRepo",
        "Comment": "",
        "DefaultDistribution": "",
        "DefaultComponent": "main"
    },
    {
        "Name": "testrepo",
        "Comment": "I'm a comment",
        "DefaultDistribution": "bookworm",
        "DefaultComponent": "main"
    }
]
	`))
	repos, err := client.ReposList()
	assert.NoError(t, err)
	assert.Equal(t, []LocalRepo{
		{
			Name:             "secondRepo",
			DefaultComponent: "main",
		},
		{
			Name:                "testrepo",
			Comment:             "I'm a comment",
			DefaultDistribution: "bookworm",
			DefaultComponent:    "main",
		},
	}, repos)
}

func TestReposCreate(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	t.Run("name only", func(t *testing.T) {
		httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/api/repos",
			tdhttpmock.JSONBody(td.JSON(`
{
	"Name": "firstRepo"
}
		`)),
			newRawJsonResponder(200, `
{
	"Name": "firstRepo",
	"Comment": "",
	"DefaultDistribution": "",
	"DefaultComponent": ""
}
	`))

		repo, err := client.ReposCreate("firstRepo", RepoCreateOptions{})
		assert.NoError(t, err)
		assert.Equal(t, LocalRepo{
			Name: "firstRepo",
		}, repo)

	})

	t.Run("all options", func(t *testing.T) {
		httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/api/repos",
			tdhttpmock.JSONBody(td.JSON(`
{
	"Name": "secondRepo",
	"DefaultComponent": "other",
	"DefaultDistribution": "dist",
	"FromSnapshot": "snap",
	"Comment": "my comment"
}
		`)),
			newRawJsonResponder(200, `
{
	"Name": "secondRepo",
	"Comment": "my comment",
	"DefaultDistribution": "dist",
	"DefaultComponent": "other"
}
	`))

		repo, err := client.ReposCreate("secondRepo", RepoCreateOptions{Comment: "my comment", DefaultComponent: "other", DefaultDistribution: "dist", FromSnapshot: "snap"})
		assert.NoError(t, err)
		assert.Equal(t, LocalRepo{
			Name:                "secondRepo",
			DefaultComponent:    "other",
			Comment:             "my comment",
			DefaultDistribution: "dist",
		}, repo)

	})
}

func TestReposEdit(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterMatcherResponder(http.MethodPut, "http://host.local/api/repos/original",
		tdhttpmock.JSONBody(td.JSON(`
{
	"Name": "edited",
	"DefaultComponent": "newComponent",
	"DefaultDistribution": "newDist",
	"Comment": "longer comment"
}
		`)),
		newRawJsonResponder(200, `
{
	"Name": "edited",
	"DefaultComponent": "newComponent",
	"DefaultDistribution": "newDist",
	"Comment": "longer comment"
}
	`))

	repo, err := client.ReposEdit("original", RepoUpdateOptions{Name: "edited", Comment: "longer comment", DefaultComponent: "newComponent", DefaultDistribution: "newDist"})
	assert.NoError(t, err)
	assert.Equal(t, LocalRepo{
		Name:                "edited",
		DefaultComponent:    "newComponent",
		Comment:             "longer comment",
		DefaultDistribution: "newDist",
	}, repo)
}

func TestReposShow(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/repos/secondRepo",
		newRawJsonResponder(200, `
{
	"Name": "secondRepo",
	"Comment": "",
	"DefaultDistribution": "",
	"DefaultComponent": "main"
}
	`))
	repo, err := client.ReposShow("secondRepo")
	assert.NoError(t, err)
	assert.Equal(t, LocalRepo{
		Name:             "secondRepo",
		DefaultComponent: "main",
	}, repo)
}

func TestReposListPackages(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/repos/testRepo/packages",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"pkg1", "pkg2", "pkg3"})
		})
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/repos/testRepo/packages",
		map[string]string{"q": "query", "withDeps": "1", "maximumVersion": "1"},
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"queried1", "queried2", "queried3"})
		})

	t.Run("without query", func(t *testing.T) {
		pkgs, err := client.ReposListPackages("testRepo", ListPackagesOptions{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"pkg1", "pkg2", "pkg3"}, pkgs)

	})
	t.Run("with query", func(t *testing.T) {
		pkgs, err := client.ReposListPackages("testRepo", ListPackagesOptions{Query: "query", WithDeps: true, MaximumVersion: true})
		assert.NoError(t, err)
		assert.Equal(t, []string{"queried1", "queried2", "queried3"}, pkgs)

	})
}

func TestReposListPackagesDetailed(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/repos/testRepo/packages",
		map[string]string{"format": "details"},
		newRawJsonResponder(200, `
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
    },
    {
        "Architecture": "any",
        "Binary": "hello",
        "Build-Depends": "build-essential, debhelper (>= 9)",
        "Checksums-Sha1": " 3f0a502de585a30e24d7c7141559602eced32858 470 hello_3.0.0-2.dsc\n 062e2e42233c6fbe058a44e3c50ef1bf454acc96 3448 hello_3.0.0-2.tar.gz\n",
        "Checksums-Sha256": " f3767c240a5221e6122e1e561bba81ab36891218a6f5471b8705e2913df9e93c 470 hello_3.0.0-2.dsc\n b84597204d5ee78dbdc9e2fe041d93aa19c444d145e21ec16bfb4602ecb36f99 3448 hello_3.0.0-2.tar.gz\n",
        "Checksums-Sha512": " 37c9da0f380303329908d00fe0c9806b215e12721faae8e6c056a3c1f0916679800f660f51ba990ca3577303a3dd982c6900959b40052afc5c88d696ee607ab2 470 hello_3.0.0-2.dsc\n caaa02e2bc9de1d7cbfdd6c7759c974c72ec0b58650e12ad34c5b7f895e67e7d4327ce4e3256e7cfcd14ee4a306ccc3f1bd5d9bf61cedf88edbfd40e7bb59243 3448 hello_3.0.0-2.tar.gz\n",
        "Files": " 58e1956baa409b0980474b33cb5a9e99 470 hello_3.0.0-2.dsc\n 30be0886385224b34c96853cf52262fe 3448 hello_3.0.0-2.tar.gz\n",
        "FilesHash": "571d33f41765ddba",
        "Format": "1.0",
        "Key": "Psource hello 3.0.0-2 571d33f41765ddba",
        "Maintainer": "John Doe <john@doe.com>",
        "Package": "hello",
        "Package-List": " hello deb devel optional arch=any\n",
        "ShortKey": "Psource hello 3.0.0-2",
        "Version": "3.0.0-2"
    }
]
	`))

	pkgs, err := client.ReposListPackagesDetailed("testRepo", ListPackagesOptions{})
	assert.NoError(t, err)
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
		{
			Architecture: "any",
			Key:          "Psource hello 3.0.0-2 571d33f41765ddba",
			ShortKey:     "Psource hello 3.0.0-2",
			FilesHash:    "571d33f41765ddba",
			Version:      "3.0.0-2",
			Package:      "hello",
		},
	}, pkgs)
}

func TestReposDrop(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodDelete, "http://host.local/api/repos/testRepo", httpmock.NewStringResponder(200, "ok"))
	httpmock.RegisterResponderWithQuery(http.MethodDelete, "http://host.local/api/repos/testForces", map[string]string{"force": "1"}, httpmock.NewStringResponder(200, "ok"))

	err := client.ReposDrop("testRepo", false)
	assert.NoError(t, err)
	err = client.ReposDrop("testForces", true)
	assert.NoError(t, err)
}

func TestReposAddFile(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName/fileName", newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName/force", map[string]string{"forceReplace": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName/noRemove", map[string]string{"noRemove": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))

	t.Run("without options", func(t *testing.T) {
		res, err := client.ReposAddFile("testRepo", "dirName", "fileName", RepoAddOptions{})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with force", func(t *testing.T) {
		res, err := client.ReposAddFile("testRepo", "dirName", "force", RepoAddOptions{ForceReplace: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with noRemove", func(t *testing.T) {
		res, err := client.ReposAddFile("testRepo", "dirName", "noRemove", RepoAddOptions{NoRemove: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
}

func TestReposAddDirectory(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName", newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/force", map[string]string{"forceReplace": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/noRemove", map[string]string{"noRemove": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))

	t.Run("without options", func(t *testing.T) {
		res, err := client.ReposAddDirectory("testRepo", "dirName", RepoAddOptions{})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with force", func(t *testing.T) {
		res, err := client.ReposAddDirectory("testRepo", "force", RepoAddOptions{ForceReplace: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with noRemove", func(t *testing.T) {
		res, err := client.ReposAddDirectory("testRepo", "noRemove", RepoAddOptions{NoRemove: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
}

func TestReposIncludeFile(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/fileName", newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/forceReplace", map[string]string{"forceReplace": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/noRemoveFiles", map[string]string{"noRemoveFiles": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/acceptUnsigned", map[string]string{"acceptUnsigned": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/ignoreSignature", map[string]string{"ignoreSignature": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))

	t.Run("without options", func(t *testing.T) {
		res, err := client.ReposIncludeFile("testRepo", "dirName", "fileName", RepoIncludeOptions{})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with forceReplace", func(t *testing.T) {
		res, err := client.ReposIncludeFile("testRepo", "dirName", "forceReplace", RepoIncludeOptions{ForceReplace: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with noRemove", func(t *testing.T) {
		res, err := client.ReposIncludeFile("testRepo", "dirName", "noRemoveFiles", RepoIncludeOptions{NoRemove: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with acceptUnsigned", func(t *testing.T) {
		res, err := client.ReposIncludeFile("testRepo", "dirName", "acceptUnsigned", RepoIncludeOptions{AcceptUnsigned: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with ignoreSignature", func(t *testing.T) {
		res, err := client.ReposIncludeFile("testRepo", "dirName", "ignoreSignature", RepoIncludeOptions{IgnoreSignature: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
}

func TestReposIncludeDirectory(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/include/fileName", newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/forceReplace", map[string]string{"forceReplace": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/noRemoveFiles", map[string]string{"noRemoveFiles": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/acceptUnsigned", map[string]string{"acceptUnsigned": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/ignoreSignature", map[string]string{"ignoreSignature": "1"}, newRawJsonResponder(200, `{"FailedFiles": []}`))

	t.Run("without options", func(t *testing.T) {
		res, err := client.ReposIncludeDirectory("testRepo", "fileName", RepoIncludeOptions{})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with forceReplace", func(t *testing.T) {
		res, err := client.ReposIncludeDirectory("testRepo", "forceReplace", RepoIncludeOptions{ForceReplace: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with noRemove", func(t *testing.T) {
		res, err := client.ReposIncludeDirectory("testRepo", "noRemoveFiles", RepoIncludeOptions{NoRemove: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with acceptUnsigned", func(t *testing.T) {
		res, err := client.ReposIncludeDirectory("testRepo", "acceptUnsigned", RepoIncludeOptions{AcceptUnsigned: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
	t.Run("with ignoreSignature", func(t *testing.T) {
		res, err := client.ReposIncludeDirectory("testRepo", "ignoreSignature", RepoIncludeOptions{IgnoreSignature: true})
		assert.NoError(t, err)
		assert.Empty(t, res.FailedFiles)
	})
}
