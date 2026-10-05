package sensitive

import (
	"strings"
	"unicode"
)

// TrieNode 字典树节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// Filter 敏感词过滤器
type Filter struct {
	root *TrieNode
}

// NewFilter 创建敏感词过滤器
func NewFilter(words []string) *Filter {
	root := &TrieNode{children: make(map[rune]*TrieNode)}
	f := &Filter{root: root}
	for _, v := range words {
		f.AddWord(v)
	}
	return f
}

// AddWord 添加敏感词
func (f *Filter) AddWord(word string) {
	word = strings.TrimSpace(word)
	if len(word) == 0 {
		return
	}
	curr := f.root
	for _, r := range word {
		r = unicode.ToLower(r)
		if _, ok := curr.children[r]; !ok {
			curr.children[r] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		curr = curr.children[r]
	}
	curr.isEnd = true
}

// isSkipRune 是否跳过干扰符号（空格、标点）
func isSkipRune(r rune) bool {
	return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// IsSensitive 判断文本是否包含敏感词
func (f *Filter) IsSensitive(text string) bool {
	runes := []rune(text)
	n := len(runes)

	for i := 0; i < n; i++ {
		curr := f.root
		for j := i; j < n; j++ {
			r := unicode.ToLower(runes[j])
			if isSkipRune(r) {
				continue
			}
			next, ok := curr.children[r]
			if !ok {
				break
			}
			if next.isEnd {
				return true
			}
			curr = next
		}
	}
	return false
}

// Replace 将文本中的敏感词替换为指定字符
func (f *Filter) Replace(text string, replaceChar rune) string {
	runes := []rune(text)
	n := len(runes)
	res := make([]rune, n)
	copy(res, runes)

	for i := 0; i < n; i++ {
		curr := f.root
		matchLen := 0
		for j := i; j < n; j++ {
			r := unicode.ToLower(runes[j])
			if isSkipRune(r) {
				continue
			}
			next, ok := curr.children[r]
			if !ok {
				break
			}
			if next.isEnd {
				matchLen = j - i + 1
			}
			curr = next
		}
		if matchLen > 0 {
			for k := i; k < i+matchLen; k++ {
				res[k] = replaceChar
			}
			i += matchLen - 1
		}
	}
	return string(res)
}
