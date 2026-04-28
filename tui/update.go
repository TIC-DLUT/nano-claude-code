package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
)

func (m *model) Init() tea.Cmd {
	return m.textarea.Focus()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(max(msg.Height-m.textarea.Height()-3, 1))
		if len(m.messages) > 0 || m.assistantDraft != "" {
			m.refreshViewport()
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.quitConfirm {
				return m, tea.Quit
			}
			m.quitConfirm = true
			return m, quitConfirmTimeout()
		case "enter":
			if m.busy {
				return m, nil
			}
			text := strings.TrimSpace(m.textarea.Value())
			if text == "" {
				return m, nil
			}
			m.messages = append(m.messages, m.sender.Render("You: ")+text)

			m.textarea.Reset()
			m.assistantDraft = ""
			m.busy = true

			m.refreshViewport()
			return m, m.startStream(text)
		default:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg, tea.PasteMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd

	case tea.MouseWheelMsg:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case streamChunkMsg:
		m.assistantDraft += msg.chunk
		m.refreshViewport()
		return m, nil

	case streamDoneMsg:
		if m.assistantDraft != "" {
			m.messages = append(m.messages, m.assistant.Render("Assistant: ")+m.assistantDraft)
		}
		m.assistantDraft = ""
		m.busy = false
		m.refreshViewport()
		return m, nil

	case quitConfirmTimeoutMsg:
		m.quitConfirm = false
		return m, nil
	}

	return m, nil
}

func quitConfirmTimeout() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return quitConfirmTimeoutMsg{}
	})
}

func (m *model) startStream(userText string) tea.Cmd {
	return func() tea.Msg {
		go func() {
			m.agent.ChatStream(userText, func(s string) {
				m.program.Send(streamChunkMsg{chunk: s})
			})
			m.program.Send(streamDoneMsg{})
		}()
		return nil
	}
}
