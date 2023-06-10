package scanner

import (
	"craftinginterpreters/errors"
	"craftinginterpreters/token"
	"fmt"
	"strconv"
)

// Workaround to represent `nil` as a byte. Equivalent of `\0` in java.
const NULL = '#'

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
	case '+':
		s.recordToken(token.PLUS)
	case ';':
		s.recordToken(token.SEMICOLON)
	case '*':
		s.recordToken(token.ASTERISK)
	case '!':
		s.recordOperator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{char, token.BANG, token.BANG_EQ},
		)
	case '=':
		s.recordOperator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{char, token.EQUAL, token.EQ_EQ},
		)
	case '<':
		s.recordOperator(
			struct {
				char     byte
				unique   token.TokenType
				twoChars token.TokenType
			}{char, token.LESS, token.LESS_EQ},
		)
	case '>':
		s.recordOperator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{char, token.GREATER, token.GREATER_EQ},
		)
	case '/':
		// To handle comments
		if s.match(char) {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.recordToken(token.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(char) {
			s.number()
		} else {
			errors.Nowhere(s.line, fmt.Sprintf("unexpected character %v", char))
		}
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

func (s *Scanner) match(expect byte) bool {

	if !s.isAtEnd() {
		if s.Source[s.current] == expect {
			s.current++
			return true
		}
	}

	return false
}

func (s *Scanner) recordOperator(props struct {
	char     byte
	unique   token.TokenType // If the lexeme has only one character, which token type should be recorded.
	twoChars token.TokenType // If the lexeme has two characters, which token type should be recorded.
}) {
	var tok token.TokenType
	if s.match(props.char) {
		tok = props.twoChars
	} else {
		tok = props.unique
	}

	s.recordToken(tok)
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return NULL
	}
	return s.Source[s.current]
}

// Scans a string literal.
func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			// Multi-line string literals are allowed
			s.line++
			s.advance()
		}
	}

	if s.isAtEnd() {
		errors.Nowhere(s.line, "Please add a double-quote at the end of the string.")
		return
	}

	s.advance()

	value := s.Source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

// Scans number literals, this handles all floating-point numbers with or without decimals.
func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	literal := s.Source[s.start:s.current]
	value, err := strconv.ParseFloat(literal, 64)
	if err != nil {
		errors.Nowhere(s.line, fmt.Sprintf("%q is an invalid %q", literal, token.NUMBER))
		return
	}
	s.addToken(token.NUMBER, value)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.Source) {
		return NULL
	}

	return s.Source[s.current+1]
}
