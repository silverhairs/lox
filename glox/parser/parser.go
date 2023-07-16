package parser

import (
	"glox/ast"
	"glox/exception"
	"glox/token"
)

type Parser struct {
	tokens   []token.Token
	errors   []error
	position int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, position: 0}
}

func (p *Parser) Parse() (ast.Expression, []error) {
	return p.expression(), p.errors
}

func (p *Parser) expression() ast.Expression {
	return p.ternary()
}

func (p *Parser) ternary() ast.Expression {
	exp := p.equality()

	for p.match(token.QUESTION_MARK) {
		left := p.previous()
		positive := p.equality()

		for p.match(token.COLON) {
			right := p.previous()
			negative := p.ternary()

			exp = ast.NewTernaryConditional(exp, left, positive, right, negative)
		}

	}

	return exp
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

		p.consume(token.R_PAREN, "Expect ')' after expression")
		return ast.NewGroupingExp(exp)

	}

	panic(captureError(p.peek(), "Expected expression"))

}

func (p *Parser) consume(tokType token.TokenType, message string) token.Token {
	if p.check(tokType) {
		return p.advance()
	}
	tok := p.peek()
	err := captureError(tok, message)
	//TODO: Maybe later capure the error instead of panicking.
	// p.errors = append(p.errors, err.Error())
	panic(err)

}

func captureError(tok token.Token, msg string) error {
	if tok.Type == token.EOF {
		return exception.Generic(tok.Line, " at end", msg)
	}

	return exception.Generic(tok.Line, "'"+tok.Lexeme+"'", msg)
}

// Discards tokens that might case cascaded errors
// func (p *Parser) synchronize() {
// 	p.advance()

// 	for !p.isAtEnd() {
// 		if p.previous().Type == token.SEMICOLON {
// 			return
// 		}

// 		for _, stmt := range statements {
// 			if p.peek().Type == stmt {
// 				return
// 			}
// 		}
// 		p.advance()
// 	}

// }

// var statements = []token.TokenType{
// 	token.CLASS,
// 	token.FUNCTION,
// 	token.LET,
// 	token.FOR,
// 	token.IF,
// 	token.WHILE,
// 	token.PRINT,
// 	token.RETURN,
// }
