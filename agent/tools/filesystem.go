package tools

import (
	"os"
	"strings"

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
