package handler

import (
	"fmt"
	"log"
	"strings"

	"chatbot/internal/model"
	"chatbot/internal/service"
	"chatbot/pkg/telegram"
)

type TelegramHandler struct {
	telegramClient *telegram.Client
	userService    *service.UserService
	llmService     *service.GroqService
}

func NewTelegramHandler(
	telegramClient *telegram.Client,
	userService *service.UserService,
	llmService *service.GroqService,
) *TelegramHandler {
	return &TelegramHandler{
		telegramClient: telegramClient,
		userService:    userService,
		llmService:     llmService,
	}
}

func (h *TelegramHandler) HandleUpdate(update telegram.Update) error {
	if update.Message == nil {
		return nil
	}

	message := update.Message
	user := message.From
	chat := message.Chat
	text := message.Text

	log.Printf("Received message from user %d (%s): %s", user.ID, user.Username, text)

	// Get or create user
	dbUser, err := h.userService.GetOrCreateUser(user.ID, user.FirstName, user.LastName, user.Username)
	if err != nil {
		log.Printf("Error getting/creating user: %v", err)
		return h.sendErrorResponse(chat.ID, "Maaf, ada error sistem. Coba lagi ya.")
	}

	// Handle commands
	if strings.HasPrefix(text, "/") {
		return h.handleCommand(chat.ID, text, dbUser.ID)
	}

	// Process regular message
	return h.handleMessage(chat.ID, text, dbUser)
}

func (h *TelegramHandler) handleCommand(chatID int64, command, userID string) error {
	switch command {
	case "/start":
		return h.handleStart(chatID)
	case "/help":
		return h.handleHelp(chatID)
	case "/summary":
		return h.handleSummary(chatID, userID)
	default:
		return h.sendMessage(chatID, "Command ga dikenal. Ketik /help untuk bantuan.")
	}
}

func (h *TelegramHandler) handleStart(chatID int64) error {
	welcomeMsg := `👋 Halo bro! Selamat datang di Finance Bot Gua.

Gua bantu lo tracking keuangan personal pake chat biasa.

🔥 Fitur-fitur:
• 📝 Input transaksi pake natural language
• 📊 Query data keuangan (mingguan/bulanan)
• 🔄 Update/delete transaksi
• 💡 Auto financial insights

📖 Contoh penggunaan:
• "bayar makan 50rb di warteg"
• "gaji masuk 8 juta kemarin"
• "brp pengeluaran gua minggu ini?"
• "eh salah, harusnya 45rb"

Langsung aja coba!`

	return h.sendMessage(chatID, welcomeMsg)
}

func (h *TelegramHandler) handleHelp(chatID int64) error {
	helpMsg := `🤖 BANTUAN FINANCE BOT

📝 INPUT TRANSAKSI:
• "bayar makan 50rb" → Pengeluaran makanan 50rb
• "gaji 8jt" → Income gaji 8 juta
• "bayar kos 2juta" → Pengeluaran kos 2 juta

📊 QUERY DATA:
• "pengeluaran minggu ini?"
• "income bulan ini berapa?"
• "total gua keluar buat makan?"

🔄 UPDATE/DELETE:
• "salah harusnya 45rb" → Update transaksi terakhir
• "hapus transaksi terakhir" → Delete

📈 INSIGHTS:
• "gimana keuangan gua bulan ini?"
• "kasi summary keuangan dong"

💡 Tips:
• Lo bisa pake format 50rb, 50.000, atau 50rbu
• Tanggal auto detect (hari ini, kemarin, dll)
• Kategori auto detect dari keywords`

	return h.sendMessage(chatID, helpMsg)
}

func (h *TelegramHandler) handleSummary(chatID int64, userID string) error {
	summary, err := h.userService.GetFinancialSummary(userID, "this_month")
	if err != nil {
		return h.sendErrorResponse(chatID, "Gagal ambil summary keuangan.")
	}

	summaryMsg := h.formatSummary(summary)
	return h.sendMessage(chatID, summaryMsg)
}

