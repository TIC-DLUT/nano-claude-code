package claude

import (
	"strings"

	"github.com/TIC-DLUT/nano-claude-code/errors"
	"resty.dev/v3"
)

type ClaudeClient struct {
	baseUrl    string
	apiKey     string
	httpClient *resty.Client
}

func NewClient(baseurl string, apikey string) (*ClaudeClient, error) {
	// 确保 baseurl 合法
	if !strings.HasPrefix(baseurl, "http://") && !strings.HasPrefix(baseurl, "https://") {
		return nil, errors.CreateClaudeClientBaseUrlError
	}
	return &ClaudeClient{
		baseUrl:    baseurl,
		apiKey:     apikey,
		httpClient: resty.New(),
	}, nil
}
