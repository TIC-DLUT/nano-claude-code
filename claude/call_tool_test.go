package claude

import (
	"fmt"
	"testing"
)

func TestCallTools(t *testing.T) {
	client := newTestClient()
	get_weatherTool, _ := NewTool("get_weather", "获取一个城市当前的天气", map[string]ToolPropertyDetail{
		"city": {
			Type:        "string",
			Description: "城市的名字",
		},
	}, []string{"city"}, func(input map[string]any) string {
		fmt.Println("天气工具被调用", input)
		return "天气良好"
	})
	message, err := client.CallTools("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("大连天气怎么样"),
		},
	}, []Tool{get_weatherTool})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(message)
}

func TestCallStreamTools(t *testing.T) {
	client := newTestClient()
	get_weatherTool, _ := NewTool("get_weather", "获取一个城市当前的天气", map[string]ToolPropertyDetail{
		"city": {
			Type:        "string",
			Description: "城市的名字",
		},
	}, []string{"city"}, func(input map[string]any) string {
		fmt.Println("天气工具被调用", input)
		return "天气良好"
	})
	message, err := client.CallStreamTools("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("大连天气怎么样"),
		},
	}, []Tool{get_weatherTool}, func(m Message) bool {
		fmt.Println(m)
		return true
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(message)
}
