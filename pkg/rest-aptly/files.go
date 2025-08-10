package aptly

import "fmt"

func (c *Client) FilesListDirs() ([]string, error) {
	var dirs []string

	resp, err := c.client.R().
		SetResult(&dirs).
		Get("api/files")

	if err != nil {
		return dirs, err
	} else if resp.IsSuccess() {
		return dirs, nil
	}
	return dirs, getError(resp)
}

func (c *Client) FilesListFiles(dir string) ([]string, error) {
	var dirs []string

	resp, err := c.client.R().
		SetResult(&dirs).
		SetPathParam("dir", dir).
		Get("api/files/{dir}")

	if err != nil {
		return dirs, err
	} else if resp.IsSuccess() {
		return dirs, nil
	}
	return dirs, getError(resp)
}

func (c *Client) FilesUpload(dir string, files []string) ([]string, error) {
	var uploaded []string

	fileMap := make(map[string]string)
	for i, file := range files {
		fileMap[fmt.Sprintf("file%d", i)] = file
	}

	resp, err := c.client.R().
		SetResult(&uploaded).
		SetFiles(fileMap).
		SetPathParam("dir", dir).
		Post("api/files/{dir}")

	if err != nil {
		return uploaded, err
	} else if resp.IsSuccess() {
		return uploaded, nil
	}
	return uploaded, getError(resp)
}

func (c *Client) FilesDeleteDir(dir string) error {
	resp, err := c.client.R().
		SetPathParam("dir", dir).
		Delete("api/files/{dir}")

	if err != nil {
		return err
	} else if resp.IsSuccess() {
		return nil
	}
	return getError(resp)
}

func (c *Client) FilesDeleteFile(dir string, file string) error {
	resp, err := c.client.R().
		SetPathParam("dir", dir).
		SetPathParam("file", file).
		Delete("api/files/{dir}/{file}")

	if err != nil {
		return err
	} else if resp.IsSuccess() {
		return nil
	}
	return getError(resp)
}
