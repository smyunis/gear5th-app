package siteverification_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"gitlab.com/gear5th/gear5th-app/internal/application/testdoubles"
	"gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"
	"gitlab.com/gear5th/gear5th-app/internal/infrastructure/siteverification"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	teardown()
}

type stubHttpClientFunc func(*http.Request) (*http.Response, error)

func (m stubHttpClientFunc) Do(req *http.Request) (*http.Response, error) {
	return m(req)
}

var httpClient stubHttpClientFunc

func setup() {
	httpClient = func(req *http.Request) (*http.Response, error) {
		switch req.URL.String() {
		case "https://www.wikihow.com/ads.txt":
			return getMockFileAsHttpResponse("testdata/wikihowads.txt")

		case "https://www.quora.com/ads.txt":
			return getMockFileAsHttpResponse("testdata/quoraads.txt")

		case "https://www.tutorialspoint.com/ads.txt":
			return getMockFileAsHttpResponse("testdata/tutspointads.txt")

		default:
			return &http.Response{
				Status: http.StatusText(http.StatusNotFound),
				Body:   io.NopCloser(bytes.NewReader([]byte("mock not implemented"))),
			}, nil
		}
	}
}

func teardown() {

}

func getMockFileAsHttpResponse(fileName string) (*http.Response, error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	responseBody := io.NopCloser(bytes.NewReader(buf))

	return &http.Response{
		Status: http.StatusText(http.StatusOK),
		Body:   responseBody,
	}, nil
}

func TestCanCheckIfRecordIsInAdsTxt(t *testing.T) {
	logger := testdoubles.ConsoleLogger{}
	service := siteverification.NewAdsTxtVerificationService(httpClient, logger)
	tuts, _ := url.Parse("https://www.tutorialspoint.com/index.htm")
	s := site.NewSite("2000970", *tuts)
	record := site.AdsTxtRecord{
		AdExchangeDomain: "consumable.com",
		PublisherId:      "2000970",
		Relation:         "DIRECT",
	}

	err := service.VerifyAdsTxt(s, record)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if !s.IsVerified() {
		t.FailNow()
	}
}
