# Semantic Search Bot 🤖🔍

An intelligent Telegram bot that enables semantic search through chat history using AI embeddings. Find conversations by meaning and context, not just keywords!

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Production%20Ready-brightgreen.svg)]()

## ✨ Features

🧠 **Semantic Understanding**: Searches by meaning and context, not exact word matches  
⚡ **Lightning Fast**: Sub-2-second search performance with real-time monitoring  
🔒 **Privacy First**: All data stored locally, no external data sharing  
🚀 **Async Processing**: Non-blocking embedding generation in background  
📊 **Performance Monitoring**: Built-in metrics and optimization tracking  
🧪 **Production Ready**: Comprehensive testing and error handling

## 🎯 How It Works

1. **Message Tracking**: Bot monitors all messages in chats where it's added
2. **AI Embeddings**: Generates semantic vectors using Ollama (local AI)
3. **Smart Storage**: Stores messages with embeddings in efficient SQLite database
4. **Semantic Search**: Uses cosine similarity to find contextually relevant results
5. **Intelligent Ranking**: Results ranked by semantic relevance with similarity scores

## 🚀 Quick Start

### Prerequisites

-   **Go 1.21+**
-   **Ollama** (local AI service)
-   **Telegram Bot Token** from [@BotFather](https://t.me/botfather)

### Installation

```bash
# 1. Clone and setup project
git clone <repository>
cd semantic-search-bot
make setup

# 2. Install and start Ollama
# Visit: https://ollama.ai/download
ollama serve
ollama pull all-minilm

# 3. Configure bot token
# Edit .env file:
TELEGRAM_TOKEN=your_bot_token_here

# 4. Run the bot
make run
```

### First Use

1. **Add bot to chat**: Find your bot and add to group or start private chat
2. **Initialize**: Send `/start` to see welcome message
3. **Test setup**: Use `/test` to verify AI service connection
4. **Chat normally**: Bot will index messages automatically
5. **Search**: Use `/search your query` to find relevant conversations

## 🔍 Search Examples

```bash
# Business & Work
/search team meeting           # Finds meeting discussions
/search project deadline       # Finds work planning conversations
/search client presentation    # Finds business-related talks

# Technical Content
/search python programming     # Finds coding discussions
/search API bug fix           # Finds technical troubleshooting
/search database optimization  # Finds technical solutions

# Daily Life
/search lunch plans           # Finds food and social arrangements
/search weekend trip          # Finds travel and personal plans
/search funny story           # Finds humorous conversations
```

## 📋 Commands

| Command           | Description                                         |
| ----------------- | --------------------------------------------------- |
| `/start`          | Welcome message and bot introduction                |
| `/help`           | Detailed help with examples and tips                |
| `/stats`          | Message count, embedding progress, search readiness |
| `/test`           | Verify AI embedding service connectivity            |
| `/perf`           | Performance metrics and system status               |
| `/search <query>` | **Semantic search through chat history**            |

## 🛠️ Development

### Available Make Commands

```bash
make run         # Run the bot (auto-checks Ollama)
make build       # Build production binary
make setup       # Initial project setup
make test        # Run unit tests
make fmt         # Format code
make clean       # Clean build artifacts and database
make clean-db    # Clean only database (fresh start)
```

### Project Structure

```
semantic-search-bot/
├── main.go                # Application entry point
├── config/                # Configuration management
│   ├── config.go          # Environment and .env handling
│   └── config_test.go     # Configuration tests
├── bot/                   # Telegram bot logic
│   ├── bot.go             # Bot initialization and lifecycle
│   ├── handlers.go        # Message and command handlers
│   └── performance.go     # Performance monitoring
├── database/              # Data persistence
│   ├── models.go          # Data models and structures
│   └── sqlite.go          # SQLite operations
├── embedding/             # AI embedding service
│   └── client.go          # Ollama API client
├── search/                # Semantic search engine
│   ├── engine.go          # Core search algorithms
│   └── engine_test.go     # Search engine tests
├── .env.example           # Environment variables template
├── Makefile               # Development automation
└── README.md              # This documentation
```

## ⚙️ Configuration

Configure via `.env` file (created automatically by `make setup`):

```bash
# Required
TELEGRAM_TOKEN=your_bot_token_here

# Optional (with defaults)
DATABASE_PATH=./messages.db
EMBEDDING_API_URL=http://localhost:11434
EMBEDDING_MODEL=all-minilm:latest
```

## 🧪 Testing

### Automated Tests

```bash
make test                  # Run all unit tests
go test -v ./search        # Test search algorithms specifically
go test -v ./config        # Test configuration handling
```

### Manual Testing

1. **Load sample data**: Start chatting or use test conversations
2. **Test search accuracy**: Try different query types and topics
3. **Performance testing**: Use `/perf` to monitor search speed
4. **Stress testing**: Test with increasing message volumes

### Success Metrics

-   **Search Relevance**: >70% of searches return relevant results in top 3
-   **Performance**: <2 seconds average search time
-   **Reliability**: Graceful handling of AI service outages
-   **Usability**: Clear, helpful responses for all user interactions

## 📊 Performance & Monitoring

### Built-in Monitoring

-   **Real-time metrics**: `/perf` command shows current performance
-   **Automatic logging**: Performance stats logged every 5 minutes
-   **Memory tracking**: Automatic garbage collection and usage monitoring
-   **Search optimization**: Query timing and similarity score tracking

### Performance Targets

| Metric               | Target                   | Status       |
| -------------------- | ------------------------ | ------------ |
| Search Speed         | < 2 seconds              | 🟢 Optimized |
| Memory Usage         | < 100MB for 1K messages  | 🟢 Efficient |
| Embedding Generation | Background, non-blocking | 🟢 Async     |
| Database Performance | Indexed queries          | 🟢 Fast      |

## 🔧 Troubleshooting

### Common Issues

**Bot not responding**

```bash
# Check token configuration
grep TELEGRAM_TOKEN .env
# Verify bot permissions in chat
```

**Search returns no results**

```bash
# Check embedding status
/stats
# Verify Ollama is running
curl http://localhost:11434/api/tags
```

**Slow performance**

```bash
# Check performance metrics
/perf
# Monitor system resources
top -p $(pgrep semantic-search-bot)
```

**Ollama connection issues**

```bash
# Start Ollama service
ollama serve
# Pull required model
ollama pull all-minilm
# Test connection
/test
```

### Debug Mode

Enable detailed logging by setting environment variable:

```bash
export BOT_DEBUG=true
make run
```

## 🏗️ Architecture

### Design Principles

-   **Modularity**: Clean separation between bot, database, embedding, and search layers
-   **Performance**: Asynchronous processing and optimized data structures
-   **Reliability**: Comprehensive error handling and graceful degradation
-   **Maintainability**: Well-documented code with comprehensive testing
-   **Privacy**: Local-first architecture with no external data dependencies

### Technology Stack

-   **Language**: Go 1.21+ (performance, concurrency, single binary deployment)
-   **AI Embeddings**: Ollama with all-minilm model (local, privacy-preserving)
-   **Database**: SQLite with JSON embedding storage (simple, efficient)
-   **Bot Framework**: go-telegram-bot-api (stable, well-maintained)
-   **Configuration**: godotenv for .env file support

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make changes and add tests
4. Run tests: `make test`
5. Format code: `make fmt`
6. Commit changes: `git commit -am 'Add feature'`
7. Push branch: `git push origin feature-name`
8. Submit pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

-   **Anthropic's Claude Sonnet 4** for AI-assisted development and code generation throughout this project
-   **Ollama** for providing excellent local AI embedding capabilities
-   **Go Telegram Bot API** for robust Telegram integration
-   **SQLite** for reliable, embedded database functionality

## 🔮 Future Enhancements

-   **Advanced Filters**: Search by user, date range, message type
-   **Export Features**: Save search results to files
-   **Multi-language Support**: Enhanced support for non-English content
-   **Web Dashboard**: Optional web interface for search analytics
-   **Advanced Analytics**: Search pattern analysis and insights
-   **Conversation Threading**: Group related messages in search results

---

**Built with ❤️ in Go | Powered by Local AI | Privacy-First Design**  
**Developed with assistance from Anthropic's Claude Sonnet 4**

Ready to make your chat history searchable and intelligent! 🚀
