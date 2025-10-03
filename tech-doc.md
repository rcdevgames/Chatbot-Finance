# Tech Spec Document: Telegram Finance Bot

## 1. Fitur-Fitur Utama

### 1.1 Conversational Transaction Tracking
- **Natural Language Input**: User chat bebas, bot auto-detect intent
  - "bayar kos 2 juta" → Auto save transaksi
  - "gaji masuk 8 juta kemarin" → Save dengan date parsing
  - "eh salah, harusnya 1.8 juta" → Update transaksi terakhir
- **Multi-format Support**: Text biasa, voice note (optional), inline buttons

### 1.2 Smart Query & Analytics
- **Flexible Queries**: Chat natural untuk tanya data
  - "brp duit gua keluar minggu ini?"
  - "total pengeluaran makan bulan ini?"
  - "bandingin bulan ini sama bulan lalu dong"
- **Auto Date Parsing**: Ngerti "hari ini", "kemarin", "minggu ini", "bulan lalu", etc
- **Contextual Response**: Bot inget chat sebelumnya

### 1.3 Financial Insights
- **Auto Summary**: Bot kasih summary berkala (mingguan/bulanan)
- **Spending Alerts**: Notif kalau spending unusual atau over budget
- **Recommendations**: Saran investasi/saving based on pattern

### 1.4 Transaction Management
- **CRUD via Chat**: 
  - Create: "bayar X untuk Y"
  - Read: "spending gua brp?"
  - Update: "eh salah, harusnya Z"
  - Delete: "hapus transaksi terakhir"
- **Categories**: Auto-detect atau tanya user kalau ambiguous

---

## 2. Alur Aplikasi

### 2.1 Bot Registration Flow
```
1. User start bot (/start)
2. Bot greeting & onboarding singkat
3. Save Telegram User ID ke database
4. Ready to use
```

### 2.2 Transaction Input Flow
```
User: "bayar makan 50rb di warteg"
  ↓
Telegram API → Webhook → Golang Handler
  ↓
Get user chat history (last 10 messages) untuk context
  ↓
Call Groq API dengan prompt:
  - System prompt (bot personality + capabilities)
  - User chat history
  - Current message
  ↓
Groq response:
  {
    "intent": "add_transaction",
    "transaction": {
      "type": "expense",
      "amount": 50000,
      "category": "Makanan",
      "description": "warteg",
      "date": "2025-10-03"
    },
    "response": "Oke gua catet makan 50rb di warteg ✅"
  }
  ↓
Save to Supabase (transactions table)
  ↓
Send response ke Telegram
```

### 2.3 Query Flow
```
User: "brp duit gua keluar minggu ini?"
  ↓
Golang Handler + Groq parse intent
  ↓
Groq response:
  {
    "intent": "query_expenses",
    "query": {
      "type": "expense",
      "date_range": "this_week"
    }
  }
  ↓
Fetch from Supabase (WHERE type='expense' AND date >= start_of_week)
  ↓
Calculate & format response
  ↓
Send ke Telegram:
  "Minggu ini lo udah keluar Rp 1.2 juta bro
   
   Breakdown:
   🍔 Makan - 600rb (50%)
   🚗 Transport - 400rb (33%)
   🎮 Hiburan - 200rb (17%)"
```

### 2.4 Update/Delete Flow
```
User: "eh salah harusnya 45rb"
  ↓
Groq detect intent: "update_last_transaction"
  ↓
Fetch last transaction from DB
  ↓
Update amount: 50000 → 45000
  ↓
Response: "Okee gua update jadi 45rb ✅"
```

### 2.5 Insight Flow (Auto-triggered)
```
Cron Job (setiap Minggu/Bulan)
  ↓
Fetch all users
  ↓
For each user:
  - Get transactions data
  - Send to Groq for analysis
  - Generate insights
  ↓
Send notification via Telegram:
  "📊 Summary Bulan Ini:
   Income: 8jt | Expense: 5.5jt | Saving: 2.5jt (31%)
   
   💡 Tips: Lo hemat bulan ini! Coba invest 1jt di..."
```

---

## 3. Database Schema (Supabase)

