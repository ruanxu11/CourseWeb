package koala

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type Session struct {
	ID         string
	CookieName string
	Values     map[string]interface{}
	UpdateTime time.Time
	ExpireTime time.Duration
	IsNew      bool
}

var Sessions = make(map[string]Session)

func NewSession() string {
	sessionID := HashString(time.Now().Format(time.UnixDate))
	s := Session{
		ID:         sessionID,
		Values:     make(map[string]interface{}),
		UpdateTime: time.Now(),
		ExpireTime: 3600 * time.Second,
		IsNew:      true,
	}
	Sessions[sessionID] = s
	return sessionID
}

func ExistSession(r *http.Request, cookieName string) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	if _, ok := Sessions[c.Value]; ok {
		return true
	}
	return false
}

func PeekSession(r *http.Request, cookieName string) *Session {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	if s, ok := Sessions[c.Value]; ok {
		s.IsNew = false
		return &s
	}
	return nil
}

func GetSession(r *http.Request, w http.ResponseWriter, cookieName string) *Session {
	c, err := r.Cookie(cookieName)
	log.Println("cookie: ", c)
	if err != nil {
		log.Println("no cookies")
		sessionID := NewSession()
		c = &http.Cookie{
			Name:  cookieName,
			Value: sessionID,
		}
		http.SetCookie(w, c)
		s := Sessions[sessionID]
		return &s
	}
	log.Println("exist cookies")
	if s, ok := Sessions[c.Value]; ok {
		log.Println("exist session")
		s.IsNew = false
		return &s
	} else {
		log.Println("no session")
		sessionID := NewSession()
		c = &http.Cookie{
			Name:  cookieName,
			Value: sessionID,
		}
		http.SetCookie(w, c)
		s := Sessions[sessionID]
		return &s
	}
	return nil
}

func (s *Session) UpdateExpireTime(r *http.Request, w http.ResponseWriter) error {
	c, _ := r.Cookie(s.CookieName)
	c.Expires = time.Now().Add(s.ExpireTime)
	http.SetCookie(w, c)
	s.UpdateTime = time.Now()
	if s.UpdateTime.Add(s.ExpireTime).Unix() < time.Now().Unix() {
		s.Destory()
		return errors.New("Session is expired")
	}
	return nil
}

func (s *Session) Destory() error {
	delete(Sessions, s.ID)
	return nil
}
