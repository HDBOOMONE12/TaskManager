package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramSender struct {
	token string
}

func NewTelegramSender(token string) *TelegramSender {
	return &TelegramSender{token: token}
}

func (s *TelegramSender) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.token)

	payload := Payload{
		ChatID: chatID,
		Text:   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("http post error: %w", err)
	}
	defer resp.Body.Close()

	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status %d, body: %+v", resp.StatusCode, respBody)
	}

	return nil
}
