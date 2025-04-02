package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"

	scanner "github.com/brysonurie/ll1-parser/tokenizer"
)

type ASTNode struct {
	Operator string
	Value    *string
	Left     *ASTNode
	Right    *ASTNode
}

func convExpr(parseNode *ParseTreeNode, exprStack *Stack[*Term]) {
	if termNode, ok := parseNode.Symbol.(*Term); ok {
		if termNode.token.Lexeme == "(" || termNode.token.Lexeme == ")" {
			return
		}
		exprStack.Push(termNode)
		return
	}
	if parseNode == nil || len(parseNode.Children) == 0 {
		return
	}
	opIndex := parseNode.findOpIndex()
	rontIndex := parseNode.findRontChild()

	processedNums := []int{}

	if opIndex != -1 {
		for i := opIndex + 1; i < len(parseNode.Children); i++ {
			curChild := parseNode.Children[i]
			processedNums = append(processedNums, i)
			convExpr(curChild, exprStack)
		}
	}
	if rontIndex != -1 {
		if !slices.Contains(processedNums, rontIndex) {
			curChild := parseNode.Children[rontIndex]
			processedNums = append(processedNums, rontIndex)
			convExpr(curChild, exprStack)
		}
	}
	if opIndex != -1 {
		if !slices.Contains(processedNums, opIndex) {
			curChild := parseNode.Children[opIndex]
			processedNums = append(processedNums, opIndex)
			convExpr(curChild, exprStack)
		}
	}
	for i := len(parseNode.Children) - 1; i >= 0; i-- {
		if !slices.Contains(processedNums, i) {
			curChild := parseNode.Children[i]
			processedNums = append(processedNums, i)
			convExpr(curChild, exprStack)
		}
	}
}

func simplifyExpr(stack *Stack[*Term]) (*Stack[[]*Term], error) {

	reverseStack := &Stack[*Term]{}

	for !stack.IsEmpty() {
		term, _ := stack.Pop()
		reverseStack.Push(term)
	}

	answerStack := &Stack[[]*Term]{}

	for !reverseStack.IsEmpty() {
		term, _ := reverseStack.Pop()
		if term.token.Type == scanner.NUMBER || term.token.Type == scanner.IDENTIFIER {
			answerStack.Push([]*Term{term})
		} else {
			top1, err := answerStack.Pop()
			top2, err := answerStack.Pop()
			if err != nil {
				return nil, err
			}
			if len(top1) == 1 &&
				len(top2) == 1 &&
				top1[0].token.Type == scanner.NUMBER &&
				top2[0].token.Type == scanner.NUMBER {
				num1, err := strconv.ParseFloat(top1[0].token.Lexeme, 64)
				num2, err := strconv.ParseFloat(top2[0].token.Lexeme, 64)
				num3 := 0.0
				if err != nil {
					return nil, err
				}
				if term.Literal() == "/" {
					num3 = num2 / num1
				} else if term.Literal() == "*" {
					num3 = num1 * num2
				} else if term.Literal() == "+" {
					num3 = num1 + num2
				} else if term.Literal() == "-" {
					num3 = num1 - num2
				} else if term.Literal() == "^" {
				}
				if math.IsNaN(num3) || math.IsInf(num3, 1) || math.IsInf(num3, 0) {
					return nil, errors.New("Error parsing numbers")
				}
				strNum3 := strconv.FormatFloat(num3, 'f', 2, 64)
				answerStack.Push([]*Term{
					{literal: strNum3, token: &scanner.Token{
						Lexeme: strNum3,
						Type:   scanner.NUMBER,
					}},
				})
			} else {
				answerStack.Push(top2)
				answerStack.Push(top1)
				answerStack.Push([]*Term{term})
			}
		}
	}
	return answerStack, nil
}

func createAST(parseNode *ParseTreeNode) {
	// Minimize an expression
	termStack := &Stack[*Term]{}
	convExpr(parseNode, termStack)
	answerStack, err := simplifyExpr(termStack)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	for !answerStack.IsEmpty() {
		answer, _ := answerStack.Pop()
		for _, item := range answer {
			fmt.Print(item.token.Lexeme + " ")
		}
	}

	fmt.Println()
}
