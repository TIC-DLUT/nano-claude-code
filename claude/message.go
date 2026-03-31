package claude

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

const (
	ClaudeMessageRoleUser      = "user"
	ClaudeMessageRoleAssistant = "assistant"
)

type SingleStringMessage string

type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolUseBlock struct {
	Type        string         `json:"type"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Input       map[string]any `json:"input"`
	PartialJson string         `json:"-"`
}

type ToolResultBlock struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	ToolUseID string `json:"tool_use_id"`
}
