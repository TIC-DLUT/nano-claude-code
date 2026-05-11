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
  
	filesysytem_writefile_tool, err := tools.NewWriteFileTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesysytem_writefile_tool)

	filesystem_bash_tool, err := tools.NewBashTool()
	if err != nil {
		return err
	}
	a.tools = append(a.tools, filesystem_bash_tool)

	return nil
}
