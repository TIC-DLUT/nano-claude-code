# day2：程序入口与简单agent

## 该day对应代码提交到commit: e7a5768c9129773c660e90b813dddad1c9d4278f 为止

经过第一天的努力，我们成功封装了一个足以满足我们后续使用的 claude 协议包。当然我们也不能忘记我们最初的目的——写一个 nano claude code。因此，接下来，我们就要正式开始从最初写下的主函数出发。

## 实现主程序入口

打开现在的主函数，我们看着空空的函数内部，一时间竟有点无从下手。不过没关系，一个程序的运行一定需要相应的环境配置，所以我们从获取环境变量开始。

### 配置

为保持项目整洁，我们将配置设置单独在根目录下提供一个 `config` 包来实现配置设置。在其内部，我们使用 `viper` 包来实现我们对配置的操作。

````go
// config/init.go
func LoadConfig() error {
	homePath, _ := os.UserHomeDir()

	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(homePath, ".nano-claude-code"))

    // 对于环境变量的绑定我们单独设置一个函数处理
    // 读者亦可直接在此处完成，具体内容见下
	BindEnv()

	err := viper.ReadInConfig()

	// 如果文件不存在，启动创建的引导
    // 该部分不是本项目重点，故不展示。默认该文件存在

	return err
}

// config/env.go
func BindEnv() {
	viper.SetEnvPrefix("ncc")
	viper.BindEnv("llm.apikey")
	viper.BindEnv("llm.baseurl")
	viper.BindEnv("llm.model")
}
````

````json
// config/config.example.json
{
    "llm": {
        "baseurl": "",
        "apikey": "",
        "model": ""
    }
}
````

如此，我们的主函数就迎来了第一个任务：通过 `LoadConfig` 读取到我们的配置文件。

### 程序入口

接下来我们的主程序还有知道使用者使用它希望它能做什么。所以，我们就需要设置一些参数来解析所需要读取的内容。

首先我们支持用户一次性调用和使用 tui 模式，所以我们定义 `TUI_Mode` 来控制程序运行模式。

在 tui 模式下，所有的信息会在进入 tui 模式后获取，但是对应单次调用的非 tui 模式，我们还需要一个 `Message` 参数来明确用户的消息。并在主函数运行前就初始化它们的值。

````go
// cmd/flags.go
var (
	TUI_Mode bool
	Message  string
)

// cmd/main.go
// 在 golang 中名为 init() 的函数会默认在 main() 之前执行
func init() {
	flag.BoolVar(&TUI_Mode, "tui", false, "是否开启tui模式")
	flag.StringVar(&Message, "message", "", "非tui模式，执行的内容")

	flag.Parse()
}

func main() {
    // 上文提到的配置文件加载
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if TUI_Mode {
		// 启动tui
	} else {
		// 直接调用
	}
}
````

今天，我们一切从简，完成非 tui 模式下的程序运行。

## 单次调用的 agent

在非 tui 模式下，我们直接将程序运行权限已将给下一级函数 `DirectRun()` 。在这个函数中我们首先要确保我们收到了 `Message` 来执行 agent。

````go
// cmd/direct.go
func DirectRun() {
	if Message == "" {
		panic("message不能为空")
	}
    
    // run agent
}
````

到此为止，我们就需要专注在创建一个基本的智能体上。

这个智能体我们需要有什么呢？

- 初始化 claudeClient 来实现请求调用。
- 在 loadTools 中将所有工具添加在 agent 中。
- 获取 systemPrompt 来明确 agent 的身份。
- 一次调用的函数。

得益于我们之前的封装工作，这一部分的内容将变得非常简单。

**agent初始化** 

agent中我们需要包含所有的工具和 client，如此，便是这样子的结构。

````go
// agent/init.go
type Agent struct {
    apiClient *claude.ClaudeClient
    tools 	  []claude.Tool
}
````

作为一个智能体，我们的工具可以想象的很多，所以我们将所有工具在 agent 包下在单独封装在一个包内来统一管理。

**tools 创建** 

我们优先实现有关文件系统和命令行的工具创建。也就是这四种工具：

- ReadFile
- WriteFile
- EditFile
- Bash

这些工具本质大同小异，只需要设计好 tool 的 name，description，参数，并提供能实现该功能的函数即可，作为教程我们将以 ReadFile 为例：

````go
// agent/tools/filesystem.go
func NewReadFileTool() (claude.Tool, error) {
	return claude.NewTool("read_file", "读一个文件，返回该文件的全部内容", map[string]claude.ToolPropertyDetail{
		"path": {
			Type:        "string",
			Description: "文件目录",
		},
	}, []string{"path"}, func(input map[string]any) string {
		path, ok := input["path"].(string)
		if !ok {
			return "path不能为空"
		}
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return "error: " + err.Error()
		}
		return string(fileContent)
	})
}
````

