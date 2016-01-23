package models

type User struct {
	Id                    int    `valid:"numeric"`
	Name                  string `valid:"alphanum,required"`
	Mobile_number         string `valid:"alphanum,required"`
	Password              string `valid:"alphanum,required"`
	Password_confirmation string `valid:"alphanum,required"`
	Devise_token          string `valid:"alphanum,required"`
}

type UserDetails struct {
	Id                    int    `valid:"numeric"`
	Name                  string `valid:"alphanum,required"`
	Mobile_number         string `valid:"alphanum,required"`
	Devise_token          string `valid:"alphanum,required"`
}

type SuccessfulSignIn struct {
	Success string
	Message string
	User    UserDetails
	Session SessionDetails
}

type SessionDetails struct {
	SessionId   int
	DeviseToken string
}

// Message struct [controllers/account]
// Common for sign_up, session and password
type Message struct {
	Success string
	Message string
	User    User
}

type ErrorMessage struct {
	Success string
	Error   string
}
