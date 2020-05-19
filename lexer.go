package bstfrm

import (
	"fmt"
	"sort"
	"strings"
)

type symbol byte

const (
	symbolSemicolon symbol = ';'
	symbolEqual symbol = '='
)

type keyword string

const (
	keywordPrint keyword = "print"
	keywordCalc keyword = "calc"
	keywordSet keyword = "set"
)

type tokenKind uint

const (
	stringKind tokenKind = iota
	symbolKind
	keywordKind
	numericKind
	identifierKind
)

type token struct {
	kind  tokenKind
	value string
}

func (t *token) equals(other *token) bool {
	return t.kind == other.kind && t.value == other.value
}

type position struct {
	col  uint
	line uint
}

type cursor struct {
	pointer uint
	pos     position
}

type lexer func(source string, cur cursor) (*token, cursor, bool)

func Lex(source string) ([]*token, error) {
	var cur cursor

	lexers := []lexer{symbolLexer, keywordLexer, stringLexer, numericLexer, identifierLexer}
	var tokens []*token

lex:
	for cur.pointer < uint(len(source)) {
		cur = eatWhitespace(source, cur)
		for _, lexer := range lexers {
			t, newCursor, ok := lexer(source, cur)
			if ok {
				tokens = append(tokens, t)
				cur = newCursor
				continue lex
			}
		}

		return nil, fmt.Errorf("error lexing at [%d:%d]", cur.pos.line, cur.pos.col)
	}

	return tokens, nil
}

func eatWhitespace(source string, cur cursor) cursor {
	for cur.pointer < uint(len(source)) && (source[cur.pointer] == ' ' || source[cur.pointer] == '\n') {
		cur.pointer++
		cur.pos.col++
		if source[cur.pointer] == '\n' {
			cur.pos.line++
			cur.pos.col = 0
		}
	}

	return cur
}

func symbolLexer(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	current := source[cur.pointer]

	for _, s := range []symbol{symbolSemicolon, symbolEqual} {
		if current == byte(s) {
			cur.pointer++
			cur.pos.col++
			return &token{
				kind: symbolKind,
				value: string(current),
			}, cur, true
		}
	}

	return nil, ic, false
}

func keywordLexer(source string, ic cursor) (*token, cursor, bool) {
	keywords := []keyword{keywordPrint, keywordCalc, keywordSet}
	cur := ic

	sort.Slice(keywords, func(a int, b int) bool {
		return len(string(keywords[a])) > len(string(keywords[b]))
	})

	for _, keyword := range keywords {
		if strings.HasPrefix(source[cur.pointer:], string(keyword)) {
			cur.pointer += uint(len(string(keyword)))
			cur.pos.col += uint(len(string(keyword)))
			return &token{
				kind:  keywordKind,
				value: string(keyword),
			}, cur, true
		}
	}

	return nil, ic, false
}

func stringLexer(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	current := source[cur.pointer]

	if current != '"' {
		return nil, ic, false
	}
	cur.pointer++

	var value []byte

	for cur.pointer < uint(len(source)) && source[cur.pointer] != '"' {
		value = append(value, source[cur.pointer])
		cur.pointer++
		cur.pos.col++
	}

	cur.pointer++
	cur.pos.col++

	return &token{
		kind:  stringKind,
		value: string(value),
	}, cur, true
}

func numericLexer(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	var value []byte

	for uint(len(source)) > cur.pointer &&  source[cur.pointer] >= '0' && source[cur.pointer] <= '9' {
		value = append(value, source[cur.pointer])
		cur.pointer++
		cur.pos.col++
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &token{
		kind:  numericKind,
		value: string(value),
	}, cur, true
}

func identifierLexer(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	if source[cur.pointer] != '#' {
		return nil, ic, false
	}
	cur.pointer++
	cur.pos.col++

	var value []byte
	for cur.pointer < uint(len(source)) && isAlphanumeric(source[cur.pointer]) {
		value = append(value, source[cur.pointer])
		cur.pointer++
		cur.pos.col++
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &token{
		kind:  identifierKind,
		value: string(value),
	}, cur, true
}

func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

