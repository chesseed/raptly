package aptly

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFilesListDirs(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/files",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"dir1", "dir2", "dir3"})
		})

	list, err := client.FilesListDirs()
	assert.NoError(t, err)
	assert.ElementsMatch(t, list, []string{"dir1", "dir2", "dir3"})
}

func TestFilesListFiles(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodGet, "http://host.local/api/files/dirTest",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"file1", "file2", "file3"})
		})

	list, err := client.FilesListFiles("dirTest")
	assert.NoError(t, err)
	assert.ElementsMatch(t, list, []string{"file1", "file2", "file3"})
}

func TestFilesDeleteDir(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodDelete, "http://host.local/api/files/dirTest",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})

	err := client.FilesDeleteDir("dirTest")
	assert.NoError(t, err)
}
func TestFilesDeleteFile(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	httpmock.RegisterResponder(http.MethodDelete, "http://host.local/api/files/dirTest/file",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, ""), nil
		})

	err := client.FilesDeleteFile("dirTest", "file")
	assert.NoError(t, err)
}

// func formFileEqual(req *http.Request, name string, data []byte) (bool, error) {
// 	file, handler, err := req.FormFile(name)
// 	if err != nil {
// 		return false, fmt.Errorf("%s does not exist", name)
// 	}
// 	defer file.Close()

// 	buffer := make([]byte, handler.Size)
// 	_, err = file.Read(buffer)
// 	if err != nil {
// 		return false, err
// 	}
// 	return bytes.Equal(data, buffer), nil
// }

// // TODO find out how to test upload with resty
// func TestFilesUpload(t *testing.T) {
// 	client := clientForTest(t, "http://host.local")

// 	data1 := "file0 data"

// 	f1, err := os.CreateTemp("", "file0_")
// 	assert.NoError(t, err)
// 	defer os.Remove(f1.Name())
// 	_, err = f1.WriteString(data1)
// 	assert.NoError(t, err)

// 	httpmock.RegisterResponder(http.MethodPost, "http://host.local/api/files/dirTest",
// 		func(req *http.Request) (*http.Response, error) {
// 			err := req.ParseMultipartForm(1024 * 1024 * 4)
// 			if err != nil {
// 				return nil, err
// 			}

// 			equal, err := formFileEqual(req, "file0", []byte(data1))
// 			if err != nil {
// 				return nil, err
// 			}
// 			if !equal {
// 				return httpmock.NewStringResponse(400, "not equal"), nil
// 			}
// 			return httpmock.NewJsonResponse(200, []string{"file0", "extra", "another"})
// 		})

// 	list, err := client.FilesUpload("dirTest", []string{f1.Name()})
// 	assert.NoError(t, err)
// 	assert.ElementsMatch(t, list, []string{"file0", "extra", "another"})
// }
