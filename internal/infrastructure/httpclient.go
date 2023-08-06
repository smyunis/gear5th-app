package infrastructure

import (
	"net/http"
	"time"
)

var HTTPClientSingleton *http.Client

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	HTTPClientSingleton = &http.Client{
		Timeout: time.Minute * 1,
	}
}
