package goweb

import (
	"net/http"
	"time"
)

const sessionKey = "sid"

type SessionIdType string

type sessionManager struct {
	sessionMaxAge int // second
	sessions      map[SessionIdType]*session
}

func newSessionManager(sessionTimeout int) *sessionManager {
	return &sessionManager{
		sessionMaxAge: sessionTimeout,
		sessions:      make(map[SessionIdType]*session),
	}
}

func (sm *sessionManager) addSession(key SessionIdType, s *session) {
	sm.sessions[key] = s
}

func (sm *sessionManager) session(key SessionIdType) *session {
	return sm.sessions[key]
}

func (sm *sessionManager) getOrCreateSession(ctx *Context) *session {
	var sessionId SessionIdType
	if c := ctx.Cookie(sessionKey); c == nil {
		return newSession(ctx, sm)
	} else {
		sessionId = SessionIdType(c.Value)
	}

	s := sm.session(sessionId)
	if s == nil {
		s = newSession(ctx, sm)
	}
	return s
}

type session struct {
	id         SessionIdType
	data       map[interface{}]interface{}
	createTime time.Time
}

func newSession(ctx *Context, sm *sessionManager) *session {
	sessionId := genSessionId()
	s := &session{
		id:   sessionId,
		data: make(map[interface{}]interface{}),
	}
	cookie := &http.Cookie{
		Name:     sessionKey,
		Value:    string(sessionId),
		HttpOnly: true,
	}
	ctx.SetRawCookie(cookie)
	sm.addSession(sessionId, s)
	return s
}

func genSessionId() SessionIdType {
	return ""
}

func (s *session) Set(key, value interface{}) {
	s.data[key] = value
}

func (s *session) Get(key interface{}) (interface{}, bool) {
	value, fd := s.data[key]
	return value, fd
}
