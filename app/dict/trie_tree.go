package dict

import (
	"strings"
)

// Trie 短语组成的Trie树.
type Trie struct {
	Root *Node
}

// Node Trie树上的一个节点.
type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
}

// newTrie 新建一棵Trie
func newTrie() *Trie {
	return &Trie{
		Root: newRootNode(0),
	}
}

// Add 添加若干个词
func (tree *Trie) Add(words ...string) {
	for _, word := range words {
		tree.add(word)
	}
}

func (tree *Trie) add(word string) {
	var current = tree.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := newNode(r)
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

// validate 验证字符串是否合法，如不合法则返回false和检测到
// 的第一个敏感词
func (tree *Trie) validate(text string) (bool, string) {
	const (
		Empty = ""
	)
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			return false, string(runes[left : position+1])
		}

		parent = current
	}

	return true, Empty
}

// searchLine 按行检索，返回行
func (tree *Trie) searchLine(text string) string {
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		found   bool
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]
		if !found {
			break
		}
		if current.IsPathEnd() && position == length-1 {
			break
		}
		parent = current
	}
	var lineBuilder strings.Builder
	for {
		lineBuilder.WriteRune(parent.Character)
		if parent.IsPathEnd() {
			break
		}
		for _, child := range parent.Children {
			parent = child
		}
	}
	return lineBuilder.String()
}

// findAll 找有所有包含在词库中的词
func (tree *Trie) findAll(text string) []string {
	var matches []string
	var (
		parent  = tree.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < length; position++ {
		current, found = parent.Children[runes[position]]

		if !found {
			parent = tree.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			matches = append(matches, string(runes[left:position+1]))
		}

		if position == length-1 {
			parent = tree.Root
			position = left
			left++
			continue
		}

		parent = current
	}

	var i = 0
	if count := len(matches); count > 0 {
		set := make(map[string]struct{})
		for i < count {
			_, ok := set[matches[i]]
			if !ok {
				set[matches[i]] = struct{}{}
				i++
				continue
			}
			count--
			copy(matches[i:], matches[i+1:])
		}
		return matches[:count]
	}

	return nil
}

// newNode 新建子节点
func newNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

// newRootNode 新建根节点
func newRootNode(character rune) *Node {
	return &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node, 0),
	}
}

// isLeafNode 判断是否叶子节点
func (node *Node) isLeafNode() bool {
	return len(node.Children) == 0
}

// IsRootNode 判断是否为根节点
func (node *Node) IsRootNode() bool {
	return node.isRootNode
}

// IsPathEnd 判断是否为某个路径的结束
func (node *Node) IsPathEnd() bool {
	return node.isPathEnd
}

// softDel 置软删除状态
func (node *Node) softDel() {
	node.isPathEnd = false
}
