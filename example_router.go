package lex

type handler func(map[string]string)

// router ルーティングを実行する構造体
type router struct {
	routes map[*Parser]handler
}

// newRouter ルーターを新規作成して返す
func newRouter() *router {
	return &router{routes: make(map[*Parser]handler)}
}

// register パスに対応する関数を登録する
func (router *router) register(path string, fn handler) {
	parser := PathParser(path)
	router.routes[&parser] = fn
}

// handle パスを検査して対応する関数が登録されていれば実行する
func (router *router) handle(path string) {
	for ptr, handler := range router.routes {
		parser := *ptr
		result := parser(path, 0)
		if result.Success {
			handler(result.Attributes)
			break
		}
	}
}
