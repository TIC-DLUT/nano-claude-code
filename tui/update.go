package tui

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
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
		m.clearSelection()
		if len(m.messages) > 0 || m.assistantDraft != "" {
			m.refreshViewport()
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			// 有选区时复制选中文本
			if m.hasSelection() {
				text := m.selectedText()
				m.clearSelection()
				m.refreshViewport()
				return m, func() tea.Msg { _ = clipboard.WriteAll(text); return nil } //通过clipboard包将文本复制到剪贴板
			}

			// 退出确认
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
			m.clearSelection()
			m.refreshViewport()
			m.viewport.GotoBottom()
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

	case tea.MouseClickMsg:
		if msg.Button == tea.MouseLeft {
			m.mouseDown = true
			m.selStartLine, m.selStartCol = m.viewport.YOffset()+msg.Y, max(msg.X, 0)
			m.selEndLine, m.selEndCol = m.selStartLine, m.selStartCol
			m.refreshViewport()
		}

	case tea.MouseMotionMsg:
		if m.mouseDown && msg.Button == tea.MouseLeft {
			m.selEndLine, m.selEndCol = m.viewport.YOffset()+msg.Y, max(msg.X, 0)
			m.refreshViewport()
		}

	case tea.MouseReleaseMsg:
		m.mouseDown = false

	case streamChunkMsg:
		m.clearSelection()
		m.assistantDraft += msg.chunk
		m.refreshViewport()
		return m, nil

	case streamDoneMsg:
		m.clearSelection()
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

func (m *model) hasSelection() bool {
	return m.selStartLine >= 0 && m.selEndLine >= 0 &&
		(m.selStartLine != m.selEndLine || m.selStartCol != m.selEndCol)
}

// 清除选区(其实更像重置或初始化 awa...)
func (m *model) clearSelection() {
	m.mouseDown = false
	m.selStartLine, m.selEndLine = -1, -1
	m.selStartCol, m.selEndCol = 0, 0
}

// 规整选区，保证选择起始点 <= 选择结束点
func (m *model) normalizedSelection() (sLine, sCol, eLine, eCol int) {
	sLine, sCol, eLine, eCol = m.selStartLine, m.selStartCol, m.selEndLine, m.selEndCol
	if sLine > eLine || (sLine == eLine && sCol > eCol) {
		sLine, eLine, sCol, eCol = eLine, sLine, eCol, sCol
	}
	return
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
