package scanner

import (
	"craftinginterpreters/errors"
	"craftinginterpreters/token"
	"fmt"
)

type Scanner struct {
	Source  string
	tokens  []*token.Token
	start   int
	current int
	line    int
}

func New(Source string) *Scanner {
	scnr := &Scanner{
		Source:  Source,
		tokens:  make([]*token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}

	return scnr
}

func (s *Scanner) Tokenize() []*token.Token {
	for !s.isAtEnd() {
		s.scanToken()
		s.start = s.current
		s.Tokenize()
	}

	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()

	switch char {
	case '(':
		s.recordToken(token.L_PAREN)
	case ')':
		s.recordToken(token.R_PAREN)
	case '{':
		s.recordToken(token.L_BRACE)
	case '}':
		s.recordToken(token.R_BRACE)
	case ',':
		s.recordToken(token.COMMA)
	case '.':
		s.recordToken(token.DOT)
	case '-':
		s.recordToken(token.MINUS)
	case '/':
		s.recordToken(token.SLASH)
	case '+':
		s.recordToken(token.PLUS)
	case ';':
		s.recordToken(token.SEMICOLON)
	case '*':
		s.recordToken(token.ASTERISK)
	default:
		errors.Report(s.line, "", fmt.Sprintf("unexpected character %v", char))
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) advance() byte {
	s.current += 1
	return s.Source[s.current]
}

func (s *Scanner) recordToken(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	lexeme := s.Source[s.start:s.current]
	tok := token.New(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, tok)
}
