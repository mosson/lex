package lex

import (
	"regexp"
	"strings"
)

type Result struct {
	Success  bool
	Target   string
	Position int
}

type Parser func(string, int) *Result

func incompatible(position int) *Result {
	return &Result{Success: false, Target: "", Position: position}
}

/*
与えられた文字に一致するパーサーを返す
*/

func Token(str string) Parser {
	l := len(str)

	return func(target string, position int) *Result {
		if len(target) < l+position {
			return incompatible(position)
		}

		if target[position:position+l] == str {
			return &Result{Success: true, Target: str, Position: position + l}
		} else {
			return incompatible(position)
		}
	}
}

/*
与えられたパーサーに０階以上合致するパーサーを返す
*/

func Many(fn Parser) Parser {
	return func(target string, position int) *Result {
		var r []string
		p := position

		for {
			var parsed *Result = fn(target, p)
			if parsed.Success {
				r = append(r, parsed.Target)
				p = parsed.Position
			} else {
				break
			}
		}

		return &Result{Success: true, Target: strings.Join(r[:], ""), Position: p}
	}
}

/*
与えられたパーサーのいづれかが成功すればよい
*/

func Choice(Parsers ...Parser) Parser {
	return func(target string, position int) *Result {
		for i := 0; i < len(Parsers); i++ {
			p := Parsers[i]

			parsed := p(target, position)
			if parsed.Success {
				return &Result{Success: true, Target: parsed.Target, Position: parsed.Position}
			}
		}

		return incompatible(position)
	}
}

/*
パーサーを連結する
*/

func Seq(Parsers ...Parser) Parser {
	return func(target string, position int) *Result {
		var r []string
		pos := position
		for i := 0; i < len(Parsers); i++ {
			p := Parsers[i]
			var parsed *Result = p(target, pos)
			if parsed.Success {
				r = append(r, parsed.Target)
				pos = parsed.Position
			} else {
				return incompatible(position)
			}
		}

		return &Result{Success: true, Target: strings.Join(r[:], ""), Position: pos}
	}
}

/*
失敗してもよいパーサー評価
*/

func Option(p Parser) Parser {
	return func(target string, position int) *Result {
		var r *Result = p(target, position)
		if r.Success {
			return r
		} else {
			return &Result{Success: true, Target: "", Position: position}
		}
	}
}

/*
入力された文字列のどれかに一致するパーサーを生成
*/

func Char(str string) Parser {
	dict := make(map[string]bool, len(str))
	for _, c := range str {
		dict[string(c)] = true
	}

	return func(target string, position int) *Result {
		targetString := target[position:1]
		if _, ok := dict[targetString]; ok {
			return &Result{Success: true, Target: targetString, Position: position + 1}
		} else {
			return incompatible(position)
		}
	}
}

/*
正規表現を使用するパーサーを生成
*/

func RegExp(pattern *regexp.Regexp) Parser {

	return func(target string, position int) *Result {
		if pattern.MatchString(target) {
			index := pattern.FindStringIndex(target)
			return &Result{Success: true, Target: pattern.FindString(target), Position: index[1]}
		} else {
			return incompatible(position)
		}
	}
}

/*
遅延評価ができるパーサーを生成
golangでは実質遅延評価はできないのであまり意味をなさない
*/

func Lazy(p Parser) Parser {
	return func(target string, position int) *Result {
		return p(target, position)
	}
}

/*
Targetを加工できるパーサーを生成
*/

func Map(p Parser, fn func(target string) string) Parser {
	return func(target string, position int) *Result {
		res := p(target, position)
		if res.Success {
			return &Result{Success: true, Target: fn(res.Target), Position: res.Position}
		} else {
			return res
		}
	}
}
