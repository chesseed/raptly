package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/maxatome/tdhttpmock"
	"github.com/stretchr/testify/assert"
)

func TestSnapshotList(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/snapshots",
		newRawJsonResponder(200, `
[
    {
        "Name": "snap1",
        "CreatedAt": "2025-07-31T19:36:01.504560598Z",
        "SourceKind": "local",
        "Description": "Snapshot from local repo [testrepo]",
        "Origin": "",
        "NotAutomatic": "",
        "ButAutomaticUpgrades": ""
    },
    {
        "Name": "snapTest",
        "CreatedAt": "2025-07-31T19:42:05.254886023Z",
        "SourceKind": "local",
        "Description": "test repo",
        "Origin": "",
        "NotAutomatic": "",
        "ButAutomaticUpgrades": ""
    }
]
	`))
	snaps, err := client.SnapshotList()
	assert.NoError(t, err)
	assert.Equal(t, []Snapshot{
		{
			Name:        "snap1",
			CreatedAt:   "2025-07-31T19:36:01.504560598Z",
			SourceKind:  "local",
			Description: "Snapshot from local repo [testrepo]",
		},
		{
			Name:        "snapTest",
			CreatedAt:   "2025-07-31T19:42:05.254886023Z",
			SourceKind:  "local",
			Description: "test repo",
		},
	}, snaps)
}

func TestSnapshotShow(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/snapshots/snapTest",
		newRawJsonResponder(200, `
{
	"Name": "snapTest",
	"CreatedAt": "2025-07-31T19:42:05.254886023Z",
	"SourceKind": "local",
	"Description": "test repo",
	"Origin": "",
	"NotAutomatic": "",
	"ButAutomaticUpgrades": ""
}
	`))
	snaps, err := client.SnapshotShow("snapTest")
	assert.NoError(t, err)
	assert.Equal(t, Snapshot{
		Name:        "snapTest",
		CreatedAt:   "2025-07-31T19:42:05.254886023Z",
		SourceKind:  "local",
		Description: "test repo",
	}, snaps)
}

func TestSnapshotPackages(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/snapshots/snapTest/packages",
		newRawJsonResponder(200, test_pkgs_simple1.Json))
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/snapshots/snapTest/packages",
		map[string]string{"q": "query"},
		newRawJsonResponder(200, test_pkgs_simple2.Json))
	httpmock.RegisterResponderWithQuery(http.MethodGet, "http://host.local/api/snapshots/snapTest/packages",
		map[string]string{"format": "details"},
		newRawJsonResponder(200, test_pkgs_detailed.Json))

	t.Run("without query", func(t *testing.T) {
		pkgs, err := client.SnapshotPackages("snapTest", ListPackagesOptions{})
		assert.NoError(t, err)
		assert.Equal(t, test_pkgs_simple1.Pkgs, pkgs)

	})
	t.Run("with query", func(t *testing.T) {
		pkgs, err := client.SnapshotPackages("snapTest", ListPackagesOptions{Query: "query"})
		assert.NoError(t, err)
		assert.Equal(t, test_pkgs_simple2.Pkgs, pkgs)

	})
	t.Run("detailed", func(t *testing.T) {
		pkgs, err := client.SnapshotPackages("snapTest", ListPackagesOptions{Detailed: true})
		assert.NoError(t, err)
		assert.Equal(t, test_pkgs_detailed.Pkgs, pkgs)
	})
}

func TestSnapshotDrop(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodDelete, "http://host.local/api/snapshots/testSnap", httpmock.NewStringResponder(200, "ok"))
	httpmock.RegisterResponderWithQuery(http.MethodDelete, "http://host.local/api/snapshots/testForce", map[string]string{"force": "1"}, httpmock.NewStringResponder(200, "ok"))

	err := client.SnapshotDrop("testSnap", false)
	assert.NoError(t, err)
	err = client.SnapshotDrop("testForce", true)
	assert.NoError(t, err)
}

func TestSnapshotFromRepo(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/api/repos/repo/snapshots",
		tdhttpmock.JSONBody(td.JSON(`
{
	"Name": "repoSnap",
	"Description": "description"
}
		`)),
		newRawJsonResponder(200, `
{
	"Name": "repoSnap",
	"Description": "",
	"SourceKind": "local",
	"LocalRepos": [
		{
			"Name": "repo",
			"Comment": "comment",
			"DefaultComponent": "component",
			"DefaultDistribution": "dist"
		}
	]
}
	`))

	snap, err := client.SnapshotFromRepo("repoSnap", "repo", "description")
	assert.NoError(t, err)
	assert.Equal(t, Snapshot{
		Name:        "repoSnap",
		Description: "",
		SourceKind:  "local",
		LocalRepos: []LocalRepo{
			{
				Name:                "repo",
				Comment:             "comment",
				DefaultComponent:    "component",
				DefaultDistribution: "dist",
			},
		},
	}, snap)
}

// func TestSnapshotFromMirror(t *testing.T) {
// 	assert.Fail(t, "todo")
// }

