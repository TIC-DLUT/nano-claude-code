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
	// 这里是工具的运行函数，不属于官方的 Tool 定义，但是我们需要其用来方便的运行 Tool 函数，不需要 json 化
	// 否则会出现 400 请求出错
	Func func(input map[string]any) string `json:"-"`
}

type ToolPropertyDetail struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func NewTool(name string, description string, properties map[string]ToolPropertyDetail, required []string, toolFunc func(input map[string]any) string) (Tool, error) {
	// name 和 description 不能为空
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
		// 将消息中的 ToolUse 请求执行，并将结果添加在 toolMessages 中
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
	realResMessages := []Message{}
	resMessages := []Message{}
	// 循环调用，因为模型可能在两次或多次请求中均需要使用 Tools
	for {
		resMessages, err = c.Call(model, messages, tools)
		if err != nil {
			return []Message{}, err
		}

		messages = append(messages, resMessages...)
		realResMessages = append(realResMessages, resMessages...)

		// 执行 ToolUse
		continueFlag := false
		toolMessages := []Message{}
		continueFlag, messages, toolMessages = toolCall(tools, messages, resMessages)
		realResMessages = append(realResMessages, toolMessages...)

		if !continueFlag {
			break
		}
	}
	return realResMessages, err
}

func (c *ClaudeClient) CallStreamTools(model string, messages []Message, tools []Tool, dealFunc func(Message) bool) ([]Message, error) {
	realResMessages := []Message{}
	// 循环调用，因为模型可能在两次或多次请求中均需要使用 Tools
	for {
		resMessages, err := c.CallStream(model, messages, tools, dealFunc)
		if err != nil {
			return resMessages, err
		}
		messages = append(messages, resMessages...)
		realResMessages = append(realResMessages, resMessages...)

		// 执行 ToolUse
		continueFlag := false
		toolMessages := []Message{}
		continueFlag, messages, toolMessages = toolCall(tools, messages, resMessages)
		realResMessages = append(realResMessages, toolMessages...)
		if !continueFlag {
			break
		}
	}
	return realResMessages, nil
}
