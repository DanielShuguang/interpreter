package parser

import (
	"fmt"
	"monkey-compiler-study/ast"
	"monkey-compiler-study/lexer"
	"monkey-compiler-study/token"
)

type Parser struct {
	l *lexer.Lexer // 指向词法分析器实例的指针

	curToken  token.Token // 当前字符的信息
	peekToken token.Token // 下一个字符的信息

	errors []string // 检测到的语法错误列表
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// 读取两个词法单元，以设置 curToken 和 peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// 往程序错误信息中插入一条错误
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// 同时移动 curToken 和 peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// 没有遇到结束符号时需要一直读取程序代码
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// 获取当前语句的准确类型
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement() // 遇到 let 关键词，判断是否为正确的声明语法
	case token.RETURN:
		return p.parseReturnStatement() // 遇到 return 关键词，判断是否为正确的返回语句
	default:
		return nil
	}
}

// 遇到 let 关键词时，继续向后窥探，判断是否为完整的声明语句
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 跳过对表达式的处理，直到遇见分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 遇到 return 关键词时，继续向后窥探，判断是否为完整的语句
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: 跳过对表达式的处理，知道遇见分号
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 当前token是否为某一类型
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 下一个token是否为某一类型
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 下一个token是否为期望类型
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		// 不是预期类型，插入一条报错信息
		p.peekError(t)
		return false
	}
}
