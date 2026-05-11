package tools

import (
	"os"
	"strings"
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
				Description: "需要运行的命令，不要包含参数内容，仅需要运行的命令的首个字段，其余字段放置在 args 中",
			},
			"args": {
				Type:        "array",
				Description: "运行命令的参数string列表",
			},
		},
		[]string{"command"},
		func(input map[string]any) string {
			// 提取 command 字段
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

func NewEditFileTool() (claude.Tool, error) {
	return claude.NewTool("edit_file",
		"编辑一个文件，返回编辑结果",
		map[string]claude.ToolPropertyDetail{
			"path": {
				Type:        "string",
				Description: "文件目录",
			},
			"old_string": {
				Type:        "string",
				Description: "需要被替换的字符串",
			},
			"new_string": {
				Type:        "string",
				Description: "替换后的字符串",
			},
		},
		[]string{"path", "old_string", "new_string"},
		func(input map[string]any) string {
			// 提取输入
			path, ok := input["path"].(string)
			if !ok {
				return "path不能为空"
			}
			oldString, ok := input["old_string"].(string)
			if !ok {
				return "old_string不能为空"
			}
			newString, ok := input["new_string"].(string)
			if !ok {
				return "new_string不能为空"
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return "error: " + err.Error()
			}
			count := strings.Count(string(content), oldString)
			if count == 0 {
				return "文件中没有找到需要替换的字符串"
			}
			if count > 1 {
				return "文件中找到多个需要替换的字符串，提供更唯一的匹配"
			}
			newContent := strings.Replace(string(content), oldString, newString, 1)
			err = os.WriteFile(path, []byte(newContent), 0777)
			if err != nil {
				return "error: " + err.Error()
			}
			return "编辑成功"
		})
}

func NewWriteFileTool() (claude.Tool, error) {
	return claude.NewTool(
		"write_file",
		"写一个文件，输入文件路径和内容，返回写入结果",
		map[string]claude.ToolPropertyDetail{
			"path": {
				Type:        "string",
				Description: "文件目录",
			},
			"content": {
				Type:        "string",
				Description: "要写入文件的内容",
			},
		},
		[]string{"path", "content"},
		func(input map[string]any) string {
      // 提取输入
			path, ok := input["path"].(string)
			if !ok {
				return "path不能为空"
			}
			content, ok := input["content"].(string)
			if !ok {
				return "content不能为空"
			}

			err := os.WriteFile(path, []byte(content), 0777)
			if err != nil {
				return "error: " + err.Error()
			}
			return "文件写入成功"
		})
}
