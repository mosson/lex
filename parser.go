package lex

import (
	"regexp"
	"strings"
)

type result struct {
	Success  bool
	Target   string
	Position int
}

type parser func(string, int) *result

func incompatible(position int) *result {
	return &result{Success: false, Target: "", Position: position}
}

/*
与えられた文字に一致するパーサーを返す
*/

func Token(str string) parser {
	l := len(str)

	return func(target string, position int) *result {
		if len(target) < l+position {
			return incompatible(position)
		}

		if target[position:position+l] == str {
			return &result{Success: true, Target: str, Position: position + l}
		} else {
			return incompatible(position)
		}
	}
}

/*
与えられたパーサーに０階以上合致するパーサーを返す
*/

func Many(fn parser) parser {
	return func(target string, position int) *result {
		var r []string
		p := position

		for {
			var parsed *result = fn(target, p)
			if parsed.Success {
				r = append(r, parsed.Target)
				p = parsed.Position
			} else {
				break
			}
		}

		return &result{Success: true, Target: strings.Join(r[:], ""), Position: p}
	}
}

/*
与えられたパーサーのいづれかが成功すればよい
*/

func Choice(parsers ...parser) parser {
	return func(target string, position int) *result {
		for i := 0; i < len(parsers); i++ {
			p := parsers[i]

			parsed := p(target, position)
			if parsed.Success {
				return &result{Success: true, Target: parsed.Target, Position: parsed.Position}
			}
		}

		return incompatible(position)
	}
}

/*
パーサーを連結する
*/

func Seq(parsers ...parser) parser {
	return func(target string, position int) *result {
		var r []string
		pos := position
		for i := 0; i < len(parsers); i++ {
			p := parsers[i]
			var parsed *result = p(target, pos)
			if parsed.Success {
				r = append(r, parsed.Target)
				pos = parsed.Position
			} else {
				return incompatible(position)
			}
		}

		return &result{Success: true, Target: strings.Join(r[:], ""), Position: pos}
	}
}

/*
失敗してもよいパーサー評価
*/

func Option(p parser) parser {
	return func(target string, position int) *result {
		var r *result = p(target, position)
		if r.Success {
			return r
		} else {
			return &result{Success: true, Target: "", Position: position}
		}
	}
}

/*
入力された文字列のどれかに一致するパーサーを生成
*/

func Char(str string) parser {
	dict := make(map[string]bool, len(str))
	for _, c := range str {
		dict[string(c)] = true
	}

	return func(target string, position int) *result {
		targetString := target[position:1]
		if _, ok := dict[targetString]; ok {
			return &result{Success: true, Target: targetString, Position: position + 1}
		} else {
			return incompatible(position)
		}
	}
}

/*
正規表現を使用するパーサーを生成
*/

func RegExp(pattern *regexp.Regexp) parser {

	return func(target string, position int) *result {
		if pattern.MatchString(target) {
			index := pattern.FindStringIndex(target)
			return &result{Success: true, Target: pattern.FindString(target), Position: index[1]}
		} else {
			return incompatible(position)
		}
	}
}

/*
遅延評価ができるパーサーを生成
golangでは実質遅延評価はできないのであまり意味をなさない
*/

func Lazy(p parser) parser {
	return func(target string, position int) *result {
		return p(target, position)
	}
}

/*
Targetを加工できるパーサーを生成
*/

func Map(p parser, fn func(target string) string) parser {
	return func(target string, position int) *result {
		res := p(target, position)
		if res.Success {
			return &result{Success: true, Target: fn(res.Target), Position: res.Position}
		} else {
			return res
		}
	}
}
