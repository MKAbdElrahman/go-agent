package main

import (
	"fmt"
	"go-agent/agent"
	"go-agent/calculator"
	"go-agent/llm"
	"go-agent/tools/toolstore"
)

func main() {
	// List of user requests
	userRequests := []string{
		"What is the sum of one , three and 6?",
		"What is the square root of 36?",
		"Divide 100 by 0",
		"Add 3 and 4.",
		"What is the square root of -24?",
		"Subtract 10 from 20.",
		"Multiply 5 by 6.",
		"What is the factorial of 5?",
		"What is 2 raised to the power of 8?",
		"What is the sine of 90 degrees?",
	}

	// Initialize the LLM engine
	ollamaEngine, err := llm.NewOllamaEngine("llama3.1:8b")
	if err != nil {
		fmt.Printf("Error initializing LLM engine: %v\n", err)
		return
	}

	// Get public functions from the calculator package
	// Create a function store for the tools
	toolStore, err := toolstore.NewFunctionStoreFromPkg("go-agent/calculator", calculator.FunctionRegistry(), nil)
	if err != nil {
		fmt.Printf("Error creating function store: %v\n", err)
		return
	}

	// Initialize the agent
	goDeveloper := agent.NewAgent(ollamaEngine, toolStore)

	// Evaluate each user request
	for _, request := range userRequests {
		fmt.Printf("User Request: %s\n", request)

		// Execute the request using the agent
		response, err := goDeveloper.Execute(request)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			fmt.Println("-----------------------------")
		} else {

			// Print the response
			fmt.Printf("Response: %+v\n", response)
			fmt.Println("-----------------------------")

		}

	}
}
