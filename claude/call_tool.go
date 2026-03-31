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

// 如果需要继续循环，返回 true
func toolCall(tools []Tool, messages []Message, resMessages []Message) (bool, []Message, []Message) {
	continueFlag := false
	toolMessages := []Message{}
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
			toolResultMessage := Message{
				Role: ClaudeMessageRoleUser,
				Content: ToolResultBlock{
					Type:      "tool_result",
					ToolUseID: toolUserItem.ID,
					Content:   content,
				},
			}
			messages = append(messages, toolResultMessage)
			toolMessages = append(toolMessages, toolResultMessage)
		}
	}
	return continueFlag, messages, toolMessages
}

func (c *ClaudeClient) CallTools(model string, messages []Message, tools []Tool) ([]Message, error) {
	var err error = nil
	realresMessages := []Message{}
	resMessages := []Message{}
	for {
		resMessages, err = c.Call(model, messages, tools)
		if err != nil {
			return []Message{}, err
		}

		messages = append(messages, resMessages...)
		realresMessages = append(realresMessages, resMessages...)

		continueFlag := false
		toolmessages := []Message{}
		continueFlag, messages, toolmessages = toolCall(tools, messages, resMessages)
		realresMessages = append(realresMessages, toolmessages...)

		if !continueFlag {
			break
		}
	}
	return realresMessages, err
}

func (c *ClaudeClient) CallStreamTools(model string, messages []Message, tools []Tool, dealFunc func(Message) bool) ([]Message, error) {
	realResMessages := []Message{}
	for {
		resMessages, err := c.CallStream(model, messages, tools, dealFunc)
		messages = append(messages, resMessages...)
		realResMessages = append(realResMessages, resMessages...)
		if err != nil {
			return resMessages, err
		}
		cotinuesFlag := false
		ToolMessages := []Message{}
		cotinuesFlag, messages, ToolMessages = toolCall(tools, messages, resMessages)
		realResMessages = append(realResMessages, ToolMessages...)
		if !cotinuesFlag {
			break
		}
	}
	return realResMessages, nil
}
