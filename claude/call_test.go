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

func TestCallWithImageUrl(t *testing.T) {
	client := newTestClient()
	message, err := client.Call("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("这个图片是啥"),
		}, {
			Role:    ClaudeMessageRoleUser,
			Content: NewImageBlockFromUrl("https://twfood.cc/img/code/A1/_.jpg"),
		}}, []Tool{})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(message)
}

func TestCallWithImageBase64(t *testing.T) {
	client := newTestClient()
	//从文件读取base64编码后的图片数据
	base64data, err := os.ReadFile("./test_data/base64.txt")
	if err != nil {
		t.Error(err.Error())
	}
	base64dataStr := string(base64data)
	message, err := client.Call("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("这个图片是啥"),
		}, {
			Role:    ClaudeMessageRoleUser,
			Content: NewImageBlockFromBase64("image/jpeg", base64dataStr),
		}}, []Tool{})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(message)
}

func TestCallStreamImageBase64(t *testing.T) {
	client := newTestClient()
	//从文件读取base64编码后的图片数据
	base64data, err := os.ReadFile("test_data/base64.txt")
	if err != nil {
		t.Error(err.Error())
	}
	base64dataStr := string(base64data)
	resMessages, err := client.CallStream("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("这个图片是啥，尽可能说长一点"),
		}, {
			Role:    ClaudeMessageRoleUser,
			Content: NewImageBlockFromBase64("image/jpeg", base64dataStr),
		}}, []Tool{}, func(m Message) bool {
		fmt.Println(m.Content)
		return true
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("总消息", resMessages)
}

func TestCallStreamImageUrl(t *testing.T) {
	client := newTestClient()
	resMessages, err := client.CallStream("claude-sonnet-4-6", []Message{
		{
			Role:    ClaudeMessageRoleUser,
			Content: SingleStringMessage("这个图片是啥，尽可能说长一点"),
		}, {
			Role:    ClaudeMessageRoleUser,
			Content: NewImageBlockFromUrl("https://twfood.cc/img/code/A1/_.jpg")},
	}, []Tool{}, func(m Message) bool {
		fmt.Println(m.Content)
		return true
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("总消息", resMessages)
}
