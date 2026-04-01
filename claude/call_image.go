package claude

import (
	"encoding/json"
	"errors"
)

type ImageReq struct {
	Type              string `json:"type"`              // 图片类型 base64 or url
	Source            string `json:"source"`            // 图片资源, base64 or url
	Cache_control_ttl string `json:"cache_control_ttl"` // 缓存控制时间, 5m or 1h or default("5m")
}

// MarshalJSON 实现 json.Marshaler 接口，自动转换为 ImageBlock 格式
func (r ImageReq) MarshalJSON() ([]byte, error) {
	block, err := r.ToImageBlock()
	if err != nil {
		return nil, err
	}
	return json.Marshal(block)
}

// ToImageBlock 将 ImageReq 转换为 ImageBlock
func (r *ImageReq) ToImageBlock() (ImageBlock, error) {
	block := ImageBlock{
		Type: "image",
		Cache_control: ImageBlockCacheControl{
			Type: "ephemeral",
			Ttl:  r.Cache_control_ttl,
		},
	}
	// 缓存控制时间默认值为 5m
	if r.Cache_control_ttl != "1h" {
		block.Cache_control.Ttl = "5m"
	}

	switch r.Type {
	case "url":
		block.Source = UrlImageSource{
			Type: "url",
			Url:  r.Source,
		}
	case "base64":
		// 解析 base64 数据，尝试提取 media type
		mediaType, data, err := parseBase64Data(r.Source)
		if err != nil {
			return block, errors.New("invalid base64 data: " + err.Error())
		}
		block.Source = Base64ImageSource{
			Type:      "string",
			MediaType: mediaType,
			Data:      data,
		}
	}

	return block, nil
}

// parseBase64Data 解析 base64 数据，提取 media type 和实际数据
// 支持格式: "data:image/jpeg;base64,xxxxx"
func parseBase64Data(source string) (mediaType, data string, err error) {
	// 默认值
	if source[:5] != "data:" {
		return "", "", errors.New("invalid base64 data")
	}

	// 检查是否是 data URI 格式
	if len(source) > 5 && source[:5] == "data:" {
		// 查找逗号分隔符
		commaIdx := findComma(source)
		if commaIdx != -1 {
			// 提取 media type 部分 (data:image/jpeg;base64 -> image/jpeg;base64)
			mediaType = source[5:commaIdx]
			data = source[commaIdx+1:]
			if mediaType[len(mediaType)-7:] != ";base64" {
				return "", "", errors.New("invalid base64 data")
			}
		} else {
			return "", "", errors.New("invalid base64 data")
		}
	}

	return mediaType, data, nil
}

// findComma 查找字符串中的逗号位置
func findComma(s string) int {
	for i, c := range s {
		if c == ',' {
			return i
		}
	}
	return -1
}

// findSemicolon 查找字符串中的分号位置
func findSemicolon(s string) int {
	for i, c := range s {
		if c == ';' {
			return i
		}
	}
	return -1
}

// 图片消息
type ImageBlock struct {
	Type          string                 `json:"type"`
	Source        any                    `json:"source"`
	Cache_control ImageBlockCacheControl `json:"cache_control"`
}

// 图片消息Base64
type Base64ImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// 图片消息URL
type UrlImageSource struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

// 图片消息缓存控制
type ImageBlockCacheControl struct {
	Type string `json:"type"`
	Ttl  string `json:"ttl"`
}
