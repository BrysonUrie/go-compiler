package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Production struct {
	lhs *NonTerm
	rhs *[]*[]Expr
}

type ProductionSet struct {
	productions *[]Production
}

func (p *Production) String() string {
	retStr := p.lhs.Literal() + " => "
	for _, exprSet := range *p.rhs {
		for _, expr := range *exprSet {
			retStr += "`" + expr.Literal() + "` "
		}
		if slices.Index(*p.rhs, exprSet) != len(*p.rhs)-1 {

			retStr += "| "
		}
	}
	return retStr
}

func (pset *ProductionSet) String() string {
	retStr := ""
	for _, prod := range *pset.productions {
		retStr += prod.String() + "\n"
	}
	return retStr
}

func (pset *ProductionSet) addProduction(production string) error {
	arrowSymbol := "→"
	splitProduction := strings.Split(production, arrowSymbol)
	if len(splitProduction) != 2 {
		return fmt.Errorf(" Invalid production" + production)
	}
	lhsString := splitProduction[0]
	rhsString := splitProduction[1]
	rhsExprs := strings.Split(rhsString, "|")

	lhs := &NonTerm{strings.TrimSpace(lhsString)}
	rhs := &[]*[]Expr{}

	for _, exprList := range rhsExprs {
		splitExpr := strings.Split(exprList, " ")
		exprArr := &[]Expr{}
		for _, exprStr := range splitExpr {
			strippedExpr := strings.TrimSpace(exprStr)
			if strippedExpr == EpsilonString {
				*exprArr = append(*exprArr, &Epsilon{})
			} else if slices.Contains(Terminals, strippedExpr) {
				*exprArr = append(*exprArr, &Term{literal: strippedExpr})
			} else if strippedExpr != "" {
				*exprArr = append(*exprArr, &NonTerm{literal: strippedExpr})
			}
		}
		*rhs = append(*rhs, exprArr)
	}

	newProd := &Production{
		lhs: lhs,
		rhs: rhs,
	}
	*pset.productions = append(*pset.productions, *newProd)
	return nil
}

func CreateProductionSet(productionFile string) (*ProductionSet, error) {
	prodSet := &ProductionSet{productions: &[]Production{}}
	file, err := os.Open(productionFile)
	if err != nil {
		return prodSet, fmt.Errorf("Could not open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		err := prodSet.addProduction(line)
		if err != nil {
			return prodSet, err
		}
	}
	if err := scanner.Err(); err != nil {
		return prodSet, err
	}

	return prodSet, nil
}
