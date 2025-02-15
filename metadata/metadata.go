package metadata

import (
	"encoding/json"
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

func ExtractMetadata(importPath, name string) (FunctionMetaData, error) {
	doc, err := getDocumentation(importPath, name)
	if err != nil {
		return FunctionMetaData{}, err
	}

	meta := parseDocumentation(name, doc)
	return meta, nil
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

// parseDocumentation parses the documentation string and extracts metadata.
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

// getDocumentation retrieves the documentation for a function or type in a package.
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
