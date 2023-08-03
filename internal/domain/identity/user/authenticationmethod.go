package user

type AuthenticationMethod int

const (
	Managed AuthenticationMethod = iota
	OAuth
)

