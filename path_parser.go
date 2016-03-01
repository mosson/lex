package lex

const availableChars = "abcdefghijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ1234567890-_."

var (
	phraseParser       = Many(Char(availableChars))
	staticPhraseParser = Seq(Token("/"), Phrase())
	paramsPhraseParser = Seq(Token("/"), Token(":"), Phrase())
)

const (
	static = "static"
	params = "params"
)

// Phrase /[a-zA-Z0-9]*/を検査するパーサーを返す
func Phrase() Parser {
	return func(target string, position int) *Result {
		return phraseParser(target, position)
	}
}

// StaticPhrase /\/[a-zA-Z0-9]*/を検査するパーサーを返す
func StaticPhrase() Parser {
	return func(target string, position int) *Result {
		result := staticPhraseParser(target, position)

		if result.Success {
			result.Attributes = assign(result.Attributes, map[string]string{
				"phrase": result.Target[1:],
				"type":   static,
			})
		}

		return result
	}
}

// ParamsPhrase /\/\:[a-zA-Z0-9]*/を検査するパーサーを返す
func ParamsPhrase() Parser {
	return func(target string, position int) *Result {
		result := paramsPhraseParser(target, position)

		if result.Success {
			result.Attributes = assign(result.Attributes, map[string]string{
				result.Target[2:]: "",
				"phrase":          result.Target[2:],
				"type":            params,
			})
		}

		return result
	}
}

// ParamsParser パラメータ表現の名前を取得するパーサーを返す
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

// PathParser 与えられた文字列を解析してパーサーを生成する
// 作成されたパーサーを実行した場合Result.Attributesにしかるべきパラメーター名と値の組を格納する
func PathParser(phrase string) Parser {
	cursor, parsers := 0, make([]Parser, 0)

	for {
		result := Choice(ParamsPhrase(), StaticPhrase())(phrase, cursor)

		if result.Success {

			cursor = result.Position
			if result.Attributes["type"] == params {
				parsers = append(parsers, ParamsParser(result.Attributes["phrase"]))
			} else {
				parsers = append(parsers, Token(result.Target))
			}
		} else {
			break
		}
	}

	return Seq(parsers...)
}
