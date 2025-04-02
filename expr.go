package main

import scanner "github.com/brysonurie/ll1-parser/tokenizer"

type Expr interface {
	Literal() string
}

type NonTerm struct {
	literal string
}

func (nonTerm *NonTerm) Literal() string {
	return nonTerm.literal
}

type Term struct {
	literal string
	token   *scanner.Token
}

func (term *Term) Literal() string {
	return term.literal
}

type Epsilon struct{}

var EpsilonString = "#"

func (term *Epsilon) Literal() string {
	return EpsilonString
}

type EOF struct{}

func (term *EOF) Literal() string {
	return "$"
}

var Terminals = []string{
	"/",    // Division operator
	"*",    // Multiplication operator
	"-",    // Minus sign or subtraction operator
	"+",    // Plus sign or addition operator
	"name", // Identifier or variable name
	"num",  // Numeric constant
	"(",    // Opening parenthesis
	")",    // Closing parenthesis
	"^",    // Exponentiation operator
}

var Operators = []string{
	"/", // Division operator
	"*", // Multiplication operator
	"-", // Minus sign or subtraction operator
	"+", // Plus sign or addition operator
	"^", // Exponentiation operator
}

func IsOperator(expr Expr) bool {
	for _, op := range Operators {
		if expr.Literal() == op {
			return true
		}
	}

	return false
}

func getTerminalExpr() []Expr {
	expr := []Expr{}
	for _, term := range Terminals {
		expr = append(expr, &Term{literal: term})

	}
	return expr
}

func removeEpsilon(exprSet []Expr) []Expr {
	newExprSet := []Expr{}
	for _, expr := range exprSet {
		if expr.Literal() != EpsilonString {
			newExprSet = append(newExprSet, expr)
		}
	}
	return newExprSet
}

func containsEpsilon(exprs []Expr) bool {
	return containsExpr(exprs, &Epsilon{})
}

func containsExpr(exprs []Expr, exprKey Expr) bool {
	for _, expr := range exprs {
		if expr.Literal() == exprKey.Literal() {
			return true
		}
	}
	return false

}

func allTermOrNonTerm(exprSet []Expr) bool {
	for _, expr := range exprSet {
		_, isTerm := expr.(*Term)
		_, isNonTerm := expr.(*NonTerm)
		_, isEps := expr.(*Epsilon)
		if !isTerm && !isNonTerm && !isEps {
			return false
		}
	}
	return true
}

func getTokenExprType(token scanner.Token) Expr {
	if token.Type == scanner.NUMBER {
		return &Term{"num", nil}
	} else if token.Type == scanner.IDENTIFIER {
		return &Term{"name", nil}
	} else if token.Type == scanner.EOF {
		return &EOF{}
	} else {
		return &Term{token.Lexeme, scanner.CreateToken(token.Type, token.Lexeme)}
	}
}
