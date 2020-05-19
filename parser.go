package bstfrm

import (
	"errors"
	"fmt"
)

type StatementKind int32

const (
	PrintKind StatementKind = iota
	SetKind
)

type PrintStatement struct {
	Tokens []*token
}

type SetStatement struct {
	Name string
	Value string
}

type Statement struct {
	Kind StatementKind
	PrintStatement *PrintStatement
	SetStatement *SetStatement
}

type Ast struct {
	Statements *[]*Statement
}

func Parse(source string) (*Ast, error) {
	tokens, err := Lex(source)
	if err != nil {
		return nil, err
	}

	cursor := uint(0)
	ast := Ast{}
	stmts, err := parseStatements(tokens, cursor)
	if err != nil {
		return nil, err
	}
	ast.Statements = stmts

	return &ast, nil
}

func parseStatements(tokens []*token, cursor uint) (*[]*Statement, error) {
	var stmts []*Statement

	for cursor < uint(len(tokens)) {
		printStmt, newCursor, ok := parsePrintStatement(tokens, cursor)
		if ok {
			cursor = newCursor
			stmts = append(stmts, &Statement{
				Kind:           PrintKind,
				PrintStatement: printStmt,
			})
			continue
		}
		setStmt, newCursor, ok, err := parseSetStatement(tokens, cursor)
		if err != nil {
			return nil, err
		}
		if ok {
			cursor = newCursor
			stmts = append(stmts, &Statement{
				Kind:           SetKind,
				SetStatement: setStmt,
			})
			continue
		}
	}

	return &stmts, nil
}

func expectToken(tokens []*token, cursor uint, token *token) bool {
	return tokens[cursor].equals(token)
}

func tokenFromKeyword(k keyword) *token {
	return &token{
		kind:  keywordKind,
		value: string(k),
	}
}

func tokenFromSymbol(s symbol) *token {
	return &token{
		kind:  symbolKind,
		value: string(s),
	}
}

func parseSetStatement(tokens []*token, ic uint) (*SetStatement, uint, bool, error) {
	if !expectToken(tokens, ic, tokenFromKeyword(keywordSet)) {
		return nil, ic, false, nil
	}
	cur := ic
	cur++

	if cur == uint(len(tokens)) || tokens[cur].kind != identifierKind {
		return nil, ic, false, errors.New("expected identifier")
	}

	name := tokens[cur].value
	cur++

	if cur == uint(len(tokens)) || !expectToken(tokens, cur, tokenFromSymbol(symbolEqual)){
		return nil, ic, false, errors.New("expected equal sign")
	}
	cur++

	if cur == uint(len(tokens)) || tokens[cur].kind != stringKind {
		return nil, ic, false, errors.New("only string values are supported currently")
	}
	value := tokens[cur].value
	cur++

	if cur == uint(len(tokens)) || !expectToken(tokens, cur, tokenFromSymbol(symbolSemicolon)) {
		return nil, ic, false, errors.New("set statements must end in a semicolon")
	}
	cur++

	return &SetStatement{
		Name:  name,
		Value: value,
	}, cur, true, nil
}

func parsePrintStatement(tokens []*token, ic uint) (*PrintStatement, uint, bool) {
	if !expectToken(tokens, ic, tokenFromKeyword(keywordPrint)) {
		return nil, ic, false
	}
	cur := ic
	cur++
	var tokensToPrint []*token
	for cur < uint(len(tokens)) && !tokens[cur].equals(tokenFromSymbol(symbolSemicolon)) {
		if tokens[cur].kind != stringKind && tokens[cur].kind != identifierKind {
			fmt.Println("Expected string or identifier to print")
			return nil, ic, false
		}
		tokensToPrint = append(tokensToPrint, tokens[cur])
		cur++
	}
	cur++

	return &PrintStatement{
		Tokens: tokensToPrint,
	}, cur, true
}
