package bstfrm

import "fmt"

type Machine struct {
	variables map[string]string
}

func NewMachine() *Machine {
	return &Machine{variables: make(map[string]string)}
}

func (m *Machine) Run(ast *Ast) {
	for _, stmt := range *ast.Statements {
		switch stmt.Kind {
		case PrintKind:
			m.print(stmt)
		case SetKind:
			m.set(stmt.SetStatement.Name, stmt.SetStatement.Value)
		case CalcKind:
			result, err := stmt.CalcStatement.Expression.Calc()
			if err != nil {
				fmt.Println(err)
			} else {
				//fmt.Println(stmt.CalcStatement)
				fmt.Println(result)
			}
		}
	}
}

func (m *Machine) print(stmt *Statement) {
	for _, t := range stmt.PrintStatement.Tokens {
		switch t.kind{
		case stringKind:
			fmt.Print(t.value)
		case identifierKind:
			value, ok := m.variables[t.value]
			if !ok {
				fmt.Print("## unidentified variable:", t.value, "##")
			} else {
				fmt.Print(value)
			}
		}
	}
}

func (m *Machine) set(name string, value string) {
	m.variables[name] = value
}


