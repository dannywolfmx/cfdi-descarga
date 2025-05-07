package client

import (
	"io"
	"net/http"
)

// Request struct use to
type Request struct {
	Body    io.Reader
	URL     string
	Method  string
	Headers map[string]string
}

// Send will perform a request to the server
func (req *Request) Send() (*http.Response, error) {
	client, err := newHTTPSClient(req.URL, req.Method, req.Body)

	if err != nil {
		return nil, err
	}

	//Set headers to the request
	for key, value := range req.Headers {
		client.SetHeader(key, value)
	}

	return client.SendRequest()
}

// RequestData represent a basig function to generate a body request
type RequestData interface {
	SendRequest() ([]byte, error)
	GetRequest() (interface{}, error)
}
