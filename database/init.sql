-- Database Schema for Telegram Finance Bot
-- Run this script in your Supabase SQL editor

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: users
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  telegram_user_id BIGINT NOT NULL UNIQUE,
  telegram_username TEXT,
  first_name TEXT,
  last_name TEXT,
  language_code TEXT DEFAULT 'id',
  timezone TEXT DEFAULT 'Asia/Jakarta',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_telegram_id ON users(telegram_user_id);

-- Table: categories
CREATE TABLE categories (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
  icon TEXT,
  keywords TEXT[],
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Seed categories
INSERT INTO categories (name, type, icon, keywords) VALUES
  ('Gaji', 'income', '=°', ARRAY['gaji', 'salary', 'income']),
  ('Investasi', 'income', '=È', ARRAY['investasi', 'dividend', 'profit', 'keuntungan']),
  ('Freelance', 'income', '=¼', ARRAY['freelance', 'project', 'client', 'proyek']),
  ('Bonus', 'income', '<', ARRAY['bonus', 'thr', 'insentif']),

  ('Makanan', 'expense', '<T', ARRAY['makan', 'food', 'resto', 'warteg', 'lunch', 'dinner']),
  ('Transport', 'expense', '=—', ARRAY['transport', 'bensin', 'grab', 'gojek', 'parkir', 'ojek']),
  ('Tagihan', 'expense', '=Ä', ARRAY['tagihan', 'listrik', 'air', 'internet', 'wifi', 'token']),
  ('Kos/Sewa', 'expense', '<à', ARRAY['kos', 'sewa', 'rent', 'kontrakan']),
  ('Hiburan', 'expense', '<®', ARRAY['hiburan', 'nonton', 'game', 'netflix', 'spotify']),
  ('Belanja', 'expense', '=Ò', ARRAY['belanja', 'shopping', 'beli', 'tokped', 'shopee']),
  ('Kesehatan', 'expense', '=Š', ARRAY['kesehatan', 'obat', 'dokter', 'hospital', 'rumah sakit']),
  ('Pendidikan', 'expense', '=Ú', ARRAY['kursus', 'buku', 'course', 'training']),
  ('Lainnya', 'expense', '=æ', ARRAY['lainnya', 'other', 'misc']);

-- Table: transactions
CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
  amount DECIMAL(12, 2) NOT NULL,
  category TEXT NOT NULL,
  description TEXT,
  date DATE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_date ON transactions(date DESC);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_category ON transactions(category);
CREATE INDEX idx_transactions_user_date ON transactions(user_id, date DESC);

-- Table: chat_history
CREATE TABLE chat_history (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  telegram_message_id BIGINT,
  role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_chat_history_user_id ON chat_history(user_id, created_at DESC);

-- Function to cleanup old chat history (keep last 50 per user)
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

-- Trigger to auto-cleanup chat history
CREATE TRIGGER trigger_cleanup_chat_history
AFTER INSERT ON chat_history
FOR EACH ROW
EXECUTE FUNCTION cleanup_old_chat_history();

-- Function to get transaction summary
CREATE OR REPLACE FUNCTION get_transaction_summary(
  p_user_id UUID,
  p_type TEXT DEFAULT NULL,
  p_start_date DATE DEFAULT NULL,
  p_end_date DATE DEFAULT NULL
)
RETURNS TABLE(total DECIMAL) AS $$
BEGIN
  RETURN QUERY
  SELECT COALESCE(SUM(t.amount), 0) as total
  FROM transactions t
  WHERE t.user_id = p_user_id
    AND (p_type IS NULL OR t.type = p_type)
    AND (p_start_date IS NULL OR t.date >= p_start_date)
    AND (p_end_date IS NULL OR t.date <= p_end_date);
END;
$$ LANGUAGE plpgsql;

-- Function to get category breakdown
CREATE OR REPLACE FUNCTION get_category_breakdown(
  p_user_id UUID,
  p_type TEXT DEFAULT NULL,
  p_start_date DATE DEFAULT NULL,
  p_end_date DATE DEFAULT NULL
)
RETURNS TABLE(
  category TEXT,
  total DECIMAL,
  percentage DECIMAL,
  transaction_count BIGINT
) AS $$
DECLARE
  total_amount DECIMAL;
BEGIN
  -- Get total amount for percentage calculation
  SELECT COALESCE(SUM(amount), 0) INTO total_amount
  FROM transactions
  WHERE user_id = p_user_id
    AND (p_type IS NULL OR type = p_type)
    AND (p_start_date IS NULL OR date >= p_start_date)
    AND (p_end_date IS NULL OR date <= p_end_date);

  -- Return category breakdown
  RETURN QUERY
  SELECT
    t.category,
    SUM(t.amount) as total,
    CASE
      WHEN total_amount > 0 THEN (SUM(t.amount) / total_amount) * 100
      ELSE 0
    END as percentage,
    COUNT(*) as transaction_count
  FROM transactions t
  WHERE t.user_id = p_user_id
    AND (p_type IS NULL OR t.type = p_type)
    AND (p_start_date IS NULL OR t.date >= p_start_date)
    AND (p_end_date IS NULL OR t.date <= p_end_date)
  GROUP BY t.category
  ORDER BY total DESC;
END;
$$ LANGUAGE plpgsql;

-- Table: budgets (optional for future features)
CREATE TABLE budgets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  category TEXT REFERENCES categories(name),
  amount DECIMAL(12, 2) NOT NULL,
  period TEXT CHECK (period IN ('daily', 'weekly', 'monthly')),
  start_date DATE NOT NULL,
  end_date DATE,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_budgets_user_id ON budgets(user_id);

-- Table: insights (cache for AI insights)
CREATE TABLE insights (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  content TEXT NOT NULL,
  metadata JSONB,
  period_start DATE,
  period_end DATE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_insights_user_period ON insights(user_id, period_start, period_end);

-- Row Level Security (RLS) policies
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE budgets ENABLE ROW LEVEL SECURITY;
ALTER TABLE insights ENABLE ROW LEVEL SECURITY;

-- Users can only access their own data
CREATE POLICY "Users can view own data" ON users
  FOR ALL USING (telegram_user_id = current_setting('app.current_user_id')::BIGINT);

CREATE POLICY "Users can view own transactions" ON transactions
  FOR ALL USING (user_id IN (
    SELECT id FROM users WHERE telegram_user_id = current_setting('app.current_user_id')::BIGINT
  ));

CREATE POLICY "Users can view own chat history" ON chat_history
  FOR ALL USING (user_id IN (
    SELECT id FROM users WHERE telegram_user_id = current_setting('app.current_user_id')::BIGINT
  ));

CREATE POLICY "Users can view own budgets" ON budgets
  FOR ALL USING (user_id IN (
    SELECT id FROM users WHERE telegram_user_id = current_setting('app.current_user_id')::BIGINT
  ));

CREATE POLICY "Users can view own insights" ON insights
  FOR ALL USING (user_id IN (
    SELECT id FROM users WHERE telegram_user_id = current_setting('app.current_user_id')::BIGINT
  ));

-- Categories are public read-only
CREATE POLICY "Categories are public read" ON categories
  FOR SELECT USING (true);

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at
  BEFORE UPDATE ON transactions
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_budgets_updated_at
  BEFORE UPDATE ON budgets
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();