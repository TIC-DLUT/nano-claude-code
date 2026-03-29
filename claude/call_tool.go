package claude

import (
	"github.com/TIC-DLUT/nano-claude-code/errors"
)

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema struct {
		Type       string                        `json:"type"`
		Properties map[string]ToolPropertyDetail `json:"properties"`
		Required   []string                      `json:"required"`
	} `json:"input_schema"`
	Func func(input map[string]any) string `json:"-"`
}

type ToolPropertyDetail struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func NewTool(name string, description string, properties map[string]ToolPropertyDetail, required []string, toolFunc func(input map[string]any) string) (Tool, error) {
	if name == "" || description == "" {
		return Tool{}, errors.ClaudeCreateToolEmptyError
	}
	return Tool{
		Name:        name,
		Description: description,
		Func:        toolFunc,
		InputSchema: struct {
			Type       string                        "json:\"type\""
			Properties map[string]ToolPropertyDetail "json:\"properties\""
			Required   []string                      "json:\"required\""
		}{
			Type:       "object",
			Properties: properties,
			Required:   required,
		},
	}, nil
}

func (c *ClaudeClient) CallTools(model string, messages []Message, tools []Tool) ([]Message, error) {
	var err error = nil
	resMessages := []Message{}
	for {
		resMessages, err = c.Call(model, messages, tools)
		if err != nil {
			return []Message{}, err
		}

		messages = append(messages, resMessages...)

		continueFlag := false

		for _, item := range resMessages {
			switch item.Content.(type) {
			case ToolUseBlock:
				continueFlag = true
				toolUserItem := item.Content.(ToolUseBlock)
				content := ""
				for _, tool := range tools {
					if tool.Name == toolUserItem.Name {
						content = tool.Func(toolUserItem.Input)
					}
				}
				messages = append(messages, Message{
					Role: ClaudeMessageRoleUser,
					Content: []any{ToolResultBlock{
						Type:      "tool_result",
						ToolUseID: toolUserItem.ID,
						Content:   content,
					}},
				})

			}
		}

		if !continueFlag {
			break
		}
	}
	return resMessages, err
}
