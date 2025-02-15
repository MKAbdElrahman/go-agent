package main

import (
	"fmt"
	"go-agent/calculator"
	"go-agent/tools"
)

func main() {
	importPath := "go-agent/calculator"

	store, err := tools.CreateFunctionStore(importPath, map[string]interface{}{"Divide": calculator.Divide})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	result := store.Evaluate("Divide", []any{1, 0})

	fmt.Println("Result: ", result.Result)
	fmt.Println("Error: ", result.Error)

}
