package aptly

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

type FileInForm struct {
	filename string
	body     []byte
}

func formFiles(expectedFiles map[string]FileInForm) httpmock.Matcher {
	return httpmock.NewMatcher("",
		func(req *http.Request) bool {
			err := req.ParseMultipartForm(1024 * 1024 * 4)
			if err != nil {
				return false
			}

			allEqual := true

			for name, formFile := range expectedFiles {
				file, handler, err := req.FormFile(name)
				if err != nil {
					fmt.Printf("file '%s': %s\n", name, err.Error())
					return false
				}
				defer file.Close()

				if handler.Filename != formFile.filename {
					fmt.Printf("filename '%s' not expected '%s'\n", handler.Filename, formFile.filename)
					return false
				}

				buffer := make([]byte, handler.Size)
				_, err = file.Read(buffer)
				if err != nil {
					fmt.Printf("could not read file '%s' from request: %s\n", name, err.Error())
					return false
				}

				if !bytes.Equal(formFile.body, buffer) {
					fmt.Printf("file '%s' from not expected value\n", name)
					allEqual = false
				}
			}
			return allEqual
		})
}

func TestFilesUpload(t *testing.T) {
	client := clientForTest(t, "http://host.local")

	data1 := "file0 data"

	f1, err := os.CreateTemp("", "file0_")
	assert.NoError(t, err)
	defer os.Remove(f1.Name())
	_, err = f1.WriteString(data1)
	assert.NoError(t, err)

	files := map[string]FileInForm{
		"file0": {
			filename: filepath.Base(f1.Name()),
			body:     []byte("file0 data"),
		},
	}

	httpmock.RegisterMatcherResponder(http.MethodPost, "http://host.local/api/files/dirTest",
		formFiles(files),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, []string{"file0", "extra", "another"})
		})

	list, err := client.FilesUpload("dirTest", []string{f1.Name()})
	assert.NoError(t, err)
	assert.ElementsMatch(t, list, []string{"file0", "extra", "another"})
}
