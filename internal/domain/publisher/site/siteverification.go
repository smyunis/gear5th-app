package site

import (
	"fmt"
	"net/url"
)

type SiteVerificationMethod int

const (
	_ SiteVerificationMethod = iota
	AdsTxt
	SiteScript
)

type SiteVerificationError struct {
	Reason string
	Inner error
}

func (e SiteVerificationError) Error() string {
	return fmt.Sprintf("%s,%s", e.Reason,e.Inner.Error())
}

type AdsTxtVerificationService interface {
	VerifyAdsTxt(site *Site) error
}

func VerifySiteHostname(site *Site, source url.URL) error {
	siteFqdn := site.url.Hostname()
	sourceFqdn := source.Hostname()

	if siteFqdn != sourceFqdn {
		return SiteVerificationError{Reason: "site hostname mismatch"}
	}

	site.Verify()
	return nil
}


