package dto

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

/**
Cursor

- 순차 탐색 가능
- 무한 스크롤 때 사용
- 일정 성능
- 구현 난이도 어려움
*/

type Cursor struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// Encode 커서를 문자열로 인코딩
func (c *Cursor) Encode() string {
	data, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(data)
}

// DecodeCursor 문자열을 커서로 디코딩
func DecodeCursor(encoded string) (*Cursor, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var cursor Cursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, err
	}

	return &cursor, nil
}

// CursorPagination 커서 페이징 요청
type CursorPagination struct {
	Cursor string `form:"cursor"`
	Size   int    `form:"size" binding:"min=1,max=100"`
}

// CursorMeta 커서 페이징 메타 정보
type CursorMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
