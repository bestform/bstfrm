package bstfrm

import (
	"errors"
	"fmt"
	"strconv"
)

type StatementKind int32

const (
	PrintKind StatementKind = iota
	SetKind
	CalcKind
)

type ExpressionType string

const (
	Times ExpressionType = "TIMES"
	Divide ExpressionType = "DIVIDE"
	Plus ExpressionType = "PLUS"
	Minus ExpressionType = "MINUS"
)

type Expression struct {
	Value string
	Left  *Expression
	Right *Expression
	Type  ExpressionType
}

func (e *Expression) Calc() (int64, error) {
	if e.Value != "" {
		i, err := strconv.ParseInt(e.Value, 10, 32)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	left, err := e.Left.Calc()
	if err != nil {
		return 0, err
	}
	if e.Right == nil {
		return left, nil
	}
	right, err := e.Right.Calc()
	if err != nil {
		return 0, err
	}

	switch e.Type {
	case Times:
		return left * right, nil
	case Divide:
		return left / right, nil
	case Plus:
		return left + right, nil
	case Minus:
		return left - right, nil
	}

	return 0, errors.New("wat")
}

func (e *Expression) String() string {
	str := ""
	str += "["
	str += "Value:"
	str += e.Value
	str += ","
	str += "Type:"
	str += string(e.Type)
	str += ","
	str += "Left:"
	if e.Left != nil {
		str += e.Left.String()
	}
	str += ","
	str += "Right:"
	if e.Right != nil {
		str += e.Right.String()
	}
	str += "]"

	return str
}

type CalcStatement struct {
	Expression *Expression
}

func (c *CalcStatement) String() string {
	return c.Expression.String()
}

type PrintStatement struct {
	Tokens []*token
}

type SetStatement struct {
	Name  string
	Value string
}

type Statement struct {
	Kind           StatementKind
	PrintStatement *PrintStatement
	SetStatement   *SetStatement
	CalcStatement  *CalcStatement
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
				Kind:         SetKind,
				SetStatement: setStmt,
			})
			continue
		}
		calcStmt, newCursor, ok, err := parseCalcStatement(tokens, cursor)
		if err != nil {
			return nil, err
		}
		if ok {
			cursor = newCursor
			stmts = append(stmts, &Statement{
				Kind:          CalcKind,
				CalcStatement: calcStmt,
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

	if cur == uint(len(tokens)) || !expectToken(tokens, cur, tokenFromSymbol(symbolEqual)) {
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

func parseCalcStatement(tokens []*token, ic uint) (*CalcStatement, uint, bool, error) {
	if !expectToken(tokens, ic, tokenFromKeyword(keywordCalc)) {
		return nil, ic, false, nil
	}
	cur := ic
	cur++

	stmt := &CalcStatement{}
	var newCursor uint
	var err error
	stmt.Expression, newCursor, err = parseExpression(tokens, cur)
	if err != nil {
		return nil, ic, false, err
	}
	cur = newCursor

	return stmt, cur, true, nil
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

func parseExpression(tokens []*token, ic uint) (*Expression, uint, error) {
	//fmt.Println("parse new expression")
	//fmt.Println("Starting with", tokens[ic].value)
	exp := &Expression{}
	cur := ic

	if cur == uint(len(tokens)) {
		return exp, ic, errors.New("incomplete expression")
	}

	var err error
	var newCursor uint
	for cur != uint(len(tokens)) && !tokens[cur].equals(tokenFromSymbol(symbolSemicolon)) {
		switch tokens[cur].kind {
		case numericKind:
			//fmt.Println("found numeric", tokens[cur].value)
			exp.Left = &Expression{Value: tokens[cur].value}
			cur++
		case symbolKind:
			//fmt.Println("found symbol")
			if tokens[cur].equals(tokenFromSymbol(symbolLeftParen)) {
				cur++
				exp.Left, newCursor, err = parseExpression(tokens, cur)
				if err != nil {
					return exp, ic, err
				}
				cur = newCursor
				break
			}
			if tokens[cur].equals(tokenFromSymbol(symbolRightParen)) {
				cur++
				return exp, cur, nil
			}
		default:
			return exp, ic, errors.New("invalid expression: Expected number or left paren")
		}

		if cur >= uint(len(tokens)) {
			return exp, cur, nil
		}

		switch tokens[cur].value {
		case string(symbolPlus):
			//fmt.Println("Found +")
			exp.Type = Plus
		case string(symbolMinus):
			//fmt.Println("Found -")
			exp.Type = Minus
		case string(symbolDivide):
			//fmt.Println("Found /")
			exp.Type = Divide
		case string(symbolTimes):
			//fmt.Println("Found *")
			exp.Type = Times
		case string(symbolSemicolon):
			//fmt.Println("Found ; Done!")
			cur++
			return exp, cur, nil
		case string(symbolRightParen):
			//fmt.Println("Found ) Done!")
			cur++
			return exp, cur, nil
		default:
			return exp, ic, errors.New("invalid expression. Expected +/-/*//, bot got " + tokens[cur].value)
		}
		cur++

		exp.Right, newCursor, err = parseExpression(tokens, cur)
		if err != nil {
			return exp, ic, err
		}
		cur = newCursor

		return exp, cur, nil
	}

	return exp, cur, nil
}
