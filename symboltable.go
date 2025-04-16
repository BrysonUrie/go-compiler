package main

import "fmt"

type SymbolTableEntry struct {
	entryType    string
	offset       int
	lineDeclared int
}
type SymbolTable struct {
	symbols   map[string]*SymbolTableEntry
	children  []*SymbolTable
	parent    *SymbolTable
	curOffset int
}

func CreateSymbolTable() (*SymbolTable, error) {
	return &SymbolTable{
		symbols:  make(map[string]*SymbolTableEntry),
		children: []*SymbolTable{},
		parent:   nil,
	}, nil
}

var typeSizeMapping = map[string]int{
	"int32": 4,
}

func (table SymbolTable) String() string {
	str := "Sign \tType \tOffset \tLine Declared"
	for key, val := range table.symbols {
		str += fmt.Sprintf("\n%s \t%s \t%d \t%d", key, val.entryType, val.offset, val.lineDeclared)
	}
	return str
}

func (table *SymbolTable) updateEntry(
	lineDeclared int,
	varName string,
) error {
	table.symbols[varName].lineDeclared = lineDeclared
	return nil
}

func (symbolTable *SymbolTable) addEntry(
	entryType string,
	lineDeclared int,
	varName string,
) (*SymbolTableEntry, error) {
	size, ok := typeSizeMapping[entryType]
	if !ok {
		return nil, fmt.Errorf("Invalid Type")
	}
	old, exists := symbolTable.symbols[varName]
	if exists {
		return old, symbolTable.updateEntry(lineDeclared, varName)
	}
	symbolTable.curOffset += size
	newEntr := &SymbolTableEntry{
		entryType:    entryType,
		offset:       symbolTable.curOffset,
		lineDeclared: lineDeclared,
	}
	symbolTable.symbols[varName] = newEntr
	return newEntr, nil
}
