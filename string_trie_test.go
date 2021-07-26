package main

import (
	"testing"

	"github.com/syple000/compiler/trie"
)

func TestStringTrie(t *testing.T) {
	tr, _ := trie.NewStringTrie([]string{"abc", "abcd", "ac", "bc"})
	if tr.Match("a") != -1 {
		t.Errorf("a is not in trie but found")
	}
	if tr.Match("b") != -1 {
		t.Errorf("b is not in trie but found")
	}
	if tr.Match("abc") != 0 {
		t.Errorf("abc is in trie but not found")
	}
	if tr.Match("abcd") != 1 {
		t.Errorf("abcd is in trie but not found")
	}
	if tr.Match("abcde") != -1 {
		t.Errorf("abcde is not in trie but found")
	}
	if tr.Match("bc") != 3 {
		t.Errorf("bc is in trie but not found")
	}
	if tr.Match("bcd") != -1 {
		t.Errorf("bcd is not in trie but found")
	}
}
