// Package aptly provides a client to access the aptly REST API in go
package aptly

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

// GetClient get resty client used for advanced use cases like testing or special auth
func (c *Client) GetClient() *resty.Client {
	return c.client
}

func (c *Client) get(url string) *resty.Request {
	return c.newRequest(resty.MethodGet, url)
}

func (c *Client) post(url string) *resty.Request {
	return c.newRequest(resty.MethodPost, url)
}

func (c *Client) put(url string) *resty.Request {
	return c.newRequest(resty.MethodPut, url)
}

func (c *Client) delete(url string) *resty.Request {
	return c.newRequest(resty.MethodDelete, url)
}

func (c *Client) newRequest(method string, url string) *resty.Request {
	r := c.client.NewRequest()
	// workaround for pre 1.6.0 servers not sending content-type on errors
	r.ExpectContentType("application/json")
	r.Method = method
	r.URL = url
	return r
}

// send prepared request and get error
func (c *Client) send(req *resty.Request) error {
	res, err := req.Send()

	if err != nil {
		return err
	} else if res.IsSuccess() {
		return nil
	}
	return getError(res)
}

func NewClient(url string) *Client {
	client := new(Client)
	client.client = resty.New()
	client.client.SetBaseURL(url)
	client.client.SetError(APIError{})

	return client
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

// common function to get errors
func getError(response *resty.Response) error {
	e := response.Error()
	if e != nil {
		apiErr := e.(*APIError)
		if apiErr.Valid() {
			return apiErr
		}
	}

	return fmt.Errorf("unexpected response code %v", response.StatusCode())
}
