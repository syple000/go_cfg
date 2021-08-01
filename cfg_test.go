package main

import (
	"fmt"
	"testing"

	"github.com/syple000/compiler/cfg"
)

type EchoAnalyzer struct{}

func (analyzer *EchoAnalyzer) Moveon(matcher *cfg.CFGMatcher, symbolId int, obj string) {
	fmt.Printf("Moveon: %s    value: %s\n", matcher.SymbolIdSymbolMap[symbolId], obj)
}

func (analyzer *EchoAnalyzer) Reduce(matcher *cfg.CFGMatcher, expIndex int) {
	fmt.Printf("Reduce: %v\n", matcher.Engine.ExpList[expIndex])
}

func TestNewCFGEngine(t *testing.T) {
	/* BEGIN -> S $
	 * S -> null
	 * S -> AS; S
	 * AS -> AS + num
	 * AS -> AS - num
	 * AS -> num
	 */
	engine, err := cfg.NewCFGEngine(
		[]string{"$", "null", ";", "+", "-", "num"},
		map[string]int{"+": 1, "-": 1},
		[]string{"BEGIN", "S", "AS"},
		[][]string{
			{"BEGIN", "S", "$"},
			{"S", "null"},
			{"S", "AS", ";", "S"},
			{"AS", "AS", "+", "AS"},
			{"AS", "AS", "-", "AS"},
			{"AS", "num"},
		},
		map[int]int{},
		"BEGIN",
		"null",
	)
	if err != nil {
		t.Errorf("NewCFGEngine fail: %v", err)
	}
	t.Logf("Engine: %p", engine)
}

func TestCFGMatcher(t *testing.T) {
	/* BEGIN -> S $
	 * S -> null
	 * S -> AS; S
	 * AS -> AS + num
	 * AS -> AS - num
	 * AS -> num
	 */
	engine, _ := cfg.NewCFGEngine(
		[]string{"$", "null", ";", "+", "-", "num"},
		map[string]int{"+": -1, "-": -1},
		[]string{"BEGIN", "S", "AS"},
		[][]string{
			{"BEGIN", "S", "$"},
			{"S", "null"},
			{"S", "AS", ";", "S"},
			{"AS", "AS", "+", "AS"},
			{"AS", "AS", "-", "AS"},
			{"AS", "num"},
		},
		map[int]int{},
		"BEGIN",
		"null",
	)

	echoAnalyzer := EchoAnalyzer{}
	matcher := cfg.NewCFGMatcher(engine, &echoAnalyzer)

	if ok, err := matcher.NextSymbol("num", "5"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("+", "+"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num", "3"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("-", "-"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num", "4"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol(";", ";"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num", "5"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("+", "+"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num", "9"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol(";", ";"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("$", "$"); !ok {
		t.Errorf("match fail, %v", err)
	}

	if ok, err := matcher.NextSymbol("num", "4"); ok {
		t.Errorf("match succ but expect fail, %v", err)
	}

	matcher = cfg.NewCFGMatcher(engine, nil)
	if ok, err := matcher.NextSymbol("$", "$"); !ok {
		t.Errorf("match fail, %v", err)
	}
}
