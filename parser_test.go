package lex

import (
	"regexp"
	"testing"
)

func TestIncompatible(t *testing.T) {
	result := incompatible(4)

	if result.Success {
		t.Error("必ず失敗しなければなりません")
	}

	if result.Target != "" {
		t.Error("必ず空文字を返さなければなりません")
	}

	if result.Position != 4 {
		t.Error("与えた位置を返さなければなりません")
	}
}

func TestToken(t *testing.T) {
	result := Token("hoge")("hoge", 0)
	if result.Success {
		if result.Target != "hoge" {
			t.Error("正確に文字列をパースできていません")
		}

		if result.Position != len("hoge") {
			t.Error("正確な長さがパースできていません")
		}
	} else {
		t.Error("正確にパースできていません")
	}

	// 解析文字が入力より長い場合
	result2 := Token("hogehoge")("hoge", 0)
	if result2.Success {
		t.Error("失敗しなければならない")
	}

	if result2.Target != "" {
		t.Error("空文字でなければならない")
	}

	if result2.Position != 0 {
		t.Error("読み出し位置を更新してはならない")
	}

	// incompatible
	result3 := Token("hoge")("fugafuga", 0)
	if result3.Success {
		t.Error("見つからない場合は失敗しなければならない")
	}

	if result3.Target != "" {
		t.Error("見つからない場合はから文字を返さなければならない")
	}

	if result3.Position != 0 {
		t.Error("見つからない場合は読み出し位置を更新しない")
	}
}

func TestMany(t *testing.T) {
	result := Many(Token("hoge"))("hogehoge", 0)
	t.Log(result)
	if !result.Success {
		t.Error("合致する場合は成功しなければならない")
	}

	if result.Target != "hogehoge" {
		t.Error("正確に文字列を回収しなければならない")
	}

	if result.Position != 8 {
		t.Error("正確な読み出し位置を返さなければならない")
	}

	result2 := Many(Token("hoge"))("hogehoge", 4)
	t.Log(result2)
	if !result2.Success {
		t.Error("合致する場合は成功しなければならない")
	}

	if result2.Target != "hoge" {
		t.Error("正確な範囲をとり出さなければならない")
	}

	if result2.Position != 8 {
		t.Error("正確な読み取り位置を更新しなければならない")
	}
}

func TestChoice(t *testing.T) {
	result := Choice(Token("foo"), Token("bar"))("bar", 0)

	if !result.Success {
		t.Error("一致するものがあればtrueを返す")
	}

	if result.Target != "bar" {
		t.Error("一致する箇所が正しくない")
	}

	if result.Position != 3 {
		t.Error("一致する読み出し位置が違う")
	}

	result2 := Choice(Token("foo"), Token("bar"))("foo", 0)

	if !result2.Success {
		t.Error("一致するものがあればtrueを返す")
	}

	if result2.Target != "foo" {
		t.Error("一致する箇所が正しくない")
	}

	if result2.Position != 3 {
		t.Error("一致する読み出し位置が違う")
	}

	result3 := Choice(Token("foo"), Token("bar"))("baz", 0)

	if result3.Success {
		t.Error("これは一致しない")
	}

	if result3.Target != "" {
		t.Error("から文字を返さないとおかしい")
	}

	if result3.Position != 0 {
		t.Error("0を返さないとおかしい")
	}
}

func TestSeq(t *testing.T) {
	result := Seq(Token("foo"), Token("bar"), Token("baz"))("foobarbaz", 0)

	if !result.Success {
		t.Error("成功しないとおかしい")
	}

	if result.Target != "foobarbaz" {
		t.Error("foobarbazを拾えないとおかしい")
	}

	if result.Position != len("foobarbaz") {
		t.Error("正しい長さをとれないとおかしい")
	}

	result2 := Seq(Token("foo"), Token("bar"), Token("baz"))("bazbarfoo", 0)
	t.Log(result2)

	if result2.Success {
		t.Error("成功するとおかしい")
	}

	if result2.Target != "" {
		t.Error("から文字でなければおかしい")
	}

	if result2.Position != 0 {
		t.Error("0でなければおかしい")
	}
}

func TestOption(t *testing.T) {
	result := Option(Token("A"))("A", 0)

	if !result.Success {
		t.Error("成功しないとおかしい")
	}

	if result.Target != "A" {
		t.Error("Aじゃないとおかしい")
	}

	if result.Position != 1 {
		t.Error("1じゃないとおかしい")
	}

	result2 := Option(Token("A"))("B", 0)

	if !result2.Success {
		t.Error("成功しないとおかしい")
	}

	if result2.Target != "" {
		t.Error("空文字じゃないとおかしい")
	}

	if result2.Position != 0 {
		t.Error("0じゃないとおかしい")
	}

}

func TestChar(t *testing.T) {
	result := Char("abc")("c", 0)

	if !result.Success {
		t.Error("成功しないとおかしい")
	}

	if result.Target != "c" {
		t.Error("cじゃないとおかしい")
	}

	if result.Position != 1 {
		t.Error("1じゃないとおかしい")
	}

	result2 := Char("abc")("d", 0)

	if result2.Success {
		t.Error("成功するとおかしい")
	}

	if result2.Target != "" {
		t.Error("空文字じゃないとおかしい")
	}

	if result2.Position != 0 {
		t.Error("0じゃないとおかしい")
	}
}

func TestRegExp(t *testing.T) {
	result := RegExp(regexp.MustCompile("\\d+"))("a2014b333", 0)

	if !result.Success {
		t.Error("成功しないとおかしい")
	}

	if result.Target != "2014" {
		t.Error("2014を抜き出せないとおかしい")
	}

	if result.Position != 5 {
		t.Error("5文字目じゃないとおかしい")
	}

	result2 := RegExp(regexp.MustCompile("\\d+"))("abcd", 0)

	if result2.Success {
		t.Error("失敗しないとおかしい")
	}

	if result2.Target != "" {
		t.Error("空文字を返さないとおかしい")
	}

	if result2.Position != 0 {
		t.Error("0を返さないとおかしい")
	}
}

func TestLazy(t *testing.T) {
	result := Lazy(Token("hello"))("hello", 0)

	if !result.Success {
		t.Error("成功しないとおかしい")
	}
}

func TestMap(t *testing.T) {
	result := Map(Token("hello"), func(str string) string {
		return str + ", world"
	})("hello", 0)

	if !result.Success {
		t.Error("成功しないとおかしい")
	}

	if result.Target != "hello, world" {
		t.Error("文字を加工できていないとおかしい")
	}

	if result.Position != 5 {
		t.Error("5じゃないとおかしい")
	}

	result2 := Map(Token("coming"), func(str string) string {
		return str + " soon"
	})("hoge", 0)

	if result2.Success {
		t.Error("成功するとおかしい")
	}

	if result2.Target != "" {
		t.Error("空文字じゃないとおかしい")
	}

	if result2.Position != 0 {
		t.Error("0じゃないとおかしい")
	}
}
