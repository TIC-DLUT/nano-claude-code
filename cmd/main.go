package main

import (
	"flag"
	"fmt"

	"github.com/TIC-DLUT/nano-claude-code/agent"
	"github.com/TIC-DLUT/nano-claude-code/config"
	"github.com/TIC-DLUT/nano-claude-code/tui"
)

func init() {
	flag.BoolVar(&TUI_Mode, "tui", false, "是否开启tui模式")
	flag.StringVar(&Message, "message", "", "非tui模式，执行的内容")
	flag.StringVar(&SessionID, "session", "", "选择从那个会话开始")

	flag.Parse()
}

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	MainAgent, err = agent.NewAgent(&SessionID)
	if err != nil {
		panic(err)
	}

	if TUI_Mode {
		// 启动tui
		if err := tui.Run(MainAgent); err != nil {
			panic(err)
		}
	} else {
		// 直接调用
		DirectRun()
	}

	fmt.Println("\nsession id: ", SessionID)
}
