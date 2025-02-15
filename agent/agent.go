package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"go-agent/evaluate"
	"go-agent/memory"
	"go-agent/tools"
	"strings"
	"text/template"
)

const promptTemplate = `You are a Go software engineer. Your task is to help users call mathematical functions in Go. 
Below are the available functions and their documentation. Respond to user requests in JSON format using the following template:

{
  "function": "<function_name>",
  "arguments": [<arg1>, <arg2>, ...]
}

Here are the functions and their documentation:
{{.Tools}}

{{if .Memory}}
### Interaction History:
{{.Memory}}
{{end}}

User Request: {{.UserRequest}}`

type LLMEngine interface {
	GenerateTokens(ctx context.Context, prompt string) (<-chan string, error)
}

type FunctionCall struct {
	Function  string `json:"function"`  // Function name (e.g., "Divide")
	Arguments []any  `json:"arguments"` // Function arguments (e.g., [4, 2])
}

type Agent struct {
	Engine        LLMEngine
	Prompt        string
	FunctionStore tools.ToolStore // Map of function names to their documentation prompts
	Memory        memory.Memory   // History of interactions
}

// NewAgent creates a new Agent instance with the specified LLM engine and prompts.
func NewAgent(engine LLMEngine, memory memory.Memory, tools tools.ToolStore) *Agent {
	return &Agent{
		Engine:        engine,
		Prompt:        promptTemplate,
		FunctionStore: tools,
		Memory:        memory,
	}
}

func (a *Agent) Execute(userRequest string) evaluate.EvaluationResult {
	functionCall, err := a.CallLLM(userRequest)
	if err != nil {
		return evaluate.EvaluationResult{Error: err}
	}
	return evaluate.Evaluate(a.FunctionStore, functionCall.Function, functionCall.Arguments)
}

func (a *Agent) CallLLM(userRequest string) (*FunctionCall, error) {
	// Execute the template to construct the final prompt
	tmpl, err := template.New("llmPrompt").Parse(a.Prompt)
	if err != nil {
		return nil, fmt.Errorf("error creating template: %w", err)
	}

	// Data for the template
	data := struct {
		Tools       string
		UserRequest string
		Memory      string
	}{
		Tools:       a.FunctionStore.CombineToolsDoc(),
		UserRequest: userRequest,
		Memory:      a.Memory.String(),
	}

	// Write the template output to a buffer (or directly to a string)
	var finalPrompt strings.Builder
	if err := tmpl.Execute(&finalPrompt, data); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}

	// Print the final prompt for debugging
	// fmt.Println("Final Prompt:\n", finalPrompt.String())

	// Generate tokens for the final prompt
	tokenCh, err := a.Engine.GenerateTokens(context.Background(), finalPrompt.String())
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	fmt.Println("------------------------------")
	fmt.Println("LLM Response")
	fmt.Println("------------------------------")

	// Collect the generated tokens
	var reply string
	for token := range tokenCh {
		fmt.Print(token)
		reply += token
	}
	fmt.Println()
	fmt.Println("------------------------------")

	var functionCall FunctionCall
	// Decode the LLM's response into the Go struct
	if err := json.Unmarshal([]byte(reply), &functionCall); err != nil {
		return nil, fmt.Errorf("error decoding LLM response: %w", err)
	}

	return &functionCall, nil
}
