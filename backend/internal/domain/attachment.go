package domain

type Attachment struct {
	URL       string `json:"url,omitempty"`
	Data      []byte `json:"-"`
	MimeType  string `json:"mimeType"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize,omitempty"`
	MediaType string `json:"mediaType"`
}
