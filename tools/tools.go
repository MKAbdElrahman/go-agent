package tools

import (
	"encoding/json"
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

type ToolStore map[string]Tool

type Tool struct {
	Metadata FunctionMetaData `json:"metadata"`
	Doc      string           `json:"doc"`      // Generated prompt
	Function interface{}      `json:"function"` // The function object (as an interface{})
}

// CreateFunctionStore creates a function store for the given import path and funcMap.
// The Doc field contains the generated prompt, and the Function field contains the function object.
func CreateFunctionStore(importPath string, funcMap map[string]interface{}) (ToolStore, error) {
	store := make(ToolStore)

	// Iterate over the funcMap to process each function
	for functionName, function := range funcMap {
		// Step 1: Get the documentation for the function
		doc, err := getDocumentation(importPath, functionName)
		if err != nil {
			return nil, fmt.Errorf("failed to get documentation for function '%s': %v", functionName, err)
		}

		// Step 2: Parse the documentation into metadata
		metadata := parseDocumentation(functionName, doc)

		// Step 3: Generate the prompt
		prompt := generatePrompt(metadata)

		// Step 4: Store the metadata, prompt, and function object in the function store
		store[functionName] = Tool{
			Metadata: metadata,
			Doc:      prompt,
			Function: function,
		}
	}

	return store, nil
}

// Evaluate evaluates a function by name using the provided arguments.
func (store ToolStore) Evaluate(functionName string, args []interface{}) ([]interface{}, error) {
	// Step 1: Look up the function in the store
	entry, exists := store[functionName]
	if !exists {
		return nil, fmt.Errorf("function '%s' not found in store", functionName)
	}

	// Step 2: Use reflection to call the function
	functionValue := reflect.ValueOf(entry.Function)
	if functionValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("entry for '%s' is not a function", functionName)
	}

	// Step 3: Convert arguments to reflect.Value
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}

	// Step 4: Call the function with the provided arguments
	results := functionValue.Call(argValues)

	// Step 5: Convert results back to []interface{}
	resultInterfaces := make([]interface{}, len(results))
	for i, result := range results {
		resultInterfaces[i] = result.Interface()
	}

	return resultInterfaces, nil
}

func (s ToolStore) CombineToolsDoc() string {
	var combinedPrompt strings.Builder

	// Add a header for the combined prompt
	combinedPrompt.WriteString("=== Combined Function Prompts ===\n\n")

	// Iterate over the function store and append each prompt
	for functionName, entry := range s {
		combinedPrompt.WriteString(fmt.Sprintf("--- Function: %s ---\n", functionName))
		combinedPrompt.WriteString(entry.Doc)
		combinedPrompt.WriteString("\n\n") // Add spacing between functions
	}

	return combinedPrompt.String()
}

// GeneratePrompt creates a human-readable prompt for an AI agent to understand how to use the function.
func generatePrompt(meta FunctionMetaData) string {
	var prompt strings.Builder

	// Add function name and description
	prompt.WriteString(fmt.Sprintf("Function: %s\n", meta.FunctionName))
	prompt.WriteString(fmt.Sprintf("Description: %s\n", meta.Description))

	// Add parameters
	if len(meta.Params) > 0 {
		prompt.WriteString("Parameters:\n")
		for _, param := range meta.Params {
			prompt.WriteString(fmt.Sprintf("  - %s: %s\n", param.Name, param.Desc))
		}
	}

	// Add return values
	if len(meta.Return) > 0 {
		prompt.WriteString("Returns:\n")
		for _, ret := range meta.Return {
			prompt.WriteString(fmt.Sprintf("  - %s: %s\n", ret.Type, ret.Description))
		}
	}

	// Add constraints
	if len(meta.Constraints) > 0 {
		prompt.WriteString("Constraints:\n")
		for _, constraint := range meta.Constraints {
			prompt.WriteString(fmt.Sprintf("  - %s: %s\n", constraint.Condition, constraint.Desc))
		}
	}

	// Add examples
	if len(meta.Examples) > 0 {
		prompt.WriteString("Examples:\n")
		for _, example := range meta.Examples {
			prompt.WriteString(fmt.Sprintf("  - %s\n", example))
		}
	}

	return prompt.String()
}

