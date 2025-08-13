package aptly

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
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

func formFileEqual(req *http.Request, name string, data []byte) (bool, error) {
	file, handler, err := req.FormFile(name)
	if err != nil {
		return false, fmt.Errorf("%s does not exist", name)
	}
	defer file.Close()

	buffer := make([]byte, handler.Size)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}
	return bytes.Equal(data, buffer), nil
}

// TODO find out how to test upload with resty
func TestFilesUpload(t *testing.T) {
	httpmock.Activate(t)
	client := NewClient("http://host.local")
	httpmock.ActivateNonDefault(client.GetClient().GetClient())

	data1 := "file 1 data"

	f1, err := os.CreateTemp("", "file0")
	assert.Nil(t, err)
	defer os.Remove(f1.Name())
	_, err = f1.WriteString(data1)
	assert.Nil(t, err)

	httpmock.RegisterResponder("POST", "http://host.local/api/files/dirTest",
		func(req *http.Request) (*http.Response, error) {
			req.ParseMultipartForm(1024 * 1024 * 4)
			equal, err := formFileEqual(req, "file0", []byte(data1))
			if err != nil {
				return nil, err
			}
			if !equal {
				return httpmock.NewStringResponse(400, "not equal"), nil
			}
			return httpmock.NewJsonResponse(200, []string{"file1", "file2", "file3"})
		})

	list, err := client.FilesUpload("dirTest", []string{f1.Name()})
	assert.Nil(t, err)
	assert.ElementsMatch(t, list, []string{"file1", "file2", "file3"})
}

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
