package claude

import (
	"bufio"
	"encoding/json"
	stdError "errors"
	"io"
	"reflect"
	"strings"

	"github.com/TIC-DLUT/nano-claude-code/errors"
	"resty.dev/v3"
)

type CallError struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	Type string `json:"type"`
}

type CallRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools"`
	System   string    `json:"system,omitempty"`
}

type CallResponse struct {
	Model        string        `json:"model"`
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Role         string        `json:"role"`
	Content      []interface{} `json:"content"`
	StopReason   string        `json:"stop_reason"`
	StopSequence interface{}   `json:"stop_sequence"`
	Usage        struct {
		InputTokens              int `json:"input_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
		CacheCreation            struct {
			Ephemeral5MInputTokens int `json:"ephemeral_5m_input_tokens"`
			Ephemeral1HInputTokens int `json:"ephemeral_1h_input_tokens"`
		} `json:"cache_creation"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type CallStreamResponse struct {
	Type         string `json:"type"`
	Index        int    `json:"index"`
	ContentBlock struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Type     string `json:"type"`
		Text     string `json:"text"`
		Thinking string `json:"thinking"`
	} `json:"content_block"`
	Delta struct {
		Type           string `json:"type"`
		Text           string `json:"text"`
		Thinking       string `json:"thinking"`
		PartialJson    string `json:"partial_json"`
		SignatureDelta string `json:"signature_delta"`
	} `json:"delta"`
}

func frontCall(httpClient *resty.Client, inBaseUrl string, apiKey string, model string, messages []Message, stream bool, tools []Tool, system string) (CallResponse, *resty.Response, error) {
	// 防止 https://example ，统一为 https://example/
	baseurl := inBaseUrl
	if inBaseUrl[len(inBaseUrl)-1] != '/' {
		baseurl += "/"
	}

	// tools 去重，防止出现重名 tool
	toolHashMap := make(map[string]bool)
	for _, item := range tools {
		_, ok := toolHashMap[item.Name]
		if ok && toolHashMap[item.Name] == true {
			return CallResponse{}, nil, errors.ClaudeToolRepeat
		}
		toolHashMap[item.Name] = true
	}

	// 合并同侧的 user 和 assistant 消息
	currentRole := ClaudeMessageRoleUser
	requestMessages := []Message{}

	if len(messages) > 0 {
		currentRole = messages[0].Role
		requestMessages = append(requestMessages, Message{Role: currentRole, Content: []any{}})
	}

	for i := 0; i < len(messages); i++ {
		if messages[i].Role != currentRole {
			currentRole = messages[i].Role
			requestMessages = append(requestMessages, Message{Role: currentRole, Content: []any{}})
		}

		newContent := requestMessages[len(requestMessages)-1].Content.([]any)
		if reflect.TypeOf(messages[i].Content) == reflect.TypeOf(SingleStringMessage("")) {
			newContent = append(newContent, TextBlock{Type: "text", Text: string(messages[i].Content.(SingleStringMessage))})
		} else {
			newContent = append(newContent, messages[i].Content)
		}
		requestMessages[len(requestMessages)-1].Content = newContent
	}

	// 构造请求发送，并获取响应
	res := CallResponse{}

	requestBody := CallRequest{
		Stream:   stream,
		Model:    model,
		Messages: requestMessages,
		System:   system,
	}

	if len(tools) != 0 {
		requestBody.Tools = tools
	}

	httpRequest := httpClient.R().
		SetHeader("x-api-key", apiKey).
		SetBody(requestBody)

	//控制是否为流式相应，若为非流式，则返回预设的响应结构，若为流式，则保留原始响应，交由上层函数处理
	if stream {
		httpRequest.SetDoNotParseResponse(true)
	} else {
		httpRequest.SetResult(&res)
	}
	httpRes, err := httpRequest.Post(baseurl + "v1/messages")

	// 解析响应错误
	if httpRes.StatusCode() != 200 {
		httpBody, _ := io.ReadAll(httpRes.Body)
		defer httpRes.Body.Close()
		errMessage := CallError{}
		json.Unmarshal(httpBody, &errMessage)
		return res, httpRes, stdError.New(errMessage.Error.Message)
	}

	return res, httpRes, err
}

func (c *ClaudeClient) Call(model string, system string, messages []Message, tools []Tool) ([]Message, error) {
	// 发送请求，获取响应信息
	res, _, err := frontCall(c.httpClient, c.baseUrl, c.apiKey, model, messages, false, tools, system)
	if err != nil {
		return []Message{}, err
	}

	// 将响应信息解析并转化为对应的 ContentBlock 返回
	resMessages := []Message{}
	for _, item := range res.Content {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return []Message{}, errors.ClaudeClientCallFormatError
		}

		messageType, ok := itemMap["type"].(string)
		if !ok {
			return []Message{}, errors.ClaudeClientCallFormatError
		}

		switch messageType {
		case "text":
			resMessages = append(resMessages, Message{
				Role: ClaudeMessageRoleAssistant,
				Content: TextBlock{
					Type: "text",
					Text: itemMap["text"].(string),
				},
			})
		case "thinking":
			resMessages = append(resMessages, Message{
				Role: ClaudeMessageRoleAssistant,
				Content: ThinkingBlock{
					Type:      "thinking",
					Thinking:  itemMap["thinking"].(string),
					Signature: itemMap["signature"].(string),
				},
			})
		case "tool_use":
			resMessages = append(resMessages, Message{
				Role: ClaudeMessageRoleAssistant,
				Content: ToolUseBlock{
					Type:  "tool_use",
					Name:  itemMap["name"].(string),
					Input: itemMap["input"].(map[string]any),
					ID:    itemMap["id"].(string),
				},
			})
		default:
			return []Message{}, errors.ClaudeClientCallFormatError
		}
	}

	return resMessages, nil
}

