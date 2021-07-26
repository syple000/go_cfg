package cfg

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

const (
	unknownExp = iota // 默认零值作为第一个值
	nullExp
	notNullExp
)

const (
	unknown = iota
	reduce
	moveon
)

type nextExpPos struct {
	ExpIndex int
	NextPos  int
}

func newNextExpPos(expIndex int, nextPos int) nextExpPos {
	return nextExpPos{ExpIndex: expIndex, NextPos: nextPos}
}

type nextExpPosList []nextExpPos

func (l nextExpPosList) Len() int {
	return len(l)
}

func (l nextExpPosList) Less(i int, j int) bool {
	if l[i].ExpIndex < l[j].ExpIndex {
		return true
	} else if l[i].ExpIndex > l[j].ExpIndex {
		return false
	} else {
		return l[i].NextPos < l[j].NextPos
	}
}

func (l nextExpPosList) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
}

type action struct {
	Act int
	// reduce时表示expIndex
	// moveon时表示下一个状态
	Arg int
}

func newAction(act int, arg int) action {
	return action{Act: act, Arg: arg}
}

type CFGEngine struct {
	ExpList         [][]string
	StatusTable     [][]action
	FinalSymbolList []string
	GenSymbolList   []string
	StartGenSymbol  string
	NullFinalSymbol string
}

type CFGMatcher struct {
	Engine      *CFGEngine
	StatusStack []int
	SymbolIdMap map[string]int
	OK          bool
}

// 返回值表示是否有变更，假如src包含dest，则返回false，否则true
// value默认设置为0
func mergeIntSet(src map[int]int, dest map[int]int) bool {
	changed := false
	for k := range dest {
		if _, ok := src[k]; !ok {
			src[k] = 0
			changed = true
		}
	}
	return changed
}

func mergeStringSet(src map[string]int, dest map[string]int) bool {
	changed := false
	for k := range dest {
		if _, ok := src[k]; !ok {
			src[k] = 0
			changed = true
		}
	}
	return changed
}

func genGenSymbolClosureMap(expList [][]string) map[string]map[int]int {
	genSymbolClosureMap := make(map[string]map[int]int)
	for index, exp := range expList {
		genSymbol := exp[0]
		if _, ok := genSymbolClosureMap[genSymbol]; !ok {
			genSymbolClosureMap[genSymbol] = make(map[int]int)
		}
		genSymbolClosureMap[genSymbol][index] = 0
	}

	// 循环给闭包添加元素直至闭包元素不再有新增
	for {
		changed := false
		for genSymbol, closure := range genSymbolClosureMap {
			newIndexSet := make(map[int]int)
			for expIndex := range closure {
				exp := expList[expIndex]
				// 表达式中的首元素，如A->B c, B就是首元素
				if s, ok := genSymbolClosureMap[exp[1]]; ok {
					mergeIntSet(newIndexSet, s)
				}
			}
			if mergeIntSet(closure, newIndexSet) {
				changed = true
				genSymbolClosureMap[genSymbol] = closure
			}
		}
		if !changed {
			break
		}
	}
	return genSymbolClosureMap
}

func isIn(symbol string, symbolSet map[string]int) bool {
	_, ok := symbolSet[symbol]
	return ok
}

