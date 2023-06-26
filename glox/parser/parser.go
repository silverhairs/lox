package parser

import (
	"glox/ast"
	"glox/token"
)

const (
	_          int = iota
	EXPRESSION     // -> equality
	EQUALITY       // Operators: [==, !=] ;  Associates: `L`
	COMPARISON     // Operators: [>, >= <=, <] ; Associates: `L`
	TERM           // Operators: [+, -] ; Associates: `L`
	FACTOR         // Opetators: [/, *] ; Associates : `L`
	UNARY          // Operators: [!, -] ; Associates: `R`
	PRIMARY        // Matches literals and grouping expressions
)

type Parser struct {
	tokens   []token.Token
	position int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, position: 0}
}

func (p *Parser) expression() ast.Expression {
	return p.equality()
}

func (p *Parser) equality() ast.Expression {
	exp := p.comparison()

	for p.match(token.BANG_EQ, token.EQ_EQ) {
		operator := p.previous()
		right := p.comparison()

		exp = ast.NewBinaryExpression(exp, operator, right)
	}
	return exp

}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, tokType := range types {
		if p.check(tokType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokType
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.position++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.position]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.position-1]
}

func (p *Parser) comparison() ast.Expression {
	exp := p.term()

	for p.match(token.GREATER, token.GREATER_EQ, token.LESS, token.LESS_EQ) {
		operator := p.previous()
		right := p.term()
		exp = ast.NewBinaryExpression(exp, operator, right)
	}

	return exp
}

func (p *Parser) term() ast.Expression {
	exp := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		exp = ast.NewBinaryExpression(exp, operator, right)
	}

	return exp
}

func (p *Parser) factor() ast.Expression {
	exp := p.unary()

	for p.match(token.SLASH, token.ASTERISK) {
		operator := p.previous()
		right := p.unary()
		exp = ast.NewBinaryExpression(exp, operator, right)
	}

	return exp
}

func (p *Parser) unary() ast.Expression {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnaryExpression(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expression {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpression(false)
	}
	if p.match(token.TRUE) {
		return ast.NewLiteralExpression(true)
	}
	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteralExpression(p.previous().Literal)
	}
	if p.match(token.L_PAREN) {
		exp := p.expression()

		p.consume(token.R_PAREN, "Expect ')' after expression.")
		return ast.NewGroupingExp(exp)
	}

	p.consume(p.peek().Type, "could not parse token")
	return p.expression()

}
