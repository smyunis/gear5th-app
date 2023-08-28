package infrastructure

import (
	"net/http"
	"time"
)

var httpClientSingleton *http.Client

func init() {
	httpClientSingleton = &http.Client{
		Timeout: time.Minute * 1,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AppHTTPClient struct{}

func (AppHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return httpClientSingleton.Do(req)
}

func NewAppHTTPClient() AppHTTPClient {
	return AppHTTPClient{}
}
