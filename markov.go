package markov

import (
	"bufio"
	"io"
	"math/rand"
	"unicode"
)

type Node struct {
	key   string
	nodes []*Node
}

type TextGenerator struct {
	rand  *rand.Rand
	nodes []Node
}

func (tg *TextGenerator) Init(seed int64) {
	tg.rand = rand.New(rand.NewSource(seed))
	tg.nodes = make([]Node, 0, 10)
}

func (tg *TextGenerator) findNode(key string) *Node {
	for i, _ := range tg.nodes {
		if tg.nodes[i].key == key {
			return &tg.nodes[i]
		}
	}
	return nil
}

func (tg *TextGenerator) appendNode(key, value string) {
	v := tg.findNode(value)
	if v == nil {
		var n Node
		n.key = value
		n.nodes = []*Node{}
		tg.nodes = append(tg.nodes, n)
		v = &tg.nodes[len(tg.nodes)-1]
	}
	k := tg.findNode(key)
	if k == nil {
		var n Node
		n.key = key
		n.nodes = []*Node{v}
		tg.nodes = append(tg.nodes, n)
	} else {
		k.nodes = append(k.nodes, v)
	}
}

func (tg *TextGenerator) Feed(in io.Reader) {
	type MatchFun func(rune) bool
	br := bufio.NewReader(in)
	read := func() rune {
		r, _, err := br.ReadRune()
		if err != nil {
			return -1
		}
		return r
	}
	accept := func(r rune) MatchFun {
		for _, mf := range []MatchFun{
			unicode.IsLetter,
			unicode.IsDigit,
			unicode.IsPunct,
		} {
			if mf(r) {
				return mf
			}
		}
		return nil
	}
	last := ""
	c := read()
next:
	for {
		mf := accept(c)
		for mf == nil && c != -1 {
			c = read()
			mf = accept(c)
		}
		if c == -1 {
			break
		}
		val := string(c)
		for {
			c = read()
			if c == -1 {
				break
			}
			if !mf(c) {
				if last != "" {
					tg.appendNode(last, val)
				}
				last = val
				continue next
			}
			val += string(c)
		}
		if val != "" {
			tg.appendNode(last, val)
			last = val
		}
	}
}

func (tg *TextGenerator) Generate(words uint) []string {
	var node *Node
	for {
		node = &tg.nodes[tg.rand.Int31n(int32(len(tg.nodes)))]
		first := []rune(node.key)[0]
		if unicode.IsLetter(first) || unicode.IsDigit(first) {
			break
		}
	}
	chain := make([]string, 0, words)
	chain = append(chain, node.key)
	var i uint
	for i = 1; i < words; i++ {
		if len(node.nodes) == 0 {
			break
		}
		if len(node.nodes) == 1 {
			node = node.nodes[0]
		} else {
			node = node.nodes[tg.rand.Int31n(int32(len(node.nodes)))]
		}
		chain = append(chain, node.key)
	}
	return chain
}
