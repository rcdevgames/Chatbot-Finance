package model

import (
	"time"
)

type User struct {
	ID                string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	TelegramUserID    int64     `gorm:"uniqueIndex;not null" json:"telegram_user_id"`
	TelegramUsername  string    `gorm:"index" json:"telegram_username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	LanguageCode      string    `json:"language_code"`
	Timezone          string    `gorm:"default:'Asia/Jakarta'" json:"timezone"`
	TrialExpiresAt    *time.Time `gorm:"index" json:"trial_expires_at"`
	IsTrialActive     bool      `gorm:"default:false" json:"is_trial_active"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type Transaction struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        string    `gorm:"type:varchar(10);not null;check:type IN ('income', 'expense')" json:"type"` // income, expense
	Amount      float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Category    string    `gorm:"type:varchar(100);not null;index" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	Date        string    `gorm:"type:date;not null;index" json:"date"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type Category struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Type      string    `gorm:"type:varchar(10);not null;check:type IN ('income', 'expense')" json:"type"` // income, expense
	Icon      string    `gorm:"type:varchar(50)" json:"icon"`
	Keywords  []string  `gorm:"type:text;serializer:json" json:"keywords"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Category) TableName() string {
	return "categories"
}

type ChatHistory struct {
	ID                string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID            string    `gorm:"type:uuid;not null;index" json:"user_id"`
	TelegramMessageID int64     `gorm:"index" json:"telegram_message_id"`
	Role              string    `gorm:"type:varchar(20);not null;check:role IN ('user', 'assistant')" json:"role"` // user, assistant
	Content           string    `gorm:"type:text;not null" json:"content"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (ChatHistory) TableName() string {
	return "chat_history"
}

type Budget struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Category  string    `gorm:"type:varchar(100);not null;index" json:"category"`
	Amount    float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Period    string    `gorm:"type:varchar(20);not null;check:period IN ('daily', 'weekly', 'monthly')" json:"period"` // daily, weekly, monthly
	StartDate string    `gorm:"type:date;not null" json:"start_date"`
	EndDate   string    `gorm:"type:date;not null" json:"end_date"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Budget) TableName() string {
	return "budgets"
}

type Insight struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        string    `gorm:"type:varchar(50);not null;index" json:"type"` // weekly_summary, monthly_summary, recommendation
	Content     string    `gorm:"type:text;not null" json:"content"`
	Metadata    string    `gorm:"type:text;serializer:json" json:"metadata"` // JSON string
	PeriodStart string    `gorm:"type:date;index" json:"period_start"`
	PeriodEnd   string    `gorm:"type:date;index" json:"period_end"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Insight) TableName() string {
	return "insights"
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