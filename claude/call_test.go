package claude

import (
	"fmt"
	"os"
	"testing"
)

func newTestClient() *ClaudeClient {
	client, _ := NewClient(os.Getenv("baseurl"), os.Getenv("apikey"))
	return client
}

func TestCall(t *testing.T) {
	client := newTestClient()
	message, err := client.Call("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("你好"),
		},
	}, []Tool{})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(message)
}

func TestCallStream(t *testing.T) {
	client := newTestClient()
	resMessages, err := client.CallStream("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("你好"),
		},
	}, []Tool{}, func(m Message) bool {
		fmt.Println(m.Content)
		return true
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("总消息", resMessages)
}

func TestCallWithTools(t *testing.T) {
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
	message, err := client.Call("claude-sonnet-4-6", []Message{
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

func TestCallStreamWithTools(t *testing.T) {
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
	message, err := client.CallStream("claude-sonnet-4-6", []Message{
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
