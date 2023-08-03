package parser

import (
	"glox/ast"
	"glox/exception"
	"glox/token"
)

type Parser struct {
	tokens   []token.Token
	position int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, position: 0}
}

func (p *Parser) Parse() ([]ast.Statement, error) {
	return p.program()
}

func (p *Parser) program() ([]ast.Statement, error) {
	stmts := []ast.Statement{}
	var err error

	for !p.isAtEnd() {
		stmt, e := p.declaration()
		if e != nil {
			err = e
			break
		}
		stmts = append(stmts, stmt)
	}

	return stmts, err

}

func (p *Parser) declaration() (ast.Statement, error) {
	if p.match(token.IF) {
		return p.ifStatement()
	} else if p.match(token.LET) {
		return p.letDeclaration()
	}
	return p.statement()
}

func (p *Parser) ifStatement() (ast.Statement, error) {
	var err error
	_, err = p.consume(token.L_PAREN, "expected '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, e := p.expression()
	if e != nil {
		return nil, e
	}
	_, err = p.consume(token.R_PAREN, "expected ')' after if condition")
	if err != nil {
		return nil, err
	}

	then, e := p.statement()
	if e != nil {
		return nil, e
	}
	var orElse ast.Statement

	if p.match(token.ELSE) {
		orElse, err = p.statement()
	}

	return ast.NewIfStmt(condition, then, orElse), err

}

func (p *Parser) letDeclaration() (ast.Statement, error) {
	tok, err := p.consume(token.IDENTIFIER, "expected variable name.")
	if err != nil {
		return nil, err
	}

	var val ast.Expression
	if p.match(token.EQUAL) {
		val, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(token.SEMICOLON, "expect ';' after variable declaration."); err != nil {
		return nil, err
	}
	return ast.NewLetStmt(tok, val), err
}

func (p *Parser) statement() (ast.Statement, error) {

	if p.match(token.PRINT) {
		return p.printStatement()
	} else if p.match(token.L_BRACE) {
		block, err := p.block()
		return ast.NewBlockStmt(block), err
	}
	return p.expressionStatement()
}

func (p *Parser) block() ([]ast.Statement, error) {
	stmts := []ast.Statement{}
	var err error
	for !p.check(token.R_BRACE) && !p.isAtEnd() {
		stmt, e := p.declaration()
		if err != nil {
			return stmts, e
		}
		stmts = append(stmts, stmt)
	}
	_, err = p.consume(token.R_BRACE, "expect '}' after block.")
	return stmts, err
}

func (p *Parser) printStatement() (ast.Statement, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.SEMICOLON, "expect ';' after value."); err != nil {
		return nil, err
	}
	return ast.NewPrintStmt(exp), err

}

func (p *Parser) expressionStatement() (ast.Statement, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(token.SEMICOLON, "expect ';' after value."); err != nil {
		return nil, err
	}
	return ast.NewExprStmt(exp), err
}

func (p *Parser) expression() (ast.Expression, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expression, error) {
	exp, err := p.ternary()

	if p.match(token.EQUAL) {
		equals := p.previous()
		val, e := p.assignment()
		if e != nil {
			return exp, e
		}

		if variable, isVar := exp.(*ast.Variable); isVar {
			return ast.NewAssignment(variable.Name, val), err
		}

		err = exception.Runtime(equals, "invalid assignment target.")
	}

	return exp, err
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

		if _, e := p.consume(token.R_PAREN, "expected ')' after expression"); e != nil {
			err = e
		}
		return ast.NewGroupingExp(exp), err

	}
	if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
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
