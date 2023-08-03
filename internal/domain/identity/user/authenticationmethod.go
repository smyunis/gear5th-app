package user

type AuthenticationMethod int

const (
	_ AuthenticationMethod = iota
	Managed 
	OAuth
)

