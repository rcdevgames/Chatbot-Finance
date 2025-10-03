-- Create tables for Telegram Finance Bot

-- Users table
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

-- Create indexes for users
CREATE INDEX idx_users_telegram_id ON users(telegram_user_id);

-- Transactions table
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

-- Create indexes for transactions
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_date ON transactions(date DESC);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_category ON transactions(category);
CREATE INDEX idx_transactions_user_date ON transactions(user_id, date DESC);

-- Categories table
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
  icon TEXT,
  keywords TEXT[],
  created_at TIMESTAMP DEFAULT NOW()
);

-- Seed categories
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

-- Chat history table
CREATE TABLE chat_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  telegram_message_id BIGINT,
  role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for chat history
CREATE INDEX idx_chat_history_user_id ON chat_history(user_id, created_at DESC);

-- Auto-cleanup function for chat history
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

-- Create trigger for auto-cleanup
CREATE TRIGGER trigger_cleanup_chat_history
AFTER INSERT ON chat_history
FOR EACH ROW
EXECUTE FUNCTION cleanup_old_chat_history();

-- Budgets table (optional - future feature)
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

-- Create indexes for budgets
CREATE INDEX idx_budgets_user_id ON budgets(user_id);

-- Insights table (optional - cache AI insights)
CREATE TABLE insights (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  content TEXT NOT NULL,
  metadata JSONB,
  period_start DATE,
  period_end DATE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for insights
CREATE INDEX idx_insights_user_period ON insights(user_id, period_start, period_end);

-- Enable Row Level Security (RLS)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE budgets ENABLE ROW LEVEL SECURITY;
ALTER TABLE insights ENABLE ROW LEVEL SECURITY;

-- Create policies for RLS
-- Users can only access their own data
CREATE POLICY "Users can view own data" ON users
  FOR SELECT USING (true);

CREATE POLICY "Users can update own data" ON users
  FOR UPDATE USING (true);

-- Transactions policies
CREATE POLICY "Users can view own transactions" ON transactions
  FOR SELECT USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own transactions" ON transactions
  FOR INSERT WITH CHECK (auth.uid()::text = user_id::text);

CREATE POLICY "Users can update own transactions" ON transactions
  FOR UPDATE USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can delete own transactions" ON transactions
  FOR DELETE USING (auth.uid()::text = user_id::text);

-- Chat history policies
CREATE POLICY "Users can view own chat history" ON chat_history
  FOR SELECT USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own chat history" ON chat_history
  FOR INSERT WITH CHECK (auth.uid()::text = user_id::text);

-- Budgets policies
CREATE POLICY "Users can view own budgets" ON budgets
  FOR SELECT USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own budgets" ON budgets
  FOR INSERT WITH CHECK (auth.uid()::text = user_id::text);

CREATE POLICY "Users can update own budgets" ON budgets
  FOR UPDATE USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can delete own budgets" ON budgets
  FOR DELETE USING (auth.uid()::text = user_id::text);

-- Insights policies
CREATE POLICY "Users can view own insights" ON insights
  FOR SELECT USING (auth.uid()::text = user_id::text);

CREATE POLICY "Users can insert own insights" ON insights
  FOR INSERT WITH CHECK (auth.uid()::text = user_id::text);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON budgets
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();