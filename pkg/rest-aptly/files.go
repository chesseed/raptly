package aptly

import "fmt"

func (c *Client) FilesListDirs() ([]string, error) {
	var dirs []string

	req := c.get("api/files").
		SetResult(&dirs)

	return dirs, c.send(req)
}

func (c *Client) FilesListFiles(dir string) ([]string, error) {
	var dirs []string

	req := c.get("api/files/{dir}").
		SetPathParam("dir", dir).
		SetResult(&dirs)

	return dirs, c.send(req)
}

func (c *Client) FilesUpload(dir string, files []string) ([]string, error) {
	var uploaded []string

	fileMap := make(map[string]string)
	for i, file := range files {
		fileMap[fmt.Sprintf("file%d", i)] = file
	}

	req := c.post("api/files/{dir}").
		SetPathParam("dir", dir).
		SetResult(&uploaded).
		SetFiles(fileMap)

	return uploaded, c.send(req)
}

func (c *Client) FilesDeleteDir(dir string) error {
	req := c.delete("api/files/{dir}").
		SetPathParam("dir", dir)

	return c.send(req)
}

func (c *Client) FilesDeleteFile(dir string, file string) error {
	req := c.delete("api/files/{dir}/{file}").
		SetPathParam("dir", dir).
		SetPathParam("file", file)

	return c.send(req)
}
