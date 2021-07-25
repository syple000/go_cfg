package main

import (
	"testing"

	"github.com/syple000/compiler/trie"
)

func TestStringTrie(t *testing.T) {
	tr, _ := trie.NewStringTrie([]string{"abc", "abcd", "ac", "bc"})
	if tr.IsIn("a") {
		t.Errorf("a is not in trie but found")
	}
	if tr.IsIn("b") {
		t.Errorf("b is not in trie but found")
	}
	if !tr.IsIn("abc") {
		t.Errorf("abc is in trie but not found")
	}
	if !tr.IsIn("abcd") {
		t.Errorf("abcd is in trie but not found")
	}
	if tr.IsIn("abcde") {
		t.Errorf("abcde is not in trie but found")
	}
	if !tr.IsIn("bc") {
		t.Errorf("bc is in trie but not found")
	}
	if tr.IsIn("bcd") {
		t.Errorf("bcd is not in trie but found")
	}
}
