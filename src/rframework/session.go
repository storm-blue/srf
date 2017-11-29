package rframework

import (
    "time"
)

var sessionContext = make(map[string]Session)

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
    delete(sessionContext, session.GetId())
}