**system prompt** 

接下来我们书写 system prompt 。在这之中我们要确保其时效性和工作边界。

````go
// agent/prompt.go
const (
	systemPrompt = `你是claude code，你需要调用工具帮助人们完成工作

当前时间是：{system_time}
当前工作地址是：{work_path}`
)

func GetNowSystemPrompt() string {
	systemTime := time.Now().Format("2006-01-02 15:04:05")
	workPath, err := os.Getwd()
	if err != nil {
		workPath = "unknown"
	}

	nowSystemPrompt := systemPrompt
	nowSystemPrompt = strings.ReplaceAll(nowSystemPrompt, "{system_time}", systemTime)
	nowSystemPrompt = strings.ReplaceAll(nowSystemPrompt, "{work_path}", workPath)

	return nowSystemPrompt
}
````

这里我们简单的使用 `strings` 中的函数实现模板，读者可以自行实现更为严谨的模板 prompt。

**chat** 

幸运的，我们第一天的封装非常有效，现在我们一次 chat 只需要使用 `apiClient` 中的函数即可。

````go
// agent/chat.go
func (a *Agent) ChatStream(message string, callback func(string)) {
	lastToolCallID := ""

	a.apiClient.CallStreamTools(viper.GetString("llm.model"), GetNowSystemPrompt(), []claude.Message{
		{
			Role:    claude.ClaudeMessageRoleUser,
			Content: claude.SingleStringMessage(message),
		},
	}, a.tools, func(m claude.Message) bool {
		switch m.Content.(type) {
		case claude.TextBlock:
			callback(m.Content.(claude.TextBlock).Text)
		case claude.ToolUseBlock:
			tooluse := m.Content.(claude.ToolUseBlock)
			if tooluse.ID != lastToolCallID {
				lastToolCallID = tooluse.ID
				callback("\n[tool_use] " + tooluse.Name + "\n")
			}
		}
		return true
	})
}

````

到此我们的一个简单 agent 所需要的内容就实现了。接下来我们只需要将这些函数在他们需要的地方使用即可。

````go
// 完善 agent 包

// agent/init.go
func NewAgent() (*Agent, error) {
	apiClient, err := claude.NewClient(viper.GetString("llm.baseurl"), viper.GetString("llm.apikey"))
	if err != nil {
		return nil, err
	}

	newAgent := &Agent{
		apiClient: apiClient,
	}

	newAgent.LoadTools()

	return newAgent, nil
}

// agent/tools.go
func (a *Agent) LoadTools() error {
	// filesystem
	filesystem_readfile_tool, err := tools.NewReadFileTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesystem_readfile_tool)

	return nil
}
````

在程序入口处，我们无论是否是 tui 模式，我们都需要的是同一个 agent，所以我们可以创建一个全局的 MainAgent 供两个模式一同使用。

````go
// cmd/main.go
func main() {
    // ...
	MainAgent, err = agent.NewAgent()
	if err != nil {
		panic(err)
	}
    // ...
}

// cmd/direct.go
func DirectRun() {
	if Message == "" {
		panic("message不能为空")
	}
	MainAgent.ChatStream(Message, func(s string) {
		fmt.Print(s)
	})
}
````

## 总结

至此，我们第二天的内容就结束了。（可喜可贺）

这一天的内容相当简单，这主要得益于我们第一天的封装足够完善，所以，我们这一天在使用 agent 相关的操作的时候，我们的逻辑十分的简单，只需要正确使用我们之前封装的内容就可以了。

由此可见，一个好的封装可以为后续节省很多工作。

## 最终文件结构

````
nano-claude-code
│ ...
├─agent
│  │  chat.go
│  │  init.go
│  │  prompt.go
│  │  tools.go
│  └─tools
│          filesystem.go
├─claude
├─cmd
│      direct.go
│      flags.go
│      main.go
│      var.go
├─config
│      config.example.json
│      env.go
│      init.go
└─errors
````



## 参考文档

- Learn Claude Code 教程：[Learn Claude Code](https://learn.shareai.run/en/s02/)
- viper 官方文档：[Viper](https://github.com/spf13/viper)

## 关联good first issues列表

- [添加工具 WriteFile](https://github.com/TIC-DLUT/nano-claude-code/issues/8)
- [添加工具 EditFile](https://github.com/TIC-DLUT/nano-claude-code/issues/9)
- [添加工具 Bash](https://github.com/TIC-DLUT/nano-claude-code/issues/10)
- [添加 tui 模式](https://github.com/TIC-DLUT/nano-claude-code/issues/6)
