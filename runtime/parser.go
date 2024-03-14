package runtime

import (
	"errors"
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
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *parser) declaration() statements.Stmt {
	if p.match(tokens.FUN) {
		return p.function("function")
	}
	if p.match(tokens.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *parser) varDeclaration() statements.Stmt {
	_, err := p.consume(tokens.IDENTIFIER, "Expect variable name.")
	if err != nil {
		// TODO: better handle
		panic(err)
	}

	name := p.previous()

	var initializer expressions.Expression = nil

	if p.match(tokens.EQUAL) {
		initializer = p.expression()
	}

	p.consume(tokens.SEMICOLON, "Expect ';' after variable declaration.")

	return statements.VarStmt{Name: name, Initializer: initializer}
}

func (p *parser) whileStatement() statements.Stmt {
	p.consume(tokens.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(tokens.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return statements.WhileStmt{Condition: condition, Body: body}
}

func (p *parser) statement() statements.Stmt {
	if p.match(tokens.FOR) {
		return p.forStatement()
	}
	if p.match(tokens.IF) {
		return p.ifStatement()
	}
	if p.match(tokens.PRINT) {
		return p.printStatement()
	}
	if p.match(tokens.WHILE) {
		return p.whileStatement()
	}
	if p.match(tokens.LEFT_BRACE) {
		return statements.Block{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *parser) forStatement() statements.Stmt {
	p.consume(tokens.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer statements.Stmt
	if p.match(tokens.SEMICOLON) {
		initializer = nil
	} else if p.match(tokens.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition expressions.Expression = nil
	if !p.check(tokens.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(tokens.SEMICOLON, "Expect ';' after loop condition")

	var increment expressions.Expression = nil
	if !p.check(tokens.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(tokens.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = statements.Block{
			Statements: []statements.Stmt{body, statements.ExpStmt{Expression: increment}},
		}
	}

	if condition == nil {
		condition = expressions.Literal{Value: true}
	}
	body = statements.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = statements.Block{Statements: []statements.Stmt{initializer, body}}
	}

	return body
}

func (p *parser) ifStatement() statements.Stmt {
	p.consume(tokens.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(tokens.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch statements.Stmt = nil
	if p.match(tokens.ELSE) {
		elseBranch = p.statement()
	}

	return statements.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *parser) block() []statements.Stmt {
	statements := []statements.Stmt{}

	for !p.check(tokens.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(tokens.RIGHT_BRACE, "Expect '}' after block.")
	return statements
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

func (p *parser) function(kind string) statements.Stmt {
	name, _ := p.consume(tokens.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))

	p.consume(tokens.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))

	params := []tokens.Token{}
	if !p.check(tokens.RIGHT_PAREN) {
		for ok := true; ok; ok = p.match(tokens.COMMA) {
			if len(params) >= 255 {
				curr := p.peek()
				p.errReporter.AddError(curr.LineNum, 0, "Can't have more than 255 parameters.")
				break
			}
			ident, _ := p.consume(tokens.IDENTIFIER, "Expect parameter name.")
			params = append(params, ident)
		}
	}

	p.consume(tokens.RIGHT_PAREN, "Expect ')' after arguments.")

	p.consume(tokens.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))

	body := p.block()

	return statements.FunctionStmt{Name: name, Params: params, Body: body}
}

func (p *parser) assignment() expressions.Expression {
	expr := p.or()

	if p.match(tokens.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if exp, ok := expr.(expressions.Variable); ok {
			name := exp.Name
			return expressions.Assign{Name: name, Value: value}
		}
		// TODO: make more accurate
		p.errReporter.AddError(0, 0, fmt.Sprintf("Invalid assignment target: %v", equals))
	}

	return expr
}

func (p *parser) or() expressions.Expression {
	expr := p.and()

	for p.match(tokens.OR) {
		operator := p.previous()
		right := p.and()
		expr = expressions.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *parser) and() expressions.Expression {
	expr := p.equality()

	for p.match(tokens.AND) {
		operator := p.previous()
		right := p.equality()
		expr = expressions.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
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
	return p.assignment()
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
	return p.call()
}

func (p *parser) call() expressions.Expression {
	expr := p.primary()

	for {
		if p.match(tokens.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *parser) finishCall(callee expressions.Expression) expressions.Expression {
	args := []expressions.Expression{}
	if !p.check(tokens.RIGHT_PAREN) {
		for ok := true; ok; ok = p.match(tokens.COMMA) {
			if len(args) >= 255 {
				curr := p.peek()
				p.errReporter.AddError(curr.LineNum, 0, "Can't have more than 255 arguments.")
				break
			}
			args = append(args, p.expression())
		}
	}

	paren, _ := p.consume(tokens.RIGHT_PAREN, "Expect ')' after arguments.")

	return expressions.Call{Callee: callee, Paren: paren, Arguments: args}
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

	if p.match(tokens.IDENTIFIER) {
		return expressions.Variable{Name: p.previous()}
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

func (p *parser) consume(t tokens.TokenType, message string) (tokens.Token, error) {
	if p.check(t) {
		t := p.advance()
		return t, nil
	}

	// err handling begins
	curr := p.peek()

	// TODO: at char count number
	p.errReporter.AddError(curr.LineNum, 0, message)

	p.synchronize()
	return tokens.Token{}, errors.New(message)
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
