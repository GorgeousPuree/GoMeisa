package usession

import "github.com/gorilla/sessions"

var Store *sessions.CookieStore

type UserSession struct {
	Email         string
	Authenticated bool
}
