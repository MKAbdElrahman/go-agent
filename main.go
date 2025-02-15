package main

import (
	"fmt"
	"go-agent/agent"
	"go-agent/calculator"
	"go-agent/llm"
	"go-agent/memory"
	"go-agent/tools"
)

func main() {
	importPath := "go-agent/calculator"

	// List of functions to generate documentation for
	functionNames := map[string]any{"Add": calculator.Add, "Subtract": calculator.Subtract, "Multiply": calculator.Multiply, "Divide": calculator.Divide}

	// Generate prompts for all specified functions
	tools, err := tools.CreateFunctionStore(importPath, functionNames)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// User request
	userRequest := "divide 4 and 3"

	// Initialize the LLM engine
	ollamaEngine, err := llm.NewOllamaEngine("llama3.1:8b")
	if err != nil {
		fmt.Printf("Error initializing LLM engine: %v\n", err)
		return
	}

	mem := memory.NewMemory()

	goDeveloper := agent.NewAgent(ollamaEngine, mem, tools)

	answer, err := goDeveloper.Execute(userRequest)
	if err != nil {
		panic(err)
	}

	fmt.Println(answer)

}
