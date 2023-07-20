package parser

import (
	"glox/ast"
	"glox/exception"
	"glox/token"
)

type Parser struct {
	tokens []token.Token

	position int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, position: 0}
}

func (p *Parser) Parse() (ast.Expression, error) {
	return p.expression()
}

func (p *Parser) expression() (ast.Expression, error) {
	return p.ternary()
}

func (p *Parser) ternary() (ast.Expression, error) {
	exp, err := p.equality()

	for p.match(token.QUESTION_MARK) {
		left := p.previous()
		positive, e := p.equality()
		if e != nil {
			err = e
		}

		for p.match(token.COLON) {
			right := p.previous()
			negative, e := p.ternary()
			if e != nil {
				err = e
			}

			exp = ast.NewTernaryConditional(exp, left, positive, right, negative)
		}

	}

	return exp, err
}

func (p *Parser) equality() (ast.Expression, error) {
	exp, err := p.comparison()

	for p.match(token.BANG_EQ, token.EQ_EQ) {
		operator := p.previous()
		right, e := p.comparison()
		if err != nil {
			err = e
		}
		exp = ast.NewBinaryExpression(exp, operator, right)
	}
	return exp, err

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

func (p *Parser) comparison() (ast.Expression, error) {
	exp, err := p.term()

	for p.match(token.GREATER, token.GREATER_EQ, token.LESS, token.LESS_EQ) {
		operator := p.previous()
		right, e := p.term()
		if e != nil {
			err = e
		}
		exp = ast.NewBinaryExpression(exp, operator, right)
	}

	return exp, err
}

func (p *Parser) term() (ast.Expression, error) {
	exp, err := p.factor()
	if err != nil {
		return exp, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, e := p.factor()
		if e != nil {
			err = e
		}
		exp = ast.NewBinaryExpression(exp, operator, right)
	}

	return exp, err
}

func (p *Parser) factor() (ast.Expression, error) {
	exp, err := p.unary()
	if err != nil {
		return exp, err
	}

	for p.match(token.SLASH, token.ASTERISK) {
		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			err = e
		}
		exp = ast.NewBinaryExpression(exp, operator, right)
	}

	return exp, err
}

func (p *Parser) unary() (ast.Expression, error) {
	var err error

	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			err = e
		}
		return ast.NewUnaryExpression(operator, right), err
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expression, error) {
	if p.match(token.FALSE) {
		return ast.NewLiteralExpression(false), nil
	}
	if p.match(token.TRUE) {
		return ast.NewLiteralExpression(true), nil
	}
	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteralExpression(p.previous().Literal), nil
	}
	if p.match(token.L_PAREN) {
		exp, err := p.expression()
		if err != nil {
			return exp, err
		}

		_, err = p.consume(token.R_PAREN, "Expect ')' after expression")
		return ast.NewGroupingExp(exp), err

	}

	return nil, exception.Parse(p.peek())

}

func (p *Parser) consume(tokType token.TokenType, message string) (token.Token, error) {
	if p.check(tokType) {
		return p.advance(), nil
	}
	tok := p.peek()
	return tok, captureError(tok, message)

}

func captureError(tok token.Token, msg string) error {
	if tok.Type == token.EOF {
		return exception.Generic(tok.Line, " at end", msg)
	}

	return exception.Generic(tok.Line, "'"+tok.Lexeme+"'", msg)
}