func TestSnapshotDiff(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/snapshots/snap1/diff/snap2",
		newRawJsonResponder(200, `
[
    {
        "Left": "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
        "Right": null
    },
    {
        "Left": "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
        "Right": null
    },
    {
        "Left": "Psource hello 3.0.0-2 571d33f41765ddba",
        "Right": null
	},
	{
		"Left": null,
		"Right": "Pamd64 nano 7.2-1+deb12u1 c5d2ac1639544e75"
	}
]
	`))

	diff, err := client.SnapshotDiff("snap1", "snap2", false)
	assert.NoError(t, err)
	assert.Equal(t, []PackageDiff{
		{
			Left: &Package{
				Key:          "Pamd64 hello 3.0.0-2 96e8a0deaf8fc95f",
				Architecture: "amd64",
				Package:      "hello",
				Version:      "3.0.0-2",
				FilesHash:    "96e8a0deaf8fc95f",
			},
		},
		{
			Left: &Package{
				Key:          "Pamd64 hello-dbgsym 3.0.0-2 185cc47ca86a934c",
				Architecture: "amd64",
				Package:      "hello-dbgsym",
				Version:      "3.0.0-2",
				FilesHash:    "185cc47ca86a934c",
			},
		},
		{
			Left: &Package{
				Key:          "Psource hello 3.0.0-2 571d33f41765ddba",
				Architecture: "source",
				Package:      "hello",
				Version:      "3.0.0-2",
				FilesHash:    "571d33f41765ddba",
			},
		},
		{
			Right: &Package{
				Key:          "Pamd64 nano 7.2-1+deb12u1 c5d2ac1639544e75",
				Architecture: "amd64",
				Package:      "nano",
				Version:      "7.2-1+deb12u1",
				FilesHash:    "c5d2ac1639544e75",
			},
		},
	}, diff)

}

func TestSnapshotUpdate(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterMatcherResponder(http.MethodPut, "http://host.local/api/snapshots/snap1",
		tdhttpmock.JSONBody(td.JSON(`
{
	"Name": "newName",
	"Description": "new"
}
		`)),
		newRawJsonResponder(200, `
{
	"Name": "newName",
	"Description": "new",
	"SourceKind": "local",
	"LocalRepos": [
		{
			"Name": "repo",
			"Comment": "comment",
			"DefaultComponent": "component",
			"DefaultDistribution": "dist"
		}
	]
}
	`))

	snap, err := client.SnapshotUpdate("snap1", SnapshotUpdateOptions{Name: "newName", Description: "new"})
	assert.NoError(t, err)
	assert.Equal(t, Snapshot{
		Name:        "newName",
		Description: "new",
		SourceKind:  "local",
		LocalRepos: []LocalRepo{
			{
				Name:                "repo",
				Comment:             "comment",
				DefaultComponent:    "component",
				DefaultDistribution: "dist",
			},
		},
	}, snap)
}

func TestSnapshotMerge(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/api/snapshots/snapMerged/merge",
		tdhttpmock.JSONBody(td.JSON(`
{
	"Sources": ["snap1","snap2","snap3"]
}
		`)),
		newRawJsonResponder(201, `
{
    "Name": "snapMerged",
    "CreatedAt": "2025-08-16T23:31:39.54837804+02:00",
    "SourceKind": "snapshot",
    "Description": "Merged from sources: 'snap1', 'snap2', 'snap3'",
    "Origin": "",
    "NotAutomatic": "",
    "ButAutomaticUpgrades": ""
}
	`))

	// invalid options
	_, err := client.SnapshotMerge("snapMerged", []string{"snap1", "snap2", "snap3"}, SnapshotMergeOptions{Latest: true, NoRemove: true})
	assert.Error(t, err)
	// empty list
	_, err = client.SnapshotMerge("snapMerged", []string{}, SnapshotMergeOptions{})
	assert.Error(t, err)

	snap, err := client.SnapshotMerge("snapMerged", []string{"snap1", "snap2", "snap3"}, SnapshotMergeOptions{})
	assert.NoError(t, err)
	assert.Equal(t, Snapshot{
		Name:        "snapMerged",
		CreatedAt:   "2025-08-16T23:31:39.54837804+02:00",
		SourceKind:  "snapshot",
		Description: "Merged from sources: 'snap1', 'snap2', 'snap3'",
	}, snap)
}

func TestSnapshotCreate(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/api/snapshots",
		tdhttpmock.JSONBody(td.JSON(`
{
	"Name": "snapEmpty"
}
		`)),
		newRawJsonResponder(201, `
{
	"Name": "snapEmpty",
	"CreatedAt": "2025-08-16T22:05:30.477336537Z",
	"SourceKind": "snapshot",
	"Description": "Created as empty",
	"Origin": "",
	"NotAutomatic": "",
	"ButAutomaticUpgrades": ""
}
	`))

	snap, err := client.SnapshotCreate("snapEmpty", SnapshotCreateOptions{})
	assert.NoError(t, err)
	assert.Equal(t, Snapshot{
		Name:        "snapEmpty",
		CreatedAt:   "2025-08-16T22:05:30.477336537Z",
		SourceKind:  "snapshot",
		Description: "Created as empty",
	}, snap)
}
