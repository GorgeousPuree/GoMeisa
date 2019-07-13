package session

import "Gomeisa"

type sessionData struct {
	email string
	specialtyId int
}

type Session struct {
	data map[string]*sessionData
}

func NewSession() *Session {
	s := new(Session)
	s.data = make(map[string]*sessionData)
	return s
}

func (s *Session) init(email string) string {
	sessionId := Gomeisa.GenerateString(32)

	data := &sessionData{email:email}
	s.data[sessionId] = data

	return sessionId
}
