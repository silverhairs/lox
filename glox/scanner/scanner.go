package scanner

import "craftinginterpreters/token"

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

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.ScanTokens()
	}

	s.tokens = append(s.tokens, token.New(token.EOF, "EOF", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}
