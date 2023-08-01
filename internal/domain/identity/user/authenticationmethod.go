package user

type AuthenticationMethod int

const (
	Email AuthenticationMethod = iota
	GoogleOAuth
)

