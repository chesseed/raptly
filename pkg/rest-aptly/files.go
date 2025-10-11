package aptly

import "fmt"

func (c *Client) FilesListDirs() ([]string, error) {
	req := c.get("api/files")
	return callAPIwithResult[[]string](c, req)
}

func (c *Client) FilesListFiles(dir string) ([]string, error) {
	req := c.get("api/files/{dir}").
		SetPathParam("dir", dir)
	return callAPIwithResult[[]string](c, req)
}

func (c *Client) FilesUpload(dir string, files []string) ([]string, error) {
	fileMap := make(map[string]string)
	for i, file := range files {
		fileMap[fmt.Sprintf("file%d", i)] = file
	}

	req := c.post("api/files/{dir}").
		SetPathParam("dir", dir).
		SetFiles(fileMap)
	return callAPIwithResult[[]string](c, req)
}

func (c *Client) FilesDeleteDir(dir string) error {
	req := c.delete("api/files/{dir}").
		SetPathParam("dir", dir)
	return callAPI(c, req)
}

func (c *Client) FilesDeleteFile(dir string, file string) error {
	req := c.delete("api/files/{dir}/{file}").
		SetPathParam("dir", dir).
		SetPathParam("file", file)
	return callAPI(c, req)
}
