# 从零开始的Claude Code(一)

本教程致力于让读者真正的明白一个AI应用是如何实现的，因此我们选择从自行撰写一个 claude 协议包开始到一个能真正实现 vibe coding 的 claude code agent。

> 本章所用内容：
>
> Claude协议官方文档： [Create a Message - Claude API Reference](https://platform.claude.com/docs/en/api/messages/create) 
>
> golang网络请求包官方文档： [Welcome | Resty](https://resty.dev/) 

## 程序入口

**nano claude code** 作为一个要真正跑起来的应用，无论我们的实际内容如何设计，我们一定需要有一个程序入口。因此，我们可以先初始化项目，并创建程序入口。然后我们再根据我们的第一需求来填充我们的后续内容。

```shell
# 读者可根据自己的需要更改项目名称
go mod init github.com/TIC-DLUT/nano-claude-code

# 创建程序入口
mkdir cmd
touch cmd/main.go
```

在文件 `cmd/main.go` 中创建 `main` 函数：

```go
package main

func main() {

}
```

此时项目结构如下：

```
nano-claude-code
├─cmd
│  └─main.go
└─go.mod
```

## 基本 claude 协议封装

由于本教程致力于让读者真正的明白一个AI应用是如何实现的，因此我们选择从自行撰写一个 claude 协议包开始到一个能真正实现 vibe coding 的 claude code agent。

> ❗注意：claude协议与openai协议不同，如果读者需要使用遵循openai协议的ai，可以参考该部分的封装思路，但不能完全照搬。
>
> ​	部分明显的不同笔者会尽可能做出标注。

### 简单初始化

那么一个调用 api 的协议应该怎么写呢，首先我们要看官方文档是怎么定义对应的协议的。

打开 claude 的官方文档，我们找到 `API Reference` 中的 `Create a Message` 这个是我们要封装的核心内容，那我们来看一下它具体是什么样子的。可以看到这里面由相当多的参数，简直眼花缭乱。该如何开始呢？我们选择从简到繁，先封装最重要的参数。使用过API调用的朋友都知道，调用一个API最基本的参数就是 `baseurl` `api-key`  这两个。这两个参数分别明确了“我们要向哪里发起请求”，“如何通过验证”这两个问题的答案。所以，我们首先实现这两个参数的封装。

作为一个封装的api协议调用包，我们将它封装在和 `cmd` 同级别的包下。

``````shell
# /nano-claude-code
mkdir claude
# 首先，毫无疑问，我们需要有一个初始化功能
touch claude/init.go
``````

在文件 `claude/init.go` 中我们首先要定义好我们的模型，将我们刚才优先选中的三个参数添加进去。我们这里直接命名为 `ClaudeClient` 。

``````go
package Claude

import (
	"strings"

	"github.com/TIC-DLUT/nano-claude-code/errors"
)

type ClaudeClient struct {
	// 这里我们选择小写来将参数私有化，来避免apiKey的泄露。
	baseUrl string
	apiKey sting
}

func NewClaudeClient(baseurl, apiKey string) (*ClaudeClient, error) {
    // ***
    if !strings.HasPrefix(baseurl, "http://") && !strings.HasPrefix(baseurl, "https://") {
		return nil, errors.CreateClaudeClientBaseUrlError // mark-1
	}
    // ***
    return &ClaudeCode{
        baseUrl: baseurl,
        apiKey: apiKey,
    }, nil
}
``````

这样一个简单的初始化函数就出现了。

但是这里需要注意一下代码中被 `***` 包含起来的三行代码。我们期望我们的 `baseUrl` 可以是由 `http` 或 `https` 开头的，所以我们通过这三行代码来验证传入的 `baseurl` 是否是合法的一个 `baseurl` ，如果没有我们就直接在 `ClaudeClient` 的创建阶段直接报错。

其次，还有被我注释为 `mark-1` 的这行代码中的 `errors.CreateClaudeClientBaseUrlError` 尽管在 `golang` 的标准库中有名为 `errors` 的标准库，但是我们这里的目的是要统一管理报错信息，所以我们将报错信息单独成包。

### errors

在根目录下创建包 `errors` ，在 `errors` 下创建文件 `claude.go` ，我们所有 `claude` 包下的错误都在这个文件中创建并管理。示例如下：

````go
// errors/claude.go
package errors

import "errors"

var (
	CreateClaudeClientBaseUrlError = errors.New("unsupport baseurl")
)
````

后续所有的报错将不再展示详细内容（其实根据变量名称也能简单判断报错内容），如有需要可直接查看该项目的源代码，或笔者在项目完结后的整理。

### 初始化httpClient

仔细观察目前的 `ClaudeClient` 模型，想一想，有没有我们遗漏的东西呢？我们基本的需要添加的调用参数已经添加好了，那么我们怎么发起请求调用呢？是的，我们还需要一个 `httpClient` 来发起请求。这里我们就不直接重写一个 `httpClient` 包了，毕竟这不是本教程的重点，我们直接使用 `golang` 的 `resty` 包来实现该功能。具体使用方法可以在文章开头给出的链接中查看。

在添加 `httpClient` 之后现在的 `claude/init.go` 是这样子的：

````go
// claude/init.go
package claude

import (
	"strings"

	"github.com/TIC-DLUT/nano-claude-code/errors"
	"resty.dev/v3"
)

type ClaudeClient struct {
	baseUrl    string
	apiKey     string
	httpClient *resty.Client
}

func NewClient(baseurl string, apikey string) (*ClaudeClient, error) {
	if !strings.HasPrefix(baseurl, "http://") && !strings.HasPrefix(baseurl, "https://") {
		return nil, errors.CreateClaudeClientBaseUrlError
	}
	return &ClaudeClient{
		baseUrl:    baseurl,
		apiKey:     apikey,
		httpClient: resty.New(),
	}, nil
}
````

这里使用的 `resty` 包是用的 `v3` 版本，但截止到笔者撰写时间（2026-3-29） `v3` 版本仍在测试版本，读者使用 `v2` 版本也是可以的，其中的具体函数暂无明显变化。（具体用法参见给出的文档，本文章不详细解释）

### 简单的 `Call` 函数实现

接下来我们就要实现一个对话功能了。创建文件 `Claude/call.go` ，我们将在这个文件中实现实现基本的非流式和流式 chat 功能。

在这里，我们要想一下我们 call 一次大模型，需要传给大模型哪些参数。那么，我们就需要再次查看官方文档。

首先，我们要传递 `model` 来明确我们要调用的大模型是哪一个，然后我们要传一个 `messages` 来告诉大模型要做什么。

- `model` 就是模型名称，这里不多赘述
- `messages` 中的内容是 `MessageParam` 的数组，这里我们详细看一下 `MessageParam` 是什么样子的。

#### MessageParam

根据官方文档所示， `MessageParam` 的字段如下图所示：

````json
{
	"role": "",
	"content": {},
}
````

**role** ：官方文档中这个字段只有两个值：“user”（即用户），“assistant”（即大模型），那我们就可以将这两个字符串定义为常量，调用时直接使用这个常量，减少错误出现。

**content** ：官方文档中这个字段可以是17种类型之一，那么我们直接使用 `any` 来传递任意值。然后先封装最简单的单条语句。

> openai协议的内容与claude协议在message上有很多区别，比较重要的是"system"是写在和claude协议"messages"同级的"input"中，而不是claude定义的单独参数。

鉴于 `Message` 中的数据类型很多，此外还需要定义返回值中的消息内容，所以我们创建文件 `claude/message.go` 来方便统一管理。基于现有的需求，我们填充文件中的内容。

````go
package claude

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // 兼容所有 Content 类型
}

// role 类型枚举
const (
	ClaudeMessageRoleUser      = "user"
	ClaudeMessageRoleAssistant = "assistant"
)

type SingleStringMessage string // 单句
````

------

这样我们就解决了当前的 `message` 的问题。现在我们就可以实现我们的 `Call` 函数了。

#### `func Call` 

接下来我们首先要发起请求，这就需要用到我们之前封装好的 `Claude.httpClient` 。

我们可以直接使用语句 `c.httpClient.R().Post(c.baseUrl + "v1/message")` 来发起一次请求，但是有一个需要注意的点，在 `baseUrl` 中，有的人会在结尾加一个 `/` 而有的人不会，我们就需要统一一下写法，统一为后面加上这个 `/` 。如下：

````go
// 防止 https://example/
baseurl := c.baseUrl
if c.baseUrl[len(c.baseUrl)-1] != '/' {
    baseurl += "/"
}
````

接下来我们完善我们的请求体部分。由于任何的请求其实都是一样的，本质只是内部的部分内容有所变化，且这些变化可以通过 `any` 来消除，所以我们可以定义一个统一的结构体来简化类型。将我们上面提到的两个字段写入。

````go
type CallRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}
````

