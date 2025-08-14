package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/maxatome/go-testdeep/td"
	"github.com/maxatome/tdhttpmock"
	"github.com/stretchr/testify/assert"
)

func TestEscapePrefix(t *testing.T) {
	assert.Equal(t, escapePrefix(""), "")
	assert.Equal(t, escapePrefix("test/test"), "test_test")
	assert.Equal(t, escapePrefix("test_test"), "test__test")
	assert.Equal(t, escapePrefix("part/slug_slug"), "part_slug__slug")
	assert.Equal(t, escapePrefix("."), ":.")
}

func TestPublishList(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("GET", "http://host.local/api/publish",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `
[
	{
		"AcquireByHash": false,
		"Architectures": [
			"amd64",
			"arm64",
			"source"
		],
		"ButAutomaticUpgrades": "",
		"Distribution": "bookworm",
		"Label": "",
		"NotAutomatic": "",
		"Origin": "",
		"Path": "repo/bookworm",
		"Prefix": "repo",
		"SkipContents": false,
		"SourceKind": "local",
		"Sources": [
			{
				"Component": "main",
				"Name": "testing"
			}
		],
		"Storage": "",
		"Suite": ""
	},
	{
		"AcquireByHash": false,
		"Architectures": [
			"amd64",
			"arm64"
		],
		"ButAutomaticUpgrades": "",
		"Distribution": "bookworm",
		"Label": "",
		"NotAutomatic": "",
		"Origin": "",
		"Path": "snap/bookworm",
		"Prefix": "snap",
		"SkipContents": false,
		"SourceKind": "snapshot",
		"Sources": [
			{
				"Component": "main",
				"Name": "testing-1"
			}
		],
		"Storage": "",
		"Suite": ""
	}
]`)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		})

	list, err := client.PublishList()
	assert.Nil(t, err)
	assert.Equal(t, list, []PublishedList{
		{
			Architectures: []string{"amd64", "arm64", "source"},
			Distribution:  "bookworm",
			Label:         "",
			Origin:        "",
			Prefix:        "repo",
			Path:          "repo/bookworm",
			SourceKind:    "local",
			Sources:       []SourceEntry{{Name: "testing", Component: "main"}},
		},
		{
			Architectures: []string{"amd64", "arm64"},
			Distribution:  "bookworm",
			Label:         "",
			Origin:        "",
			Prefix:        "snap",
			Path:          "snap/bookworm",
			SourceKind:    "snapshot",
			Sources:       []SourceEntry{{Name: "testing-1", Component: "main"}},
		},
	})
}

func TestPublishShow(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("GET", "http://host.local/api/publish/snap/bookworm",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `
{
	"AcquireByHash": false,
	"Architectures": [
		"amd64",
		"arm64"
	],
	"ButAutomaticUpgrades": "",
	"Distribution": "bookworm",
	"Label": "",
	"NotAutomatic": "",
	"Origin": "",
	"Path": "snap/bookworm",
	"Prefix": "snap",
	"SkipContents": false,
	"SourceKind": "snapshot",
	"Sources": [
		{
			"Component": "main",
			"Name": "testing-1"
		}
	],
	"Storage": "",
	"Suite": ""
}
`)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		})

	published, err := client.PublishShow("bookworm", "snap")
	assert.Nil(t, err)
	assert.Equal(t, PublishedList{
		Architectures: []string{"amd64", "arm64"},
		Distribution:  "bookworm",
		Label:         "",
		Origin:        "",
		Prefix:        "snap",
		Path:          "snap/bookworm",
		SourceKind:    "snapshot",
		Sources:       []SourceEntry{{Name: "testing-1", Component: "main"}},
	}, published)
}

func TestPublishDrop(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	t.Run("without parameters", func(t *testing.T) {
		httpmock.Activate(t)
		client := NewClient("http://host.local")
		httpmock.ActivateNonDefault(client.GetClient().GetClient())

		httpmock.RegisterResponderWithQuery("DELETE", "http://host.local/api/publish/simple/bookworm", map[string]string{},
			httpmock.NewStringResponder(200, "ok").Once())

		err := client.PublishDrop("bookworm", "simple", PublishDropOptions{})
		assert.Nil(t, err)
	})

	t.Run("with parameters", func(t *testing.T) {
		httpmock.Activate(t)
		client := NewClient("http://host.local")
		httpmock.ActivateNonDefault(client.GetClient().GetClient())

		httpmock.RegisterResponderWithQuery("DELETE", "http://host.local/api/publish/params/bookworm", map[string]string{"force": "1", "skipCleanup": "1"},
			httpmock.NewStringResponder(200, "ok").Once())
		err := client.PublishDrop("bookworm", "params", PublishDropOptions{Force: true, SkipCleanup: true})
		assert.Nil(t, err)
	})
}

func TestPublishRepo(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterMatcherResponder("POST", "http://host.local/api/publish/prefix",
		httpmock.Matcher{}.And(
			tdhttpmock.JSONBody(td.JSONPointer("/SourceKind", "local")),
			tdhttpmock.JSONBody(td.JSONPointer("/Sources/0/Name", "testing")),
			tdhttpmock.JSONBody(td.JSONPointer("/Signing/Skip", true)),
		),
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, `
{
	"AcquireByHash": false,
	"Architectures": [
		"amd64",
		"arm64"
	],
	"ButAutomaticUpgrades": "",
	"Distribution": "bookworm",
	"Label": "",
	"NotAutomatic": "",
	"Origin": "",
	"Path": "prefix/bookworm",
	"Prefix": "prefix",
	"SkipContents": false,
	"SourceKind": "local",
	"Sources": [
		{
			"Component": "main",
			"Name": "testing"
		}
	],
	"Storage": "",
	"Suite": ""
}
`)
			resp.Header.Add("Content-Type", "application/json")
			return resp, nil
		},
	)

	published, err := client.PublishRepo("testing", "prefix", PublishOptions{}, WithoutSigning())
	assert.Nil(t, err)
	assert.Equal(t, PublishedList{
		Architectures: []string{"amd64", "arm64"},
		Distribution:  "bookworm",
		Label:         "",
		Origin:        "",
		Prefix:        "prefix",
		Path:          "prefix/bookworm",
		SourceKind:    "local",
		Sources:       []SourceEntry{{Name: "testing", Component: "main"}},
	}, published)

	// TODO more complicated options
}

// func TestPublishSnapshot(t *testing.T) {
// 	assert.Fail(t, "TODO")
// }

// func TestPublishUpdateOrSwitch(t *testing.T) {
// 	assert.Fail(t, "TODO")
// }
