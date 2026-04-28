package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *model) refreshViewport() {
	m.viewport.SetContent(
		lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderMessages()),
	)
	m.viewport.GotoBottom()
	//TODO: 目前的话，在llm流式输出的时候会一直滚动到最底部，需要优化为在输出未完成时也可以用户进行滚动
}

func (m *model) renderMessages() string {
	var b strings.Builder

	if len(m.messages) > 0 {
		b.WriteString(strings.Join(m.messages, "\n"))
	}

	if m.busy {
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(m.assistant.Render("Assistant: "))
		/*
			TODO:
			目前：You: 你好
				 Assistant: 你好
				 Assistant: 新回复内容
			期望： 你好(用户输出整行高光对齐claude-code样式)
				  你好(llm输出不做处理)
			进阶： 使其完成markdown格式渲染
		*/
		b.WriteString(m.assistantDraft)
	}
	return b.String()
}

func (m *model) View() tea.View {
	viewportView := m.viewport.View()
	separator := strings.Repeat("─", max(m.viewport.Width(), 1))

	footer := ""
	if m.quitConfirm {
		footer = "Press Ctrl-C again to exit"
	}

	v := tea.NewView(viewportView + "\n" + separator + "\n" + m.textarea.View() + "\n" + separator + "\n" + footer)
	v.MouseMode = tea.MouseModeCellMotion //tui接管鼠标控制！会导致无法进行鼠标选择复制，之后需处理
	if c := m.textarea.Cursor(); c != nil {
		c.Y += lipgloss.Height(viewportView) + 1
		v.Cursor = c
	}
	v.AltScreen = true
	return v
}
