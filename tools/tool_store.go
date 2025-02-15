package tools

import (
	"errors"
	"fmt"
	"go-agent/metadata"
	"log/slog"
	"strings"
	"sync"
)

var (
	ErrNotAFunction       = errors.New("entry is not a function")
	ErrArgumentMismatch   = errors.New("argument count mismatch")
	ErrArgumentType       = errors.New("argument type mismatch")
	ErrFunctionPanic      = errors.New("function execution panicked")
	ErrMetadataExtraction = errors.New("failed to extract metadata")
)

// ToolStore is a thread-safe collection of tools indexed by their names.
type ToolStore struct {
	mu     sync.RWMutex
	tools  map[string]Tool
	logger *slog.Logger
}

// NewToolStore creates a new ToolStore with an optional logger.
func NewToolStore(logger *slog.Logger) *ToolStore {
	if logger == nil {
		logger = slog.Default() // Use the default logger if none is provided
	}
	return &ToolStore{
		tools:  make(map[string]Tool),
		logger: logger,
	}
}

// AddTool adds a new tool to the ToolStore.
func (ts *ToolStore) AddTool(name string, tool Tool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tools[name] = tool
	ts.logger.Info("Tool added", "name", name)
}

// GetTool retrieves a tool from the ToolStore by name.
func (ts *ToolStore) GetTool(name string) (Tool, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	tool, exists := ts.tools[name]
	return tool, exists
}

// CreateFunctionStore creates a ToolStore from the given import path and function map.
func NewFunctionStoreFromPkg(importPath string, funcMap map[string]interface{}, logger *slog.Logger) (*ToolStore, error) {
	store := NewToolStore(logger)

	for functionName, function := range funcMap {
		metadata, err := metadata.ExtractMetadata(importPath, functionName)
		if err != nil {
			logger.Error("Failed to extract metadata", "function", functionName, "error", err)
			return nil, fmt.Errorf("%w: %v", ErrMetadataExtraction, err)
		}

		store.AddTool(functionName, Tool{
			Metadata: metadata,
			Doc:      generatePrompt(metadata),
			Function: function,
		})
	}

	return store, nil
}

// CombineToolsDoc combines the documentation of all tools in the ToolStore.
func (ts *ToolStore) CombineToolsDoc() string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var combinedPrompt strings.Builder
	combinedPrompt.WriteString("=== Combined Function Prompts ===\n\n")

	for functionName, entry := range ts.tools {
		combinedPrompt.WriteString(fmt.Sprintf("--- Function: %s ---\n", functionName))
		combinedPrompt.WriteString(entry.Doc)
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
