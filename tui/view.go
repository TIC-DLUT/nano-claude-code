package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func (m *model) refreshViewport() {
	m.viewport.SetContent(
		m.applyHighlight(lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderMessages())),
	)
	if m.viewport.AtBottom() {
		m.viewport.GotoBottom()
	}
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

// 对选区范围内的字符加reverse反白；用ansi.Cut处理ANSI/CJK/emoji
func (m *model) applyHighlight(content string) string {
	if !m.hasSelection() {
		return content
	}
	sLine, sCol, eLine, eCol := m.normalizedSelection()
	lines := strings.Split(content, "\n")
	if sLine >= len(lines) {
		return content
	}
	eLine = min(eLine, len(lines)-1)
	reverse := lipgloss.NewStyle().Reverse(true)
	for i := sLine; i <= eLine; i++ {
		line := lines[i]
		w := ansi.StringWidth(line)
		s, e := selectionCols(i, sLine, eLine, sCol, eCol, w)
		if s < e {
			lines[i] = ansi.Cut(line, 0, s) + reverse.Render(ansi.Cut(line, s, e)) + ansi.Cut(line, e, w)
		}
	}
	return strings.Join(lines, "\n")
}

// 提取选择的纯文本字符串
func (m *model) selectedText() string {
	if !m.hasSelection() {
		return ""
	}
	sLine, sCol, eLine, eCol := m.normalizedSelection()
	lines := strings.Split(lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderMessages()), "\n")
	if sLine >= len(lines) {
		return ""
	}
	eLine = min(eLine, len(lines)-1)
	var b strings.Builder
	for i := sLine; i <= eLine; i++ {
		s, e := selectionCols(i, sLine, eLine, sCol, eCol, ansi.StringWidth(lines[i]))
		if i > sLine {
			b.WriteByte('\n')
		}
		if s < e {
			b.WriteString(ansi.Strip(ansi.Cut(lines[i], s, e)))
		}
	}
	return b.String()
}

// 计算第i行的选中的列范围 [s, e)
func selectionCols(i, sLine, eLine, sCol, eCol, w int) (s, e int) {
	s, e = 0, w
	if i == sLine {
		s = min(sCol, w)
	}
	if i == eLine {
		e = min(eCol, w)
	}
	return min(s, e), e
}

func (m *model) View() tea.View {
	viewportView := m.viewport.View()
	separator := strings.Repeat("─", max(m.viewport.Width(), 1))

	footer := ""
	if m.quitConfirm {
		footer = "Press Ctrl-C again to exit"
	}

	v := tea.NewView(viewportView + "\n" + separator + "\n" + m.textarea.View() + "\n" + separator + "\n" + footer)
	v.MouseMode = tea.MouseModeCellMotion
	if c := m.textarea.Cursor(); c != nil {
		c.Y += lipgloss.Height(viewportView) + 1
		v.Cursor = c
	}
	v.AltScreen = true
	return v
}
