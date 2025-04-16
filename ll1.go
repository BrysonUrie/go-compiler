package main

import (
	"errors"
	"fmt"

	scanner "github.com/brysonurie/ll1-parser/tokenizer"
)

type first map[Expr][]Expr
type follow map[Expr][]Expr
type firstPlus map[Production][]Expr
type table map[Expr]map[Expr]*Production

type ll1Table struct {
	first     first
	follow    follow
	firstPlus firstPlus
	table     table

	prodSet ProductionSet
}

func (table *ll1Table) String() string {
	ret := "\033[31m===First===\n"
	for key, val := range table.first {
		ret += "\033[34m" + key.Literal() + ":\033[32m "
		for _, expr := range val {
			ret += expr.Literal() + " "
		}
		ret += "\033[0m\n"
	}
	ret += "\033[31m===Follow===\n"
	for key, val := range table.follow {
		ret += "\033[34m" + key.Literal() + "::\033[32m "
		for _, expr := range val {
			ret += expr.Literal() + " "
		}
		ret += "\033[0m\n"
	}
	ret += "\033[31m===First Plus===\n"
	for key, val := range table.firstPlus {
		ret += "\033[34m" + key.lhs.Literal() + " => "
		firstRhs := (*key.rhs)[0]
		for _, expr := range *firstRhs {
			ret += expr.Literal() + " "
		}
		ret += "::\033[32m "
		for _, expr := range val {
			ret += expr.Literal() + " "
		}
		ret += "\033[0m\n"
	}
	ret += "\033[31m===LL1 Table===\n"
	ret += "\033[34mNon-Terminal\tTerminal\tProduction Rule\n"
	for nonTerminal, innerMap := range table.table {
		for terminal, prod := range innerMap {
			ret += "\033[34m" + nonTerminal.Literal() + "\t\t" + terminal.Literal() + "\t\t"
			if prod != nil {
				ret += "\033[32m" + prod.String()
			} else {
				ret += "\033[32mNIL"
			}
			ret += "\n"
		}
	}
	ret += "\033[0m\n"
	return ret
}

func createTable(prodSet ProductionSet) (*ll1Table, error) {
	table := &ll1Table{
		first:     make(first),
		follow:    make(follow),
		firstPlus: make(firstPlus),
		table:     make(map[Expr]map[Expr]*Production),
		prodSet:   prodSet,
	}
	table.populateFirst()
	table.populateFollow()
	table.populateFirstPlus()
	table.populateTable()
	return table, nil
}

func (table ll1Table) getNonTerm(exprKey Expr) (Expr, error) {
	for expr, _ := range table.follow {
		if exprKey.Literal() == expr.Literal() {
			return expr, nil
		}
	}
	return nil, fmt.Errorf("Could not find non-terminal")
}

func (table ll1Table) addToFirst(exprKey Expr, newExprSet []Expr) (bool, error) {
	changed := false
	for _, expr := range newExprSet {
		tableAtKey, err := table.getFirstAtExpr(exprKey)
		if err != nil {
			return changed, err
		}
		if !containsExpr(tableAtKey, expr) {
			table.first[exprKey] = append(table.first[exprKey], expr)
			changed = true
		}

	}
	return changed, nil
}

func (table ll1Table) addToFollow(exprKey Expr, newExprSet []Expr) (bool, error) {
	changed := false
	exprKey, err := table.getNonTerm(exprKey)
	if err != nil {
		return changed, err
	}
	for _, expr := range newExprSet {
		tableAtKey, err := table.getFollowAtExpr(exprKey)
		if err != nil {
			return changed, err
		}
		if !containsExpr(tableAtKey, expr) {
			table.follow[exprKey] = append(table.follow[exprKey], expr)
			changed = true
		}

	}
	return changed, nil
}
func (table ll1Table) getFirstAtExpr(expr Expr) ([]Expr, error) {
	for key, val := range table.first {
		if key.Literal() == expr.Literal() {
			return val, nil
		}
	}
	return []Expr{}, fmt.Errorf("Expr not found")
}

func (table ll1Table) getFollowAtExpr(expr Expr) ([]Expr, error) {
	for key, val := range table.follow {
		if key.Literal() == expr.Literal() {
			return val, nil
		}
	}
	return []Expr{}, fmt.Errorf("Expr not found")
}

func (table ll1Table) populateFirst() error {
	nonTerms := append([]Expr{}, getTerminalExpr()...)
	nonTerms = append(nonTerms, &Epsilon{})
	nonTerms = append(nonTerms, &EOF{})
	for _, expr := range nonTerms {
		table.first[expr] = []Expr{expr}
	}
	for _, prod := range *table.prodSet.productions {
		table.first[prod.lhs] = []Expr{}
	}

	stillChanging := true
	for stillChanging {
		stillChanging = false
		for _, prod := range *table.prodSet.productions {
			for _, curExprSet := range *prod.rhs {
				isValidSet := allTermOrNonTerm(*curExprSet)
				if isValidSet {
					k := len(*curExprSet) - 1
					firstExpr := (*curExprSet)[0]
					rhs, err := table.getFirstAtExpr(firstExpr)
					if err != nil {
						return err
					}
					rhs = removeEpsilon(rhs)

					i := 0
					for true {
						exprAtI := (*curExprSet)[i]
						firstAtI, err := table.getFirstAtExpr(exprAtI)
						if err != nil {
							return err
						}
						if !containsEpsilon(firstAtI) || i > k-1 {
							break
						}
						nextElem := (*curExprSet)[i+1]
						firstAtNext, err := table.getFirstAtExpr(nextElem)
						rhs = append(rhs, removeEpsilon(firstAtNext)...)
						i++
					}

					exprAtK := (*curExprSet)[k]
					firstAtK, err := table.getFirstAtExpr(exprAtK)
					if i == k && containsEpsilon(firstAtK) {
						rhs = append(rhs, &Epsilon{})
					}
					changed, err := table.addToFirst(prod.lhs, rhs)
					if err != nil {
						return err
					}
					if changed {
						stillChanging = true
					}
				}
			}
		}
	}

	return nil
}

