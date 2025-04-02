package main

import (
	"errors"
	"fmt"

	scanner "github.com/brysonurie/ll1-parser/tokenizer"
)

type ParseTreeNode struct {
	Symbol   Expr
	Children []*ParseTreeNode
	IsRont   bool
}

func printParseTreeFancy(node *ParseTreeNode, prefix string, isLast bool) {
	if node == nil {
		return
	}
	// Determine the tree branch characters
	fmt.Print(prefix)
	if isLast {
		fmt.Print("└── ")
		prefix += "    "
	} else {
		fmt.Print("├── ")
		prefix += "│   "
	}

	// Print node data
	printStr := node.Symbol.Literal()
	if termNode, ok := node.Symbol.(*Term); ok {
		printStr += " : " + termNode.token.Lexeme
	} else if node.IsRont {
		printStr += " (RONT)"
	}
	fmt.Println(printStr)
	// Recursively print children
	for i, child := range node.Children {
		printParseTreeFancy(child, prefix, i == len(node.Children)-1)
	}
}

func CreateParseTree(table ll1Table, tokens []scanner.Token) (*ParseTreeNode, error) {
	if tokens[len(tokens)-1].Lexeme != "$" {
		tokens = append(tokens, *scanner.CreateToken(scanner.EOF, "$"))
	}
	curWordIndex := 0
	word := tokens[curWordIndex]

	stack := Stack[*ParseTreeNode]{}
	stack.Push(&ParseTreeNode{Symbol: &EOF{}, IsRont: false})
	root := &ParseTreeNode{Symbol: &NonTerm{"Goal"}, IsRont: false}
	stack.Push(root)
	//stack.Push(&EOF{})
	//stack.Push(&NonTerm{"Goal"})

	for true {
		focusNode, err := stack.Peek()
		if err != nil {
			return nil, err
		}
		focus := focusNode.Symbol
		_, focusIsEof := focus.(*EOF)
		_, focusIsTerm := focusNode.Symbol.(*Term)
		if focusIsEof && word.Type == scanner.EOF {
			return root, nil
		} else if focusIsEof || focusIsTerm {
			wordExprType := getTokenExprType(word)
			if focus.Literal() == wordExprType.Literal() {
				focusNode.Symbol = &Term{
					literal: focus.Literal(),
					token:   &word,
				}
				if termNode, ok := focusNode.Symbol.(*Term); ok { // Ensure it's a pointer
					token := scanner.CreateToken(word.Type, word.Lexeme) // Create token
					termNode.token = token                               // Store token as a pointer
				}

				curWordIndex++
				word = tokens[curWordIndex]
				stack.Pop()
			} else {
				return nil, errors.New("Error looking for symbol at top of stack " + word.Lexeme)
			}
		} else {
			currentProd, err := table.getTableValue(focus, word)
			if err != nil {
				return nil, err
			}
			stack.Pop()
			rhsLen := len(*(*currentProd.rhs)[0]) - 1
			newChildren := []*ParseTreeNode{}
			for i := rhsLen; i >= 0; i-- {
				curExpr := (*(*currentProd.rhs)[0])[i]
				isRont := false
				if curExpr.Literal() != "#" {
					if i == 1 {
						prevExpr := (*(*currentProd.rhs)[0])[0]
						if IsOperator(prevExpr) {
							isRont = true
						}
					}
					childNode := &ParseTreeNode{Symbol: curExpr, IsRont: isRont}
					newChildren = append(newChildren, childNode)
					stack.Push(childNode)
				}
			}
			focusNode.Children = newChildren
		}
	}
	return nil, errors.New("Should not break")
}

func (node *ParseTreeNode) findRontChild() int {
	for i, child := range node.Children {
		if child.IsRont {
			return i
		}
	}
	return -1
}

func (node *ParseTreeNode) findOpIndex() int {
	for i, child := range node.Children {
		if IsOperator(child.Symbol) {
			return i
		}
	}
	return -1
}
