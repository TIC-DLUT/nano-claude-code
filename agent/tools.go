package agent

import "github.com/TIC-DLUT/nano-claude-code/agent/tools"

func (a *Agent) LoadTools() error {
	// filesystem
	filesystem_readfile_tool, err := tools.NewReadFileTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesystem_readfile_tool)
	filesysytem_writefile_tool, err := tools.NewWriteFileTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesysytem_writefile_tool)

	return nil
}
