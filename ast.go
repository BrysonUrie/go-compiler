package main

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"

	scanner "github.com/brysonurie/ll1-parser/tokenizer"
)

type ASTNode interface {
	Literal() string
}

type ExprNode struct {
	expr Stack[[]*Term]
}

func (exprNode *ExprNode) Literal() string {
	copyStack := exprNode.expr.Copy()
	str := ""
	for !copyStack.IsEmpty() {
		fmt.Println("another stack itme")
		answer, _ := copyStack.Pop()
		for _, item := range answer {
			str += (item.token.Lexeme + " ")
		}
	}
	return str
}

type ExprToNasmNode interface {
	loadNum(register string) string
}

type NumNasmNode struct {
	num int
}

func (numNode *NumNasmNode) loadNum(register string) string {
	return fmt.Sprintf("mov %s, %d", register, numNode.num)
}

type VarNasmNode struct {
	stackPos int
}

func (varNode *VarNasmNode) loadNum(register string) string {
	if varNode.stackPos == 0 {
		return fmt.Sprintf("mov %s dword [rbp]", register)
	}
	return fmt.Sprintf("mov %s, dword [rbp-%d]", register, varNode.stackPos)
}

type StackNasmNode struct{}

func (stackNode StackNasmNode) loadNum(register string) string {
	// this will probably not work when I have two pushed into stack
	return fmt.Sprintf("pop %s", register)
}

func (assNode *AssignNode) ConvertToNasm(symbolTable SymbolTable) ([]string, error) {
	copyStack := assNode.expr.expr.Copy()
	// reverse for post-fix
	copyStack.Reverse()
	termStack := Stack[ExprToNasmNode]{}
	newCode := []string{}

	for !copyStack.IsEmpty() {
		symbol, _ := copyStack.Pop()
		if len(symbol) != 1 {
			return nil, fmt.Errorf("I think sometihng is wrong")
		}
		trueSymbol := symbol[0]
		if IsOperator(trueSymbol) {
			op := trueSymbol.token.Lexeme
			if termStack.IsEmpty() {
				termStack.Push(&StackNasmNode{})
			}
			pop1, _ := termStack.Pop()
			if termStack.IsEmpty() {
				termStack.Push(&StackNasmNode{})
			}
			pop2, _ := termStack.Pop()

			// Is an op
			if op == "-" {

			} else if op == "+" {

			} else if op == "/" {
				dividendReg := "RAX"
				newCode = append(newCode, pop1.loadNum(dividendReg))
				divisorReg := "RCX"
				newCode = append(newCode, pop2.loadNum(divisorReg))
				newCode = append(newCode, "xor RDX, RDX")
				newCode = append(newCode, "div RCX")
				newCode = append(newCode, fmt.Sprintf("push RAX"))
			} else if op == "*" {

			} else {
				return nil, fmt.Errorf("Invalid Operator %s", op)
			}

		} else {
			numConversion, err := strconv.ParseInt(trueSymbol.token.Lexeme, 10, 64)
			if err == nil {
				// this means that it is a number
				termStack.Push(&NumNasmNode{
					num: int(numConversion),
				})
				continue
			}
			// this means that it is likely a var
			varName := trueSymbol.token.Lexeme
			symbolEntry, exists := symbolTable.symbols[varName]
			if exists {
				termStack.Push(&VarNasmNode{
					stackPos: symbolEntry.offset,
				})
			} else {
				return nil, fmt.Errorf("No var with name %s", varName)
			}
		}
	}
	entry, err := symbolTable.addEntry(assNode.nodeType, 0, assNode.name)
	if err != nil {
		return nil, err
	}
	newCode = append(newCode, []string{
		"pop RAX",
		fmt.Sprintf("mov dword [edp-%d], RAX", entry.offset),
	}...)
	return newCode, nil
}

type AssignNode struct {
	nodeType string
	name     string
	expr     ExprNode
}

func (assignNode AssignNode) Literal() string {
	return "AssignNode: " + assignNode.nodeType + " " + assignNode.name + " = " + assignNode.expr.Literal()
}

// func (assignNode AssignNode) ConvertToNasm() string[] {
//
// }

type PrintExprNode struct {
	expr ExprNode
}

func (printExprNode PrintExprNode) Literal() string {
	return " << " + printExprNode.expr.Literal()
}

type PrintVarNode struct {
	name string
}

func (printVar PrintVarNode) Literal() string {
	return "PrintVarNode: << " + printVar.name
}

type PrintNumNode struct {
	num int
}

func (printNum PrintNumNode) Literal() string {
	return "PrintNumNode: << " + strconv.Itoa(printNum.num)
}

type PrintNameNode struct {
	name string
}

