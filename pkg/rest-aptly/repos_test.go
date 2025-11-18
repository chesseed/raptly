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
		newRawJSONResponder(200, `
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
			newRawJSONResponder(200, `
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
			newRawJSONResponder(200, `
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
		newRawJSONResponder(200, `
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
		newRawJSONResponder(200, `
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
		newRawJSONResponder(200, testPkgsSimple1.JSON))
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/repos/testRepo/packages",
		map[string]string{"q": "query", "withDeps": "1", "maximumVersion": "1"},
		newRawJSONResponder(200, testPkgsSimple2.JSON))
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/repos/testRepo/packages",
		map[string]string{"format": "details"},
		newRawJSONResponder(200, testPkgsDetailed.JSON))

	t.Run("without query", func(t *testing.T) {
		pkgs, err := client.ReposListPackages("testRepo", ListPackagesOptions{})
		assert.NoError(t, err)
		assert.Equal(t, testPkgsSimple1.Pkgs, pkgs)
	})
	t.Run("with query", func(t *testing.T) {
		pkgs, err := client.ReposListPackages("testRepo", ListPackagesOptions{Query: "query", WithDeps: true, MaximumVersion: true})
		assert.NoError(t, err)
		assert.Equal(t, testPkgsSimple2.Pkgs, pkgs)
	})
	t.Run("detailed", func(t *testing.T) {
		pkgs, err := client.ReposListPackages("testRepo", ListPackagesOptions{Detailed: true})
		assert.NoError(t, err)
		assert.Equal(t, testPkgsDetailed.Pkgs, pkgs)
	})
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

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName/fileName", newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName/force", map[string]string{"forceReplace": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName/noRemove", map[string]string{"noRemove": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))

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

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/file/dirName", newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/force", map[string]string{"forceReplace": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/file/noRemove", map[string]string{"noRemove": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))

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

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/fileName", newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/forceReplace", map[string]string{"forceReplace": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/noRemoveFiles", map[string]string{"noRemoveFiles": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/acceptUnsigned", map[string]string{"acceptUnsigned": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/dirName/ignoreSignature", map[string]string{"ignoreSignature": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))

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

	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/repos/testRepo/include/fileName", newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/forceReplace", map[string]string{"forceReplace": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/noRemoveFiles", map[string]string{"noRemoveFiles": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/acceptUnsigned", map[string]string{"acceptUnsigned": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))
	httpmock.RegisterResponderWithQuery(http.MethodPost, "http://host.local/api/repos/testRepo/include/ignoreSignature", map[string]string{"ignoreSignature": "1"}, newRawJSONResponder(200, `{"FailedFiles": []}`))

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
