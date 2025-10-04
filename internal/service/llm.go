package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"chatbot/internal/model"
)

type GroqService struct {
	apiKey  string
	baseURL string
	model   string
}

type GroqRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage   `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewGroqService(apiKey string) *GroqService {
	return &GroqService{
		apiKey:  apiKey,
		baseURL: "https://api.groq.com/openai/v1/chat/completions",
		model:   "llama-3.3-70b-versatile",
	}
}

func (g *GroqService) ProcessUserInput(request *model.LLMRequest) (*model.LLMResponse, error) {
	systemPrompt := g.buildSystemPrompt()

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: g.buildUserPrompt(request)},
	}

	groqReq := GroqRequest{
		Model:       g.model,
		Messages:    messages,
		Temperature: 0.3,
		MaxTokens:   1000,
	}

	jsonData, err := json.Marshal(groqReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Groq request: %w", err)
	}

	req, err := http.NewRequest("POST", g.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Groq: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Groq response: %w. Response body: %s", err, string(body))
	}

	if len(groqResp.Choices) == 0 {
		// Log the full response for debugging
		return nil, fmt.Errorf("no choices in Groq response. Status: %d, Body: %s", resp.StatusCode, string(body))
	}

	content := groqResp.Choices[0].Message.Content
	return g.parseLLMResponse(content)
}

func (g *GroqService) buildSystemPrompt() string {
	return `Kamu adalah temen gua yang bantu tracking keuangan personal.

PERSONALITY:
- Ngomong santai, informal, pake "gua/lo"
- Friendly & helpful
- Boleh pake slang Indonesia

CAPABILITIES:
1. Parse transaksi dari natural language
2. Query data keuangan (filter by date, category, type)
3. Update/delete transaksi terakhir
4. Kasih insight & financial advice
5. Deteksi tanggal relatif (hari ini, kemarin, minggu ini, bulan lalu, etc)

OUTPUT FORMAT (JSON):
{
  "intent": "add_transaction|query_expenses|query_income|update_last|delete_last|get_summary|ask_advice",
  "transaction": {
    "type": "income|expense",
    "amount": number,
    "category": "string",
    "description": "string",
    "date": "YYYY-MM-DD"
  },
  "query": {
    "type": "income|expense|all",
    "date_range": "today|this_week|this_month|last_month|custom",
    "start_date": "YYYY-MM-DD",
    "end_date": "YYYY-MM-DD",
    "category": "string"
  },
  "response": "string (what to reply to user)"
}

RULES:
- Always respond in Indonesian
- Format currency: Rp 1.200.000 (pake titik separator)
- Kalau ambiguous, tanya balik
- Default date adalah hari ini kalau ga disebutin
- Kalau user bilang "salah" atau "update", refer to last transaction
- Amount dalam rupiah, tanpa format (contoh: 50000 untuk 50rb)

USER INFO:
- Timezone: Asia/Jakarta
- Currency: IDR`
}

func (g *GroqService) buildUserPrompt(request *model.LLMRequest) string {
	prompt := fmt.Sprintf("USER INPUT: %s\n", request.UserInput)

	if request.Context != "" {
		prompt += fmt.Sprintf("CONTEXT:\n%s\n", request.Context)
	}

	if request.SystemPrompt != "" {
		prompt += fmt.Sprintf("ADDITIONAL CONTEXT:\n%s\n", request.SystemPrompt)
	}

	prompt += "\nParse user input dan return dalam format JSON."

	return prompt
}

func (g *GroqService) parseLLMResponse(content string) (*model.LLMResponse, error) {
	// Clean up response - extract JSON if it's wrapped in code blocks
	jsonContent := content
	if strings.Contains(content, "```json") {
		start := strings.Index(content, "```json") + 7
		end := strings.LastIndex(content, "```")
		if end > start {
			jsonContent = content[start:end]
		}
	} else if strings.Contains(content, "```") {
		start := strings.Index(content, "```") + 3
		end := strings.LastIndex(content, "```")
		if end > start {
			jsonContent = content[start:end]
		}
	}

	var response model.LLMResponse
	if err := json.Unmarshal([]byte(jsonContent), &response); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response as JSON: %w\nContent: %s", err, jsonContent)
	}

	return &response, nil
}

func (g *GroqService) GenerateFinancialInsights(transactions []model.Transaction, period string) (string, error) {
	prompt := g.buildInsightPrompt(transactions, period)

	messages := []Message{
		{Role: "system", Content: "Kamu adalah financial advisor yang friendly. Berikan insight keuangan dalam bahasa Indonesia yang santai dan actionable."},
		{Role: "user", Content: prompt},
	}

	groqReq := GroqRequest{
		Model:       g.model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   800,
	}

	jsonData, err := json.Marshal(groqReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Groq request: %w", err)
	}

	req, err := http.NewRequest("POST", g.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Groq: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal Groq response: %w. Response body: %s", err, string(body))
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in Groq response. Status: %d, Body: %s", resp.StatusCode, string(body))
	}

	return groqResp.Choices[0].Message.Content, nil
}

func (g *GroqService) buildInsightPrompt(transactions []model.Transaction, period string) string {
	prompt := fmt.Sprintf(`Analisis data transaksi berikut untuk periode %s dan berikan insight keuangan:

`, period)

	// Group transactions by type and category
	incomeByCategory := make(map[string]float64)
	expenseByCategory := make(map[string]float64)
	totalIncome := 0.0
	totalExpense := 0.0

	for _, t := range transactions {
		if t.Type == "income" {
			incomeByCategory[t.Category] += t.Amount
			totalIncome += t.Amount
		} else {
			expenseByCategory[t.Category] += t.Amount
			totalExpense += t.Amount
		}
	}

	prompt += fmt.Sprintf("Total Income: Rp %.0f\n", totalIncome)
	prompt += fmt.Sprintf("Total Expense: Rp %.0f\n", totalExpense)
	prompt += fmt.Sprintf("Balance: Rp %.0f\n\n", totalIncome-totalExpense)

	if len(incomeByCategory) > 0 {
		prompt += "Income Breakdown:\n"
		for cat, amount := range incomeByCategory {
			percentage := (amount / totalIncome) * 100
			prompt += fmt.Sprintf("- %s: Rp %.0f (%.1f%%)\n", cat, amount, percentage)
		}
		prompt += "\n"
	}

	if len(expenseByCategory) > 0 {
		prompt += "Expense Breakdown:\n"
		for cat, amount := range expenseByCategory {
			percentage := (amount / totalExpense) * 100
			prompt += fmt.Sprintf("- %s: Rp %.0f (%.1f%%)\n", cat, amount, percentage)
		}
		prompt += "\n"
	}

	prompt += `Berikan insight dalam format:
1. Summary singkat
2. 3-5 highlights/pola menarik
3. 2-3 saran actionable
4. Format dengan emoji dan friendly tone

Gunakan bahasa Indonesia yang santai (gua/lo).`

	return prompt
}