### 3.1 Table: `users`
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  telegram_user_id BIGINT NOT NULL UNIQUE,
  telegram_username TEXT,
  first_name TEXT,
  last_name TEXT,
  language_code TEXT DEFAULT 'id',
  timezone TEXT DEFAULT 'Asia/Jakarta',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_telegram_id ON users(telegram_user_id);
```

### 3.2 Table: `transactions`
```sql
CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
  amount DECIMAL(12, 2) NOT NULL,
  category TEXT NOT NULL,
  description TEXT,
  date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_date ON transactions(date DESC);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_category ON transactions(category);

-- Composite index untuk query umum
CREATE INDEX idx_transactions_user_date ON transactions(user_id, date DESC);
```

### 3.3 Table: `categories`
```sql
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
  icon TEXT, -- emoji
  keywords TEXT[], -- untuk auto-detection
  created_at TIMESTAMP DEFAULT NOW()
);

-- Seed data
INSERT INTO categories (name, type, icon, keywords) VALUES
  ('Gaji', 'income', '💰', ARRAY['gaji', 'salary', 'income']),
  ('Investasi', 'income', '📈', ARRAY['investasi', 'dividend', 'profit']),
  ('Freelance', 'income', '💼', ARRAY['freelance', 'project', 'client']),
  ('Bonus', 'income', '🎁', ARRAY['bonus', 'thr', 'insentif']),
  
  ('Makanan', 'expense', '🍔', ARRAY['makan', 'food', 'resto', 'warteg', 'lunch', 'dinner']),
  ('Transport', 'expense', '🚗', ARRAY['transport', 'bensin', 'grab', 'gojek', 'parkir']),
  ('Tagihan', 'expense', '📄', ARRAY['tagihan', 'listrik', 'air', 'internet', 'wifi', 'token']),
  ('Kos/Sewa', 'expense', '🏠', ARRAY['kos', 'sewa', 'rent', 'kontrakan']),
  ('Hiburan', 'expense', '🎮', ARRAY['hiburan', 'nonton', 'game', 'netflix', 'spotify']),
  ('Belanja', 'expense', '🛒', ARRAY['belanja', 'shopping', 'beli', 'tokped', 'shopee']),
  ('Kesehatan', 'expense', '💊', ARRAY['kesehatan', 'obat', 'dokter', 'hospital']),
  ('Pendidikan', 'expense', '📚', ARRAY['kursus', 'buku', 'course', 'training']),
  ('Lainnya', 'expense', '📦', ARRAY['lainnya', 'other', 'misc']);
```

### 3.4 Table: `chat_history`
```sql
CREATE TABLE chat_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  telegram_message_id BIGINT,
  role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_chat_history_user_id ON chat_history(user_id, created_at DESC);

-- Auto-cleanup old messages (keep last 50 per user)
CREATE OR REPLACE FUNCTION cleanup_old_chat_history()
RETURNS trigger AS $$
BEGIN
  DELETE FROM chat_history
  WHERE user_id = NEW.user_id
  AND id NOT IN (
    SELECT id FROM chat_history
    WHERE user_id = NEW.user_id
    ORDER BY created_at DESC
    LIMIT 50
  );
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cleanup_chat_history
AFTER INSERT ON chat_history
FOR EACH ROW
EXECUTE FUNCTION cleanup_old_chat_history();
```

### 3.5 Table: `budgets` (Optional - Future feature)
```sql
CREATE TABLE budgets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  category TEXT REFERENCES categories(name),
  amount DECIMAL(12, 2) NOT NULL,
  period TEXT CHECK (period IN ('daily', 'weekly', 'monthly')),
  start_date DATE NOT NULL,
  end_date DATE,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_budgets_user_id ON budgets(user_id);
