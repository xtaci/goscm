package goscm

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"unicode"
)

// scope definiton
type scope struct {
	vars   map[string]interface{}
	parent *scope
}

func (s *scope) find(symbol string) *scope {
	if _, ok := s.vars[symbol]; ok {
		return s
	} else {
		return s.parent.find(symbol)
	}
}

var global scope

// proc
type proc struct {
	params interface{}
	body   interface{}
	scope  *scope
}

// token definition
const (
	TK_UNKNOWN    = iota
	TK_FORM_BEGIN // (
	TK_FORM_END   // )
	TK_INTEGER
	TK_FLOAT
	TK_STRING
	TK_BOOL // #f #t
	TK_SYMBOL
)

var chars = map[rune]bool{
	'[': true,
	'!': true,
	'$': true,
	'%': true,
	'&': true,
	'*': true,
	'+': true,
	'-': true,
	'/': true,
	':': true,
	'<': true,
	'>': true,
	'=': true,
	'?': true,
	'@': true,
	'^': true,
	'_': true,
	'~': true,
	']': true,
}

type token struct {
	typ     int
	i       int
	f       float64
	s       string
	b       bool
	literal string
}

type Lexer struct {
	reader *bytes.Buffer
	lines  []string
	lineno int
}

func (lex *Lexer) init(r io.Reader) {
	bts, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	lex.reader = bytes.NewBuffer(bts)
	lex.lineno = 1
}

func (lex *Lexer) next() (t *token) {
	defer func() {
		//log.Println(t)
	}()
	var r rune
	var err error
	for {
		r, _, err = lex.reader.ReadRune()
		if err == io.EOF {
			return &token{typ: TK_UNKNOWN}
		} else if unicode.IsSpace(r) {
			if r == '\n' {
				lex.lineno++
			}
			continue
		}
		break
	}

	if unicode.IsLetter(r) || chars[r] {
		var runes []rune
		for {
			runes = append(runes, r)
			r, _, err = lex.reader.ReadRune()
			if err == io.EOF {
				break
			} else if unicode.IsLetter(r) || unicode.IsNumber(r) || chars[r] {
				continue
			} else {
				lex.reader.UnreadRune()
				break
			}
		}

		t := &token{}
		t.literal = string(runes)
		t.typ = TK_SYMBOL
		return t
	} else if r == '"' { // quoted string
		var runes []rune
		for {
			r, _, err = lex.reader.ReadRune()
			if err == io.EOF {
				break
			} else if r != '"' { // read until '"'
				runes = append(runes, r)
				continue
			} else {
				break
			}
		}
		t := &token{}
		t.literal = string(runes)
		t.typ = TK_STRING
		return t
	} else if unicode.IsDigit(r) {
		var runes []rune
		for {
			runes = append(runes, r)
			r, _, err = lex.reader.ReadRune()
			if err == io.EOF {
				break
			} else if unicode.IsDigit(r) {
				continue
			} else {
				lex.reader.UnreadRune()
				break
			}
		}

		t := &token{}
		t.i, _ = strconv.Atoi(string(runes))
		t.typ = TK_INTEGER
		return t
	} else if r == '(' {
		return &token{typ: TK_FORM_BEGIN}
	} else if r == ')' {
		return &token{typ: TK_FORM_END}
	} else {
		return &token{typ: TK_UNKNOWN}
	}
}

//////////////////////////////////////////////////////////////
type form struct {
	list []interface{}
}

type Parser struct {
	lexer *Lexer
	forms []*form
}

func (p *Parser) init(lex *Lexer) {
	p.lexer = lex
}
func (p *Parser) match(typ int) *token {
	t := p.lexer.next()
	if t.typ != typ {
		panic("syntax error")
	}
	return t
}

func (p *Parser) parse() {
	for tk := p.lexer.next(); tk.typ != TK_UNKNOWN; tk = p.lexer.next() {
		switch tk.typ {
		case TK_FORM_BEGIN:
			p.forms = append(p.forms, p.parseform())
		default:
			panic("syntax error")
		}
	}
}

func (p *Parser) parseform() *form {
	var f form
	for tk := p.lexer.next(); tk.typ != TK_UNKNOWN; tk = p.lexer.next() {
		switch tk.typ {
		case TK_FORM_BEGIN:
			f.list = append(f.list, p.parseform())
		case TK_FORM_END:
			return &f
		case TK_INTEGER,
			TK_FLOAT,
			TK_STRING,
			TK_SYMBOL,
			TK_BOOL:
			f.list = append(f.list, tk)
		default:
			panic("syntax error")
		}
	}

	return &f
}

func eval(f *form) interface{} {
	var partial []interface{}
	for _, v := range f.list {
		partial := append(partial, evalv(v))
	}
	return apply(partial)
}

// apply
func apply(partial []interface{}) interface{} {
}

func evalv(op interface{}) interface{} {
	switch op.(type) {
	case (*form):
	case (*token):
		t := op.(*token)
		switch t.typ {
		case '+':
		case '-':
		case '*':
		case '/':
		}
	}
}
