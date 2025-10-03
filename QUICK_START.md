# 🚀 Quick Start Guide

## 📋 Prerequisites

- Telegram Bot token dari [@BotFather](https://t.me/botfather)
- Supabase account (gratis di [supabase.com](https://supabase.com))
- Groq API key (gratis di [groq.com](https://groq.com))

## ⚡ 5 Menit Setup

### 1. Create Telegram Bot
```
/start /newbot
Nama bot: Finance Bot Gua
Username: finance_bot_gua
```
Copy token yang diberikan.

### 2. Setup Supabase
1. Go to [supabase.com](https://supabase.com)
2. Create new project
3. Copy URL dan anon key
4. Buka SQL editor
5. Copy-paste isi `database/init.sql`
6. Run SQL

### 3. Get Groq API Key
1. Go to [groq.com](https://groq.com)
2. Sign up
3. Dashboard → API Keys
4. Create new key
5. Copy key

### 4. Configure Environment
```bash
# Copy .env template
cp .env.example .env

# Edit .env
nano .env
```

Isi dengan credentials kamu:
```bash
TELEGRAM_BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your_supabase_anon_key
GROQ_API_KEY=gsk_1234567890abcdef
```

### 5. Run Bot
```bash
# Download dependencies
go mod tidy

# Build & run
go run cmd/bot/main.go
```

## 🎱 Test Bot

Chat dengan bot kamu di Telegram:
- `/start` - Welcome message
- `bayar makan 50rb` - Add transaction
- `pengeluaran hari ini?` - Query expenses
- `/help` - Lihat semua commands

## 🚀 Deploy ke Railway

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login & deploy
railway login
railway init
railway up
```

Set environment variables di Railway dashboard, then deploy!

## ✅ Check Status

Health check endpoint:
```bash
curl http://localhost:8080/health
```

## 🔧 Commands Lengkap

| Command | Deskripsi |
|---------|-----------|
| `/start` | Mulai bot |
| `/help` | Bantuan |
| `/summary` | Summary keuangan |

## 💬 Natural Language Examples

**Input Transaksi:**
- `bayar makan 50rb di warteg`
- `gaji masuk 8jt kemarin`
- `bayar kos 2juta`
- `belanja pulsa 25rb`

**Query Data:**
- `pengeluaran minggu ini?`
- `income bulan ini berapa?`
- `total gua keluar buat makan?`
- `brp duit gua sisa?`

**Update/Delete:**
- `salah harusnya 45rb` (update last)
- `hapus transaksi terakhir`

**Insights:**
- `gimana keuangan gua bulan ini?`
- `kasi summary keuangan dong`
- `ada tips ga?`

## 🎯 Tips & Tricks

1. **Format Angka**: 50rb, 50.000, 50000, 5jt
2. **Auto-detect**: Hari ini, kemarin, minggu ini, bulan lalu
3. **Kategori**: Auto-detect dari keywords (makan, transport, dll)
4. **Context**: Bot ingat chat sebelumnya
5. **Error Handling**: Bot kasih saran kalau input tidak jelas

## 🆘 Troubleshooting

### Bot tidak merespon
- Check `TELEGRAM_BOT_TOKEN` benar
- Pastikan tidak ada error di console
- Test ping ke Telegram API

### Database error
- Verify Supabase credentials
- Pastikan SQL schema sudah di-run
- Check RLS policies

### Groq tidak bekerja
- Verify `GROQ_API_KEY`
- Check quota tidak habis
- Test network connection

## 📞 Support

Jika ada masalah:
1. Check console logs
2. Verify environment variables
3. Test dengan examples di atas
4. Cek README.md untuk detail

Enjoy using your Finance Bot! 🎉