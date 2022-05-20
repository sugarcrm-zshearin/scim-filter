package main

import (
	"fmt"
	"log"

	"errors"

	// "github.com/imulab/go-scim/pkg/v2/crud"
	"github.com/imulab/go-scim/pkg/v2/crud/expr"
	// "github.com/imulab/go-scim/pkg/v2/prop"
)

var (
	ErrUnsupportedDBFilter   = errors.New("unsupported database filter - will only use post filtering")
	EqualsRelationalOperator = "eq"
)

func main() {

	// myFilter := "((value eq true) and (dog eq true)) or (primary eq true)"

	// resource := &prop.Resource{}
	// meetsCriteria, err := crud.Evaluate(resource, myFilter)
	// if err != nil {
	// 	log.Fatalf("error parsing expression")
	// }

	// if meetsCriteria {
	// 	fmt.Printf("resource meets criteria, include in response\n")
	// }

	//	myFilter := "(value eq true) and (primary ne true)"
	//myFilter := "username ne \"foo\""
	//myFilter := "username eq foo"
	myFilter := "((value eq true) and (dog eq true)) or (primary eq true)"

	expression, err := getExpressionForFilter(myFilter)
	if err != nil {
		log.Fatalf("error parsing expression")
	}

	theMap, err := transformFilter(expression)
	if err != nil {
		fmt.Printf("error converting expression to filter string: %+v\n", err.Error())
	}

	if len(theMap) == 0 {
		fmt.Printf("no values in map \n")
	}

	for key, value := range theMap {
		fmt.Printf("key: %v, value: %v\n", key, value)
	}

}
func getExpressionForFilter(filter string) (*expr.Expression, error) {
	expression, err := expr.CompileFilter(filter)
	if err != nil {
		return nil, err
	}
	return expression, nil
}

//TODO  1 MAKE JUST ONE LAYER
// 		2 CONCATENATE AND STATEMENTS
func transformFilter(expression *expr.Expression) (map[string]string, error) {
	myMap := make(map[string]string, 1)
	left := expression.Left()
	right := expression.Right()

	if expression.IsRelationalOperator() {
		newMap, err := processBasicEqFilter(expression)
		//I think will still return map to use validly parsed filters if any
		if err != nil {
			return myMap, err
		}
		myMap = mergeMaps(myMap, newMap)
	} else if expression.IsLogicalOperator() {

		if left != nil {
			newMap1, err1 := transformFilter(left)
			if err1 != nil {
				return myMap, err1
			}
			myMap = mergeMaps(myMap, newMap1)
		}

		if right != nil {
			newMap2, err2 := transformFilter(right)
			if err2 != nil {
				return myMap, err2
			}
			myMap = mergeMaps(myMap, newMap2)
		}

	}

	return myMap, nil
}

// utility for merging two maps
func mergeMaps(a, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func processBasicEqFilter(filter *expr.Expression) (map[string]string, error) {
	myMap := make(map[string]string, 1)

	if filter.Token() != EqualsRelationalOperator {
		return nil, ErrUnsupportedDBFilter
	}

	left := filter.Left()
	right := filter.Right()

	if left.IsPath() && right.IsLiteral() {
		myMap[left.Token()] = right.Token()
	} else if right.IsPath() && left.IsLiteral() {
		myMap[right.Token()] = left.Token()
	} else {
		return nil, ErrUnsupportedDBFilter
	}

	return myMap, nil
}
