package toolstore

import (
	"errors"
	"fmt"
	"go-agent/metadata"
	"go-agent/tools/evaluation"
	"log/slog"
)

var (
	ErrToolNotFound       = errors.New("tool not found")
	ErrToolExists         = errors.New("tool already exists")
	ErrMetadataExtraction = errors.New("failed to extract metadata")
)

// ToolStore is a thread-safe collection of tools indexed by their names.
type ToolStore struct {
	tools  map[string]evaluation.Tool
	logger *slog.Logger
}

// NewToolStore creates a new ToolStore with an optional logger.
func NewToolStore(logger *slog.Logger) *ToolStore {
	if logger == nil {
		logger = slog.Default() // Use the default logger if none is provided
	}
	return &ToolStore{
		tools:  make(map[string]evaluation.Tool),
		logger: logger,
	}
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

		store.AddTool(functionName, evaluation.Tool{
			Metadata: metadata,
			Function: function,
		})
	}

	return store, nil
}

// AddTool adds a new tool to the ToolStore.
func (ts *ToolStore) AddTool(name string, tool evaluation.Tool) error {
	if _, exists := ts.tools[name]; exists {
		ts.logger.Error("Tool already exists", "name", name)
		return ErrToolExists
	}

	ts.tools[name] = tool
	ts.logger.Info("Tool added", "name", name)
	return nil
}

// GetTool retrieves a tool from the ToolStore by name.
func (ts *ToolStore) GetTool(name string) (evaluation.Tool, error) {
	tool, exists := ts.tools[name]
	if !exists {
		ts.logger.Error("Tool not found", "name", name)
		return evaluation.Tool{}, ErrToolNotFound
	}
	return tool, nil
}

// RemoveTool removes a tool from the ToolStore by name.
func (ts *ToolStore) RemoveTool(name string) error {

	if _, exists := ts.tools[name]; !exists {
		ts.logger.Error("Tool not found", "name", name)
		return ErrToolNotFound
	}

	delete(ts.tools, name)
	ts.logger.Info("Tool removed", "name", name)
	return nil
}

// ListTools returns a list of all tool names in the ToolStore.
func (ts *ToolStore) ListToolNames() []string {

	toolNames := make([]string, 0, len(ts.tools))
	for name := range ts.tools {
		toolNames = append(toolNames, name)
	}

	return toolNames
}

func (ts *ToolStore) Tools() map[string]evaluation.Tool {
	return ts.tools
}
