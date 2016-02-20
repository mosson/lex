package lex

import (
	"regexp"
	"testing"
)

func TestAssign(t *testing.T) {
	a, b, c := map[string]string{
		"a": "a",
		"b": "b",
	}, map[string]string{
		"c": "c",
		"a": "d",
	}, map[string]string{
		"d": "d",
	}

	result := assign(a, b, c)

	if result["a"] != "d" {
		t.Errorf("expected d, actual %v", result["a"])
	}

	if result["b"] != "b" {
		t.Errorf("expected b, actual %v", result[""])
	}

	if result["c"] != "c" {
		t.Errorf("expected c, actual %v", result["c"])
	}

	if result["d"] != "d" {
		t.Errorf("expected d, actual %v", result["d"])
	}

}

func TestIncompatible(t *testing.T) {
	result := incompatible(4)

	if result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "" {
		t.Errorf("expected \"\", actual %v", result.Target)
	}

	if result.Position != 4 {
		t.Errorf("expected 4, actual %v", result.Position)
	}
}

func TestToken(t *testing.T) {
	result := Token("hoge")("hoge", 0)
	if result.Success {
		if result.Target != "hoge" {
			t.Errorf("expected hoge, actual %v", result.Target)
		}

		if result.Position != len("hoge") {
			t.Errorf("expected %v, actual %v", len("hoge"), result.Position)
		}
	} else {
		t.Errorf("expected true, actual %v", result.Success)
	}

	// 解析文字が入力より長い場合
	result2 := Token("hogehoge")("hoge", 0)
	if result2.Success {
		t.Errorf("expected true, actual %v", result2.Success)
	}

	if result2.Target != "" {
		t.Errorf("expected \"\", actual %v", result2.Target)
	}

	if result2.Position != 0 {
		t.Errorf("expected 0, actual %v", result2.Position)
	}

	// incompatible
	result3 := Token("hoge")("fugafuga", 0)
	if result3.Success {
		t.Errorf("expected true, actual %v", result3.Success)
	}

	if result3.Target != "" {
		t.Errorf("expected \"\", actual %v", result3.Target)
	}

	if result3.Position != 0 {
		t.Errorf("expected 0, actual %v", result3.Position)
	}
}

func TestMany(t *testing.T) {
	result := Many(Token("hoge"))("hogehoge", 0)
	t.Log(result)
	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "hogehoge" {
		t.Errorf("expected hogehoge, actual %v", result.Target)
	}

	if result.Position != 8 {
		t.Errorf("expected 8, actual %v", result.Position)
	}

	result2 := Many(Token("hoge"))("hogehoge", 4)
	t.Log(result2)
	if !result2.Success {
		t.Errorf("expected true, actual %v", result2.Success)
	}

	if result2.Target != "hoge" {
		t.Errorf("expected hoge, actual %v", result2.Target)
	}

	if result2.Position != 8 {
		t.Errorf("expected 8, actual %v", result2.Position)
	}
}

