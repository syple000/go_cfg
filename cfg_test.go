package main

import (
	"testing"

	"github.com/syple000/compiler/cfg"
)

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

func TestCFGEngine(t *testing.T) {
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

	matcher := cfg.NewCFGMatcher(engine)

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
}
