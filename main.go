package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	scanner "github.com/brysonurie/ll1-parser/tokenizer"
)

func main() {
	prodSet, err := CreateProductionSet("replacements.txt")
	if err != nil {
		panic("Error creating production" + err.Error())
	}
	// fmt.Println(prodSet.String())
	table, err := createTable(*prodSet)
	if err != nil {
		panic("Error creating table" + err.Error())

	}

	file, err := os.Open("examples.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when done

	// Create a fileScanner to read the file line by line
	fileScanner := bufio.NewScanner(file)

	// Read and print each line
	lineNum := 0
	symbolTable, err := CreateSymbolTable()
	if err != nil {
		panic("Error creating table" + err.Error())

	}

	code := []string{
		`section .data`,
		`fmtstr: db "%s", 10, 0`,
		`fmtint: db "%d", 10, 0`,
		`fmtfloat: db "%f", 10, 0`,
		`fmtintin: db "%d", 0`,
		`fmtfloatin: db "%f", 0`,
		`float1: dd 0.0`,
		`section .text`,
		`extern printf`,
		`global main`,
		`main:`,
		// preserve stack
		`push RBP`,
		`mov RBP, RSP`,
		`sub RSP, 64`,
	}

	for fileScanner.Scan() {
		line := fileScanner.Text() // Get the current line as a string
		scanner := scanner.CreateScanner(line)
		tokens := scanner.Scan()

		parseTree, err := CreateParseTree(*table, tokens)
		if err != nil {
			fmt.Println(line)
			fmt.Println(err.Error())
			continue
		}

		node, err := createAST(parseTree)
		if err != nil {
			fmt.Println(line)
			printParseTreeFancy(parseTree, "", true)
			fmt.Println(err)
			continue
		}

		fmt.Println(strconv.Itoa(lineNum) + " " + node.Literal())
		/// AT THIS POINT THE AST NODE HAS BEEN BUILT
		if assNode, ok := node.(*AssignNode); ok {
			// THIS IS FOR ASSIGN NODES
			newCode, err := assNode.ConvertToNasm(*symbolTable)
			if err != nil {
				panic(err)
			}
			fmt.Println(newCode)
			code = append(code, newCode...)
			if err != nil {
				panic(err)
			}
			// tableEntry, _ := symbolTable.symbols[assNode.name]

		}

		lineNum++
	}
	fmt.Println(symbolTable)

	err = os.WriteFile("output.asm", []byte(fmt.Sprintf("%s\n", joinLines(code))), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ NASM code written to output.asm")

}

func joinLines(lines []string) string {
	return fmt.Sprint(join(lines, "\n"))
}

func join(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		result += s
		if i < len(strs)-1 {
			result += sep
		}
	}
	return result
}