func TestChoice(t *testing.T) {
	result := Choice(Token("foo"), Token("bar"))("bar", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "bar" {
		t.Errorf("expected bar, actual %v", result.Target)
	}

	if result.Position != 3 {
		t.Errorf("expected 3, actual %v", result.Position)
	}

	result2 := Choice(Token("foo"), Token("bar"))("foo", 0)

	if !result2.Success {
		t.Errorf("expected true, actual %v", result2.Success)
	}

	if result2.Target != "foo" {
		t.Errorf("expected foo, actual %v", result2.Target)
	}

	if result2.Position != 3 {
		t.Errorf("expected 3, actual %v", result2.Position)
	}

	result3 := Choice(Token("foo"), Token("bar"))("baz", 0)

	if result3.Success {
		t.Errorf("expected false, actual %v", result.Success)
	}

	if result3.Target != "" {
		t.Errorf("expected \"\", actual %v", result3.Target)
	}

	if result3.Position != 0 {
		t.Errorf("expected 0, actual %v", result3.Position)
	}
}

func TestSeq(t *testing.T) {
	result := Seq(Token("foo"), Token("bar"), Token("baz"))("foobarbaz", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "foobarbaz" {
		t.Errorf("expected foobarbaz, actual %v", result.Target)
	}

	if result.Position != len("foobarbaz") {
		t.Errorf("expected foobarbaz, actual %v", result.Position)
	}

	result2 := Seq(Token("foo"), Token("bar"), Token("baz"))("bazbarfoo", 0)

	if result2.Success {
		t.Errorf("expected false, actual %v", result2.Success)
	}

	if result2.Target != "" {
		t.Errorf("expected \"\", actual %v", result2.Target)
	}

	if result2.Position != 0 {
		t.Errorf("expected 0, actual %v", result2.Position)
	}
}

func TestOption(t *testing.T) {
	result := Option(Token("A"))("A", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "A" {
		t.Errorf("expected A, actual %v", result.Target)
	}

	if result.Position != 1 {
		t.Errorf("expected 1, actual %v", result.Position)
	}

	result2 := Option(Token("A"))("B", 0)

	if !result2.Success {
		t.Errorf("expected true, actual %v", result2.Success)
	}

	if result2.Target != "" {
		t.Errorf("expected \"\", actual %v", result2.Target)
	}

	if result2.Position != 0 {
		t.Errorf("expected 0, actual %v", result2.Position)
	}

}

func TestChar(t *testing.T) {
	result := Char("abc")("c", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "c" {
		t.Errorf("expected c, actual %v", result.Target)
	}

	if result.Position != 1 {
		t.Errorf("expected 1, actual %v", result.Position)
	}

	result2 := Char("abc")("d", 0)

	if result2.Success {
		t.Errorf("expected false, actual %v", result2.Success)
	}

	if result2.Target != "" {
		t.Errorf("expected \"\", actual %v", result2.Target)
	}

	if result2.Position != 0 {
		t.Errorf("expected 0, actual %v", result2.Position)
	}

	result3 := Many(Char("a"))("aaabb", 0)

	if !result3.Success {
		t.Errorf("expected true, actual %v", result3.Success)
	}

	if result3.Target != "aaa" {
		t.Errorf("expected aaa, actual %v", result3.Target)
	}
}

func TestRegExp(t *testing.T) {
	result := RegExp(regexp.MustCompile("\\d+"))("a2014b333", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "2014" {
		t.Errorf("expected 2014, actual %v", result.Target)
	}

	if result.Position != 5 {
		t.Errorf("expected 5, actual %v", result.Position)
	}

	result2 := RegExp(regexp.MustCompile("\\d+"))("abcd", 0)

	if result2.Success {
		t.Errorf("expected false, actual %v", result2.Success)
	}

	if result2.Target != "" {
		t.Errorf("expected \"\", actual %v", result2.Target)
	}

	if result2.Position != 0 {
		t.Errorf("expected 0, actual %v", result2.Position)
	}

	result3 := RegExp(regexp.MustCompile("hoge"))("hoge", 5)

	if result3.Success {
		t.Errorf("expected false, actual %v", result3.Success)
	}
}

func TestLazy(t *testing.T) {
	fn := Token("hello")
	result := Lazy(&fn)("hello", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	var parser Parser
	parser = Option(Seq(Token("hoge"), Lazy(&parser)))

	result2 := parser("hogehogehoge", 0)

	if !result2.Success {
		t.Errorf("expected true, actual %v", result2.Success)
	}

	if result2.Target != "hogehogehoge" {
		t.Errorf("expected hogehogehoge, actual %v", result2.Target)
	}
}

func TestMap(t *testing.T) {
	result := Map(Token("hello"), func(str string) string {
		return str + ", world"
	})("hello", 0)

	if !result.Success {
		t.Errorf("expected true, actual %v", result.Success)
	}

	if result.Target != "hello, world" {
		t.Errorf("expected hello, world, actual %v", result.Target)
	}

	if result.Position != 5 {
		t.Errorf("expected 5, actual %v", result.Position)
	}

	result2 := Map(Token("coming"), func(str string) string {
		return str + " soon"
	})("hoge", 0)

	if result2.Success {
		t.Errorf("expected false, actual %v", result2.Success)
	}

	if result2.Target != "" {
		t.Errorf("expected \"\", actual %v", result2.Target)
	}

	if result2.Position != 0 {
		t.Errorf("expected 0, actual %v", result2.Position)
	}
}
