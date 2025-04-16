package main

import (
	"bufio"
	"os"
)

// import (
// 	scanner "github.com/brysonurie/ll1-parser/tokenizer"
// )

type Compiler struct {
	table       *ll1Table
	fileScanner *bufio.Scanner
}

func createCompiler(table *ll1Table, filePath string) (*Compiler, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure the file is closed when done
	return &Compiler{
		table:       table,
		fileScanner: bufio.NewScanner(file),
	}, nil
}

func (compiler *Compiler) readLine() bool {
	if compiler.fileScanner.Scan() {

	} else {
		return false
	}
	return true
}
