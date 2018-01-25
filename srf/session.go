package srf

import (
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

var sessionContext sync.Map

type Session interface {
	GetId() string
	SetAttribute(key string, value interface{})
	GetAttribute(key string) interface{}
	Invalid()
}

type defaultSession struct {
	id         string
	createTime time.Time
	attributes map[string]interface{}
}

func (session *defaultSession) GetId() string {
	return session.id
}

func (session *defaultSession) SetAttribute(key string, value interface{}) {
	session.attributes[key] = value
}

func (session *defaultSession) GetAttribute(key string) interface{} {
	return session.attributes[key]
}

func (session *defaultSession) Invalid() {
	sessionContext.Delete(session.GetId())
}

func CreateSession() Session {
	id := getUuid()
	session := &defaultSession{id: id, createTime: time.Now(), attributes: make(map[string]interface{})}
	sessionContext.Store(id, session)
	return session
}

func GetSession(id string) Session {
	session, _ := sessionContext.Load(id)
	if session == nil {
		return nil
	}
	return session.(Session)
}

func getUuid() string {
	return uuid.NewV4().String()
}
