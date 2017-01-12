package goweb

import (
	"bytes"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// a-z A-Z 0-9 in random order
const words = "MaSHgPl3jh64EzrpqLstV5XOnF1BcTufDiIWKCbyZwRQ98JGkUNx0A2m7veYod"

//const words = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//const words = "0123456789"

var random = rand.New(rand.NewSource(time.Now().Unix() + 987653))

type SessionIdType string

type sessionManagerConfig struct {
	sessionMaxAge int // second
	sessionKey    string
	sessionIdLen  int
	cleanInterval time.Duration
}

type sessionManager struct {
	*sessionManagerConfig
	sync.RWMutex
	sessions    map[SessionIdType]*session
	cleanTicker *time.Ticker
}

func newSessionManager(config *sessionManagerConfig) *sessionManager {
	sm := &sessionManager{
		sessionManagerConfig: config,
		sessions:             make(map[SessionIdType]*session),
		cleanTicker:          time.NewTicker(config.cleanInterval),
	}
	sm.goRunClean()
	return sm
}

// delete expired session
func (sm *sessionManager) goRunClean() {
	go func() {
		for {
			select {
			case t, ok := <-sm.cleanTicker.C:
				if !ok {
					return
				}
				sm.RLock()
				for sid, s := range sm.sessions {
					if int(t.Sub(s.lastAccessTime)/time.Second) >= sm.sessionMaxAge {
						sm.RUnlock()
						sm.Lock()
						delete(sm.sessions, sid)
						sm.Unlock()
						sm.RLock()
					}
				}
				sm.RUnlock()
			}
		}
	}()
}

func (sm *sessionManager) Close() {
	sm.cleanTicker.Stop()
}

func (sm *sessionManager) session(ctx *Context) *session {
	sm.Lock()
	defer sm.Unlock()
	var sessionId SessionIdType
	if c := ctx.Cookie(sm.sessionKey); c == nil {
		return newSession(ctx, sm)
	} else {
		sessionId = SessionIdType(c.Value)
	}
	s := sm.sessions[sessionId]
	if s == nil {
		s = newSession(ctx, sm)
	}
	s.lastAccessTime = time.Now()
	return s
}

// get reversed len(words) scale string
func sessionIdEncodeInt64(buf *bytes.Buffer, i int64) {
	if i == 0 {
		buf.WriteByte(words[0])
		return
	}
	wl := int64(len(words))
	for i > 0 {
		buf.WriteByte(words[i%wl])
		i /= wl
	}
}

type session struct {
	sync.RWMutex
	id             SessionIdType
	data           map[interface{}]interface{}
	lastAccessTime time.Time
}

func newSession(ctx *Context, sm *sessionManager) *session {
	sessionId := GenSessionId(sm.sessionIdLen)
	s := &session{
		id:             sessionId,
		data:           make(map[interface{}]interface{}),
		lastAccessTime: time.Now(),
	}
	cookie := &http.Cookie{
		Name:     sm.sessionKey,
		Value:    string(sessionId),
		HttpOnly: true,
	}
	ctx.SetRawCookie(cookie)
	sm.sessions[sessionId] = s
	return s
}

func GenSessionId(length int) SessionIdType {
	var buf bytes.Buffer
	sessionIdEncodeInt64(&buf, time.Now().UnixNano())
	if buf.Len() >= length {
		return SessionIdType(buf.Bytes()[:length])
	}
	wl := len(words)
	for i := 0; i < length-buf.Len(); i++ {
		buf.WriteByte(words[random.Intn(wl)])
	}
	return SessionIdType(buf.String())
}

func (s *session) Set(key, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *session) Get(key interface{}) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	value, fd := s.data[key]
	return value, fd
}
