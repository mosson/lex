package lex

import "testing"

func TestStaticParser(t *testing.T) {
	result := StaticParser()("/api/v-a1", 0)
	if !result.Success {
		t.Error("成功しないとおかしい")
	}

	if result.Target != "/api" {
		t.Log(result.Target)
		t.Error("最初のトークンを取得できないとおかしい")
	}

	if result.Position != 4 {
		t.Error("最初のトークン位置まで読み取らないとおかしい")
	}

	if result.Key != "" {
		t.Error("空文字じゃないとおかしい")
	}

	if result.Value != "" {
		t.Error("空文字じゃないとおかしい")
	}
}
