package basic

import (
	"sync"

	"github.com/zeyadyasser/autocom/complete"
)

type BasicEngine struct {
	sync.RWMutex
	completer complete.AutoCompleter
}

func (engine *BasicEngine) Set(key string, value interface{}) error {
	engine.Lock()
	defer engine.Unlock()
	return engine.completer.Set(key, value)
}

func (engine *BasicEngine) Remove(key string) error {
	engine.Lock()
	defer engine.Unlock()
	return engine.completer.Remove(key)
}

func (engine *BasicEngine) TopN(key string, n int) (complete.Map, error) {
	engine.RLock()
	defer engine.RUnlock()
	return engine.completer.TopN(key, n)
}

func NewBasicEngine(c complete.AutoCompleter) *BasicEngine {
	return &BasicEngine{
		completer: c,
	}
}