```

### 3.6 Table: `insights` (Optional - Cache AI insights)
```sql
CREATE TABLE insights (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL, -- 'weekly_summary', 'monthly_summary', 'recommendation'
  content TEXT NOT NULL,
  metadata JSONB, -- store additional data
  period_start DATE,
  period_end DATE,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_insights_user_period ON insights(user_id, period_start, period_end);
```

---

## 4. Tech Stack

### 4.1 Backend (Golang)
```
/cmd
  /bot
    main.go                    # Entry point
/internal
  /handler
    telegram.go                # Handle Telegram updates
    transaction.go             # Transaction CRUD logic
    query.go                   # Query & analytics logic
  /service
    llm.go                     # Groq API integration
    finance.go                 # Financial calculations
    category.go                # Auto-categorization
  /repository
    supabase.go                # Database operations
    user.go
    transaction.go
  /model
    types.go                   # Structs & types
  /config
    config.go                  # Environment config
  /util
    date.go                    # Date parsing utilities
    formatter.go               # Response formatting
/pkg
  /telegram
    client.go                  # Telegram Bot API wrapper
```

### 4.2 Key Dependencies
```go
require (
    github.com/go-telegram-bot-api/telegram-bot-api/v5
    github.com/supabase-community/supabase-go
    github.com/joho/godotenv
    github.com/gin-gonic/gin // untuk webhook
)
```

### 4.3 Environment Variables
```bash
# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token

# Supabase
SUPABASE_URL=https://xxx.supabase.co
SUPABASE_KEY=your_anon_key

# Groq
GROQ_API_KEY=gsk_xxx

# Server
PORT=8080
WEBHOOK_URL=https://yourdomain.com/webhook
```

---

## 5. LLM Integration

### 5.1 System Prompt Template
```
Kamu adalah temen gua yang bantu tracking keuangan personal.

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

USER INFO:
- Telegram ID: {telegram_user_id}
- Timezone: Asia/Jakarta
- Currency: IDR

RECENT TRANSACTIONS (for context):
{recent_transactions}
```

### 5.2 Groq Model Config
```go
type GroqConfig struct {
    Model       string  // "llama-3.1-70b-versatile"
    Temperature float64 // 0.3 (lebih deterministic)
    MaxTokens   int     // 1000
}
```

---

## 6. Deployment

### 6.1 Options
- **Railway**: One-click deploy Golang + webhook
- **Fly.io**: Free tier, auto-scale
- **Google Cloud Run**: Serverless, pay-per-use
- **VPS**: Full control (DigitalOcean, Linode)

### 6.2 Webhook Setup
```bash
# Set webhook ke Telegram
curl -X POST \
  "https://api.telegram.org/bot<TOKEN>/setWebhook" \
  -d "url=https://yourdomain.com/webhook"
```

---

## 7. Estimasi & Biaya

**Development Time**: 1-2 minggu (solo developer)

**Monthly Cost**: 
- Telegram Bot: **FREE**
- Supabase: **FREE** (up to 500MB database, 2GB bandwidth)
- Groq: **FREE** (14,400 requests/day)
- Hosting: **$5-10/month** (Railway/Fly.io)

**Total**: ~$5-10/bulan untuk hosting aja

---

## 8. Contoh Interaksi

### Scenario 1: Input Transaksi
```
User: "bayar kos 2 juta tadi pagi"
Bot: "Oke gua catet pengeluaran kos 2 juta hari ini ✅
     Total pengeluaran bulan ini jadi 3.2 juta"
```

### Scenario 2: Query Data
```
User: "brp duit gua keluar minggu ini?"
Bot: "Minggu ini lo udah keluar Rp 1.2 juta bro
     
     Breakdown:
     🍔 Makan - 600rb (50%)
     🚗 Transport - 400rb (33%)
     🎮 Hiburan - 200rb (17%)"
```

### Scenario 3: Update Transaksi
```
User: "eh salah, harusnya 1.8 juta bukan 2 juta"
Bot: "Okee gua update kos jadi 1.8 juta ✅
     Total pengeluaran bulan ini jadi 3 juta"
```

### Scenario 4: Financial Insight
```
User: "gimana keuangan gua bulan ini?"
Bot: "Bulan ini:
     💰 Masuk: 8 juta
     💸 Keluar: 3 juta  
     💵 Sisa: 5 juta (62.5%)
     
     Lo hemat banget bulan ini! Pengeluaran turun 30% dari bulan lalu.
     
     💡 Saran: Dari 5 juta sisanya, coba:
     - 3 juta masuk reksadana/deposito
     - 1.5 juta emergency fund
     - 500rb bebas dipake"
```

### Scenario 5: Ambiguous Input
```
User: "bayar 100rb"
Bot: "100rb buat apa nih bro? 
     Pilih kategori atau jelasin aja:"
     [🍔 Makan] [🚗 Transport] [🎮 Hiburan] [💊 Kesehatan]
```
