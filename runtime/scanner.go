package runtime

import (
	"fmt"
	"strconv"

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
	s.Tokens = append(s.Tokens, tokens.NewToken(tokens.EOF, "", nil, s.line))
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(tokens.LEFT_PAREN, nil)
		break
	case ')':
		s.addToken(tokens.RIGHT_PAREN, nil)
		break
	case '{':
		s.addToken(tokens.LEFT_BRACE, nil)
		break
	case '}':
		s.addToken(tokens.RIGHT_BRACE, nil)
		break
	case ',':
		s.addToken(tokens.COMMA, nil)
		break
	case '.':
		s.addToken(tokens.DOT, nil)
		break
	case '-':
		s.addToken(tokens.MINUS, nil)
		break
	case '+':
		s.addToken(tokens.PLUS, nil)
		break
	case ';':
		s.addToken(tokens.SEMICOLON, nil)
		break
	case '*':
		s.addToken(tokens.STAR, nil)
		break
	case '!':
		if s.match('=') {
			s.addToken(tokens.BANG_EQUAL, nil)
			break
		}
		s.addToken(tokens.BANG, nil)
		break
	case '=':
		if s.match('=') {
			s.addToken(tokens.EQUAL_EQUAL, nil)
			break
		}
		s.addToken(tokens.EQUAL, nil)
		break
	case '<':
		if s.match('=') {
			s.addToken(tokens.LESS_EQUAL, nil)
			break
		}
		s.addToken(tokens.LESS, nil)
		break
	case '>':
		if s.match('=') {
			s.addToken(tokens.GREATER_EQUAL, nil)
			break
		}
		s.addToken(tokens.GREATER, nil)
		break
	case '/':
		if s.match('/') {
			// this is a comment line
			// continue until new line or EOF
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
			break
		}
		if s.match('*') {
			s.advance()
			s.handleMultiLineComment()
			break
		}
		s.addToken(tokens.SLASH, nil)
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

func (s *Scanner) handleMultiLineComment() {
	startLine := s.line
	startPos := s.pos

	for {
		if s.isAtEnd() {
			s.errReporter.AddError(startLine, startPos, "Unterminated multi-line comment")
			break
		}

		next := s.peek()
		if next == '\n' {
			s.setNewLine()
			s.advance()
			continue
		}

		if next == '*' && s.peekNext() == '/' {
			s.advance()
			s.advance()
			break
		}

		s.advance()
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
		s.addToken(tt, nil)
		return
	}

	s.addToken(tokens.IDENTIFIER, nil)
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

	num, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 32)
	// TODO: come back to this
	if err != nil {
		panic("unable to parse what should be a number")
	}

	s.addToken(tokens.NUMBER, num)
}

func (s *Scanner) handleString() {
	strStartpos := s.pos
	strStartline := s.line

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.setNewLine()
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errReporter.AddError(
			strStartline,
			strStartpos,
			"Unterminated string.",
		)
		return
	}

	s.advance() // the closing "

	// trim surrounding quotes
	val := string(s.source[s.start+1 : s.current-1])
	s.addToken(tokens.STRING, val)
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

func (s *Scanner) match(exp rune) bool {
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

func (s *Scanner) addToken(tt tokens.TokenType, literal interface{}) {
	str := string(s.source[s.start:s.current])
	s.Tokens = append(s.Tokens, tokens.NewToken(tt, str, literal, s.line))
}
