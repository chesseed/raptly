package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFilesListDirs(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("GET", "http://host.local/api/files",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"dir1", "dir2", "dir3"})
		})

	list, err := client.FilesListDirs()
	assert.Nil(t, err)
	assert.ElementsMatch(t, list, []string{"dir1", "dir2", "dir3"})
}

func TestFilesListFiles(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("GET", "http://host.local/api/files/dirTest",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"file1", "file2", "file3"})
		})

	list, err := client.FilesListFiles("dirTest")
	assert.Nil(t, err)
	assert.ElementsMatch(t, list, []string{"file1", "file2", "file3"})
}

// TODO find out how to test upload with resty
// func TestFilesUpload(t *testing.T) {
// 	httpmock.Activate(t)
// 	client := NewClient("http://host.local")
// 	httpmock.ActivateNonDefault(client.GetClient().GetClient())
// 	assert.Fail(t, "todo")
// }

func TestFilesDeleteDir(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("DELETE", "http://host.local/api/files/dirTest",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})

	err := client.FilesDeleteDir("dirTest")
	assert.Nil(t, err)
}
func TestFilesDeleteFile(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	httpmock.RegisterResponder("DELETE", "http://host.local/api/files/dirTest/file",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})

	err := client.FilesDeleteFile("dirTest", "file")
	assert.Nil(t, err)
}
