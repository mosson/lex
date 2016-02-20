package lex

import "regexp"

func Phrase() Parser {
	parser := RegExp(regexp.MustCompile(`[a-zA-Z0-9]+`))
	return func(target string, position int) *Result {
		return parser(target, position)
	}
}

func StaticPhrase() Parser {
	return func(target string, position int) *Result {
		parser := Seq(Token("/"), Phrase())
		result := parser(target, position)

		if result.Success {
			result.Attributes = assign(result.Attributes, map[string]string{
				"phrase": result.Target[1:],
				"type":   "static",
			})
		}

		return result
	}
}

func ParamsPhrase() Parser {
	return func(target string, position int) *Result {
		parser := Seq(Token("/"), Token(":"), Phrase())

		result := parser(target, position)

		if result.Success {
			result.Attributes = assign(result.Attributes, map[string]string{
				result.Target[2:]: "",
				"phrase":          result.Target[2:],
				"type":            "params",
			})
		}

		return result
	}
}

func ParamsParser(phrase string) Parser {
	return func(target string, position int) *Result {
		parser := Seq(Token("/"), Phrase())

		result := parser(target, position)

		if result.Success {
			result.Attributes = assign(result.Attributes, map[string]string{
				phrase: result.Target[1:],
			})
		}

		return result
	}
}

func PathParser(phrase string) Parser {
	pos, frag := 0, make([]Parser, 0)

	for {
		result := Choice(ParamsPhrase(), StaticPhrase())(phrase, pos)

		if result.Success {

			pos = result.Position
			if result.Attributes["type"] == "params" {
				frag = append(frag, ParamsParser(result.Attributes["phrase"]))
			} else {
				frag = append(frag, Token(result.Target))
			}
		} else {
			break
		}
	}

	return Seq(frag...)
}