package infrastructure

import (
	"net/http"
	"time"
)

var HttpClientSingleton *http.Client

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	HttpClientSingleton = &http.Client{
		Timeout: time.Minute * 1,
	}
}
