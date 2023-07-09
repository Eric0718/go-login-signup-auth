package config

import "github.com/gorilla/sessions"

const SESSION_ID = "go-login-signup-auth-session"

var Store = sessions.NewCookieStore([]byte("mhd123ian456rush"))
// that's a random number


