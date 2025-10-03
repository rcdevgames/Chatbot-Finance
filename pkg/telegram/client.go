package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Token   string
}

type SendMessageRequest struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

type SendMessageResponse struct {
	OK          bool `json:"ok"`
	Result      Message `json:"result"`
	ErrorCode   int  `json:"error_code"`
	Description string `json:"description"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
	Date      int    `json:"date"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Update struct {
	UpdateID int          `json:"update_id"`
	Message  *Message     `json:"message,omitempty"`
	Callback *CallbackQuery `json:"callback_query,omitempty"`
}

type CallbackQuery struct {
	ID      string `json:"id"`
	From    User   `json:"from"`
	Message Message `json:"message"`
	Data    string `json:"data"`
}

type SetWebhookRequest struct {
	URL             string `json:"url"`
	MaxConnections  int    `json:"max_connections,omitempty"`
	AllowedUpdates  []string `json:"allowed_updates,omitempty"`
}

type SetWebhookResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
}

func NewClient(token string) *Client {
	return &Client{
		BaseURL: "https://api.telegram.org/bot" + token,
		Token:   token,
	}
}

func (c *Client) SendMessage(chatID int64, text string, parseMode ...string) (*SendMessageResponse, error) {
	req := SendMessageRequest{
		ChatID:                chatID,
		Text:                  text,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}

	if len(parseMode) > 0 {
		req.ParseMode = parseMode[0]
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.OK {
		return &response, fmt.Errorf("telegram API error: %s", response.Description)
	}

	return &response, nil
}

func (c *Client) SetWebhook(webhookURL string) (*SetWebhookResponse, error) {
	req := SetWebhookRequest{
		URL:            webhookURL,
		MaxConnections: 40,
		AllowedUpdates: []string{"message", "callback_query"},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/setWebhook", c.BaseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response SetWebhookResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.OK {
		return &response, fmt.Errorf("telegram API error: %s", response.Description)
	}

	return &response, nil
}

func (c *Client) GetWebhookInfo() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/getWebhookInfo", c.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}