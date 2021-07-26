package trie

import "fmt"

type stringTrieNode struct {
	ChildNodes map[byte]*stringTrieNode
	Match      int
}

type nextStrPos struct {
	Str     string
	NextPos int
}

type StringTrie struct {
	Root *stringTrieNode
}

func newstringTrieNode() *stringTrieNode {
	return &stringTrieNode{
		ChildNodes: make(map[byte]*stringTrieNode),
		Match:      -1,
	}
}

func NewStringTrie(strList []string) (*StringTrie, error) {
	// 初始化
	stringTrie := &StringTrie{
		Root: newstringTrieNode(),
	}
	nodeStrPosMap := make(map[*stringTrieNode][]nextStrPos)
	nodeStrPosMap[stringTrie.Root] = make([]nextStrPos, 0, len(strList))
	strIndexMap := make(map[string]int)
	for index, str := range strList {
		if len(str) == 0 {
			return nil, fmt.Errorf("string with index: %d is null", index)
		}
		strIndexMap[str] = index
		nodeStrPosMap[stringTrie.Root] = append(nodeStrPosMap[stringTrie.Root],
			nextStrPos{Str: str, NextPos: 0})
	}

	// 进行树的深建立
	nextLevelNodes := make([]*stringTrieNode, 0, 1)
	nextLevelNodes = append(nextLevelNodes, stringTrie.Root)
	for len(nextLevelNodes) > 0 {
		newNextLevelNodes := make([]*stringTrieNode, 0)
		for _, node := range nextLevelNodes {
			nextStrPosList := nodeStrPosMap[node]
			// 将所有下一个可能的字符记录
			nextByteMap := make(map[byte][]int)
			for index, value := range nextStrPosList {
				if _, ok := nextByteMap[value.Str[value.NextPos]]; !ok {
					nextByteMap[value.Str[value.NextPos]] = make([]int, 0)
				}
				nextByteMap[value.Str[value.NextPos]] = append(nextByteMap[value.Str[value.NextPos]],
					index)
			}
			// 该节点的下一层子节点的建立
			for b, indexList := range nextByteMap {
				newNode := newstringTrieNode()
				nodeStrPosMap[newNode] = make([]nextStrPos, 0, len(indexList))
				for _, index := range indexList {
					value := nextStrPosList[index]
					if value.NextPos+1 == len(value.Str) {
						newNode.Match = strIndexMap[value.Str]
					} else {
						value.NextPos += 1
						nodeStrPosMap[newNode] = append(nodeStrPosMap[newNode], value)
					}
				}
				node.ChildNodes[b] = newNode
				newNextLevelNodes = append(newNextLevelNodes, newNode)
			}
		}
		nextLevelNodes = newNextLevelNodes
	}

	return stringTrie, nil
}

func (trie *StringTrie) Match(str string) int {
	if trie.Root == nil {
		return -1
	}
	node := trie.Root
	for _, b := range []byte(str) {
		if nextNode, ok := node.ChildNodes[b]; ok {
			node = nextNode
		} else {
			return -1
		}
	}
	return node.Match
}
