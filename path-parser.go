package lex

import "regexp"

type PathToken struct {
	Result
	Key   string
	Value string
}

type PathParser func(string, int) *PathToken

const AvailableChar = `[\-\_\.\!\*\'\(\)a-zA-Z0-9]+`

func StaticParser() PathParser {
	p := Seq(Token("/"), RegExp(regexp.MustCompile(AvailableChar)))

	return func(target string, position int) *PathToken {
		res := p(target, position)
		return &PathToken{
			Result: *res,
			Key:    "",
			Value:  "",
		}
	}
}
