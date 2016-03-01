package lex

import (
	"regexp"
	"strings"
)

// Result パースの結果を表す構造体
type Result struct {
	Success    bool
	Target     string
	Position   int
	Attributes map[string]string
}

// Parser パースを行う関数の型のシグニチャ
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
	return &Result{
		Success:    false,
		Target:     "",
		Position:   position,
		Attributes: map[string]string{},
	}
}

// Token 与えられた文字に一致するパーサーを返す
func Token(sample string) Parser {
	length := len(sample)

	return func(target string, position int) *Result {
		if len(target) < length+position {
			return incompatible(position)
		}

		if target[position:position+length] == sample {
			return &Result{
				Success:    true,
				Target:     sample,
				Position:   position + length,
				Attributes: map[string]string{},
			}
		}

		return incompatible(position)
	}
}

// Many 与えられたパーサーに0回以上合致するパーサーを返す
func Many(parser Parser) Parser {
	return func(target string, position int) *Result {
		var targets []string
		cursor := position
		targets = make([]string, 0)
		attributes := make(map[string]string, 0)

		for {
			result := parser(target, cursor)
			if result.Success {
				cursor = result.Position
				targets = append(targets, result.Target)
				attributes = assign(attributes, result.Attributes)
			} else {
				break
			}
		}

		return &Result{
			Success:    true,
			Target:     strings.Join(targets[:], ""),
			Position:   cursor,
			Attributes: attributes,
		}
	}
}

// Choice 与えられたパーサーのいづれかが成功すればよいパーサーを返す
func Choice(parsers ...Parser) Parser {
	return func(target string, position int) *Result {
		for i := 0; i < len(parsers); i++ {
			parser := parsers[i]

			result := parser(target, position)
			if result.Success {
				return &Result{
					Success:    true,
					Target:     result.Target,
					Position:   result.Position,
					Attributes: result.Attributes,
				}
			}
		}

		return incompatible(position)
	}
}

// Seq 引数すべてのパーサーを順番通りに検査するパーサーを返す
func Seq(parsers ...Parser) Parser {
	return func(target string, position int) *Result {
		var targets []string
		cursor := position
		targets = make([]string, 0)
		attributes := make(map[string]string, 0)
		for i := 0; i < len(parsers); i++ {
			parser := parsers[i]
			result := parser(target, cursor)
			if result.Success {
				cursor = result.Position
				targets = append(targets, result.Target)
				attributes = assign(attributes, result.Attributes)
			} else {
				return incompatible(position)
			}
		}

		return &Result{
			Success:    true,
			Target:     strings.Join(targets[:], ""),
			Position:   cursor,
			Attributes: attributes,
		}
	}
}

// Option 引数のパーサーの成否を判断しないパーサーを返す
func Option(parser Parser) Parser {
	return func(target string, position int) *Result {
		result := parser(target, position)
		if result.Success {
			return result
		}

		return &Result{
			Success:    true,
			Target:     "",
			Position:   position,
			Attributes: map[string]string{},
		}
	}
}

// Char 入力された文字列のどれかに一致するパーサーを返す
func Char(src string) Parser {
	dict := make(map[string]bool, len(src))
	for _, r := range src {
		dict[string(r)] = true
	}

	return func(target string, position int) *Result {
		if len(target) < position+1 {
			return incompatible(position)
		}

		targetChar := target[position : position+1]

		if _, ok := dict[targetChar]; ok {
			return &Result{
				Success:    true,
				Target:     targetChar,
				Position:   position + 1,
				Attributes: map[string]string{},
			}
		}

		return incompatible(position)
	}
}

// RegExp 正規表現を使用するパーサーを返す
func RegExp(pattern *regexp.Regexp) Parser {
	return func(target string, position int) *Result {
		if position > len(target) {
			return incompatible(position)
		}

		sample := target[position:]
		if pattern.MatchString(sample) {
			index := pattern.FindStringIndex(sample)
			return &Result{
				Success:    true,
				Target:     pattern.FindString(sample),
				Position:   position + index[1],
				Attributes: map[string]string{},
			}
		}

		return incompatible(position)
	}
}

// Lazy 遅延評価ができるパーサーを生成
func Lazy(ptr *Parser) Parser {
	return func(target string, position int) *Result {
		parser := *ptr
		return parser(target, position)
	}
}

// Map Result.Targetを加工できるパーサーを返す
func Map(parser Parser, processFn func(target string) string) Parser {
	return func(target string, position int) *Result {
		result := parser(target, position)
		if result.Success {
			return &Result{
				Success:    true,
				Target:     processFn(result.Target),
				Position:   result.Position,
				Attributes: map[string]string{},
			}
		}
		return result
	}
}
