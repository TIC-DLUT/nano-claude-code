package tools

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/TIC-DLUT/nano-claude-code/claude"
)

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

func NewBashTool() (claude.Tool, error) {
	return claude.NewTool("bash",
		"运行命令行工具，并返回运行结果。使用 golang 的 exec 包实现，不需要对平台做特殊处理",
		map[string]claude.ToolPropertyDetail{
			"command": {
				Type:        "string",
				Description: "需要运行的命令",
			},
			"args": {
				Type:        "array",
				Description: "运行命令的参数string列表",
			},
		},
		[]string{"command"},
		func(input map[string]any) string {
			// 提取 command 字段
			fmt.Println(input["command"].(string))
			command, ok := input["command"].(string)
			if !ok {
				return "command 不能为空"
			}

			// 若存在 args 字段，则提取
			var args = []string{}
			argsFlag, ok := input["args"]
			if ok {
				argsInAny := argsFlag.([]any)
				for _, argInAny := range argsInAny {
					arg, ok := argInAny.(string)
					if !ok {
						return "如果需要使用参数，请使用 []string 类型"
					}
					args = append(args, arg)
				}
			}

			// 运行完整指令，并获取结果
			cmd := exec.Command(command, args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return "error: " + err.Error()
			}

			return string(output)
		})
}
