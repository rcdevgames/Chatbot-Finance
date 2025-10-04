package service

import (
	"chatbot/internal/model"
	"chatbot/internal/repository"
	"chatbot/internal/util"
)

type UserService struct {
	userRepo         *repository.UserRepository
	transactionRepo  *repository.TransactionRepository
	chatHistoryRepo  *repository.ChatHistoryRepository
	categoryRepo     *repository.CategoryRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
	transactionRepo *repository.TransactionRepository,
	chatHistoryRepo *repository.ChatHistoryRepository,
	categoryRepo *repository.CategoryRepository,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		chatHistoryRepo: chatHistoryRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *UserService) GetOrCreateUser(telegramUserID int64, firstName, lastName, username string) (*model.User, error) {
	// Check if user exists
	user, err := s.GetUserByTelegramID(telegramUserID)
	if err == nil {
		// Update user info if needed
		if user.FirstName != firstName || user.LastName != lastName || user.TelegramUsername != username {
			user.FirstName = firstName
			user.LastName = lastName
			user.TelegramUsername = username
			s.userRepo.UpdateUser(user)
		}
		return user, nil
	}

	// Create new user
	newUser := &model.User{
		TelegramUserID:    telegramUserID,
		TelegramUsername:  username,
		FirstName:         firstName,
		LastName:          lastName,
		LanguageCode:      "id",
		Timezone:          "Asia/Jakarta",
	}

	err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetUserByTelegramID(telegramUserID int64) (*model.User, error) {
	return s.userRepo.GetUserByTelegramID(telegramUserID)
}

func (s *UserService) AddTransaction(userID string, transaction *model.Transaction) error {
	transaction.UserID = userID
	return s.transactionRepo.CreateTransaction(transaction)
}

func (s *UserService) QueryTransactions(userID string, query *model.QueryParams) ([]model.Transaction, error) {
	startDate, endDate := getDatesFromQuery(query)

	if query.Type == "all" {
		income, err := s.transactionRepo.GetTransactionsByType(userID, "income")
		if err != nil {
			return nil, err
		}
		expense, err := s.transactionRepo.GetTransactionsByType(userID, "expense")
		if err != nil {
			return nil, err
		}
		return append(income, expense...), nil
	}

	return s.transactionRepo.GetTransactionsByUserAndDateRange(userID, startDate, endDate)
}

func (s *UserService) UpdateLastTransaction(userID string, transaction *model.Transaction) error {
	last, err := s.transactionRepo.GetLastTransaction(userID)
	if err != nil {
		return err
	}

	transaction.ID = last.ID
	return s.transactionRepo.UpdateTransaction(transaction)
}

func (s *UserService) DeleteLastTransaction(userID string) error {
	last, err := s.transactionRepo.GetLastTransaction(userID)
	if err != nil {
		return err
	}

	return s.transactionRepo.DeleteTransaction(last.ID)
}

func (s *UserService) GetRecentTransactions(userID string, limit int) ([]model.Transaction, error) {
	return s.transactionRepo.GetTransactionsByUserID(userID, limit)
}

func (s *UserService) SaveChatHistory(history *model.ChatHistory) error {
	return s.chatHistoryRepo.SaveChatHistory(history)
}

func (s *UserService) GetRecentChatHistory(userID string, limit int) ([]model.ChatHistory, error) {
	return s.chatHistoryRepo.GetRecentChatHistory(userID, limit)
}

func (s *UserService) GetFinancialSummary(userID string, period string) (*model.TransactionSummary, error) {
	startDate, endDate := util.GetDateRangeFromPeriod(period)

	income, err := s.transactionRepo.GetTransactionsByUserAndDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	expense, err := s.transactionRepo.GetTransactionsByUserAndDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	summary := &model.TransactionSummary{
		CategoryBreakdown: make(map[string]float64),
	}

	// Calculate income
	for _, t := range income {
		if t.Type == "income" {
			summary.TotalIncome += t.Amount
		}
	}

	// Calculate expense and breakdown
	for _, t := range expense {
		if t.Type == "expense" {
			summary.TotalExpense += t.Amount
			summary.CategoryBreakdown[t.Category] += t.Amount
		}
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpense
	summary.TransactionCount = len(income) + len(expense)

	return summary, nil
}

func (s *UserService) DeleteDataByPeriod(userID string, period string, year int, month int) error {
	switch period {
	case "month":
		// Hapus transaksi per bulan
		if err := s.transactionRepo.DeleteTransactionsByUserIDAndMonth(userID, year, month); err != nil {
			return err
		}
		// Hapus chat history per bulan
		if err := s.chatHistoryRepo.DeleteChatHistoryByUserIDAndMonth(userID, year, month); err != nil {
			return err
		}
	case "year":
		// Hapus transaksi per tahun
		if err := s.transactionRepo.DeleteTransactionsByUserIDAndYear(userID, year); err != nil {
			return err
		}
		// Hapus chat history per tahun
		if err := s.chatHistoryRepo.DeleteChatHistoryByUserIDAndYear(userID, year); err != nil {
			return err
		}
	case "all":
		// Hapus semua transaksi
		if err := s.transactionRepo.DeleteAllTransactionsByUserID(userID); err != nil {
			return err
		}
		// Hapus semua chat history
		if err := s.chatHistoryRepo.DeleteChatHistoryByUserID(userID); err != nil {
			return err
		}
	}
	return nil
}

func getDatesFromQuery(query *model.QueryParams) (string, string) {
	if query.StartDate != "" && query.EndDate != "" {
		return query.StartDate, query.EndDate
	}
	return util.GetDateRangeFromPeriod(query.DateRange)
}