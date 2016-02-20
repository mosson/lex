package lex

const availableChars = "abcdefghijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ1234567890-_."

var (
	phraseParser       Parser = Many(Char(availableChars))
	staticPhraseParser Parser = Seq(Token("/"), Phrase())
	paramsPhraseParser Parser = Seq(Token("/"), Token(":"), Phrase())
)

func Phrase() Parser {
	return func(target string, position int) *Result {
		return phraseParser(target, position)
	}
}

func StaticPhrase() Parser {
	return func(target string, position int) *Result {
		result := staticPhraseParser(target, position)

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
		result := paramsPhraseParser(target, position)

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
		result := staticPhraseParser(target, position)

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
