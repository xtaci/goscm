package goscm

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestParse(t *testing.T) {
	buf := bytes.NewBufferString("(+ (+ 1 3) (+ 2 4))")
	lexer := Lexer{}
	lexer.init(buf)
	parser := Parser{}
	parser.init(&lexer)
	parser.parse()
	spew.Dump(parser.forms)
}
