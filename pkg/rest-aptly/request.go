// Package aptly provides a client to access the aptly REST API in go
package aptly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type request struct {
	PathTemplate string
	Method       string
	QueryParams  url.Values
	PathParams   map[string]string
	Files        map[string]string
	Body         any
}

func initRequest(method string, reqURL string) *request {
	return &request{
		PathTemplate: reqURL,
		Method:       method,
		PathParams:   make(map[string]string),
		Files:        make(map[string]string),
		QueryParams:  url.Values{},
	}
}

// SetQueryParams set query parameters
func (r *request) SetQueryParams(params map[string]string) *request {
	for k, v := range params {
		r.QueryParams.Add(k, v)
	}
	return r
}

// SetPathParam set named path parameter
func (r *request) SetPathParam(name string, value string) *request {
	r.PathParams[name] = value
	return r
}

// SetFiles add file to upload
func (r *request) SetFiles(files map[string]string) *request {
	for k, v := range files {
		r.Files[k] = v
	}
	return r
}

// SetBody set request json body
func (r *request) SetBody(body any) *request {
	r.Body = body
	return r
}

// GetURL get the resolved URL with path and query parameters applied
func (r *request) GetURL(baseURL string) (string, error) {
	path, err := r.GetPath()
	if err != nil {
		return "", err
	}

	query := ""
	// Add Query Params
	if len(r.QueryParams) > 0 {
		query = "?" + r.QueryParams.Encode()
	}

	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[0 : len(baseURL)-1]
	}

	return fmt.Sprintf("%s/%s%s", baseURL, path, query), nil
}

// GetPath resolves all placeholders
func (r *request) GetPath() (string, error) {
	params := make(map[string]string, len(r.PathParams))
	for k, v := range r.PathParams {
		params[k] = url.PathEscape(v)
	}

	type placeholder struct {
		start int
		end   int
	}

	placeholders := make([]placeholder, 0, len(params))
	currStart := -1

	for index, char := range r.PathTemplate {
		switch char {
		case '{':
			if currStart == -1 {
				currStart = index
			} else {
				return "", fmt.Errorf("missing closing bracket at '%s'", r.PathTemplate[currStart:index+1])
			}
		case '}':
			if currStart != -1 {
				placeholders = append(placeholders, placeholder{start: currStart, end: index})
				currStart = -1
			} else {
				return "", fmt.Errorf("closing bracket without opening bracket")
			}
		case '/':
			if currStart != -1 {
				return "", fmt.Errorf("missing closing bracket at '%s'", r.PathTemplate[currStart:index+1])
			}
		}
	}

	if len(placeholders) == 0 {
		if r.PathTemplate[0] != '/' {
			return r.PathTemplate, nil
		} else {
			return r.PathTemplate[1:], nil
		}
	}

	lastIndex := -1
	buf := bytes.Buffer{}

	for _, p := range placeholders {
		// copy preceeding non template part
		if lastIndex+1 != p.start {
			buf.WriteString(r.PathTemplate[lastIndex+1 : p.start])
		}
		key := r.PathTemplate[p.start+1 : p.end]

		fragment, ok := params[key]
		if !ok {
			return "", fmt.Errorf("path parameter '{%s}' not set", key)
		}
		buf.WriteString(fragment)
		lastIndex = p.end
	}

	// copy remainder
	if lastIndex != len(r.PathTemplate) {
		buf.WriteString(r.PathTemplate[lastIndex+1:])
	}

	path := buf.String()
	if path[0] != '/' {
		return path, nil
	} else {
		return path[1:], nil
	}
}

// SetResult set the result object
func (r *request) GetBodyReader() io.Reader {
	if r.Body == nil {
		return nil
	}

	requestByte, _ := json.Marshal(r.Body)
	requestReader := bytes.NewReader(requestByte)
	return requestReader
}

// SetResult set the result object
func (r *request) GetRawRequest(baseUrl string) (*http.Request, error) {
	url, err := r.GetURL(baseUrl)
	if err != nil {
		return nil, err
	}

	contentType := ""
	payload := new(bytes.Buffer)
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
		mpw := multipart.NewWriter(payload)

		for name, path := range r.Files {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			fileWriter, err := mpw.CreateFormFile(name, filepath.Base(path)) // path get filename
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
