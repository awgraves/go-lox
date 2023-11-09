package runtime

import (
	"fmt"

	"github.com/awgraves/go-lox/tokens"
)

type Scanner struct {
	source      []rune
	Tokens      []*tokens.Token
	start       int
	current     int
	line        int
	pos         int
	errReporter ErrorReporter
}

func newScanner(source string, errReporter ErrorReporter) *Scanner {
	return &Scanner{
		source:      []rune(source),
		Tokens:      []*tokens.Token{},
		line:        1,
		pos:         0,
		errReporter: errReporter,
	}
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.Tokens = append(s.Tokens, tokens.NewToken(tokens.EOF, "", s.line))
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(tokens.LEFT_PAREN)
		break
	case ')':
		s.addToken(tokens.RIGHT_PAREN)
		break
	case '{':
		s.addToken(tokens.LEFT_BRACE)
		break
	case '}':
		s.addToken(tokens.RIGHT_BRACE)
		break
	case ',':
		s.addToken(tokens.COMMA)
		break
	case '.':
		s.addToken(tokens.DOT)
		break
	case '-':
		s.addToken(tokens.MINUS)
		break
	case '+':
		s.addToken(tokens.PLUS)
		break
	case ';':
		s.addToken(tokens.SEMICOLON)
		break
	case '*':
		s.addToken(tokens.STAR)
		break
	case '!':
		if s.matc('=') {
			s.addToken(tokens.BANG_EQUAL)
			break
		}
		s.addToken(tokens.BANG)
		break
	case '=':
		if s.matc('=') {
			s.addToken(tokens.EQUAL_EQUAL)
			break
		}
		s.addToken(tokens.EQUAL)
		break
	case '<':
		if s.matc('=') {
			s.addToken(tokens.LESS_EQUAL)
			break
		}
		s.addToken(tokens.LESS)
		break
	case '>':
		if s.matc('=') {
			s.addToken(tokens.GREATER_EQUAL)
			break
		}
		s.addToken(tokens.GREATER)
		break
	case '/':
		if s.matc('/') {
			// this is a comment line
			// continue until new line or EOF
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
			break
		}
		s.addToken(tokens.SLASH)
		break
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		s.setNewLine()
		break

	case '"':
		s.handleString()
		break
	default:
		if s.isDigit(c) {
			s.handleNum()
			break
		}

		if s.isAlpha(c) {
			s.handleIdentifier()
			break
		}
		s.errReporter.AddError(s.line, s.pos, fmt.Sprintf("Unexpected character '%v'", string(c)))
	}
}

func (s *Scanner) setNewLine() {
	s.line++
	s.pos = 0
}

func (s *Scanner) handleIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := string(s.source[s.start:s.current])

	tt, ok := tokens.KeywordsMap[text]
	if ok {
		s.addToken(tt)
		return
	}

	s.addToken(tokens.IDENTIFIER)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) handleNum() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	s.addToken(tokens.NUMBER)
}

func (s *Scanner) handleString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.setNewLine()
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errReporter.AddError(
			s.line,
			s.pos,
			"Unterminated string.",
		)
		return
	}

	s.advance() // the closing "

	// trim surrounding quotes
	val := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithVal(tokens.STRING, val)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\r'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\n'
	}
	return s.source[s.current+1]
}

func (s *Scanner) matc(exp rune) bool {
	if s.isAtEnd() || s.source[s.current] != exp {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	s.pos++
	return c
}

func (s *Scanner) addTokenWithVal(tt tokens.TokenType, val string) {
	s.Tokens = append(s.Tokens, tokens.NewToken(tt, val, s.line))
}

func (s *Scanner) addToken(tt tokens.TokenType) {
	str := string(s.source[s.start:s.current])
	s.Tokens = append(s.Tokens, tokens.NewToken(tt, str, s.line))
}
