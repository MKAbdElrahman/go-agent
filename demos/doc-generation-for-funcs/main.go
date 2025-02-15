package main

import (
	"fmt"
	"go-agent/calculator"
	"go-agent/tools"
)

func main() {
	importPath := "go-agent/calculator"

	store, err := tools.CreateFunctionStore(importPath, map[string]interface{}{"Add": calculator.Add})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	result, err := store.Evaluate("Add", []any{1., 2.})

	fmt.Println(result...)
}
