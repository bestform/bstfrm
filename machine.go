package bstfrm

import "fmt"

type Machine struct {
	ast *Ast
}

func NewMachine(ast *Ast) *Machine {
	return &Machine{ast: ast}
}

func (m *Machine) Run() {
	for _, stmt := range *m.ast.Statements {
		switch stmt.Kind {
		case PrintKind:
			runPrint(stmt)
		}
	}
}

func runPrint(stmt *Statement) {
	for _, s := range stmt.PrintStatement.Strings {
		fmt.Print(s)
	}
}


