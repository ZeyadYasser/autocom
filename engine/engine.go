package engine

import (
	"github.com/zeyadyasser/autocom/complete"
)

type Engine interface {
	Set(string, interface{}) error
	Remove(string) error
	TopN(string, int) (complete.Map, error)
}