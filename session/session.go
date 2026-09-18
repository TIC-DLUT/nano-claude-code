package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	stdError "errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TIC-DLUT/nano-claude-code/claude"
	"github.com/TIC-DLUT/nano-claude-code/errors"
)

type head struct {
	SessionID   string `json:"session_id"`
	Cwd         string `json:"cwd"`
	SessionDir  string `json:"session_dir"`
	SessionFile string `json:"-"`
}

type entryDetail struct {
	Type      string         `json:"type"`
	EntryID   string         `json:"entry_id"`
	ParentID  string         `json:"parent_id"`
	Message   claude.Message `json:"message"`
	TimeStamp string         `json:"time_stamp"`
}

type SessionManager struct {
	head
	entryDetails map[string]entryDetail
	nowEntryID   string

	encoder *json.Encoder
}

func NewSessionManager(sessionID string) (*SessionManager, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	homePath, err := os.UserHomeDir()
	if err != nil {
		homePath = ""
	}
	if sessionID == "" {
		sessionID = newID()
	}
	s := &SessionManager{
		head: head{
			SessionID:   sessionID,
			Cwd:         cwd,
			SessionDir:  homePath + "/.nano-claude-code/sessions",
			SessionFile: fmt.Sprintf(homePath+"/.nano-claude-code/sessions/%s.jsonl", sessionID),
		},
		entryDetails: make(map[string]entryDetail),
		nowEntryID:   "",
	}
	err = os.MkdirAll(s.SessionDir, 0755)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(s.SessionFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	// 想办法解决文件关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(
			sigChan,
			os.Interrupt,
			syscall.SIGTERM,
		)
		<-sigChan
		if err := file.Close(); err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
		os.Exit(0)
	}()

	info, _ := file.Stat()
	s.encoder = json.NewEncoder(file)
	decoder := json.NewDecoder(file)

	// 文件为空则先写入 head 信息
	if info.Size() == 0 {
		s.encoder.Encode(struct {
			Type       string `json:"type"`
			HeadDetail head   `json:"head_detail"`
		}{
			Type:       "head",
			HeadDetail: s.head,
		})
	}

	// 读取数据
	for {
		rawMessage := json.RawMessage{}
		err := decoder.Decode(&rawMessage)
		if stdError.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		typeInfo := struct {
			Type string `json:"type"`
		}{}
		json.Unmarshal(rawMessage, &typeInfo)

		switch typeInfo.Type {
		case "head":
			json.Unmarshal(rawMessage, &s.head)
		default:
			entrydetail := entryDetail{}
			json.Unmarshal(rawMessage, &entrydetail)

			decodeEntry(&entrydetail.Message)

			s.entryDetails[entrydetail.EntryID] = entrydetail
			s.nowEntryID = entrydetail.EntryID
		}
	}

	return s, nil
}

func (s *SessionManager) Append(messages ...claude.Message) {
	// 解析 messages 到 entry 树
	for _, message := range messages {
		content := message.Content
		entrydetail := entryDetail{
			Type:      "message",
			EntryID:   newID(),
			ParentID:  s.nowEntryID,
			Message:   message,
			TimeStamp: time.Now().Format("2006-01-02 15:04:05")}
		switch content.(type) {
		case claude.TextBlock, claude.SingleStringMessage, claude.ToolUseBlock, claude.ToolResultBlock, claude.ThinkingBlock:
			entrydetail.Type = "message"
		}
		s.entryDetails[entrydetail.EntryID] = entrydetail
		s.nowEntryID = entrydetail.EntryID

		// 写入文件
		s.encoder.Encode(entrydetail)
	}
}

func (s *SessionManager) BuildSessionContext() ([]claude.Message, error) {
	path := []entryDetail{}
	for i := s.nowEntryID; i != ""; {
		nowEntryDetail, ok := s.entryDetails[i]
		if !ok {
			return []claude.Message{}, errors.UnknowEntry
		}
		path = append(path, nowEntryDetail)
		i = nowEntryDetail.ParentID
	}
	messages := []claude.Message{}
	for i := len(path) - 1; i >= 0; i-- {
		messages = append(messages, path[i].Message)
	}
	return messages, nil
}

func (s *SessionManager) Fork(entryID string) error {
	if _, ok := s.entryDetails[entryID]; !ok {
		return errors.UnknowEntry
	}
	s.nowEntryID = entryID
	return nil
}

func newID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func decodeEntry(message *claude.Message) {
	switch c := message.Content.(type) {
	case string:
		message.Content = claude.SingleStringMessage(c)
	case map[string]any:
		switch c["type"] {
		case "text":
			message.Content = claude.TextBlock{
				Type: "text",
				Text: c["text"].(string),
			}
		case "thinking":
			message.Content = claude.ThinkingBlock{
				Type:      "thinking",
				Thinking:  c["thinking"].(string),
				Signature: c["signature"].(string),
			}
		case "tool_use":
			message.Content = claude.ToolUseBlock{
				Type:  "tool_use",
				Name:  c["name"].(string),
				Input: c["input"].(map[string]any),
				ID:    c["id"].(string),
			}
		case "tool_result":
			message.Content = claude.ToolResultBlock{
				Type:      "tool_result",
				Content:   c["content"].(string),
				ToolUseID: c["tool_use_id"].(string),
			}
		default:
		}
	default:
	}
}