接下来我们就可以填充我们的 Body 部分了。

````go
c.httpClient.R().
	SetHeader("x-api-key", apikey).
	SetBody(CallRequest{
        Stream: false,
        Model: model,
        Message: message,
    }).
	Post(baseurl + "v1/message")
````

> ❗注意此处的 **url** ，openai协议下baseurl的**后缀**与 claude 协议是不同的，读者务必注意。

为了接收数据，我们根据官方文档的定义创建一个用于接受数据的结构体 `CallResponse` ，然后 `SetResult(&CallResponse{})` 即可自动解析结果。

````go
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
````

其中我们最重要的部分就是 `Content` 中的东西。那就直接取出来 `res.Content` 。这里这个内容是用 `[]map[string]interface{}` 来表示的。并且我们可以明确的是每个 `map[string]interface{}` 中就有一个叫 `"type"` 的字段来明确其他字段的类型。那么我们就可以取出这个字段，针对每个 `Type` 来专门实现对应的处理逻辑。

同理我们先实现最简单的 `TextBlock` :

````go
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
````

依次取出 `Content` 的每一个然后再取出 `type` ，针对不同的内容处理，即为下面的代码逻辑：

````go
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
````

然后我们加上数据接收处理和错误处理，就得到了这样的代码：

````go
func (c *ClaudeClient) Call(model string, messages []Message) ([]Message, error) {
    baseurl := c.baseUrl
    if c.baseUrl[len(c.baseUrl)-1] != '/' {
        baseurl += "/"
    }
    
    res, err := c.httpClient.R().
        SetHeader("x-api-key", c.apiKey). // 不要忘了加上鉴权密钥
        SetBody(CallRequest{
            Stream: false,
            Model: model,
            Message: message,
        }).
        Post(baseurl + "v1/message")
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
````

这样一个简单的 `Call` 函数就实现了。

但是考虑到我们还需要写一个 `CallStream` 的函数，那么我们可以将网络请求的内容再次封装一次为 `frontCall` 。

````go
// claude/call.go

// 返回三个参数分别是直接响应，流式响应，错误
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

    // 如果是流式响应，就设置不处理返回，
    // 反之则自动解析
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
````

#### `func CallStream` 

`CallStream` 的内容我们可以通过我们再上一小节中提前封装的 `frontCall` 中的返回值 `*resty.Response` 获取，并新建一个 `reader` 来不断地读取返回的信息，当我们发现已经没有消息时，就可以停止读取操作了。也就是下面的代码

````go
// dealFun 调用时输入，通过这种方式将信息同步传到函数外使用。
func (c *ClaudeClient) CallStream(model string, messages []Message, dealFunc func(Message) bool) error {
	_, originHttpRes, err := frontCall(c.httpClient, c.baseUrl, c.apiKey, model, messages, true)
	if err != nil {
		return err
	}

    // 新建 reader 读取数据
	reader := bufio.NewReader(originHttpRes.Body)
	defer originHttpRes.Body.Close() // 切记关闭

	for {
		eventStr, err := reader.ReadString('\n')
		if stdError.Is(err, io.EOF) { // 数据为空，说明已经结束
			break
		}
		if err != nil {
			return err
		}

		if strings.Trim(eventStr, " ") == "" { // 如果数据无意义则跳过
			continue
		}
	}
	return nil
}
````

这样我们就实现了不断读取的操作。读者可以自行添加相应的测试文件以及测试语句将这些数据输出查看。

接下来我们进行数据处理环节。resty 的 SSE 每次返回会返回两个字符串字段，即 `event` 和 `data` ，观察输出，发现在 `data` 返回的内容是下面这样格式的 json 字符串，其中最外层的 `type` 字段与 `event` 中的内容是相同的。因此，我们可以直接将 `event` 字段的内容忽略。

````json
{
    "type": "",
    "index": 0,
    "Content": {
        "type": "",
        "text": "",
    },
    "delta": {
        "type": "",
        "text": "",
    }
}
````

那么我们可以通过解析 `data` 中的这段字符串，提取在如下结构体中来方便我们后续使用。

````go
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
````

现在，我们先关注 `CallStreamResponse` 中的 `Type` 字段，其中的值主要包括：

"message_start"，"message_delta"，"content_block_start"，"content_block_delta"，"content_block_end"，"message_end"这些，其中 message 中的内容我们并不关心，我们真正需要使用的内容实际上是 content 中的内容，所以，我们只针对"content_block_start"(新增 content)和"content_block_delta"(后续 content 拼接)。就有了下面的处理结构：

````go
switch dataDetail.Type { // dataDetail 即为上文提到的取出来的 data 结构体
    case "content_block_start":

    case "content_block_delta":
    
}
````

因为流式响应可能使用在中间过程，这时我们更希望能够直接获取一次性的信息来执行，所以我们额外定义一个返回值 `[]Message` 和变量 `resMessage` 将完整的信息返回。

**`content_block_start`** 

对于新增的开始信息，我们为 `resMessage` append 一个新的元素，并根据其中的 `Content.Type` 创建对应的 Content 类型。（目前只有 TextBlock 一种）。

````go
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
````

**`content_block_delta`** 

对于新的补充信息，我们将其拼接在上一个信息内。

````go
case "content_block_delta":
	// 默认调用未完成
    continueFlag := true
    switch resMessages[len(resMessages)-1].Content.(type) {
        case TextBlock:
        resMessages[len(resMessages)-1].Content = TextBlock{
            Type: "text",
            Text: resMessages[len(resMessages)-1].Content.(TextBlock).Text + dataDetail.Delta.Text,
        }
        // 可以根据实际返回内容通过 dealFunc 来控制是否需要来停止该次 CallStream 调用
        continueFlag = dealFunc(Message{
            Role: ClaudeMessageRoleAssistant,
            Content: TextBlock{
                Type: "text",
                Text: dataDetail.Delta.Text,
            },
        })
    }
	// 停止调用，直接返回
    if !continueFlag {
        return resMessages, nil
    }
````

最终 `CallStream` 如下：

````go
// claude/call.go
func (c *ClaudeClient) CallStream(model string, messages []Message, dealFunc func(Message) bool) ([]Message, error) {
	_, originHttpRes, err := frontCall(c.httpClient, c.baseUrl, c.apiKey, model, messages, true, []Tool{})
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
````



### 封装 `ToolCall` 

接下来，作为一个 claude code 的智能体，我们的封装绝对不能停止在简单的 Call 函数，我们还需要将赋予大模型直接影响世界的能力——tool call。

#### `Tool` 封装

在官方的 API 中，我们可以看到，关于 tools 的关键定义：

````go
// claude/call_tool.go
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema struct {
		Type       string                        `json:"type"`
		Properties map[string]ToolPropertyDetail `json:"properties"`
		Required   []string                      `json:"required"`
	} `json:"input_schema"`
    // 这里是工具的运行函数，不属于官方的 Tool 定义，但是我们需要其用来方便的运行 Tool 函数，不需要 json 化
	Func func(input map[string]any) string `json:"-"` 
}
// 输入参数的具体内容
type ToolPropertyDetail struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
````

创建新的文件，将上面的定义存入。我们就有了一个工具的定义。

现在我们提供一个接口创建新的工具。对于用户来说，一个新的 `Tool` 创建，需要有工具的名称，描述，需要传入的参数，以及运行的 Tool 函数。即如下结果。
````go
// claude/call_tool.go
func NewTool(name string, description string, properties map[string]ToolPropertyDetail, required []string, toolFunc func(input map[string]any) string) (Tool, error) {
    // name 和 description 是必须不能为空的。
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
````

到此为止，我们有了一个完整的 `Tool` 模型。

#### 补充底层封装

接下来，我们就需要将tool的定义添加到我们之前在写好的 Call 模型中。

````go
// claude/message.go
type ToolUseBlock struct {
	Type  string         `json:"type"`
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}
type ToolResultBlock struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	ToolUseID string `json:"tool_use_id"`
}

// claude/call.go
type CallRequest struct {
    // ...
	Tools    []Tool    `json:"tools"`
    // ...
}
type CallStreamResponse struct {
    // ...
	ContentBlock struct {
		ID   string `json:"id"`
		Name string `json:"name"`
        // ...
	} `json:"content_block"`
	Delta struct {
        // ...
		PartialJson string `json:"partial_json"`
	} `json:"delta"`
}
func frontCall(httpClient *resty.Client, inBaseUrl string, apiKey string, model string, messages []Message, stream bool, tools []Tool) (CallResponse, *resty.Response, error) {
    // ❗ 注意上面的参数传入了 tools
    // ...

    // ❗
    // 此处将所有的同 Role 的 Content 合并，同时避免出现在 Message 中出现单个 ContentBlock 的错误情况
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
        // 统一格式，将 SingleStringMessage 也转化为 TextBlock
		if reflect.TypeOf(messages[i].Content) == reflect.TypeOf(SingleStringMessage("")) {
			newContent = append(newContent, TextBlock{Type: "text", Text: string(messages[i].Content.(SingleStringMessage))})
		} else {
			newContent = append(newContent, messages[i].Content)
		}
		requestMessages[len(requestMessages)-1].Content = newContent
	}

	res := CallResponse{}

	requestBody := CallRequest{
		Stream:   stream,
		Model:    model,
		Messages: messages,
	}

	if len(tools) != 0 {
		requestBody.Tools = tools
	}
    
    // ...
}
func (c *ClaudeClient) Call(model string, messages []Message, tools []Tool) ([]Message, error) {
// ❗ 注意上面的参数传入了 tools，在 CallStream 中也要进行相应的修改，
// ...
    switch messageType {
    // ...
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
	}
    // ...
// ...
}
func (c *ClaudeClient) CallStream(model string, messages []Message, tools []Tool, dealFunc func(Message) bool) ([]Message, error) {
	// ...
    switch dataDetail.Type {
    // ...
    case "content_block_start":
        // ...
        switch dataDetail.ContentBlock.Type {
        // ...
        case "tool_use":
            // 新建一个 ToolUseBlock，并记录 ID 和 Name（仅在 start 时才出现）
            content = ToolUseBlock{
                Type: "tool_use",
                ID:   dataDetail.ContentBlock.ID,
                Name: dataDetail.ContentBlock.Name,
            }
        }

        resMessages[len(resMessages)-1].Content = content

    case "content_block_delta":
        continueFlag := true
        switch resMessages[len(resMessages)-1].Content.(type) {
            // ...
        case ToolUseBlock:
            // 拼接参数string
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
    // 由于我们只进行了输入参数的字符串拼接，现在 Tool 的输入仍是一段 string，无法在 Func 中使用
    // 所以我们在这里将所有的 PartialJson 转化为 map[string]any 的类型传入 Tool 函数
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

````

#### `func CallTools` 

函数 `CallTools` 我们的设想是向模型发送请求，然后模型会返回相应的Content Block，这些Block在处理过后 `CallTools` 会根据其中的内容调用相应的 Tool。发送请求我们可以使用刚刚补充封装的 `Call` 来实现。此外，注意到，我们收到的tooluse请求未必只有一次，他可能会依次调用多次或多个 tool，所以我们需要在一个 `for` 循环中不断地调用 `Call` ，当不需要使用 Tools 时，我们再停止循环。这样，一个基本的框架就有了。

````go
func (c *ClaudeClient) CallTools(model string, messages []Message, tools []Tool) ([]Message, error) {
	var err error = nil
    // 新的消息
    realMessage := []Message{}
    // 新的单次返回的新消息
	resMessages := []Message{}
	for {
		resMessages, err = c.Call(model, messages, tools)
		if err != nil {
			return []Message{}, err
		}

        // 将新消息拼接在原来的消息之后
		messages = append(messages, resMessages...)
		realresMessages = append(realresMessages, resMessages...)

        // 控制循环，仍在 tool call 过程则持续循环
		continueFlag := false

		for _, item := range resMessages {
			switch item.Content.(type) {
				// 具体 tool use 内容
                messages = append(messages, toolResultMessage)
                realresMessages = append(realresMessages, toolResultMessage)
			}
		}

		if !continueFlag {
			break
		}
	}
	return realMessage, err
}
````

而 tool use 的过程实际也很简单，我们提取出其中的 `ToolUseBlock` 内容，然后在所有的 tools 中找到对应的 tool 并运行即可。

````go
switch item.Content.(type) {
    case ToolUseBlock:
    continueFlag = true
    toolUserItem := item.Content.(ToolUseBlock)
    content := ""
    // 找到对应函数并运行获取结果
    for _, tool := range tools {
        if tool.Name == toolUserItem.Name {
            content = tool.Func(toolUserItem.Input)
        }
    }
    // 返回实际运行结果
    toolResultMessage := Message{
        Role: ClaudeMessageRoleUser,
        Content: ToolResultBlock{
            Type:      "tool_result",
            ToolUseID: toolUserItem.ID,
            Content:   content,
        },
    }
}
````

非常简单，我们的 `CallTools` 函数就封装好了。完整代码如下：

````go
// claude/call_tool.go
func (c *ClaudeClient) CallTools(model string, messages []Message, tools []Tool) ([]Message, error) {
	var err error = nil
    realMessage := []Message{}
	resMessages := []Message{}
	for {
		resMessages, err = c.Call(model, messages, tools)
		if err != nil {
			return []Message{}, err
		}

		messages = append(messages, resMessages...)
		realresMessages = append(realresMessages, resMessages...)

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
                toolResultMessage := Message{
                    Role: ClaudeMessageRoleUser,
                    Content: ToolResultBlock{
                        Type:      "tool_result",
                        ToolUseID: toolUserItem.ID,
                        Content:   content,
                    },
                }

                messages = append(messages, toolResultMessage)
                realresMessages = append(realresMessages, toolResultMessage)
			}
		}

		if !continueFlag {
			break
		}
	}
	return resMessages, err
}
````

#### `func CallStreamTools` 

流式的CallTool思路与非流式相同，甚至没有任何区别。并且在运行Tool时的逻辑完全一样，那么我们可以先将这一部分再次封装，用来专门运行工具函数并返回循环是否继续，下次请求的 message，新增的 toolResultMessage。

````go
// claude/call_tool.go
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
````

这样封装之后，原先的 `CallTools` `CallStreamTools` 的逻辑就只剩下循环 `Call` / `CallStream` ，然后拼接所有新增信息，并控制循环进行。

````go
// claude/call_tool.go
// 原先的 CallTools 与该函数基本相同，只需要将 CallStream 换为 Call，调整相关参数即可。
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
````



### 解析http错误

到这里为止，我们已经将我们需要用来做 agent 的所有协议封装好了。但是，在 `frontCall` 中，我们仅仅提取了返回的响应体，如果我们的请求出错了，我们无法获取到任何信息。因此，我们需要在请求出错时解析返回的错误。当响应码不为200时，就会解析错误信息并报错。

````go
// claude/call.go
type CallError struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	Type string `json:"type"`
}
func frontCall(httpClient *resty.Client, inBaseUrl string, apiKey string, model string, messages []Message, stream bool, tools []Tool) (CallResponse, *resty.Response, error) {
	// ...
	if httpRes.StatusCode() != 200 {
		httpBody, _ := io.ReadAll(httpRes.Body)
		defer httpRes.Body.Close()
		errMessage := CallError{}
		json.Unmarshal(httpBody, &errMessage)
		return res, httpRes, stdError.New(errMessage.Error.Message)
	}
    return res, httpRes, err
}
````

## 总结

至此，我们第一天的内容就结束了。（可喜可贺）

这一天的内容相当充实，从网络请求封装，到一次 `Message` 调用封装，再到工具调用封装，我们封装了一个真正属于我们自己的 claude 包，也迈出了 claude code 的第一步。

## 最终文件结构

````
nano-claude-code
│  go.mod
│  go.sum
├─claude
│      call.go
│      call_test.go
│      call_tool.go
│      call_tool_test.go
│      init.go
│      message.go
├─cmd
│      main.go
└─errors
        claude.go
````

