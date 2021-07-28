package main

import (
	"fmt"
	"testing"

	"github.com/syple000/compiler/cfg"
)

type EchoAnalyzer struct{}

func (analyzer *EchoAnalyzer) Moveon(matcher *cfg.CFGMatcher, symbolId int) {
	fmt.Printf("Moveon: %s\n", matcher.SymbolIdSymbolMap[symbolId])
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
		[]string{"BEGIN", "S", "AS"},
		[][]string{
			{"BEGIN", "S", "$"},
			{"S", "null"},
			{"S", "AS", ";", "S"},
			{"AS", "AS", "+", "num"},
			{"AS", "AS", "-", "num"},
			{"AS", "num"},
		},
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
		[]string{"BEGIN", "S", "AS"},
		[][]string{
			{"BEGIN", "S", "$"},
			{"S", "null"},
			{"S", "AS", ";", "S"},
			{"AS", "AS", "+", "num"},
			{"AS", "AS", "-", "num"},
			{"AS", "num"},
		},
		"BEGIN",
		"null",
	)

	echoAnalyzer := EchoAnalyzer{}
	matcher := cfg.NewCFGMatcher(engine, &echoAnalyzer)

	if ok, err := matcher.NextSymbol("num"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("+"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("-"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol(";"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("+"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("num"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol(";"); !ok {
		t.Errorf("match fail, %v", err)
	}
	if ok, err := matcher.NextSymbol("$"); !ok {
		t.Errorf("match fail, %v", err)
	}

	if ok, err := matcher.NextSymbol("num"); ok {
		t.Errorf("match succ but expect fail, %v", err)
	}

	matcher = cfg.NewCFGMatcher(engine, nil)
	if ok, err := matcher.NextSymbol("$"); !ok {
		t.Errorf("match fail, %v", err)
	}
}