// FunctionMetaData represents structured metadata extracted from the function documentation.
type FunctionMetaData struct {
	FunctionName string       `json:"function_name"`
	Description  string       `json:"description"`
	Params       []Param      `json:"params"`
	Return       []ReturnType `json:"return"`
	Examples     []string     `json:"examples"`
	Constraints  []Constraint `json:"constraints"`
}

// Constraint represents a constraint on the function or its parameters.
type Constraint struct {
	Condition string `json:"condition"`
	Desc      string `json:"desc"`
}

// Param represents a function parameter.
type Param struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

// ReturnType represents the return type and its description.
type ReturnType struct {
	Type        string `json:"type"`        // The return type (e.g., "float64")
	Description string `json:"description"` // A description of the return value
}

// ToJSON converts the FunctionMetaData struct to a JSON-formatted string.
func (meta FunctionMetaData) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ParseDocumentation parses the documentation string and extracts metadata.
func parseDocumentation(functionName, doc string) FunctionMetaData {
	meta := FunctionMetaData{
		FunctionName: functionName,
	}

	// Extract description (the first line of the doc string)
	lines := strings.Split(doc, "\n")
	if len(lines) > 0 {
		meta.Description = strings.TrimSpace(lines[0])
	}

	// Regex patterns
	paramRegex := regexp.MustCompile(`@param (\w+): (.+)`)
	returnRegex := regexp.MustCompile(`@return (\w+): (.+)`) // Updated to capture type and description
	constraintRegex := regexp.MustCompile(`@constraint (.+): (.+)`)
	exampleRegex := regexp.MustCompile(`@example:\s*(.+)`)

	// Parse the doc string line by line
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Extract parameters
		if matches := paramRegex.FindStringSubmatch(line); len(matches) == 3 {
			meta.Params = append(meta.Params, Param{
				Name: matches[1],
				Desc: matches[2],
			})
		}

		// Extract return values
		if matches := returnRegex.FindStringSubmatch(line); len(matches) == 3 {
			meta.Return = append(meta.Return, ReturnType{
				Type:        matches[1], // Capture the return type
				Description: matches[2], // Capture the return description
			})
		}

		// Extract constraints
		if matches := constraintRegex.FindStringSubmatch(line); len(matches) == 3 {
			meta.Constraints = append(meta.Constraints, Constraint{
				Condition: matches[1],
				Desc:      matches[2],
			})
		}

		// Extract examples
		if matches := exampleRegex.FindStringSubmatch(line); len(matches) == 2 {
			meta.Examples = append(meta.Examples, matches[1])
		}
	}

	return meta
}

// GetDocumentation retrieves the documentation for a function or type in a package.
func getDocumentation(importPath, name string) (string, error) {
	// Create a new file set.
	fset := token.NewFileSet()

	// Locate the package directory using go/build.
	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("failed to locate package: %v", err)
	}

	// Parse the package directory.
	pkgs, err := parser.ParseDir(fset, pkg.Dir, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse package: %v", err)
	}

	// Iterate over the packages (usually just one).
	for _, pkg := range pkgs {
		// Create a new doc.Package from the parsed package.
		docPkg := doc.New(pkg, importPath, doc.AllDecls)

		// Search for the type or function in the package.
		for _, t := range docPkg.Types {
			if t.Name == name {
				return t.Doc, nil
			}
			for _, method := range t.Methods {
				if method.Name == name {
					return method.Doc, nil
				}
			}
		}

		for _, fun := range docPkg.Funcs {
			if fun.Name == name {
				return fun.Doc, nil
			}
		}
	}

	return "", fmt.Errorf("function or type '%s' not found in package '%s'", name, importPath)
}
