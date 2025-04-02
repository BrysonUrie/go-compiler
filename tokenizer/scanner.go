package scanner

import (
	"unicode"
)

type Scanner struct {
	Tokens  []Token
	Source  string
	Start   int
	Current int
	Line    int
}

func CreateScanner(text string) *Scanner {
	return &Scanner{
		Tokens:  []Token{},
		Source:  text,
		Start:   0,
		Current: 0,
		Line:    0,
	}
}

func (s *Scanner) Scan() []Token {
	for !s.IsAtEnd() {
		s.Start = s.Current
		s.ScanToken()
	}
	s.AddToken(EOF, "")
	return s.Tokens
}

func (s *Scanner) ScanToken() {
	c := s.Advance()
	switch c {
	// Single-character tokens.
	case '(':
		s.AddToken(LEFT_PAREN, "(")
	case ')':
		s.AddToken(RIGHT_PAREN, ")")
	case '{':
		s.AddToken(LEFT_BRACE, "{")
	case '}':
		s.AddToken(RIGHT_BRACE, "}")
	case ',':
		s.AddToken(COMMA, ",")
	case '.':
		s.AddToken(DOT, ".")
	case '-':
		s.AddToken(MINUS, "-")
	case '+':
		s.AddToken(PLUS, "+")
	case ';':
		s.AddToken(SEMICOLON, ";")
	case '^':
		s.AddToken(CARET, "^")
	// One or two character tokens.
	case '/':
		if s.Match('/') {
			for s.Peek() != '\n' && !s.IsAtEnd() {
				s.Advance()
			}
		} else {
			s.AddToken(SLASH, "/")
		}
	case '*':
		s.AddToken(STAR, "*")
	case '!':
		if s.Match('=') {
			s.AddToken(BANG_EQUAL, "!=")
		} else {
			s.AddToken(EQUAL, "=")
		}
	case '=':
		if s.Match('=') {
			s.AddToken(EQUAL_EQUAL, "==")
		} else {
			s.AddToken(EQUAL, "=")
		}
	case '>':
		if s.Match('=') {
			s.AddToken(GREATER_EQUAL, ">=")
		} else {
			s.AddToken(GREATER, ">")
		}
	case '<':
		if s.Match('=') {
			s.AddToken(LESS_EQUAL, "<=")
		} else {
			s.AddToken(LESS, "<")
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.Line++
	case '"':
		s.String()
	default:
		if unicode.IsDigit(c) {
			s.Digit()
		} else if unicode.IsLetter(c) {
			s.Identifier()
		}
	}
}

func (s *Scanner) Digit() {
	for unicode.IsDigit(s.Peek()) {
		s.Advance()
	}
	if s.Peek() == '.' && unicode.IsDigit(s.PeekNext()) {
		s.Advance()
		for unicode.IsDigit(s.Peek()) {
			s.Advance()
		}
	}
	number := s.Source[s.Start:s.Current]
	s.AddToken(NUMBER, number)
}
func (s *Scanner) Identifier() {
	for unicode.IsDigit(s.Peek()) || unicode.IsLetter(s.Peek()) {
		s.Advance()
	}

	str := s.Source[s.Start:s.Current]
	keyWord, exists := keywordMap[str]
	if exists {
		s.AddToken(keyWord, str)
	} else {
		s.AddToken(IDENTIFIER, str)
	}
}
func (s *Scanner) PeekNext() rune {
	if s.Current+1 >= len(s.Source) {
		return rune(0)
	}
	return rune(s.Source[s.Current+1])
}
func (s *Scanner) String() {
	for s.Peek() != '"' && !s.IsAtEnd() {
		if s.Peek() == '\n' {
			s.Line++
		}
		s.Advance()
	}

	if s.IsAtEnd() {
		panic("Unterminated string")
	}
	s.Advance()

	value := s.Source[s.Start+1 : s.Current]
	s.AddToken(STRING, value)
}

func (s *Scanner) Advance() rune {
	ret := rune(s.Source[s.Current])
	s.Current++
	return ret
}
func (s *Scanner) Match(expected rune) bool {
	if s.IsAtEnd() {
		return false
	}
	next := rune(s.Source[s.Current])
	if next != expected {
		return false
	}
	s.Current++
	return true
}
func (s *Scanner) Peek() rune {
	if s.IsAtEnd() {
		return rune(0)
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) IsAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) AddToken(tokenType TokenType, lexeme string) {
	token := CreateToken(tokenType, lexeme)
	s.Tokens = append(s.Tokens, *token)
}
