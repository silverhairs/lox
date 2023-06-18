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
	return &Scanner{
		Source:  Source,
		tokens:  make([]*token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) Tokenize() []*token.Token {

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()

	switch char {
	case '(':
		s.addTokenType(token.L_PAREN)
	case ')':
		s.addTokenType(token.R_PAREN)
	case '{':
		s.addTokenType(token.L_BRACE)
	case '}':
		s.addTokenType(token.R_BRACE)
	case ',':
		s.addTokenType(token.COMMA)
	case '.':
		s.addTokenType(token.DOT)
	case '-':
		s.addTokenType(token.MINUS)
	case '+':
		s.addTokenType(token.PLUS)
	case ';':
		s.addTokenType(token.SEMICOLON)
	case '*':
		s.addTokenType(token.ASTERISK)
	case '!':
		s.operator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{'=', token.BANG, token.BANG_EQ},
		)
	case '=':
		s.operator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{'=', token.EQUAL, token.EQ_EQ},
		)
	case '<':
		s.operator(
			struct {
				char     byte
				unique   token.TokenType
				twoChars token.TokenType
			}{'=', token.LESS, token.LESS_EQ},
		)
	case '>':
		s.operator(struct {
			char     byte
			unique   token.TokenType
			twoChars token.TokenType
		}{'=', token.GREATER, token.GREATER_EQ},
		)
	case '/':
		s.slash()
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
		} else if isAlpha(char) {
			s.identifier()
		} else {
			errors.Nowhere(s.line, fmt.Sprintf("Unexpected character %q", char))
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) advance() byte {
	prev := s.current
	s.current++
	return s.Source[prev]
}

func (s *Scanner) addTokenType(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	lexeme := s.Source[s.start:s.current]
	tok := token.New(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, tok)
}

func (s *Scanner) match(expect byte) bool {
	if s.isAtEnd() || s.Source[s.current] != expect {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) operator(props struct {
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

	s.addTokenType(tok)
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

func isAlpha(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func isAlphaNumeric(char byte) bool {
	return isAlpha(char) || isDigit(char)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	literal := s.Source[s.start:s.current]
	tok := token.LookupIdentifier(literal)

	s.addTokenType(tok)
}

func (s *Scanner) slash() {
	if s.match('/') {
		for s.peek() != '\n' && !s.isAtEnd() {
			s.advance()
		}
		literal := s.Source[s.start+2 : s.current]
		s.addToken(token.COMMENT_L, literal)

	} else if s.match('*') {
		for s.peek() != '*' && !s.isAtEnd() {
			s.advance()
		}

		if s.match('/') {
			literal := s.Source[s.start+2 : s.current-2]
			s.addToken(token.COMMENT_B, literal)
		} else {
			errors.Nowhere(s.line, "opened multi-line comment has not been closed.")
		}

	} else {
		s.addTokenType(token.SLASH)
	}
}
