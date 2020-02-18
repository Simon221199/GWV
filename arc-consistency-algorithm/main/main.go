package main

import (
	"fmt"
	"github.com/patrickz98/arc-consistency-algorithm/cartesian"
	"github.com/patrickz98/arc-consistency-algorithm/set"
	"strings"
)

// type constrain struct {
// 	variable1 string
// }

type toDoElem struct {
	variable  string
	constrain string
}

func matchedConstrains(variable string, constrains *set.Set) *set.Set {

	scopes := set.New()

	constrains.Do(func(value interface{}) {

		constrain := value.(string)

		if strings.Contains(constrain, variable) {
			scopes.Insert(constrain)
		}
	})

	return scopes
}

func scope(constrain string) (variables []string) {

	parts := strings.Split(constrain, " ")

	return []string{parts[0], parts[2]}
}

func toSlice(dom *set.Set) []interface{} {

	vars := make([]interface{}, 0)

	dom.Do(func(val interface{}) {
		elem := val.(int)
		vars = append(vars, elem)
	})

	return vars
}

func GAC() {

	variables := set.New("A", "B", "C", "D", "E")

	dom := make(map[string]*set.Set)
	dom["A"] = set.New(1, 2, 3, 4)
	dom["B"] = set.New(1, 2, 4)
	dom["C"] = set.New(1, 3, 4)
	dom["D"] = set.New(1, 2, 3, 4)
	dom["E"] = set.New(1, 2, 3, 4)

	constrains := set.New("A != B", "A == D", "E < A", "B != D", "C < D", "E < D", "B != C", "E < B", "E < C")

	todo := set.New()

	variables.Do(func(val1 interface{}) {

		variable := val1.(string)
		scopes := matchedConstrains(variable, constrains)

		scopes.Do(func(val2 interface{}) {
			constrain := val2.(string)
			todo.Insert(toDoElem{variable, constrain})
		})
	})

	todo.Do(func(val interface{}) {
		elem := val.(toDoElem)
		fmt.Printf("var=%s, const=(%s)\n", elem.variable, elem.constrain)
	})

	GAC2(variables, dom, constrains, todo)
}

func GAC2(variables *set.Set, dom map[string]*set.Set, constrains *set.Set, todo *set.Set) {

	for todo.Len() > 0 {
		elem := todo.Get().(toDoElem)
		todo.Remove(elem)

		fmt.Println("variable: " + elem.variable)
		fmt.Println("constrain: " + elem.constrain)

		nd := set.New()

		doms := make([][]interface{}, 0)
		// doms = append(doms, toSlice(dom[ elem.variable ]))

		ys := scope(elem.constrain)

		for _, scpe := range ys {

			fmt.Printf("doms: %s\n", scpe)
			doms = append(doms, toSlice(dom[scpe]))
		}

		fmt.Printf("len(doms): %d\n", len(doms))

		allCombinations := cartesian.Iter(doms...)

		for combination := range allCombinations {

			var1 := 0
			var2 := 0

			for inx, part := range combination {

				if inx == 0 {
					var1 = part.(int)
				}

				if inx == 1 {
					var2 = part.(int)
				}
			}

			fmt.Printf("%d, %d\n", var1, var2)

			if strings.Contains(elem.constrain, "!=") {
				if var1 != var2 {
					nd.Insert(var1)
				}
			}

			if strings.Contains(elem.constrain, "==") {
				if var1 == var2 {
					nd.Insert(var1)
				}
			}

			if strings.Contains(elem.constrain, "<") {
				if var1 < var2 {
					nd.Insert(var1)
				}
			}

			if strings.Contains(elem.constrain, ">") {
				if var1 > var2 {
					nd.Insert(var1)
				}
			}
		}

		fmt.Println(nd)
		fmt.Printf("done.")

		if dom[elem.variable].SubsetOf(nd) && nd.SubsetOf(dom[elem.variable]) {
			continue
		}

		fmt.Printf("dom[X] != nd\n")

	}

}

func main() {
	fmt.Println("Hello")

	GAC()
}
