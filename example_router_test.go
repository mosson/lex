package lex

import "testing"

func TestHandle(t *testing.T) {
	router := newRouter()

	done := false

	router.register("/api/v1/entries/:id/:name", func(params map[string]string) {
		if params["id"] != "123" {
			t.Errorf("expected 123, actual %v", params["id"])
		}

		if params["name"] != "hoge" {
			t.Errorf("expected hoge, actual %v", params["name"])
		}

		done = true
	})

	router.handle("/api/v1/entries/123/hoge")

	if done != true {
		t.Errorf("not invoked handler")
	}
}