func (c *ClaudeClient) CallStream(model string, system string, messages []Message, tools []Tool, dealFunc func(Message) bool) ([]Message, error) {
	// 发送请求，获取响应信息
	_, originHttpRes, err := frontCall(c.httpClient, c.baseUrl, c.apiKey, model, messages, true, tools, system)
	if err != nil {
		return []Message{}, err
	}

	// 构建 reader 读取 SSE 响应内容
	reader := bufio.NewReader(originHttpRes.Body)
	defer originHttpRes.Body.Close()

	resMessages := []Message{}

	// 解析 SSE 响应内容，并将其转化为对应的 ContentBlock 返回
	for {
		eventStr, err := reader.ReadString('\n')
		if stdError.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return []Message{}, err
		}

		if strings.Trim(eventStr, " ") == "" {
			continue
		}

		if strings.HasPrefix(eventStr, "data: ") {
			data := eventStr[6:]
			dataDetail := CallStreamResponse{}
			err := json.Unmarshal([]byte(data), &dataDetail)
			if err != nil {
				continue
			}

			switch dataDetail.Type {
			// ContentBlock 的开始，根据类型创建对应的 ContentBlock，并添加在 resMessages 中
			case "content_block_start":
				resMessages = append(resMessages, Message{
					Role: ClaudeMessageRoleAssistant,
				})
				var content any
				switch dataDetail.ContentBlock.Type {
				case "text":
					content = TextBlock{
						Type: "text",
						Text: "",
					}
				case "thinking":
					content = ThinkingBlock{
						Type:      "thinking",
						Thinking:  "",
						Signature: "",
					}
				case "tool_use":
					content = ToolUseBlock{
						Type: "tool_use",
						ID:   dataDetail.ContentBlock.ID,
						Name: dataDetail.ContentBlock.Name,
					}
				case "":
					continue
				}

				resMessages[len(resMessages)-1].Content = content

			// ContentBlock 的增量更新，根据类型将增量内容添加在 resMessages 中对应的 ContentBlock 上
			case "content_block_delta":
				continueFlag := true
				switch resMessages[len(resMessages)-1].Content.(type) {
				case TextBlock:
					resMessages[len(resMessages)-1].Content = TextBlock{
						Type: "text",
						Text: resMessages[len(resMessages)-1].Content.(TextBlock).Text + dataDetail.Delta.Text,
					}
					continueFlag = dealFunc(Message{
						Role: ClaudeMessageRoleAssistant,
						Content: TextBlock{
							Type: "text",
							Text: dataDetail.Delta.Text,
						},
					})
				case ThinkingBlock:
					resMessages[len(resMessages)-1].Content = ThinkingBlock{
						Type:      "thinking",
						Signature: dataDetail.Delta.SignatureDelta,
						Thinking:  resMessages[len(resMessages)-1].Content.(ThinkingBlock).Thinking + dataDetail.Delta.Thinking,
					}
					continueFlag = dealFunc(Message{
						Role: ClaudeMessageRoleAssistant,
						Content: ThinkingBlock{
							Type:      "thinking",
							Thinking:  dataDetail.Delta.Thinking,
							Signature: dataDetail.Delta.SignatureDelta,
						},
					})
				case ToolUseBlock:
					changeContent := resMessages[len(resMessages)-1].Content.(ToolUseBlock)
					changeContent.PartialJson += dataDetail.Delta.PartialJson
					resMessages[len(resMessages)-1].Content = changeContent
					continueFlag = dealFunc(Message{
						Role: ClaudeMessageRoleAssistant,
						Content: ToolUseBlock{
							Type:        "tool_use",
							ID:          changeContent.ID,
							Name:        changeContent.Name,
							PartialJson: dataDetail.Delta.PartialJson,
						},
					})
				}
				if !continueFlag {
					return resMessages, nil
				}
			}
		}
	}
	// 将 ToolUseBlock 中的 PartialJson 解析为 Input
	for i := 0; i < len(resMessages); i++ {
		if reflect.TypeOf(resMessages[i].Content) == reflect.TypeOf(ToolUseBlock{}) {
			changeBlock := resMessages[i].Content.(ToolUseBlock)
			inputMap := make(map[string]any)
			err := json.Unmarshal([]byte(changeBlock.PartialJson), &inputMap)
			if err != nil {
				return resMessages, errors.ClaudeToolStreamPartParseError
			}
			changeBlock.Input = inputMap
			resMessages[i].Content = changeBlock
		}
	}
	return resMessages, nil
}
