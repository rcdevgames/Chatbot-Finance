# Deployment Guide

## 🚀 Quick Deploy

### 1. Docker Compose (Recommended)

```bash
# Create external network
docker network create waw_bridge

# Setup environment variables
cp .env.example .env
# Edit .env dengan credentials yang benar

# Deploy
docker-compose up -d --build

# Check logs
docker-compose logs -f app

# Check health
curl http://localhost:8080/health
```

### 2. Railway Deployment

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login to Railway
railway login

# Initialize project
railway init

# Deploy
railway up
```

Environment variables di Railway:
- `TELEGRAM_BOT_TOKEN`: Token dari BotFather
- `DATABASE_URL`: PostgreSQL connection string
- `GROQ_API_KEY`: API key dari Groq
- `PORT`: 8080 (auto-set oleh Railway)

### 3. Manual VPS Deployment

```bash
# Build aplikasi
go build -o bot cmd/bot/main.go

# Copy ke server
scp bot user@server:/path/to/deploy/

# Setup systemd service
sudo nano /etc/systemd/system/finance-bot.service
```

Content untuk systemd service:
```ini
[Unit]
Description=Telegram Finance Bot
After=network.target

[Service]
Type=simple
User=financebot
WorkingDirectory=/path/to/deploy
ExecStart=/path/to/deploy/bot
Restart=always
RestartSec=5
Environment=TELEGRAM_BOT_TOKEN=your_token
Environment=DATABASE_URL=postgres://user:password@localhost:5432/chatbot?sslmode=disable
Environment=GROQ_API_KEY=your_groq_key

[Install]
WantedBy=multi-user.target
```

```bash
# Enable dan start service
sudo systemctl enable finance-bot
sudo systemctl start finance-bot
sudo systemctl status finance-bot
```

### 4. Docker Deployment

```bash
# Build image
docker build -t finance-bot .

# Run container
docker run -d \
  --name finance-bot \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  --network waw_bridge \
  finance-bot
```

## 🔧 Setup Webhook

Setelah deployment, setup webhook:

```bash
# Set webhook (ganti URL dengan deployment URL)
curl -X POST "https://api.telegram.org/bot<TOKEN>/setWebhook" \
  -d "url=https://yourdomain.com/webhook"
```

## ✅ Health Check

Test jika bot berjalan:
```bash
curl https://yourdomain.com/health
```

Response harusnya:
```json
{
  "status": "ok",
  "message": "Telegram Finance Bot is running",
  "database": "healthy"
}
```

## 📊 Monitoring

### Logs
- Railway: Dashboard → Logs
- VPS: `journalctl -u finance-bot -f`
- Docker: `docker logs -f finance-bot`

### Metrics
- Monitor webhook response time
- Track user activity di database
- Monitor Groq API usage
- Database connection pool status
- Error rate tracking

## 🔄 CI/CD (Opsional)

GitHub Actions example:
```yaml
name: Deploy to Railway
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: railway-app/railway-action@v1
        with:
          api-key: ${{ secrets.RAILWAY_TOKEN }}
```

## 🚨 Troubleshooting

### Bot tidak merespon
1. Check environment variables
2. Verify webhook URL
3. Check PostgreSQL connection
4. Test Groq API key

### Database errors
1. Pastikan PostgreSQL running dan accessible
2. Run schema migration: check auto-migration logs
3. Verify UUID extension enabled
4. Check connection string format
5. Monitor connection pool status

### GORM errors
1. Check model definitions untuk GORM tags
2. Verify foreign key relationships
3. Check database migration logs
4. Test manual query di database

### Webhook errors
1. Pastikan HTTPS enabled
2. Check SSL certificate
3. Verify port accessibility

## 📈 Scaling

- Auto scaling di Railway
- Load balancer untuk multiple instances
- Database connection pooling
- Redis untuk caching (future)