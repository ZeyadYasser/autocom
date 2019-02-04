package skip

import (
	"strings"
	"sync"
	"log"

	"github.com/zeyadyasser/autocom/complete"
	"github.com/zeyadyasser/autocom/complete/tst"
)

type completerFactory func() complete.AutoCompleter

type result struct {
	Key string
	Value interface{}
}

type Options struct {
	MaxLevels int
	SkipBegin bool
	ToLower bool
}

type SkipEngine struct {
	sync.RWMutex
	opts Options
	levels []complete.AutoCompleter
}

func (E *SkipEngine) Set(key string, value interface{}) error {
	E.Lock()
	defer E.Unlock()

	normalizedKey := key
	if E.opts.ToLower {
		normalizedKey = strings.ToLower(normalizedKey)
	}

	// All variations of the key point to this result
	res := &result{
		Key: key, 			// Keep Original key without normalization
		Value: value,
	}

	E.levels[0].Set(normalizedKey, res)
	
	words, size := E.splitKey(normalizedKey)
	if E.opts.SkipBegin {
		for i := 1; i < size; i++ {
			err := E.levels[i].Set(
				strings.Join(words[i:]," "),
				res,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (E *SkipEngine) Remove(key string) error {
	E.Lock()
	defer E.Unlock()

	if E.opts.ToLower {
		key = strings.ToLower(key)
	}

	err := E.levels[0].Remove(key)
	if err != nil {
		return err
	}

	words, size := E.splitKey(key)
	if E.opts.SkipBegin {
		for i := 1; i < size; i++ {
			err := E.levels[i].Remove(
				strings.Join(words[i:]," "),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (E *SkipEngine) TopN(key string, n int) (complete.Map, error) {
	E.RLock()
	defer E.RUnlock()

	if E.opts.ToLower {
		key = strings.ToLower(key)
	}

	topN := make(complete.Map, n)

	res, err := E.levels[0].TopN(key, n)
	if err != nil {
		return topN, err
	}
	for _, v := range res {
		topN[v.(*result).Key] = v.(*result).Value
	}

	if E.opts.SkipBegin {
		for i := 1; i < E.opts.MaxLevels; i++ {
			if len(topN) == n {
				break
			}
			res, err := E.levels[i].TopN(key, n - len(topN))
			if err != nil {
				return topN, err
			}
			for _, v := range res {
				topN[v.(*result).Key] = v.(*result).Value
			}
		}
	}

	return topN, nil
}

func (E *SkipEngine) splitKey(key string) (words []string, size int) {
	words = strings.Fields(key)
	switch {
	case E.opts.MaxLevels > len(words):
		size = len(words)
	default:
		size = E.opts.MaxLevels
	}
	return
}

func NewSkipEngine(opts Options, factory completerFactory) *SkipEngine{
	if factory == nil {
		factory = tst.NewTSTCompleter // Use default factory
	}

	if opts.MaxLevels == 0 {
		log.Fatalf("MaxLevels should be at least 1, zero provieded.")
	}

	engine := &SkipEngine{
		opts: opts,
	}

	engine.levels = make([]complete.AutoCompleter, opts.MaxLevels)
	for idx := 0; idx < opts.MaxLevels; idx++ {
		engine.levels[idx] = factory()
	}
	return engine
}
