package site

import (
	"errors"
	"fmt"
	"net/url"

)

// type SiteVerificationError struct {
// 	message string
// 	inner   error
// }

// func NewSiteVerificationError(message string, inner error) error {
// 	return SiteVerificationError{
// 		message: message,
// 		inner:   inner,
// 	}
// }

// func (e SiteVerificationError) Error() string {
// 	return e.message
// }

// func (e SiteVerificationError) Unwrap() error {
// 	return e.inner
// }

var ErrSiteVerification = errors.New("site verification error")

type AdsTxtRecord struct {
	AdExchangeDomain string
	PublisherId      string
	Relation         string
	CertAuthTag      string
}

func (a AdsTxtRecord) String() string {
	record := fmt.Sprintf("%s, %s, %s", a.AdExchangeDomain, a.PublisherId, a.Relation)

	if a.CertAuthTag != "" {
		record = fmt.Sprintf("%s, %s", record, a.CertAuthTag)
	}
	return record
}

type AdsTxtVerificationService interface {
	VerifyAdsTxt(site Site, record AdsTxtRecord) error
}

func VerifySiteHostname(site *Site, source url.URL) error {
	siteFqdn := site.url.Hostname()
	sourceFqdn := source.Hostname()

	if siteFqdn != sourceFqdn {
		return ErrSiteVerification
	}

	site.Verify()
	return nil
}



func GetAdsTxtRecord(s Site) AdsTxtRecord {
	record := AdsTxtRecord{
		AdExchangeDomain: "gear5th.com",
		PublisherId:      s.publisherId.String(),
		Relation:         "DIRECT",
	}
	return record
}