func genGenSymbolNullInfoMap(expList [][]string,
	finalSymbolSet map[string]int,
	genSymbolSet map[string]int,
	nullFinalSymbol string) map[string]int {

	genSymbolExpMap := make(map[string]map[int]int)
	for index, exp := range expList {
		if _, ok := genSymbolExpMap[exp[0]]; !ok {
			genSymbolExpMap[exp[0]] = make(map[int]int)
		}
		genSymbolExpMap[exp[0]][index] = 0
	}

	genSymbolNullInfoMap := make(map[string]int)
	for genSymbol := range genSymbolExpMap {
		genSymbolNullInfoMap[genSymbol] = unknownExp
	}

	for {
		changed := false
		for genSymbol, expIndexSet := range genSymbolExpMap {
			curNullInfo := genSymbolNullInfoMap[genSymbol]
			if curNullInfo == nullExp || curNullInfo == notNullExp {
				continue
			}
			isGenSymbolNull := notNullExp
			for expIndex := range expIndexSet {
				isNull := nullExp
				exp := expList[expIndex]
				for i := 1; i < len(exp); i++ {
					if isIn(exp[i], finalSymbolSet) {
						if exp[i] != nullFinalSymbol {
							isNull = notNullExp
							break
						}
					} else {
						if genSymbolNullInfoMap[exp[i]] == notNullExp {
							isNull = notNullExp
							break
						} else if genSymbolNullInfoMap[exp[i]] == unknownExp {
							isNull = unknownExp
						}
					}
				}
				if isNull == nullExp {
					isGenSymbolNull = nullExp
					break
				}
				if isNull == unknownExp {
					isGenSymbolNull = unknownExp
				}
			}
			if isGenSymbolNull != unknownExp {
				genSymbolNullInfoMap[genSymbol] = isGenSymbolNull
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return genSymbolNullInfoMap
}

func genGenSymbolFirstFinalSymbolSetMap(expList [][]string,
	genSymbolClosureMap map[string]map[int]int,
	genSymbolNullInfoMap map[string]int,
	finalSymbolSet map[string]int,
	genSymbolSet map[string]int,
	nullFinalSymbol string) map[string]map[string]int {

	genSymbolFirstFinalSymbolSetMap := make(map[string]map[string]int)

	for {
		changed := false
		for genSymbol, expIndexSet := range genSymbolClosureMap {
			firstFinalSymbolSet := make(map[string]int)
			for expIndex := range expIndexSet {
				exp := expList[expIndex]
				for i := 1; i < len(exp); i++ {
					if isIn(exp[i], finalSymbolSet) {
						if exp[i] != nullFinalSymbol {
							firstFinalSymbolSet[exp[i]] = 0
							break
						}
					} else {
						if v, ok := genSymbolFirstFinalSymbolSetMap[exp[i]]; ok {
							mergeStringSet(firstFinalSymbolSet, v)
						}
						if genSymbolNullInfoMap[exp[i]] != nullExp {
							break
						}
					}
				}
			}
			if v, ok := genSymbolFirstFinalSymbolSetMap[genSymbol]; ok {
				if mergeStringSet(v, firstFinalSymbolSet) {
					changed = true
				}
				genSymbolFirstFinalSymbolSetMap[genSymbol] = v
			} else {
				changed = true
				genSymbolFirstFinalSymbolSetMap[genSymbol] = firstFinalSymbolSet
			}
		}
		if !changed {
			break
		}
	}
	return genSymbolFirstFinalSymbolSetMap
}

func genGenSymbolNextFinalSymbolSetMap(expList [][]string,
	genSymbolFirstFinalSymbolSetMap map[string]map[string]int,
	genSymbolNullInfoMap map[string]int,
	finalSymbolSet map[string]int,
	genSymbolSet map[string]int,
	nullFinalSymbol string) map[string]map[string]int {

	genSymbolNextFinalSymbolSetMap := make(map[string]map[string]int)

	for {
		changed := false
		tmpGenSymbolNextFinalSymbolSetMap := make(map[string]map[string]int)
		for _, exp := range expList {
			nextFinalSymbolSet := make(map[string]int)
			// 拷贝当前表达式的产生符号的next final symbol 数组用于初始化
			if v, ok := genSymbolNextFinalSymbolSetMap[exp[0]]; ok {
				mergeStringSet(nextFinalSymbolSet, v)
			}
			for i := len(exp) - 1; i > 0; i-- {
				if isIn(exp[i], genSymbolSet) {
					if v, ok := tmpGenSymbolNextFinalSymbolSetMap[exp[i]]; ok {
						mergeStringSet(v, nextFinalSymbolSet)
						tmpGenSymbolNextFinalSymbolSetMap[exp[i]] = v
					} else {
						m := make(map[string]int)
						mergeStringSet(m, nextFinalSymbolSet)
						tmpGenSymbolNextFinalSymbolSetMap[exp[i]] = m
					}
					if genSymbolNullInfoMap[exp[i]] == nullExp {
						mergeStringSet(nextFinalSymbolSet, genSymbolFirstFinalSymbolSetMap[exp[i]])
					} else {
						nextFinalSymbolSet = make(map[string]int)
						mergeStringSet(nextFinalSymbolSet, genSymbolFirstFinalSymbolSetMap[exp[i]])
					}
				} else {
					if exp[i] != nullFinalSymbol {
						nextFinalSymbolSet = make(map[string]int)
						nextFinalSymbolSet[exp[i]] = 0
					}
				}
			}
		}
		for genSymbol, s := range tmpGenSymbolNextFinalSymbolSetMap {
			if _, ok := genSymbolNextFinalSymbolSetMap[genSymbol]; !ok {
				genSymbolNextFinalSymbolSetMap[genSymbol] = make(map[string]int)
			}
			cs := genSymbolNextFinalSymbolSetMap[genSymbol]
			if mergeStringSet(cs, s) {
				changed = true
				genSymbolNextFinalSymbolSetMap[genSymbol] = cs
			}
		}

		if !changed {
			break
		}
	}
	return genSymbolNextFinalSymbolSetMap
}

func serializeStatusMap(m map[nextExpPos]int) string {
	keys := make([]nextExpPos, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Sort(nextExpPosList(keys))

	var buffer bytes.Buffer
	for _, key := range keys {
		buffer.WriteString("{" + strconv.FormatInt(int64(key.ExpIndex), 10) + ":" +
			strconv.FormatInt(int64(key.NextPos), 10) + "}")
	}
	return buffer.String()
}

// 默认规约具有最高优先级
func genStatusTable(expList [][]string,
	genSymbolClosureMap map[string]map[int]int,
	genSymbolNextFinalSymbolSetMap map[string]map[string]int,
	finalSymbolList []string,
	genSymbolList []string,
	startGenSymbol string) [][]action {

	initStatus := make(map[nextExpPos]int)
	for expIndex := range genSymbolClosureMap[startGenSymbol] {
		initStatus[newNextExpPos(expIndex, 0)] = 0
	}
	statusArray := make([]map[nextExpPos]int, 0)
	statusMap := make(map[string]int)
	statusArray = append(statusArray, initStatus)
	statusMap[serializeStatusMap(initStatus)] = 0

	transTable := make(map[string]map[int]action)
	statusCount := 0

	for {
		index := statusCount
		statusCount = len(statusArray)
		if index == statusCount {
			break
		}
		for index < statusCount {
			newStatusMap := make(map[string]map[nextExpPos]int)
			status := statusArray[index]
			for k := range status {
				expIndex := k.ExpIndex
				expPos := k.NextPos
				exp := expList[expIndex]
				if expPos == len(exp)-1 {
					// 规约
					firstFinalSymbolSet := genSymbolNextFinalSymbolSetMap[exp[0]]
					for firstFinalSymbol := range firstFinalSymbolSet {
						if _, ok := transTable[firstFinalSymbol]; !ok {
							transTable[firstFinalSymbol] = make(map[int]action)
						}
						if v, ok := transTable[firstFinalSymbol][index]; ok {
							fmt.Printf("status: %d reduce by following: %s, but conflict with reduce exp: %d, rewrite",
								index, firstFinalSymbol, v.Arg)
						}
						transTable[firstFinalSymbol][index] = newAction(reduce, expIndex)
					}
				} else {
					symbol := exp[expPos+1]
					if _, ok := newStatusMap[symbol]; !ok {
						newStatusMap[symbol] = make(map[nextExpPos]int)
					}
					newStatusMap[symbol][newNextExpPos(expIndex, expPos+1)] = 0
					expPos++
					// 判断下一个元素是否是generate symbol，假如是，将该元素的所有closure的表达式加入集合
					if expPos+1 < len(exp) {
						nextSymbol := exp[expPos+1]
						if closure, ok := genSymbolClosureMap[nextSymbol]; ok {
							for expIndex := range closure {
								newStatusMap[symbol][newNextExpPos(expIndex, 0)] = 0
							}
						}
					}
				}
			}
			for symbol, info := range newStatusMap {
				if _, ok := transTable[symbol]; !ok {
					transTable[symbol] = make(map[int]action)
				}
				if v, ok := transTable[symbol][index]; ok {
					fmt.Printf("status: %d moveon by %s, but conflict with reduce exp: %d, skip",
						index, symbol, v.Arg)
					continue
				}
				infoStr := serializeStatusMap(info)
				if v, ok := statusMap[infoStr]; !ok {
					// 新状态
					statusMap[infoStr] = len(statusArray)
					statusArray = append(statusArray, info)
					transTable[symbol][index] = newAction(moveon, statusMap[infoStr])
				} else {
					transTable[symbol][index] = newAction(moveon, v)
				}
			}
			index++
		}
	}
	fmt.Printf("trans table:\n%v\nstatus array:\n%v\n", transTable, statusArray)
	// 整理状态表为
	// --\--    symbol_id
	// status
	symbolIdMap := make(map[string]int)
	for i := 0; i < len(finalSymbolList); i++ {
		symbolIdMap[finalSymbolList[i]] = i
	}
	for i := len(finalSymbolList); i < len(finalSymbolList)+len(genSymbolList); i++ {
		symbolIdMap[genSymbolList[i-len(finalSymbolList)]] = i
	}
	table := make([][]action, len(statusArray))
	for i := 0; i < len(statusArray); i++ {
		table[i] = make([]action, len(symbolIdMap))
	}
	for symbol, statusActMap := range transTable {
		for status, act := range statusActMap {
			table[status][symbolIdMap[symbol]] = act
		}
	}
	return table
}

func NewCFGEngine(finalSymbolList []string,
	genSymbolList []string,
	expList [][]string,
	startGenSymbol string,
	nullFinalSymbol string) (*CFGEngine, error) {
	// 校验初始化数据是否完全正确
	finalSymbolSet := make(map[string]int)
	genSymbolSet := make(map[string]int)
	for _, v := range finalSymbolList {
		finalSymbolSet[v] = 0
	}
	for _, v := range genSymbolList {
		genSymbolSet[v] = 0
	}
	if _, ok := genSymbolSet[startGenSymbol]; !ok {
		return nil, fmt.Errorf("start generate symbol: %s not found in generate symbol list", startGenSymbol)
	}
	if _, ok := finalSymbolSet[nullFinalSymbol]; !ok {
		return nil, fmt.Errorf("null final symbol: %s not found in final symbol list", nullFinalSymbol)
	}
	// 校验exp list的同时，将同一个生成符号的数据进行归类
	for _, exp := range expList {
		if len(exp) < 2 {
			return nil, fmt.Errorf("exp: %v invalid", exp)
		}
		if _, ok := genSymbolSet[exp[0]]; !ok {
			return nil, fmt.Errorf("exp: %s should start with generate symbol", exp)
		}
		for _, symbol := range exp {
			_, inFinalSymbolList := finalSymbolSet[symbol]
			_, inGenSymbolList := genSymbolSet[symbol]
			if !inFinalSymbolList && !inGenSymbolList {
				return nil, fmt.Errorf("symbol: %s not found in symbol list", symbol)
			}
		}
	}

	// 要建立的是lr语法分析器，需要找出所有产生符号的开始终止符号与是否可空用于归约。
	// 大致流程：通过当前符号对应的表达式闭包，多次循环中将自身直接开始终止符号集合与其它产生符号的
	// 开始终止符号集合（初始化为空）加入到自身开始终止符号集合并根据自身表达式判断是否可空并结合表
	// 达式中其它产生符号的是否可空（初始为unknown，后可以是yes or no）。等到所有状态稳定即可
	// 闭包是指当前产生符号下，下一个可能匹配符号可以出现的表达式的集合

	// 解析逻辑依赖状态跳转表，状态跳转的关键在于状态的完整性，找到当前状态的闭包在进行下一个状态的
	// 计算。该算法基于上述算法

	// 未在该map中的generate symbol都是无效符号，忽略
	genSymbolClosureMap := genGenSymbolClosureMap(expList)
	// 校验产生符号都有对应表达式
	for _, genSymbol := range genSymbolList {
		if _, ok := genSymbolClosureMap[genSymbol]; !ok {
			return nil, fmt.Errorf("genarate symbol: %s has no generate expression", genSymbol)
		}
	}

	genSymbolNullInfoMap := genGenSymbolNullInfoMap(expList, finalSymbolSet, genSymbolSet, nullFinalSymbol)
	for genSymbol, nullInfo := range genSymbolNullInfoMap {
		if nullInfo == unknownExp {
			return nil, fmt.Errorf("generate symbol: %s null info unknown", genSymbol)
		}
	}

	genSymbolFirstFinalSymbolSetMap := genGenSymbolFirstFinalSymbolSetMap(expList, genSymbolClosureMap,
		genSymbolNullInfoMap, finalSymbolSet, genSymbolSet, nullFinalSymbol)

	genSymbolNextFinalSymbolSetMap := genGenSymbolNextFinalSymbolSetMap(expList, genSymbolFirstFinalSymbolSetMap,
		genSymbolNullInfoMap, finalSymbolSet, genSymbolSet, nullFinalSymbol)

	fmt.Printf("closure:\n%v\nsymbol null info map:\n%v\nfirst final symbol set:\n%v\nnext final symbol set:\n%v\n",
		genSymbolClosureMap, genSymbolNullInfoMap, genSymbolFirstFinalSymbolSetMap, genSymbolNextFinalSymbolSetMap)

	statusTable := genStatusTable(expList, genSymbolClosureMap, genSymbolNextFinalSymbolSetMap, finalSymbolList,
		genSymbolList, startGenSymbol)

	return &CFGEngine{
		ExpList:         expList,
		StatusTable:     statusTable,
		FinalSymbolList: finalSymbolList,
		GenSymbolList:   genSymbolList,
		StartGenSymbol:  startGenSymbol,
		NullFinalSymbol: nullFinalSymbol,
	}, nil
}

func NewCFGMatcher(engine *CFGEngine) *CFGMatcher {
	matcher := CFGMatcher{
		Engine:      engine,
		StatusStack: make([]int, 0, 16),
		OK:          true,
	}
	matcher.StatusStack = append(matcher.StatusStack, 0)

	matcher.SymbolIdMap = make(map[string]int)
	for i := 0; i < len(engine.FinalSymbolList); i++ {
		matcher.SymbolIdMap[engine.FinalSymbolList[i]] = i
	}
	for i := len(engine.FinalSymbolList); i < len(engine.FinalSymbolList)+len(engine.GenSymbolList); i++ {
		matcher.SymbolIdMap[engine.GenSymbolList[i-len(engine.FinalSymbolList)]] = i
	}

	return &matcher
}

func (matcher *CFGMatcher) NextSymbol(symbol string) (bool, error) {
	if !matcher.OK {
		return false, fmt.Errorf("matcher is not ok")
	}
	if id, ok := matcher.SymbolIdMap[symbol]; !ok {
		return false, fmt.Errorf("symbol: %s not found", symbol)
	} else {
		return matcher.nextSymbolId(id)
	}
}

func (matcher *CFGMatcher) NextSymbolId(symbolId int) (bool, error) {
	if !matcher.OK {
		return false, fmt.Errorf("matcher is not ok")
	}
	if symbolId < 0 || symbolId >= len(matcher.SymbolIdMap) {
		return false, fmt.Errorf("symbol id: %d is invalid", symbolId)
	}
	return matcher.nextSymbolId(symbolId)
}

// 空符号永远不会被作为参数，所以在匹配异常时，需要尝试空符号
func (matcher *CFGMatcher) nextSymbolId(symbolId int) (bool, error) {
	for {
		curSatus := matcher.StatusStack[len(matcher.StatusStack)-1]
		action := matcher.Engine.StatusTable[curSatus][symbolId]
		if action.Act == unknown {
			// 考虑空符号
			action = matcher.Engine.StatusTable[curSatus][matcher.SymbolIdMap[matcher.Engine.NullFinalSymbol]]
			if action.Act != moveon {
				matcher.OK = false
				return false, fmt.Errorf("symbol id: %d match fail", symbolId)
			}
			matcher.StatusStack = append(matcher.StatusStack, action.Arg)
			continue
		}
		if action.Act == moveon {
			matcher.StatusStack = append(matcher.StatusStack, action.Arg)
			return true, nil
		} else {
			// reduce
			exp := matcher.Engine.ExpList[action.Arg]
			for i := 1; i < len(exp); i++ {
				matcher.StatusStack = matcher.StatusStack[0 : len(matcher.StatusStack)-1]
			}
			curSatus = matcher.StatusStack[len(matcher.StatusStack)-1]
			action = matcher.Engine.StatusTable[curSatus][matcher.SymbolIdMap[exp[0]]]
			if action.Act != moveon {
				matcher.OK = false
				return false, fmt.Errorf("symbol id: %d reduce by: %v but moveon fail", symbolId, exp)
			}
			matcher.StatusStack = append(matcher.StatusStack, action.Arg)
			continue
		}
	}
}
