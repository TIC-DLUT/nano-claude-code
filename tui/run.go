package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/TIC-DLUT/nano-claude-code/agent"
)

func Run(agent *agent.Agent) error {
	m := newModel(agent)
	program := tea.NewProgram(m)
	m.program = program

	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}
