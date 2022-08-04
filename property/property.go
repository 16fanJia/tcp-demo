package property

import (
	"sync"
)

type Property struct {
	//锁
	mu sync.RWMutex
	//连接的属性值
	property map[string]interface{}
}

var PropertyInstance *Property

func init() {
	PropertyInstance = &Property{property: map[string]interface{}{}}
}

func (p *Property) SetProperty(key string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.property[key] = value
}

func (p *Property) GetProperty(key string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if val, ok := p.property[key]; ok {
		return val
	} else {
		return nil
	}
}

func (p *Property) RemoveProperty(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.property, key)
}
