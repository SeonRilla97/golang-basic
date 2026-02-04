package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackMessage struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Color  string `json:"color"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
}

func SendToSlack(webhookURL string) func(err any, stack string) {
	return func(err any, stack string) {
		msg := SlackMessage{
			Text: "서버 패닉 발생",
			Attachments: []SlackAttachment{
				{
					Color:  "danger",
					Title:  fmt.Sprintf("Panic: %v", err),
					Text:   truncate(stack, 500),
					Footer: "Go Board API",
				},
			},
		}

		body, _ := json.Marshal(msg)
		http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
