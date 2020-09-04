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
func calc(inputStr string) int {
	input := antlr.NewInputStream(inputStr)
	lexer := newLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewEarthParser(stream)

	errorListener := antlr.NewConsoleErrorListener()
	errorStrategy := antlr.NewBailErrorStrategy()

	p.AddErrorListener(errorListener)
	p.SetErrorHandler(errorStrategy)
	p.BuildParseTrees = true

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

////////////////////////

type lexer struct {
	*parser.EarthLexer
	prevIndentLevel                              int
	indentLevel                                  int
	afterNewLine                                 bool
	tokenQueue                                   []antlr.Token
	wsChannel, wsStart, wsStop, wsLine, wsColumn int
}

func newLexer(input antlr.CharStream) antlr.Lexer {
	l := new(lexer)
	l.EarthLexer = parser.NewEarthLexer(input)
	return l
}

func (l *lexer) NextToken() antlr.Token {
	peek := l.EarthLexer.NextToken()
	ret := peek
	switch peek.GetTokenType() {
	case parser.EarthLexerWS:
		if l.afterNewLine {
			l.indentLevel++
		}
		l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn =
			peek.GetChannel(), peek.GetStart(), peek.GetStop(), peek.GetLine(), peek.GetColumn()
	case parser.EarthLexerNL:
		l.indentLevel = 0
		l.afterNewLine = true
	default:
		if l.afterNewLine {
			if l.prevIndentLevel < l.indentLevel {
				l.tokenQueue = append(l.tokenQueue, l.GetTokenFactory().Create(
					l.GetTokenSourceCharStreamPair(), parser.EarthLexerINDENT, "",
					l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn))
			} else if l.prevIndentLevel > l.indentLevel {
				l.tokenQueue = append(l.tokenQueue, l.GetTokenFactory().Create(
					l.GetTokenSourceCharStreamPair(), parser.EarthLexerDEDENT, "",
					l.wsChannel, l.wsStart, l.wsStop, l.wsLine, l.wsColumn))
				l.PopMode() // Pop RECIPE mode.
			}
		}
		l.prevIndentLevel = l.indentLevel
		l.afterNewLine = false
	}
	if len(l.tokenQueue) > 0 {
		l.tokenQueue = append(l.tokenQueue, peek)
		ret = l.tokenQueue[0]
		l.tokenQueue = l.tokenQueue[1:]
	}
	return ret
}
