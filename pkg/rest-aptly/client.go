// Package aptly provides a client to access the aptly REST API in go
package aptly

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Client struct {
	BaseURL string

	client *http.Client
}

// GetClient get http client used for advanced use cases like testing
func (c *Client) GetClient() *http.Client {
	return c.client
}

func (c *Client) get(url string) *request {
	return c.newRequest(http.MethodGet, url)
}

func (c *Client) post(url string) *request {
	return c.newRequest(http.MethodPost, url)
}

func (c *Client) put(url string) *request {
	return c.newRequest(http.MethodPut, url)
}

func (c *Client) delete(url string) *request {
	return c.newRequest(http.MethodDelete, url)
}

func (c *Client) SetTLSClientConfig(config *tls.Config) *Client {
	// TODO
	return c
}

func (c *Client) SetBasicAuth(user string, password string) *Client {
	// TODO
	return c
}

func (c *Client) newRequest(method string, url string) *request {
	r := initRequest(method, url, c.client)
	// // workaround for pre 1.6.0 servers not sending content-type on errors
	// r.ExpectContentType("application/json")
	// r.Method = method
	// r.URL = url
	return r
}

// send prepared request and get error
func (c *Client) Send(r *request) (*http.Response, error) {

	req, err := makeRequest(c, r)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func makeRequest(c *Client, r *request) (*http.Request, error) {
	url, err := r.GetURL(c.BaseURL)
	if err != nil {
		return nil, err
	}

	contentType := ""
	var payload io.Reader = nil
	if r.Body != nil {
		// send JSON body
		b, err := json.Marshal(r.Body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewBuffer(b)
		contentType = "application/json"
	} else if len(r.Files) > 0 {
		// send files body
		buf := &bytes.Buffer{}
		mpw := multipart.NewWriter(buf)

		for name, path := range r.Files {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			fileWriter, err := mpw.CreateFormFile("upload", name)
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(fileWriter, f)

			if err != nil {
				return nil, err
			}
		}
		err = mpw.Close()
		if err != nil {
			return nil, err
		}
		contentType = mpw.FormDataContentType()
	}

	req, err := http.NewRequest(r.Method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func NewClient(url string) *Client {
	client := new(Client)
	client.BaseURL = url
	client.client = &http.Client{}
	return client
}

func checkResponseForError(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}
	var apiErr APIError
	decodeErr := json.NewDecoder(res.Body).Decode(&apiErr)
	if decodeErr != nil {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return &apiErr

}

func callAPIWithResponse(c *Client, r *request) (*http.Response, error) {
	res, err := c.Send(r)
	if err != nil {
		return nil, err
	}
	err = checkResponseForError(res)
	if err != nil {
		res.Body.Close()
		return nil, err
	}
	return res, nil
}

func callAPI(c *Client, r *request) error {
	res, err := c.Send(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = checkResponseForError(res)
	if err != nil {
		return err
	}
	return nil
}

// T can be any type
func callAPIwithResult[T any](c *Client, r *request) (T, error) {
	var ret T

	res, err := c.Send(r)
	if err != nil {
		return ret, err
	}
	defer res.Body.Close()
	err = checkResponseForError(res)
	if err != nil {
		return ret, err
	}

	decodeErr := json.NewDecoder(res.Body).Decode(&ret)
	if decodeErr != nil {
		return ret, err
	}
	return ret, nil
}

type APIError struct {
	// as pointer to distinguish between valid error and failed parsing errors (like empty bodies)
	ErrorMsg *string `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	if e.Valid() {
		return *e.ErrorMsg
	}
	return ""
}

func (e *APIError) Valid() bool {
	return e.ErrorMsg != nil
}
