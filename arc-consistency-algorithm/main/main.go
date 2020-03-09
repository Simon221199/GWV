package main

import (
	"fmt"
	"strings"
)

// type constrain struct {
// 	variable1 string
// }

type arc struct {
	variable   string
	constraint constraint
}

type set map[string]bool

type constraint struct {
	aVariable string
	aChar     int
	dVariable string
	dChar     int
}

func (con constraint) check(wordA, wordD string) bool {

	return wordA[con.aChar] == wordD[con.dChar]

}

func (con constraint) getScope() []string {
	return []string{con.aVariable, con.dVariable}
}

func (con constraint) has(search string) bool {

	for _, variable := range con.getScope() {
		if variable == search {
			return true
		}
	}

	return false
}

func (con constraint) getOther(except string) string {

	if except == con.aVariable {
		return con.dVariable
	} else {
		return con.aVariable
	}
}

func newSet(variables ...string) set {

	set := make(map[string]bool)

	for _, variable := range variables {
		set[variable] = true
	}

	return set
}

func eq(aaa, bbb map[string]bool) bool {

	if len(aaa) != len(bbb) {
		return false
	}

	for key1 := range aaa {
		if _, ok := bbb[key1]; !ok {
			return false
		}
	}

	return true
}

func GAC() {

	variables := []string{
		"add", "ado", "age", "ago", "aid", "ail", "aim", "air", "and", "any",
		"ape", "apt", "arc", "are", "ark", "arm", "art", "ash", "ask", "auk",
		"awe", "awl", "aye", "bad", "bag", "ban", "bat", "bee", "boa", "ear",
		"eel", "eft", "far", "fat", "fit", "lee", "oaf", "rat", "tar", "tie",
	}

	dom := make(map[string]set)
	// dom["A0"] = newSet("eel")
	dom["A0"] = newSet(variables...)
	dom["A1"] = newSet(variables...)
	dom["A2"] = newSet(variables...)
	// dom["A2"] = newSet("art")
	dom["D0"] = newSet(variables...)
	dom["D1"] = newSet(variables...)
	dom["D2"] = newSet(variables...)

	constraints := make(map[constraint]bool, 0)

	for inx := 0; inx < 3; inx++ {

		a := fmt.Sprintf("A%d", inx)

		for iny := 0; iny < 3; iny++ {

			d := fmt.Sprintf("D%d", iny)

			fmt.Printf("%s[%d] = %s[%d]\n", a, iny, d, inx)
			constrain := constraint{
				aVariable: a,
				aChar:     iny,
				dVariable: d,
				dChar:     inx,
			}

			constraints[constrain] = true
		}
	}

	todo := make(map[arc]bool)

	for constraint := range constraints {

		for _, X := range constraint.getScope() {
			nArc := arc{
				variable:   X,
				constraint: constraint,
			}

			todo[nArc] = true
		}
	}

	fmt.Println(todo)

	GAC2(variables, dom, constraints, todo)

	for X, vars := range dom {
		fmt.Println()
		fmt.Println(X)

		for word := range vars {
			fmt.Println(word)
		}
		// fmt.Println(vars)
	}
}

func GAC2(variables []string, dom map[string]set, constraints map[constraint]bool, todo map[arc]bool) {

	for elem := range todo {

		delete(todo, elem)

		yk := elem.constraint.getOther(elem.variable)

		// fmt.Println("variable: " + elem.variable)
		// fmt.Println("yk:       " + yk)
		// fmt.Println("constrain: " + elem.constraint)

		nd := make(set)

		for word1 := range dom[elem.variable] {
			for word2 := range dom[yk] {

				var wordA string
				var wordD string

				if strings.HasPrefix(elem.variable, "A") {
					wordA = word1
					wordD = word2
				} else {
					wordA = word2
					wordD = word1
				}

				if elem.constraint.check(wordA, wordD) {
					nd[word1] = true
				}
			}
		}

		// ND == dom[X]
		if eq(dom[elem.variable], nd) {
			continue
		}

		dom[elem.variable] = nd

		for constraint := range constraints {

			if constraint == elem.constraint {
				continue
			}

			if constraint.has(elem.variable) {
				nArc := arc{
					variable:   constraint.getOther(elem.variable),
					constraint: constraint,
				}

				todo[nArc] = true
			}
		}
	}
}

func main() {

	fmt.Println("Hello")

	GAC()
}
