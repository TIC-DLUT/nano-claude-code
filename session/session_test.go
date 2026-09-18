package session

import (
	"testing"

	"github.com/TIC-DLUT/nano-claude-code/claude"
)

func TestSessionManager(t *testing.T) {
	sessionManager, err := NewSessionManager("test")
	if err != nil {
		t.Error(err)
	}
	sessionManager.Append(claude.Message{
		Role: claude.ClaudeMessageRoleAssistant,
		Content: claude.ToolUseBlock{
			Type: "tool_use",
			ID:   "tool_use_id",
			Name: "tool_use_test",
			Input: map[string]any{
				"input": "input",
			},
		},
	})

	messages, err := sessionManager.BuildSessionContext()
	if err != nil {
		t.Error(err)
	}

	t.Log(messages)
}