func (h *TelegramHandler) handleMessage(chatID int64, text string, user *model.User) error {
	// Get chat history for context
	chatHistory, err := h.userService.GetRecentChatHistory(user.ID, 5)
	if err != nil {
		log.Printf("Error getting chat history: %v", err)
		// Continue without history
	}

	// Get recent transactions for context
	recentTransactions, err := h.userService.GetRecentTransactions(user.ID, 3)
	if err != nil {
		log.Printf("Error getting recent transactions: %v", err)
		// Continue without transactions
	}

	// Build context for LLM
	context := h.buildContext(chatHistory, recentTransactions)

	// Process with LLM
	llmRequest := &model.LLMRequest{
		UserInput:    text,
		Context:      context,
		SystemPrompt: fmt.Sprintf("Telegram ID: %d", user.TelegramUserID),
	}

	llmResponse, err := h.llmService.ProcessUserInput(llmRequest)
	if err != nil {
		log.Printf("Error processing with LLM: %v", err)
		return h.sendErrorResponse(chatID, "Maaf, gua bingung. Coba diulang dengan kata-kata yang berbeda.")
	}

	// Save user message to chat history
	userHistory := &model.ChatHistory{
		UserID: user.ID,
		Role:   "user",
		Content: text,
	}
	h.userService.SaveChatHistory(userHistory)

	// Process the intent
	response := h.processIntent(chatID, llmResponse, user)

	// Save bot response to chat history
	botHistory := &model.ChatHistory{
		UserID: user.ID,
		Role:   "assistant",
		Content: response,
	}
	h.userService.SaveChatHistory(botHistory)

	return h.sendMessage(chatID, response)
}

func (h *TelegramHandler) processIntent(chatID int64, response *model.LLMResponse, user *model.User) string {
	switch response.Intent {
	case "add_transaction":
		if response.Transaction != nil {
			err := h.userService.AddTransaction(user.ID, response.Transaction)
			if err != nil {
				log.Printf("Error adding transaction: %v", err)
				return "Maaf, gagal simpan transaksi. Coba lagi ya."
			}
			return response.Response
		}
	case "query_expenses", "query_income":
		if response.Query != nil {
			transactions, err := h.userService.QueryTransactions(user.ID, response.Query)
			if err != nil {
				log.Printf("Error querying transactions: %v", err)
				return "Maaf, gagal ambil data. Coba lagi ya."
			}
			return h.formatQueryResults(transactions, response.Query.Type)
		}
	case "update_last":
		if response.Transaction != nil {
			err := h.userService.UpdateLastTransaction(user.ID, response.Transaction)
			if err != nil {
				log.Printf("Error updating transaction: %v", err)
				return "Maaf, gagal update transaksi. Coba lagi ya."
			}
			return response.Response
		}
	case "delete_last":
		err := h.userService.DeleteLastTransaction(user.ID)
		if err != nil {
			log.Printf("Error deleting transaction: %v", err)
			return "Maaf, gagal hapus transaksi. Coba lagi ya."
		}
		return response.Response
	case "get_summary":
		summary, err := h.userService.GetFinancialSummary(user.ID, "this_month")
		if err != nil {
			log.Printf("Error getting summary: %v", err)
			return "Maaf, gagal ambil summary. Coba lagi ya."
		}
		return h.formatSummary(summary)
	case "ask_advice":
		transactions, _ := h.userService.GetRecentTransactions(user.ID, 10)
		insight, err := h.llmService.GenerateFinancialInsights(transactions, "recent")
		if err != nil {
			log.Printf("Error generating insights: %v", err)
			return response.Response
		}
		return insight
	default:
		return response.Response
	}

	return response.Response
}

