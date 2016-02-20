package lex

import (
	"regexp"
	"strings"
)

type Result struct {
	Success    bool
	Target     string
	Position   int
	Attributes map[string]string
}

type Parser func(string, int) *Result

func assign(assigns ...map[string]string) map[string]string {
	result := make(map[string]string)

	for i := 0; i < len(assigns); i++ {
		target := assigns[i]
		for k, v := range target {
			result[k] = v
		}
	}

	return result
}

func incompatible(position int) *Result {
	return &Result{Success: false, Target: "", Position: position, Attributes: map[string]string{}}
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
			return &Result{Success: true, Target: str, Position: position + l, Attributes: map[string]string{}}
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
		r, a := make([]string, 0), make(map[string]string, 0)
		p := position

		for {
			var parsed *Result = fn(target, p)
			if parsed.Success {
				r = append(r, parsed.Target)
				p = parsed.Position
				a = assign(a, parsed.Attributes)
			} else {
				break
			}
		}

		return &Result{Success: true, Target: strings.Join(r[:], ""), Position: p, Attributes: a}
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
				return &Result{Success: true, Target: parsed.Target, Position: parsed.Position, Attributes: parsed.Attributes}
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
		r, a := make([]string, 0), make(map[string]string, 0)
		pos := position
		for i := 0; i < len(Parsers); i++ {
			p := Parsers[i]
			var parsed *Result = p(target, pos)
			if parsed.Success {
				r = append(r, parsed.Target)
				pos = parsed.Position
				a = assign(a, parsed.Attributes)
			} else {
				return incompatible(position)
			}
		}

		return &Result{Success: true, Target: strings.Join(r[:], ""), Position: pos, Attributes: a}
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
			return &Result{Success: true, Target: "", Position: position, Attributes: map[string]string{}}
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
		if len(target) < position+1 {
			return incompatible(position)
		}

		targetString := target[position : position+1]

		if _, ok := dict[targetString]; ok {
			return &Result{Success: true, Target: targetString, Position: position + 1, Attributes: map[string]string{}}
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
		if position > len(target) {
			return incompatible(position)
		}

		t := target[position:]
		if pattern.MatchString(t) {
			index := pattern.FindStringIndex(t)
			return &Result{Success: true, Target: pattern.FindString(t), Position: position + index[1], Attributes: map[string]string{}}
		} else {
			return incompatible(position)
		}
	}
}

/*
遅延評価ができるパーサーを生成
主に再帰を表現するのに使う
*/

func Lazy(p *Parser) Parser {
	return func(target string, position int) *Result {
		var fn Parser = *p
		return fn(target, position)
	}
}

/*
Targetを加工できるパーサーを生成
*/

func Map(p Parser, fn func(target string) string) Parser {
	return func(target string, position int) *Result {
		res := p(target, position)
		if res.Success {
			return &Result{Success: true, Target: fn(res.Target), Position: res.Position, Attributes: map[string]string{}}
		} else {
			return res
		}
	}
}
