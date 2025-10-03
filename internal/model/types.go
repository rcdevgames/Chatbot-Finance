package model

import (
	"time"
)

type User struct {
	ID             string    `json:"id"`
	TelegramUserID int64     `json:"telegram_user_id"`
	TelegramUsername string  `json:"telegram_username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	LanguageCode   string    `json:"language_code"`
	Timezone       string    `json:"timezone"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"` // income, expense
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"` // income, expense
	Icon      string   `json:"icon"`
	Keywords  []string `json:"keywords"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatHistory struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	TelegramMessageID int64     `json:"telegram_message_id"`
	Role              string    `json:"role"` // user, assistant
	Content           string    `json:"content"`
	CreatedAt         time.Time `json:"created_at"`
}

type Budget struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount"`
	Period    string    `json:"period"` // daily, weekly, monthly
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Insight struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"` // weekly_summary, monthly_summary, recommendation
	Content     string    `json:"content"`
	Metadata    string    `json:"metadata"` // JSON string
	PeriodStart string    `json:"period_start"`
	PeriodEnd   string    `json:"period_end"`
	CreatedAt   time.Time `json:"created_at"`
}

// LLM Related Types
type LLMRequest struct {
	SystemPrompt string `json:"system_prompt"`
	UserInput    string `json:"user_input"`
	Context      string `json:"context"`
}

type LLMResponse struct {
	Intent       string       `json:"intent"`
	Transaction  *Transaction `json:"transaction,omitempty"`
	Query        *QueryParams `json:"query,omitempty"`
	Response     string       `json:"response"`
	Suggestions  []string     `json:"suggestions,omitempty"`
}

type QueryParams struct {
	Type      string `json:"type"`       // income, expense, all
	DateRange string `json:"date_range"`  // today, this_week, this_month, last_month, custom
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Category  string `json:"category"`
}

// Telegram Types
type TelegramUpdate struct {
	UpdateID int                `json:"update_id"`
	Message  *TelegramMessage   `json:"message,omitempty"`
}

type TelegramMessage struct {
	MessageID int            `json:"message_id"`
	From      *TelegramUser  `json:"from"`
	Chat      *TelegramChat  `json:"chat"`
	Text      string         `json:"text"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type TelegramChat struct {
	ID int64 `json:"id"`
}

// Response Types
type TransactionSummary struct {
	TotalIncome    float64            `json:"total_income"`
	TotalExpense   float64            `json:"total_expense"`
	Balance        float64            `json:"balance"`
	CategoryBreakdown map[string]float64 `json:"category_breakdown"`
	TransactionCount int              `json:"transaction_count"`
}