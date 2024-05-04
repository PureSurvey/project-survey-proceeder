package pools

import (
	"github.com/mssola/useragent"
	"sync"
)

type UserAgentPool struct {
	pool sync.Pool
}

func (uap *UserAgentPool) Get() *useragent.UserAgent {
	v := uap.pool.Get()
	if v == nil {
		return &useragent.UserAgent{}
	}
	return v.(*useragent.UserAgent)
}

func (uap *UserAgentPool) Put(ua *useragent.UserAgent) {
	uap.pool.Put(ua)
}
