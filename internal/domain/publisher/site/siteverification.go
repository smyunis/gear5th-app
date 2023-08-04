package site

import (
	"net/url"
)

type SiteVerificationMethod int

const (
	_ SiteVerificationMethod = iota
	AdsTxt
	SiteScript
)

type SiteVerificationError struct {
	message string
	inner   error
}

func NewSiteVerificationError(message string, inner error) error {
	return SiteVerificationError{
		message: message,
		inner:   inner,
	}
}

func (e SiteVerificationError) Error() string {
	return e.message
}

func (e SiteVerificationError) Unwrap() error {
	return e.inner
}

type AdsTxtVerificationService interface {
	VerifyAdsTxt(site *Site) error
}

func VerifySiteHostname(site *Site, source url.URL) error {
	siteFqdn := site.url.Hostname()
	sourceFqdn := source.Hostname()

	if siteFqdn != sourceFqdn {
		return NewSiteVerificationError("site hostname mismatch", nil)
	}

	site.Verify()
	return nil
}
