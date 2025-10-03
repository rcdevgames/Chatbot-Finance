package service

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"chatbot/internal/model"
	"chatbot/internal/repository"
	"chatbot/pkg/telegram"
)

type CronService struct {
	cron           *cron.Cron
	userRepo       *repository.UserRepository
	telegramClient *telegram.Client
	llmService     *GroqService
}

func NewCronService(
	userRepo *repository.UserRepository,
	telegramClient *telegram.Client,
	llmService *GroqService,
) *CronService {
	return &CronService{
		cron:           cron.New(),
		userRepo:       userRepo,
		telegramClient: telegramClient,
		llmService:     llmService,
	}
}

func (s *CronService) Start() {
	// Schedule weekly summary every Sunday at 9 AM
	s.cron.AddFunc("0 9 * * 0", s.sendWeeklySummaries)

	// Schedule monthly summary on the 1st of each month at 9 AM
	s.cron.AddFunc("0 9 1 * *", s.sendMonthlySummaries)

	// Schedule spending alerts check every day at 8 PM
	s.cron.AddFunc("0 20 * * *", s.checkSpendingAlerts)

	s.cron.Start()
	log.Println("Cron service started")
}

func (s *CronService) Stop() {
	s.cron.Stop()
	log.Println("Cron service stopped")
}

func (s *CronService) sendWeeklySummaries() {
	log.Println("Sending weekly summaries...")

	// Get all users
	users, err := s.getAllUsers()
	if err != nil {
		log.Printf("Error getting users for weekly summary: %v", err)
		return
	}

	for _, user := range users {
		go func(u model.User) {
			if err := s.sendWeeklySummaryToUser(u); err != nil {
				log.Printf("Error sending weekly summary to user %d: %v", u.TelegramUserID, err)
			}
		}(user)
	}
}

func (s *CronService) sendMonthlySummaries() {
	log.Println("Sending monthly summaries...")

	// Get all users
	users, err := s.getAllUsers()
	if err != nil {
		log.Printf("Error getting users for monthly summary: %v", err)
		return
	}

	for _, user := range users {
		go func(u model.User) {
			if err := s.sendMonthlySummaryToUser(u); err != nil {
				log.Printf("Error sending monthly summary to user %d: %v", u.TelegramUserID, err)
			}
		}(user)
	}
}

func (s *CronService) checkSpendingAlerts() {
	log.Println("Checking spending alerts...")
	// Implementation for spending alerts
}

func (s *CronService) sendWeeklySummaryToUser(user model.User) error {
	// Get transactions for this week
	startDate, endDate := getWeekDateRange(time.Now())

	transactions, err := s.getTransactionsForUser(user.ID, startDate, endDate)
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil // No transactions this week
	}

	// Generate weekly insight
	insight, err := s.llmService.GenerateFinancialInsights(transactions, "weekly")
	if err != nil {
		return err
	}

	// Format message
	message := s.formatWeeklyMessage(insight, transactions)

	// Send to user
	_, err = s.telegramClient.SendMessage(user.TelegramUserID, message, "HTML")
	if err != nil {
		return err
	}

	// Save insight to database
	s.saveInsight(user.ID, "weekly_summary", message, startDate, endDate)

	return nil
}

func (s *CronService) sendMonthlySummaryToUser(user model.User) error {
	// Get transactions for this month
	startDate, endDate := getMonthDateRange(time.Now())

	transactions, err := s.getTransactionsForUser(user.ID, startDate, endDate)
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil // No transactions this month
	}

	// Generate monthly insight
	insight, err := s.llmService.GenerateFinancialInsights(transactions, "monthly")
	if err != nil {
		return err
	}

	// Format message
	message := s.formatMonthlyMessage(insight, transactions)

	// Send to user
	_, err = s.telegramClient.SendMessage(user.TelegramUserID, message, "HTML")
	if err != nil {
		return err
	}

	// Save insight to database
	s.saveInsight(user.ID, "monthly_summary", message, startDate, endDate)

	return nil
}

func (s *CronService) getAllUsers() ([]model.User, error) {
	// This would need to be implemented in UserRepository
	// For now, return empty slice
	return []model.User{}, nil
}

func (s *CronService) getTransactionsForUser(userID, startDate, endDate string) ([]model.Transaction, error) {
	// This would need to be implemented in TransactionRepository
	// For now, return empty slice
	return []model.Transaction{}, nil
}

func (s *CronService) saveInsight(userID, insightType, content, startDate, endDate string) {
	// Save insight to database for future reference
	_ = &model.Insight{
		UserID:      userID,
		Type:        insightType,
		Content:     content,
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}

	// Save to repository (implementation needed)
}

func (s *CronService) formatWeeklyMessage(insight string, transactions []model.Transaction) string {
	return `📊 <b>WEEKLY FINANCE SUMMARY</b>

` + insight + `

💡 <i>Weekly insights generated automatically</i>`
}

func (s *CronService) formatMonthlyMessage(insight string, transactions []model.Transaction) string {
	return `📈 <b>MONTHLY FINANCE REPORT</b>

` + insight + `

🎯 <i>Monthly report generated automatically</i>`
}

func getWeekDateRange(t time.Time) (string, string) {
	// Get Monday of current week
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}

	startOfWeek := t.AddDate(0, 0, -(weekday - 1))
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	return startOfWeek.Format("2006-01-02"), endOfWeek.Format("2006-01-02")
}

func getMonthDateRange(t time.Time) (string, string) {
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	return startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02")
}