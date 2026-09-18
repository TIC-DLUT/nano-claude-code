package agent

import (
	"github.com/TIC-DLUT/nano-claude-code/claude"
	"github.com/spf13/viper"
)

func (a *Agent) ChatStream(message string, callback func(string)) {
	lastToolCallID := ""
	messages, err := a.sessionManager.BuildSessionContext()
	if err != nil {
		panic(err)
	}
	newMessage := claude.Message{
		Role:    claude.ClaudeMessageRoleUser,
		Content: claude.SingleStringMessage(message),
	}

	messages = append(messages, newMessage)
	a.sessionManager.Append(newMessage)

	resMessages, err := a.apiClient.CallStreamTools(viper.GetString("llm.model"), GetNowSystemPrompt(), messages,
		a.tools, func(m claude.Message) bool {
			switch m.Content.(type) {
			case claude.TextBlock:
				callback(m.Content.(claude.TextBlock).Text)
			case claude.ToolUseBlock:
				tooluse := m.Content.(claude.ToolUseBlock)
				if tooluse.ID != lastToolCallID {
					lastToolCallID = tooluse.ID
					callback("\n[tool_use] " + tooluse.Name + "\n")
				}
			}
			return true
		})
	if err != nil {
		panic(err)
	}

	a.sessionManager.Append(resMessages...)
}
