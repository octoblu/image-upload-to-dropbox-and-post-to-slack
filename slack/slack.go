package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Slack is the interface for interacting with the Slack API
type Slack interface {
	Post(text string) error
}

type webhookSlack struct {
	webhookURI string
}

type slackMessage struct {
	Text string `json:"text"`
}

// New constructs a new slack instance using a webhook
func New(webhookURI string) Slack {
	return &webhookSlack{webhookURI}
}

func (slack *webhookSlack) Post(text string) error {
	message := &slackMessage{Text: text}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(slack.webhookURI, "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Non 200 status received from slack: %v, %v", resp.StatusCode, resp.Body)
	}

	return nil
}
