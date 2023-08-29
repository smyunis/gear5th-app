package testdoubles

import "gitlab.com/gear5th/gear5th-app/internal/domain/publisher/site"

type AdsTxtSiteVerificaitonStub struct{}

// VerifyAdsTxt implements site.AdsTxtVerificationService.
func (AdsTxtSiteVerificaitonStub) VerifyAdsTxt(site site.Site, record site.AdsTxtRecord) error {
	return nil
}

