package testdoubles

import "gitlab.com/gear5th/gear5th-api/internal/domain/shared"

type JwtAccessTokenGeneratorStub struct{}

func (j JwtAccessTokenGeneratorStub) Generate(id shared.ID) (string, error) {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkuZ2VhcjV0aC5jb20iLCJzdWIiOiJzdHViLWlkLXh4eCIsImF1ZCI6WyJhcGkuZ2VhcjV0aC5jb20iXSwiZXhwIjoxNjkzOTEwMjU4LCJpYXQiOjE2OTEzMTgyNTh9.6Aj-eIuYbU_06YN4vRDhb9zTWSB2nyxHKnz9Y2NH3uo", nil

}
