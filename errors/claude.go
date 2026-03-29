package errors

import "errors"

var (
	CreateClaudeClientBaseUrlError = errors.New("unsupport baseurl")
	ClaudeClientCallFormatError    = errors.New("response format parse error")
	ClaudeCreateToolEmptyError     = errors.New("name or description can't be empty")
)
