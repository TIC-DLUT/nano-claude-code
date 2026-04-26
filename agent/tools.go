package agent

import "github.com/TIC-DLUT/nano-claude-code/agent/tools"

func (a *Agent) LoadTools() error {
	// filesystem
	filesystem_readfile_tool, err := tools.NewReadFileTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesystem_readfile_tool)
	filesystem_editfile_tool, err := tools.NewEditFileTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesystem_editfile_tool)

	return nil
}
