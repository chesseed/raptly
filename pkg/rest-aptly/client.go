// Package aptly provides a client to access the aptly REST API in go
package aptly

import (
	"encoding/json"
	"errors"
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

func NewClient(url string) *Client {
	client := new(Client)
	client.client = resty.New()
	client.client.SetBaseURL(url)
	client.client.SetError(APIError{})

	return client
}

type APIError struct {
	Error string `json:"error"`
}

// common function to get errors
func getError(response *resty.Response) error {
	e := response.Error()
	if e != nil {
		return errors.New(e.(*APIError).Error)
	}

	// workaround for pre 1.6.x version not sending json content-type header
	var msg APIError
	me := json.Unmarshal(response.Body(), &msg)
	if me == nil {
		return errors.New(msg.Error)
	}
	return fmt.Errorf("unexpected response code %v", response.StatusCode())
}
