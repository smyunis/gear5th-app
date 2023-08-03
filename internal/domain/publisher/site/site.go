package site

import (
	"net/url"
	"time"

	"gitlab.com/gear5th/gear5th-api/internal/domain/shared"
)

type Site struct {
	id           shared.Id
	publisherId  shared.Id
	url          url.URL
	isVerified   bool
	lastVerified time.Time
}

func NewSite(publisherId shared.Id, url url.URL) Site {
	return Site{
		id:          shared.NewId(),
		publisherId: publisherId,
		url:         url,
	}
}

func (s *Site) Verify() {
	s.isVerified = true
	s.lastVerified = time.Now()
}

func (s *Site) IsVerified() bool {
	return s.isVerified
}

func (s *Site) PublisherId() shared.Id {
	return s.publisherId
}

func (s *Site) Url() url.URL {
	return s.url
}
