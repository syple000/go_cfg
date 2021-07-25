/*
 * 正则解析器：仅支持简单正则表达式规则包括:
 * 连接. 或| 左右括号() 星号* 以及特殊的转义\
 * 正则上下文无关语法:
 * 1. E0 -> any_char E1
 * 2. E1 -> E2
 *    E1 -> E3
 * 3.
 */
package re

type ReEngine struct{}

func NewReEngine(expList []string) *ReEngine {
	return nil
}
