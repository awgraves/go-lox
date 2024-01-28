package runtime

import (
	"fmt"

	"github.com/awgraves/go-lox/expressions"
	"github.com/awgraves/go-lox/statements"
	"github.com/awgraves/go-lox/tokens"
)

type parser struct {
	source      []*tokens.Token
	current     int
	errReporter ErrorReporter
}

func newParser(source []*tokens.Token, errReporter ErrorReporter) *parser {
	return &parser{
		source:      source,
		errReporter: errReporter,
	}
}

// parse will attempt to parse the tokens into statements.
// it is the caller's responsibility to check the err reporter as to whether this list of statements is usable.
func (p *parser) parse() []statements.Stmt {
	statements := []statements.Stmt{}

	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}
	return statements
}

func (p *parser) statement() statements.Stmt {
	if p.match(tokens.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *parser) printStatement() statements.Stmt {

	value := p.expression()
	p.consume(tokens.SEMICOLON, "Expect ';' after value.")
	return statements.PrintStmt{Expression: value}
}

func (p *parser) expressionStatement() statements.Stmt {
	expr := p.expression()
	p.consume(tokens.SEMICOLON, "Expect ';' after expression.")
	return statements.ExpStmt{Expression: expr}
}

func (p *parser) peek() tokens.Token {
	return *p.source[p.current]
}

func (p *parser) isAtEnd() bool {
	return p.peek().TokenType == tokens.EOF
}

func (p *parser) previous() tokens.Token {
	return *p.source[p.current-1]
}

func (p *parser) advance() tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) match(types ...tokens.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(t tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *parser) expression() expressions.Expression {
	return p.equality()
}

func (p *parser) equality() expressions.Expression {
	expr := p.comparison()

	for p.match(tokens.BANG_EQUAL, tokens.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) comparison() expressions.Expression {
	expr := p.term()

	for p.match(tokens.GREATER, tokens.GREATER_EQUAL, tokens.LESS, tokens.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) term() expressions.Expression {
	expr := p.factor()

	for p.match(tokens.MINUS, tokens.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) factor() expressions.Expression {
	expr := p.unary()

	for p.match(tokens.SLASH, tokens.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = expressions.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) unary() expressions.Expression {
	if p.match(tokens.BANG, tokens.MINUS) {
		operator := p.previous()
		right := p.unary()
		return expressions.Unary{Operator: operator, Right: right}
	}
	return p.primary()
}

func (p *parser) primary() expressions.Expression {
	if p.match(tokens.FALSE) {
		return expressions.Literal{Value: false} // TODO: better solution that interface
	}
	if p.match(tokens.TRUE) {
		return expressions.Literal{Value: true}
	}
	if p.match(tokens.NIL) {
		return expressions.Literal{Value: nil}
	}

	if p.match(tokens.NUMBER, tokens.STRING) {
		return expressions.Literal{Value: p.previous().Literal}
	}

	if p.match(tokens.LEFT_PAREN) {
		expr := p.expression()
		p.consume(tokens.RIGHT_PAREN, "Expect ')' after expression.")
		return expressions.Grouping{Expression: expr}
	}

	// TODO: revist this, not totally sure yet
	curr := p.peek()
	p.errReporter.AddError(curr.LineNum, 0, fmt.Sprintf("unknown primary expression: %s", curr.Lexeme))
	p.synchronize()

	// important to tell error reporter and avoid executing the expression tree.
	// otherwise might nil pointer dereference
	return nil
}

func (p *parser) consume(t tokens.TokenType, message string) {
	if p.check(t) {
		p.advance()
		return
	}

	// err handling begins
	curr := p.peek()

	// TODO: at char count number
	p.errReporter.AddError(curr.LineNum, 0, message)

	p.synchronize()
}

// synchronize moves the parser along to the next statement after an error was found
func (p *parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().TokenType == tokens.SEMICOLON {
			return
		}

		curr := p.peek()

		for _, t := range []tokens.TokenType{tokens.CLASS, tokens.FOR, tokens.FUN, tokens.IF, tokens.PRINT, tokens.RETURN, tokens.VAR, tokens.WHILE} {
			if t == curr.TokenType {
				return
			}
		}
		p.advance()
	}
}
