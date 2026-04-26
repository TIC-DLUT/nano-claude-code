package tools

import (
	"os"

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
