package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"go-agent/metadata"
	"go-agent/tools/toolstore"
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
	FunctionStore *toolstore.ToolStore // Map of function names to their documentation prompts
}

// NewAgent creates a new Agent instance with the specified LLM engine and prompts.
func NewAgent(engine LLMEngine, tools *toolstore.ToolStore) *Agent {
	return &Agent{
		Engine:        engine,
		Prompt:        promptTemplate,
		FunctionStore: tools,
	}
}

func (a *Agent) Execute(userRequest string) ([]any, error) {
	functionCall, err := a.CallLLM(userRequest)
	if err != nil {
		return nil, err
	}

	tool, err := a.FunctionStore.GetTool(functionCall.Function)
	if err != nil {
		return nil, fmt.Errorf("function '%s' not found in tool store", functionCall.Function)
	}

	return tool.Evaluate(functionCall.Arguments)
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
	}{
		Tools:       combineToolsDoc(a.FunctionStore),
		UserRequest: userRequest,
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

// CombineToolsDoc combines the documentation of all tools in the ToolStore.
func combineToolsDoc(ts *toolstore.ToolStore) string {

	var combinedPrompt strings.Builder
	combinedPrompt.WriteString("=== Combined Function Prompts ===\n\n")

	for functionName, entry := range ts.Tools() {
		combinedPrompt.WriteString(fmt.Sprintf("--- Function: %s ---\n", functionName))
		combinedPrompt.WriteString(generatePrompt(entry.Metadata))
		combinedPrompt.WriteString("\n\n")
	}

	return combinedPrompt.String()
}

// generatePrompt creates a human-readable prompt for a function based on its metadata.
func generatePrompt(meta metadata.FunctionMetaData) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Function: %s\nDescription: %s\n", meta.FunctionName, meta.Description))

	if len(meta.Params) > 0 {
		prompt.WriteString("Parameters:\n")
		for _, param := range meta.Params {
			prompt.WriteString(fmt.Sprintf("  - %s: %s\n", param.Name, param.Desc))
		}
	}

	if len(meta.Return) > 0 {
		prompt.WriteString("Returns:\n")
		for _, ret := range meta.Return {
			prompt.WriteString(fmt.Sprintf("  - %s: %s\n", ret.Type, ret.Description))
		}
	}

	if len(meta.Constraints) > 0 {
		prompt.WriteString("Constraints:\n")
		for _, constraint := range meta.Constraints {
			prompt.WriteString(fmt.Sprintf("  - %s: %s\n", constraint.Condition, constraint.Desc))
		}
	}

	if len(meta.Examples) > 0 {
		prompt.WriteString("Examples:\n")
		for _, example := range meta.Examples {
			prompt.WriteString(fmt.Sprintf("  - %s\n", example))
		}
	}

	return prompt.String()
}
