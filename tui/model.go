package tui

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/TIC-DLUT/nano-claude-code/agent"
)

type (
	streamChunkMsg        struct{ chunk string }
	streamDoneMsg         struct{}
	quitConfirmTimeoutMsg struct{}
)

type model struct {
	agent     *agent.Agent
	program   *tea.Program
	viewport  viewport.Model
	textarea  textarea.Model
	messages  []string
	sender    lipgloss.Style //用户发送的消息的前缀样式
	assistant lipgloss.Style

	assistantDraft string //流式返回的文本
	busy           bool   //agent是否正在处理请求
	quitConfirm    bool   //二次退出确认

	// 鼠标选区坐标
	mouseDown    bool
	selStartLine int
	selStartCol  int
	selEndLine   int
	selEndCol    int
}

func newModel(agent *agent.Agent) *model {
	ta := textarea.New()
	ta.SetVirtualCursor(true)
	ta.Prompt = "❯ "
	ta.CharLimit = 4096
	ta.SetWidth(30)
	ta.SetHeight(1)

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to nano-claude-code`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)
	vp.MouseWheelDelta = 1 //鼠标滚轮灵敏度

	return &model{
		agent:        agent,
		textarea:     ta,
		viewport:     vp,
		messages:     []string{},
		sender:       lipgloss.NewStyle(),
		assistant:    lipgloss.NewStyle(),
		selStartLine: -1, //-1 表示无选区
		selEndLine:   -1,
	}
}
