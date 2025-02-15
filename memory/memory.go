package memory

import (
	"fmt"
	"strings"
)

type Memory interface {
	AppendInteraction(request string, response string)
	String() string
}

// memoryImpl is the concrete implementation of the Memory interface.
type memoryImpl struct {
	interactions []interaction
}

// interaction represents a single exchange between the user and the agent.
type interaction struct {
	Request  string
	Response string
}

// NewMemory creates a new Memory instance with empty interaction history.
func NewMemory() Memory {
	return &memoryImpl{
		interactions: make([]interaction, 0),
	}
}

// AppendInteraction adds a new interaction to the memory.
func (m *memoryImpl) AppendInteraction(request, response string) {
	m.interactions = append(m.interactions, interaction{
		Request:  request,
		Response: response,
	})
}

// String returns a formatted string of the interaction history.
func (m *memoryImpl) String() string {
	if len(m.interactions) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, interaction := range m.interactions {
		sb.WriteString(fmt.Sprintf("User Request: %s\nAgent Response: %s\n\n",
			interaction.Request, interaction.Response))
	}

	// Remove the last newline to avoid trailing whitespace
	return strings.TrimSpace(sb.String())
}