func convExpr(parseNode *ParseTreeNode, exprStack *Stack[*Term]) {
	if termNode, ok := parseNode.Symbol.(*Term); ok {
		parent := parseNode.Parent
		grandparent := parent.Parent
		negative := false
		if grandparent != nil && grandparent.Symbol.Literal() == "GTermSign" {
			if len(grandparent.Children) > 1 {
				negative = true
			}
		}
		if termNode.token.Lexeme == "(" || termNode.token.Lexeme == ")" {
			return
		}
		if negative {

			termNode.token.Lexeme = "-" + termNode.token.Lexeme
		}
		if parent.Symbol.Literal() != "GTermSign" {
			exprStack.Push(termNode)
		}
		return
	}
	if parseNode == nil || len(parseNode.Children) == 0 {
		return
	}
	opIndex := parseNode.findOpIndex()
	rontIndex := parseNode.findRontChild()

	processedNums := []int{}

	// parent := parseNode.Parent
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

func getExpr(parseNode *ParseTreeNode) (*Stack[[]*Term], error) {
	if parseNode == nil {
		return nil, fmt.Errorf("No children")
	}
	if parseNode.Symbol.Literal() == "Expr" {
		termStack := &Stack[*Term]{}
		convExpr(parseNode, termStack)
		simpleExpr, err := simplifyExpr(termStack)
		if err != nil {
			return nil, err
		}
		return simpleExpr, nil
	}
	for _, child := range parseNode.Children {
		stack, err := getExpr(child)
		if err == nil {
			return stack, nil
		}
	}
	return nil, fmt.Errorf("No Valid Children")
}

func getTermVal(parseNode *ParseTreeNode, nodeType string) (string, error) {
	if parseNode == nil {
		return "", fmt.Errorf("Reached end of branch")
	}
	if parseNode.Symbol.Literal() == nodeType {
		if termNode, ok := parseNode.Symbol.(*Term); ok {
			return termNode.token.Lexeme, nil
		}
	}
	for _, child := range parseNode.Children {
		ret, err := getTermVal(child, nodeType)
		if err == nil {
			return ret, nil
		}
	}
	return "", fmt.Errorf("Node not found")
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
				strNum3 := strconv.FormatFloat(num3, 'f', 0, 64)
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

func createAST(parseNode *ParseTreeNode) (ASTNode, error) {
	// expr, err := getExpr(parseNode)
	// if err != nil {
	// 	fmt.Println("bad job")
	// } else {
	// 	newNode := ExprNode{
	// 		expr: *expr,
	// 	}
	// 	fmt.Println(newNode.Literal())
	// }
	if parseNode == nil {
		return nil, fmt.Errorf("Reached end of node tree")
	}

	if parseNode.Symbol.Literal() == "LineFull" {
		if len(parseNode.Children) == 1 {
			firstChild := parseNode.Children[0]
			if firstChild.Symbol.Literal() == "ExprWithoutName" {
				if len(firstChild.Children) == 1 {
					stack, err := getExpr(parseNode)
					if err != nil {
						return nil, err
					}
					exprNode := &ExprNode{
						expr: *stack,
					}
					// this is void expression
					return exprNode, nil
				} else if len(firstChild.Children) == 2 {
					//this is print expr
					gTermSign := firstChild.Children[0]
					gTerm := gTermSign.Children[0]
					gTermChild := gTerm.Children[0]
					if gTermChild.Symbol.Literal() == "Parens" {
						// Print Expr
						stack, err := getExpr(parseNode)
						if err != nil {
							return nil, err
						}
						exprNode := &ExprNode{
							expr: *stack,
						}
						return &PrintExprNode{expr: *exprNode}, nil
					} else if gTermChild.Symbol.Literal() == "name" {
						name, err := getTermVal(gTermChild, "name")
						if err != nil {
							return nil, fmt.Errorf("No name found")
						}
						return &PrintVarNode{
							name: name,
						}, nil
					} else if gTermChild.Symbol.Literal() == "num" {
						num, err := getTermVal(gTermChild, "num")
						if err != nil {
							return nil, fmt.Errorf("No num found")
						}
						int, err := strconv.Atoi(num)
						if err != nil {
							return nil, fmt.Errorf("Not a valid number")
						}
						return &PrintNumNode{
							num: int,
						}, nil

					}

				}
			} else {
				fmt.Println("bad job didnt find expr iwhtout name")
				fmt.Println("probably a LineVarName" + firstChild.Symbol.Literal())
			}
		} else if len(parseNode.Children) == 2 {
			// this is var declaration
			stack, err := getExpr(parseNode)
			if err != nil {
				return nil, err
			}
			exprNode := &ExprNode{
				expr: *stack,
			}
			varTypeAfter := parseNode.Children[0]
			lineVarName := varTypeAfter.Children[0]
			varNameProp := lineVarName.Children[1]

			varName, ok := varNameProp.Symbol.(*Term)
			if !ok {
				return nil, fmt.Errorf("Bad name")
			}

			return &AssignNode{
				nodeType: "int32",
				name:     varName.token.Lexeme,
				expr:     *exprNode,
			}, nil

		}
	}

	for _, child := range parseNode.Children {
		node, err := createAST(child)
		if err == nil {
			return node, nil
		}
	}

	return nil, fmt.Errorf("Couldnt find suitable AST conversion")
}
