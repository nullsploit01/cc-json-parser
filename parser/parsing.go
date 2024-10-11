package parser

import "fmt"

type Parser struct {
	lexer    *Lexer
	curToken Token
}

func NewParser(data string) *Parser {
	lexer := NewLexer(data)
	p := &Parser{lexer: lexer}
	p.curToken = p.lexer.NextToken() // Read the first token
	return p
}

func (p *Parser) Parse() error {
	if p.curToken.Type != TknLeftBrace {
		return fmt.Errorf("expected left brace, found %s", p.curToken.Literal)
	}

	p.curToken = p.lexer.NextToken()

	if p.curToken.Type != TknRightBrace {
		return fmt.Errorf("expected right brace, found %s", p.curToken.Literal)
	}

	return nil
}
