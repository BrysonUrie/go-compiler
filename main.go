package main

import (
	"bufio"
	"fmt"
	"os"

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
	table.populateFirst()
	table.populateFollow()
	table.populateFirstPlus()
	table.populateTable()
	fmt.Println(table.String())

	file, err := os.Open("examples.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // Ensure the file is closed when done

	// Create a fileScanner to read the file line by line
	fileScanner := bufio.NewScanner(file)

	// Read and print each line
	for fileScanner.Scan() {
		line := fileScanner.Text() // Get the current line as a string
		scanner := scanner.CreateScanner(line)
		tokens := scanner.Scan()
		parseTree, err := CreateParseTree(*table, tokens)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// printParseTreeFancy(parseTree, "", true)

		fmt.Print("Successfully parsed: ")
		for _, token := range tokens {
			fmt.Print(token.String())
		}
		fmt.Println()

		createAST(parseTree)
	}

	// Check for errors during scanning
	if err := fileScanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

}
