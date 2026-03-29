package claude

import (
	"bufio"
	"encoding/json"
	stdError "errors"
	"io"
	"strings"

	"github.com/TIC-DLUT/nano-claude-code/errors"
	"resty.dev/v3"
)

type CallRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
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
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

func frontCall(httpClient *resty.Client, inBaseUrl string, apiKey string, model string, messages []Message, stream bool) (CallResponse, *resty.Response, error) {
	// 防止 https://example/
	baseurl := inBaseUrl
	if inBaseUrl[len(inBaseUrl)-1] != '/' {
		baseurl += "/"
	}

	res := CallResponse{}

	httpRequest := httpClient.R().
		SetHeader("x-api-key", apiKey).
		SetBody(CallRequest{
			Stream:   stream,
			Model:    model,
			Messages: messages,
		})

	if stream {
		httpRequest.SetDoNotParseResponse(true)
	} else {
		httpRequest.SetResult(&res)
	}
	httpRes, err := httpRequest.Post(baseurl + "v1/messages")

	return res, httpRes, err
}

func (c *ClaudeClient) Call(model string, messages []Message) ([]Message, error) {
	res, _, err := frontCall(c.httpClient, c.baseUrl, c.apiKey, model, messages, false)
	if err != nil {
		return []Message{}, err
	}

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
		default:
			return []Message{}, errors.ClaudeClientCallFormatError
		}
	}

	return resMessages, nil
}

func (c *ClaudeClient) CallStream(model string, messages []Message, dealFunc func(Message) bool) ([]Message, error) {
	_, originHttpRes, err := frontCall(c.httpClient, c.baseUrl, c.apiKey, model, messages, true)
	if err != nil {
		return []Message{}, err
	}

	reader := bufio.NewReader(originHttpRes.Body)
	defer originHttpRes.Body.Close()

	resMessages := []Message{}

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
				case "":
					continue
				}

				resMessages[len(resMessages)-1].Content = content

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
				}
				if !continueFlag {
					return resMessages, nil
				}
			}
		}
	}
	return resMessages, nil
}
