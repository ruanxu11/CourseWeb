package utilKL

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
	sessionID := hashString(time.Now().Format(time.UnixDate))
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

func GetSession(r *http.Request, w http.ResponseWriter, cookieName string) *Session {
	c, err := r.Cookie(cookieName)
	log.Println(c)
	log.Println(err)
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

// func GetSession(sessionId string, key string) (interface{}, error) {
// 	data, ok := Session[sessionId]
// 	if !ok {
// 		return nil, errors.New("SessionId not found")
// 	}
// 	if data.UpdateTime+SESSION_EXPIRE_TIME < time.Now().Unix() {
// 		delete(Session, sessionId)
// 		return nil, errors.New("Session is expired")
// 	}
// 	value, ok := data.Values[key]
// 	if !ok {
// 		return nil, errors.New("error key")
// 	}
// 	data.UpdateTime = time.Now().Unix()
// 	return value, nil
// }

// func SetSession(sessionId string, key string, value interface{}) error {
// 	data, ok := Session[sessionId]
// 	if !ok {
// 		data = Session{
// 			Values:     make(map[string]interface{}),
// 			UpdateTime: time.Now().Unix(),
// 		}
// 		Session[sessionId] = data
// 	}
// 	data.Values[key] = value
// 	data.UpdateTime = time.Now().Unix()
// 	return nil
// }

// func DestorySession(sessionId string) error {
// 	data, ok := Session[sessionId]
// 	if !ok {
// 		return errors.New("SessionId not found")
// 	}
// 	if data.UpdateTime+SESSION_EXPIRE_TIME < time.Now().Unix() {
// 		delete(Session, sessionId)
// 		return errors.New("Session is expired")
// 	}
// 	delete(Session, sessionId)
// 	return nil
// }

// func getSessionBycookieName(r *http.Request, cookieName string, key string) (interface{}, error) {
// 	c, err := r.Cookie(cookieName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return GetSession(c.Value, key)
// }

// func createSessionBycookieName(r *http.Request, w http.ResponseWriter, cookieName string) (interface{}, error) {
// 	sessionID := NewSession()
// 	c := &http.Cookie{
// 		Name:  cookieName,
// 		Value: sessionID,
// 	}
// 	http.SetCookie(w, c)
// }
