package main

import (
	"fmt"
	"io/ioutil"

	"github.com/alexcb/antlrcalc/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type calcListener struct {
	*parser.BaseEarthParserListener

	stack []int
}

func (l *calcListener) ExitFromStmt(c *parser.FromStmtContext) {
	fmt.Printf("ExitFromStmt %v\n", c.GetText())
}

// calc takes a string expression and returns the evaluated result.
func calc(input string) int {
	// Setup the input
	is := antlr.NewInputStream(input)

	// Create the Lexer
	lexer := parser.NewEarthLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewEarthParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener calcListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, p.EarthFile())

	return 0

}

func main() {

	s, err := ioutil.ReadFile("test")
	if err != nil {
		panic(err)
	}

	fmt.Printf("got %s\n", string(s))

	calc(string(s))
}
