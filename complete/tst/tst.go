package tst

import (
	"github.com/zeyadyasser/autocom/complete"
)

type node struct{
	left *node
	right *node
	mid *node
	midByte byte
	isEnd bool
	value interface{}
}

type TSTCompleter struct {
	root *node
}

func (t *TSTCompleter) Set(key string, value interface{}) error {
	curr := t.root
	idx := 0
	bKey := []byte(key)
	for idx < len(bKey) {
		if curr.mid == nil {
			curr.midByte = bKey[idx]
			curr.mid = &node{}
			curr = curr.mid
			idx++
		} else if bKey[idx] < curr.midByte {
			if curr.left == nil {
				curr.left = &node{}
			}
			curr = curr.left
		} else if bKey[idx] > curr.midByte {
			if curr.right == nil {
				curr.right = &node{}
			}
			curr = curr.right
		} else {
			curr = curr.mid
			idx++
		}
	}
	curr.isEnd = true
	curr.value = value

	return nil
}

func (t *TSTCompleter) Remove(key string) error {
	curr := t.root
	idx := 0
	bKey := []byte(key)
	for idx < len(bKey) {
		if curr == nil {
			return nil
		}
		if bKey[idx] == curr.midByte {
			curr = curr.mid
			idx++
		} else if bKey[idx] < curr.midByte {
			curr = curr.left
		} else if bKey[idx] > curr.midByte {
			curr = curr.right
		}
	}
	if curr != nil {
		curr.isEnd = false
	}

	return nil
}

func (t *TSTCompleter) TopN(key string, n int) (complete.Map, error) {
	res := make(complete.Map, n);
	curr := t.root
	idx := 0
	bKey := []byte(key)
	for idx < len(bKey) {
		if curr == nil {
			return res, nil
		}
		if bKey[idx] == curr.midByte {
			curr = curr.mid
			idx++
		} else if bKey[idx] < curr.midByte {
			curr = curr.left
		} else if bKey[idx] > curr.midByte {
			curr = curr.right
		}
	}
	if curr == nil {
		return res, nil
	}
	res = dfs(curr, res, key, n)
	return res, nil
}

func dfs(src *node, res complete.Map, currstr string, n int) complete.Map {
	if len(res) == n {
		return res
	}
	if src.isEnd {
		res[currstr] = src.value
	}

	if src.left != nil {
		res = dfs(src.left, res, currstr, n)
	}
	if src.mid != nil {
		newstr := currstr + string(src.midByte)
		res = dfs(src.mid, res, newstr, n)
	}
	if src.right != nil {
		res = dfs(src.right, res, currstr, n)
	}

	return res
}

func NewTSTCompleter() complete.AutoCompleter {
	return &TSTCompleter{
		root: &node{},
	}
}
