package bstfrm

import "fmt"

type StatementKind int32

const (
	PrintKind StatementKind = iota
)

type PrintStatement struct {
	Strings []string
}

type Statement struct {
	Kind StatementKind
	PrintStatement *PrintStatement
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
		stmt, newCursor, ok := parsePrintStatement(tokens, cursor)
		if ok {
			cursor = newCursor
			stmts = append(stmts, &Statement{
				Kind:           PrintKind,
				PrintStatement: stmt,
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

func parsePrintStatement(tokens []*token, ic uint) (*PrintStatement, uint, bool) {
	if !expectToken(tokens, ic, tokenFromKeyword(keywordPrint)) {
		return nil, ic, false
	}
	cur := ic
	cur++
	var values []string
	for cur < uint(len(tokens)) && !tokens[cur].equals(tokenFromSymbol(symbolSemicolon)) {
		if tokens[cur].kind != stringKind {
			fmt.Println("Expected string to print")
			return nil, ic, false
		}
		values = append(values, tokens[cur].value)
		cur++
	}
	cur++

	return &PrintStatement{
		Strings: values,
	}, cur, true
}
