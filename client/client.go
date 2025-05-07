package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type httpClient struct {
	client *http.Client
	req    *http.Request
}

func (c *httpClient) SetHeader(key, value string) {
	c.req.Header.Set(key, value)
}

func (c *httpClient) SendRequest() (*http.Response, error) {
	return c.client.Do(c.req)
}

func newHTTPSClient(url, method string, body io.Reader) (*httpClient, error) {
	req, err := getDefautlRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	return &httpClient{
		client: getClient(),
		req:    req,
	}, nil
}

func getDefautlRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", `text/xml; charset="utf-8"`)
	req.Header.Set("Accept", "text/xml")
	req.Header.Set("Cache-Control", "no-cache")
	return req, nil
}

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(20 * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Renegotiation:      tls.RenegotiateOnceAsClient,
				InsecureSkipVerify: true,
			},
		},
	}
}
