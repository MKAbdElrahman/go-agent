package main

import (
	"fmt"
	"go-agent/calculator"
	"go-agent/tools/toolstore"
)

func main() {
	importPath := "go-agent/calculator"

	store, err := toolstore.NewFunctionStoreFromPkg(importPath, calculator.FunctionRegistry(), nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	f, err := store.GetTool("Divide")
	if err != nil {
		fmt.Printf("function '%s' not found in tool store", f.Function)
		return
	}
	result, err := f.Evaluate([]any{1, 0})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Result: ", result)

}
