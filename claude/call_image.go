package claude

// 从 URL 创建 ImageBlock
func NewImageBlockFromUrl(url string) ImageBlock {
	return ImageBlock{
		Type: "image",
		Source: ImageBlockSourceUrl{
			Type: "url",
			Url:  url,
		},
	}
}

// base64Data为base64编码后的图片数据（不包含data:image/jpeg;base64,前缀）
// mediaType为图片的media type ，例如image/jpeg
func NewImageBlockFromBase64(mediaType string, base64Data string) ImageBlock {
	return ImageBlock{
		Type: "image",
		Source: ImageBlockSourceBase64{
			Type:      "base64",
			MediaType: mediaType,
			Data:      base64Data,
		},
	}
}

// 图片消息
type ImageBlock struct {
	Type   string `json:"type"`
	Source any    `json:"source"`
}

// 图片消息Base64
type ImageBlockSourceBase64 struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

// 图片消息URL
type ImageBlockSourceUrl struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}
