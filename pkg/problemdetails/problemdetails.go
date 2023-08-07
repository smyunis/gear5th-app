package problemdetails

import (
	"fmt"
	"net/http"
)

// reference https://datatracker.ietf.org/doc/html/rfc7807

type ProblemDetails struct {
	Status   int    `json:"status,omitempty"`
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func NewProblemDetails(status int) ProblemDetails {
	return ProblemDetails{
		Status: status,
		Type:   "about:blank",
		Title:  http.StatusText(status),
	}
}

func (p ProblemDetails) Error() string {
	return fmt.Sprintf("%s: %s", p.Title, p.Detail)
}