func (table ll1Table) populateFollow() error {
	for _, prod := range *table.prodSet.productions {
		if prod.lhs.Literal() == "Goal" {
			table.follow[prod.lhs] = []Expr{&EOF{}}
		} else {
			table.follow[prod.lhs] = []Expr{}
		}
	}
	stillChanging := true
	for stillChanging {
		stillChanging = false
		for _, production := range *table.prodSet.productions {
			for _, exprSet := range *production.rhs {
				trailer, err := table.getFollowAtExpr(production.lhs)
				if err != nil {
					return err
				}
				k := len(*exprSet) - 1
				i := k
				for i >= 0 {
					exprAtI := (*exprSet)[i]
					_, isNonTerm := exprAtI.(*NonTerm)
					if isNonTerm {
						followAtI, err := table.getFollowAtExpr(exprAtI)
						if err != nil {
							return err
						}
						followAtIPlusTrailer := append(followAtI, trailer...)
						stillChanging, err = table.addToFollow(exprAtI, followAtIPlusTrailer)
						if err != nil {
							return err
						}
						firstAtI, err := table.getFirstAtExpr(exprAtI)
						if err != nil {
							return err
						}
						if containsEpsilon(firstAtI) {
							firstAtILessEps := removeEpsilon(firstAtI)
							trailer = append(trailer, firstAtILessEps...)
						} else {
							trailer, err = table.getFirstAtExpr(exprAtI)
							if err != nil {
								return err
							}
						}

					} else {
						trailer, err = table.getFirstAtExpr(exprAtI)
						if err != nil {
							return err
						}
					}
					i--
				}
			}
		}
	}
	return nil
}

func (table ll1Table) populateFirstPlus() error {
	for _, productions := range *table.prodSet.productions {
		for _, exprSet := range *productions.rhs {
			newExprs := &[]*[]Expr{exprSet}
			newProd := Production{lhs: productions.lhs, rhs: newExprs}
			firstExpr := (*exprSet)[0]
			firstAtFirstExpr, err := table.getFirstAtExpr(firstExpr)
			if err != nil {
				return err
			}
			if !containsEpsilon(firstAtFirstExpr) {
				table.firstPlus[newProd] = firstAtFirstExpr
			} else {
				followAtLhs, err := table.getFollowAtExpr(productions.lhs)
				if err != nil {
					return err
				}
				combinedFollowAndFirst := append(followAtLhs, firstAtFirstExpr...)
				table.firstPlus[newProd] = combinedFollowAndFirst
			}
		}
	}

	return nil
}

func findExprFromSlice(exprSet []Expr, exprKey Expr) (Expr, error) {
	for _, expr := range exprSet {
		if exprKey.Literal() == expr.Literal() {
			return expr, nil
		}
	}
	return nil, fmt.Errorf("Could not find in slice")
}

func (table ll1Table) populateTable() error {
	nonTerms := []Expr{}
	for _, expr := range *table.prodSet.productions {
		nonTerms = append(nonTerms, expr.lhs)
	}
	terms := getTerminalExpr()
	terms = append(terms, &EOF{})
	for _, nonTerm := range nonTerms {
		table.table[nonTerm] = make(map[Expr]*Production)
		for _, term := range terms {
			table.table[nonTerm][term] = nil
		}
	}
	for prod, exprSet := range table.firstPlus {
		for _, expr := range exprSet {
			_, isTerm := expr.(*Term)
			nonTerm, err := findExprFromSlice(nonTerms, prod.lhs)
			if err != nil {
				return err
			}
			if isTerm {
				term, err := findExprFromSlice(terms, expr)
				if err != nil {
					return err
				}
				table.table[nonTerm][term] = &prod
			}
			if containsExpr(exprSet, &EOF{}) {
				eofPos, err := findExprFromSlice(terms, &EOF{})
				if err != nil {
					return err
				}
				table.table[nonTerm][eofPos] = &prod
			}
		}
	}
	return nil
}

func (table ll1Table) getTableValue(focus Expr, word scanner.Token) (*Production, error) {
	tokenExprType := getTokenExprType(word)
	for nonTerminal, innerMap := range table.table {
		if focus.Literal() == nonTerminal.Literal() {
			for terminal, prod := range innerMap {
				if terminal.Literal() == tokenExprType.Literal() && prod != nil {
					return prod, nil
				}
			}
		}
	}
	return nil, errors.New("No value found in table for [" + focus.Literal() + "][" + word.Lexeme + " : " + tokenExprType.Literal() + "]")
}
