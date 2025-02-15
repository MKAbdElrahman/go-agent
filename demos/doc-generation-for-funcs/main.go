package main

import (
	"fmt"
	"go-agent/calculator"
	"go-agent/tools"
)

func main() {
	importPath := "go-agent/calculator"

	store, err := tools.NewFunctionStoreFromPkg(importPath, calculator.FunctionRegistry(), nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	f, ok := store.GetTool("Divide")
	if !ok {
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
