package agent

import (
	"github.com/TIC-DLUT/nano-claude-code/claude"
	"github.com/TIC-DLUT/nano-claude-code/session"
	"github.com/spf13/viper"
)

type Agent struct {
	apiClient      *claude.ClaudeClient
	tools          []claude.Tool
	sessionManager *session.SessionManager
}

func NewAgent(sessionID *string) (*Agent, error) {
	apiClient, err := claude.NewClient(viper.GetString("llm.baseurl"), viper.GetString("llm.apikey"))
	if err != nil {
		return nil, err
	}

	newAgent := &Agent{
		apiClient: apiClient,
	}

	newAgent.LoadTools()

	newAgent.sessionManager, err = session.NewSessionManager(*sessionID)
	if err != nil {
		return nil, err
	}

	*sessionID = newAgent.sessionManager.SessionID

	return newAgent, nil
}
