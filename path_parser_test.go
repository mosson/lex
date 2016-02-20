package lex

import "testing"

func TestPhrase(t *testing.T) {
	result := Phrase()("id/:uid", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "id" {
		t.Errorf("expected id, actual %v", result.Target)
	}

	if result.Position != 2 {
		t.Errorf("expected 2, actual %v", result.Position)
	}
}

func TestStaticPhrase(t *testing.T) {
	result := StaticPhrase()("/foo/bar", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "/foo" {
		t.Errorf("expected /foo, actual %v", result.Target)
	}

	if result.Position != 4 {
		t.Errorf("expected 4, actual %v", result.Position)
	}
}

func TestParamsPhrase(t *testing.T) {
	result := ParamsPhrase()("/:id/:uid", 0)

	if !result.Success {
		t.Errorf("exptected true, actual %v", result.Success)
	}

	if result.Target != "/:id" {
		t.Errorf("expected /:id, actual %v", result.Target)
	}

	if result.Position != 4 {
		t.Errorf("expected 4, actual %v", result.Position)
	}

	if _, ok := result.Attributes["id"]; !ok {
		t.Errorf("expected true, actual %v", ok)
	}
}

func TestParamsParser(t *testing.T) {
	result := ParamsParser("id")("/123/456", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Attributes["id"] != "123" {
		t.Errorf("expected 123, actual %v", result.Attributes["id"])
	}
}

func TestPathParser(t *testing.T) {
	result := PathParser("/api/v1/entries/:id/:query")("/api/v1/entries/123/hello", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Attributes["id"] != "123" {
		t.Errorf("expected 123, actual %v", result.Attributes["id"])
	}

	if result.Attributes["query"] != "hello" {
		t.Errorf("expected hello, actual %v", result.Attributes["query"])
	}
}
