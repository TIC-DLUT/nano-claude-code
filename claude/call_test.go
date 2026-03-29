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
	})
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
	}, func(m Message) bool {
		fmt.Println(m.Content)
		return true
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println("总消息", resMessages)
}