func (h *TelegramHandler) buildContext(history []model.ChatHistory, transactions []model.Transaction) string {
	var context strings.Builder

	if len(transactions) > 0 {
		context.WriteString("RECENT TRANSACTIONS:\n")
		for i, t := range transactions {
			if i >= 3 { break }
			context.WriteString(fmt.Sprintf("- %s %s: %s (%s)\n",
				strings.ToUpper(t.Type),
				t.Category,
				h.formatCurrency(t.Amount),
				t.Date))
		}
		context.WriteString("\n")
	}

	if len(history) > 0 {
		context.WriteString("RECENT CHAT:\n")
		for i := len(history) - 1; i >= 0 && i >= len(history)-3; i-- {
			h := history[i]
			role := "User"
			if h.Role == "assistant" {
				role = "Bot"
			}
			context.WriteString(fmt.Sprintf("%s: %s\n", role, h.Content))
		}
	}

	return context.String()
}

func (h *TelegramHandler) formatSummary(summary *model.TransactionSummary) string {
	response := fmt.Sprintf(`📊 SUMMARY KEUANGAN BULAN INI

💰 Income: %s
💸 Expense: %s
💵 Balance: %s (%.1f%%)

`,
		h.formatCurrency(summary.TotalIncome),
		h.formatCurrency(summary.TotalExpense),
		h.formatCurrency(summary.Balance),
		(summary.Balance/summary.TotalIncome)*100,
	)

	if len(summary.CategoryBreakdown) > 0 {
		response += "\n📈 BREAKDOWN:\n"
		for category, amount := range summary.CategoryBreakdown {
			percentage := (amount / summary.TotalExpense) * 100
			response += fmt.Sprintf("%s %s: %s (%.1f%%)\n",
				h.getCategoryIcon(category),
				category,
				h.formatCurrency(amount),
				percentage)
		}
	}

	return response
}

func (h *TelegramHandler) formatQueryResults(transactions []model.Transaction, queryType string) string {
	if len(transactions) == 0 {
		return "Ga ada transaksi yang ditemukan."
	}

	total := 0.0
	response := fmt.Sprintf("📋 Daftar Transaksi (%s):\n\n", strings.ToUpper(queryType))

	for _, t := range transactions {
		total += t.Amount
		response += fmt.Sprintf("%s %s: %s\n%s: %s\n📅: %s\n\n",
			h.getTypeIcon(t.Type),
			strings.ToUpper(t.Type),
			h.formatCurrency(t.Amount),
			h.getCategoryIcon(t.Category),
			t.Category,
			t.Date)
	}

	response += fmt.Sprintf("💵 Total: %s", h.formatCurrency(total))
	return response
}

func (h *TelegramHandler) formatCurrency(amount float64) string {
	return fmt.Sprintf("Rp %.0f", amount)
}

func (h *TelegramHandler) getCategoryIcon(category string) string {
	icons := map[string]string{
		"Makanan":     "🍔",
		"Transport":   "🚗",
		"Tagihan":     "📄",
		"Kos/Sewa":    "🏠",
		"Hiburan":     "🎮",
		"Belanja":     "🛒",
		"Kesehatan":   "💊",
		"Pendidikan":  "📚",
		"Lainnya":     "📦",
		"Gaji":        "💰",
		"Investasi":   "📈",
		"Freelance":   "💼",
		"Bonus":       "🎁",
	}

	if icon, exists := icons[category]; exists {
		return icon
	}
	return "💰"
}

func (h *TelegramHandler) getTypeIcon(transactionType string) string {
	if transactionType == "income" {
		return "💰"
	}
	return "💸"
}

func (h *TelegramHandler) sendMessage(chatID int64, text string) error {
	_, err := h.telegramClient.SendMessage(chatID, text, "HTML")
	return err
}

func (h *TelegramHandler) sendErrorResponse(chatID int64, message string) error {
	errorMsg := fmt.Sprintf("❌ %s", message)
	return h.sendMessage(chatID, errorMsg)
}