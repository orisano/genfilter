package main

import (
	"encoding/gob"
	"os"

	"github.com/anknown/darts"
	"golang.org/x/xerrors"
)

const (
	FailState = -1
	RootState = 1
)

type Filter struct {
	Trie    *godarts.DoubleArrayTrie
	Failure []int
	Output  map[int]struct{}
}

func buildFilter(paths [][]rune) (*Filter, error) {
	var d godarts.Darts
	dat, llt, err := d.Build(paths)
	if err != nil {
		return nil, xerrors.Errorf("build double array: %w", err)
	}

	output := make(map[int]struct{}, len(d.Output))
	for state := range d.Output {
		output[state] = struct{}{}
	}

	failure := make([]int, len(dat.Base))
	for _, c := range llt.Root.Children {
		failure[c.Base] = godarts.ROOT_NODE_BASE
	}

	m := &Filter{
		Trie:    dat,
		Failure: failure,
		Output:  output,
	}

	queue := make([]*godarts.LinkedListTrieNode, len(llt.Root.Children))
	copy(queue, llt.Root.Children)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, n := range node.Children {
			if n.Base == godarts.END_NODE_BASE {
				continue
			}
			input := n.Code - godarts.ROOT_NODE_BASE
			outState := FailState
			for inState := node.Base; outState == FailState; {
				inState = m.Failure[inState]
				outState = m.g(inState, input)
			}
			if _, ok := m.Output[outState]; ok {
				m.Output[n.Base] = struct{}{}
			}
			m.Failure[n.Base] = outState
		}
		queue = append(queue, node.Children...)
	}
	return m, nil
}

func loadFilter(path string) (*Filter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, xerrors.Errorf("open filter: %w", err)
	}
	defer f.Close()

	var m Filter
	if err := gob.NewDecoder(f).Decode(&m); err != nil {
		return nil, xerrors.Errorf("decode filter: %w", err)
	}
	return &m, nil
}

func (m *Filter) g(inState int, input rune) int {
	if inState == FailState {
		return RootState
	}
	t := inState + int(input) + godarts.ROOT_NODE_BASE
	if t < len(m.Trie.Base) && inState == m.Trie.Check[t] {
		return m.Trie.Base[t]
	}
	if inState == RootState {
		return RootState
	}
	return FailState
}

func (m *Filter) Contains(r []rune) bool {
	state := RootState
	for _, c := range r {
		for {
			next := m.g(state, c)
			if next != FailState {
				state = next
				break
			}
			state = m.Failure[state]
		}
		if _, ok := m.Output[state]; ok {
			return true
		}
	}
	return false
}
