package bot

import (
	"fmt"
	"log"
	"semantic-search-bot/database"
	"semantic-search-bot/search"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	// Handle regular messages
	if update.Message != nil {
		b.handleMessage(update.Message)
		return
	}

	// Handle edited messages (optional - treat as new messages)
	if update.EditedMessage != nil {
		b.handleMessage(update.EditedMessage)
		return
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// Skip empty messages
	if message.Text == "" {
		return
	}

	// Handle commands
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	// Store regular messages
	b.storeMessage(message)
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	command := message.Command()
	args := message.CommandArguments()

	switch command {
	case "start":
		b.handleStartCommand(message)
	case "help":
		b.handleHelpCommand(message)
	case "stats":
		b.handleStatsCommand(message)
	case "test":
		b.handleTestCommand(message)
	case "perf":
		b.handlePerfCommand(message)
	case "search":
		b.handleSearchCommand(message, args)
	default:
		b.sendReply(message, fmt.Sprintf("Unknown command: /%s", command))
	}

	log.Printf("Command /%s executed by %s in chat %d", command, message.From.UserName, message.Chat.ID)
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	welcomeText := `🤖 *Welcome to Semantic Search Bot!*

I'm your AI-powered chat search assistant! I understand conversations by *meaning*, not just keywords.

🧠 *What makes me special?*
• I learn from every message in this chat
• I understand context and intent behind your words
• I find relevant conversations even with different wording

⚡ *Quick Start:*
1️⃣ Just chat normally - I'm already learning!
2️⃣ When you need to find something: ` + "`/search your question`" + `
3️⃣ I'll show you the most relevant conversations

🔍 *Try these searches:*
• ` + "`/search meeting plans`" + ` - finds scheduling discussions
• ` + "`/search technical issue`" + ` - finds troubleshooting talks  
• ` + "`/search funny moment`" + ` - finds humorous conversations

*Ready to make your chat history searchable!* 🚀

Use /help for detailed instructions or /search to start exploring!`

	b.sendReply(message, welcomeText)
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) {
	helpText := `🔍 *How to Use Semantic Search*

I'm an AI that understands the *meaning* behind your words, not just exact matches!

🎯 *Search Examples:*

*📅 Find Planning & Meetings:*
• ` + "`/search team meeting`" + ` → finds scheduling, agenda discussions
• ` + "`/search deadline project`" + ` → finds work planning conversations
• ` + "`/search client call`" + ` → finds business communications

*💻 Find Technical Discussions:*
• ` + "`/search bug fix`" + ` → finds troubleshooting conversations  
• ` + "`/search code review`" + ` → finds development discussions
• ` + "`/search API problem`" + ` → finds technical issues

*🎉 Find Social & Fun:*
• ` + "`/search lunch plans`" + ` → finds food and social arrangements
• ` + "`/search funny story`" + ` → finds humorous moments
• ` + "`/search weekend trip`" + ` → finds travel discussions

💡 *Pro Tips:*
✅ Use natural language - "when is the meeting" works great!
✅ Try different phrasings if first search doesn't work
✅ I get smarter as more messages are added to chat
✅ Check /stats to see how many messages I've learned from

🛠️ *Available Commands:*
• ` + "`/search <your question>`" + ` - Find relevant conversations
• ` + "`/stats`" + ` - See my learning progress  
• ` + "`/test`" + ` - Check if my AI brain is working
• ` + "`/perf`" + ` - View performance metrics

*Happy searching!* 🚀`

	b.sendReply(message, helpText)
}

func (b *Bot) handleStatsCommand(message *tgbotapi.Message) {
	count, err := b.db.GetStats(message.Chat.ID)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		b.sendReply(message, "❌ Oops! I couldn't retrieve the statistics right now. Please try again.")
		return
	}

	// Count messages with embeddings
	countWithEmbeddings, err := b.db.GetStatsWithEmbeddings(message.Chat.ID)
	if err != nil {
		log.Printf("Error getting embedding stats: %v", err)
		countWithEmbeddings = 0
	}

	readinessPercent := float64(countWithEmbeddings) / float64(max(count, 1)) * 100

	var statusEmoji string
	var statusText string

	if readinessPercent >= 80 {
		statusEmoji = "🟢"
		statusText = "Excellent - Ready for great search results!"
	} else if readinessPercent >= 50 {
		statusEmoji = "🟡"
		statusText = "Good - Search quality improving as I learn"
	} else if readinessPercent >= 10 {
		statusEmoji = "🟠"
		statusText = "Getting started - Keep chatting for better results"
	} else {
		statusEmoji = "🔴"
		statusText = "Just beginning - I need more messages to learn from"
	}

	statsText := fmt.Sprintf(`📊 *My Learning Progress*

💬 *Messages Collected:* %d
🧠 *Messages I've Learned From:* %d
📈 *Search Readiness:* %.1f%%

%s *Status:* %s

🔍 *Search Quality:*
%s

*What's Next?*
• Keep chatting naturally - I learn from every message!
• Try `+"`/search`"+` to find conversations by meaning
• Use `+"`/test`"+` to check my AI connection

*Model:* %s | *Chat ID:* %d`,
		count,
		countWithEmbeddings,
		readinessPercent,
		statusEmoji,
		statusText,
		getSearchQualityTips(countWithEmbeddings),
		b.config.EmbeddingModel,
		message.Chat.ID)

	b.sendReply(message, statsText)
}

func getSearchQualityTips(embeddingCount int) string {
	if embeddingCount >= 100 {
		return "🎯 Excellent search quality expected!"
	} else if embeddingCount >= 50 {
		return "👍 Good search quality - results should be relevant"
	} else if embeddingCount >= 20 {
		return "📚 Fair search quality - improving with more messages"
	} else if embeddingCount >= 5 {
		return "🌱 Basic search available - quality will improve"
	} else {
		return "⏳ Need more messages for meaningful search results"
	}
}

func (b *Bot) handleTestCommand(message *tgbotapi.Message) {
	b.sendReply(message, "🧪 *Testing My AI Brain...*")

	// Test embedding generation
	testText := "Testing AI connection for semantic understanding"
	startTime := time.Now()
	embedding, err := b.embedding.GetEmbedding(testText)
	testDuration := time.Since(startTime)

	if err != nil {
		errorMsg := fmt.Sprintf(`❌ *AI Connection Failed*

*Problem:* %s

🔧 *How to Fix:*
1️⃣ Make sure Ollama is running: `+"`ollama serve`"+`
2️⃣ Install the AI model: `+"`ollama pull %s`"+`
3️⃣ Check the service: `+"`curl %s/api/tags`"+`

💡 *Need Help?*
• Restart Ollama service and try again
• Verify model installation with `+"`ollama list`"+`
• Check if port 11434 is available

Once fixed, I'll be ready to understand your conversations!`,
			err.Error(), b.config.EmbeddingModel, b.config.EmbeddingAPIURL)

		b.sendReply(message, errorMsg)
		return
	}

	var performanceEmoji string
	var performanceText string

	if testDuration < 1*time.Second {
		performanceEmoji = "🚀"
		performanceText = "Lightning fast!"
	} else if testDuration < 3*time.Second {
		performanceEmoji = "⚡"
		performanceText = "Great speed!"
	} else if testDuration < 5*time.Second {
		performanceEmoji = "✅"
		performanceText = "Good performance"
	} else {
		performanceEmoji = "🐌"
		performanceText = "A bit slow, but working"
	}

	successMsg := fmt.Sprintf(`✅ *AI Brain Test Successful!*

🧠 *Test Results:*
• Response time: %v %s %s
• AI dimensions: %d vectors
• Model: %s
• Service: %s

🎯 *What this means:*
I can understand the meaning behind your messages and find relevant conversations when you search!

*Ready to help you explore your chat history!* 🔍`,
		testDuration, performanceEmoji, performanceText,
		len(embedding),
		b.config.EmbeddingModel,
		b.config.EmbeddingAPIURL)

	b.sendReply(message, successMsg)
}

func (b *Bot) handlePerfCommand(message *tgbotapi.Message) {
	searchAvg, embeddingAvg, memUsage := b.perf.GetStats()

	perfMsg := fmt.Sprintf(`⚡ *Performance Dashboard*

🔍 *Search Performance:*
• Average speed: %v
• Target: < 2 seconds
• Status: %s

🧠 *AI Processing:*
• Embedding speed: %v  
• Processing: Background (non-blocking)
• Status: %s

💾 *System Health:*
• Memory usage: %s
• Optimization: %s

📊 *Performance Notes:*
• Search speed depends on chat history size
• AI processing runs automatically in background  
• Memory usage scales efficiently with message count

*Everything running smoothly!* 🎯`,
		formatDuration(searchAvg),
		getPerformanceStatus(searchAvg),
		formatDuration(embeddingAvg),
		getEmbeddingStatus(embeddingAvg),
		memUsage,
		getMemoryStatus(memUsage))

	b.sendReply(message, perfMsg)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "No data yet"
	}
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d.Nanoseconds())/1000000)
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func getEmbeddingStatus(embeddingAvg time.Duration) string {
	if embeddingAvg == 0 {
		return "🟡 Waiting for messages"
	} else if embeddingAvg < 2*time.Second {
		return "🟢 Fast processing"
	} else if embeddingAvg < 5*time.Second {
		return "🟡 Normal speed"
	} else {
		return "🔴 Consider checking Ollama performance"
	}
}

func getMemoryStatus(memUsage string) string {
	// Simple heuristic based on memory string
	if strings.Contains(memUsage, "GB") {
		return "🟡 Higher usage - consider restart if issues occur"
	} else {
		return "🟢 Efficient memory usage"
	}
}

func (b *Bot) handleSearchCommand(message *tgbotapi.Message, query string) {
	if strings.TrimSpace(query) == "" {
		b.sendReply(message, `🔍 *Semantic Search Help*

*How to search:* `+"`/search <your question or keywords>`"+`

💡 *Search Ideas:*

📅 *Find Planning:*
• `+"`/search meeting next week`"+`
• `+"`/search project deadline`"+`
• `+"`/search team lunch plans`"+`

💻 *Find Technical Stuff:*
• `+"`/search bug in code`"+`
• `+"`/search API not working`"+`
• `+"`/search database issue`"+`

🎉 *Find Fun Conversations:*
• `+"`/search funny story`"+`  
• `+"`/search weekend plans`"+`
• `+"`/search restaurant recommendation`"+`

✨ *Remember:* I understand meaning, not just exact words! Try natural language like you're asking a friend.

*Ready to explore your chat history?* Just add your question after /search!`)
		return
	}

	// Show searching indicator with friendly message
	b.sendReply(message, fmt.Sprintf("🔍 *Searching for:* \"%s\"\n⏳ *Let me find the most relevant conversations...*", query))

	// Start performance timing
	startTime := time.Now()

	// Perform search
	results, err := b.search.Search(query, message.Chat.ID)

	// Record search performance
	searchDuration := time.Since(startTime)
	b.perf.RecordSearchTime(searchDuration)

	if err != nil {
		log.Printf("Search error: %v", err)
		b.sendReply(message, fmt.Sprintf(`❌ *Search Error*

Something went wrong while searching: %s

💡 *Try:*
• Checking /stats to see if I have enough messages to learn from
• Using /test to verify my AI connection  
• Rephrasing your search query

*I'm ready to help once the issue is resolved!*`, err.Error()))
		return
	}

	// Handle no results with helpful suggestions
	if len(results) == 0 {
		totalMessages, withEmbeddings, _ := b.search.SearchStats(message.Chat.ID)

		var suggestionText string
		if withEmbeddings < 10 {
			suggestionText = "I need more conversations to learn from! Keep chatting and try again soon."
		} else if withEmbeddings < 50 {
			suggestionText = "Try broader search terms or different keywords. I'm still learning from this chat!"
		} else {
			suggestionText = "Try rephrasing your search or using different keywords. Sometimes a slight change helps!"
		}

		noResultsMsg := fmt.Sprintf(`🤷‍♂️ *No Matching Conversations Found*

*Your search:* "%s"

💭 *Why this might happen:*
• This topic hasn't been discussed yet
• Try different keywords or phrasing
• I might need more messages to understand better

📊 *My Knowledge:*
• Total messages: %d
• Messages I've learned from: %d

💡 *Suggestion:* %s

*Keep chatting - I get smarter with every message!* 🧠`,
			query, totalMessages, withEmbeddings, suggestionText)

		b.sendReply(message, noResultsMsg)
		return
	}

	// Format and send results with encouraging message
	resultMsg := b.formatSearchResults(query, results, searchDuration)
	b.sendReply(message, resultMsg)

	log.Printf("Search completed: query='%s', results=%d, duration=%v, chat=%d",
		query, len(results), searchDuration, message.Chat.ID)
}

func (b *Bot) formatSearchResults(query string, results []search.SearchResult, searchDuration time.Duration) string {
	var msg strings.Builder

	// Header with performance indicator
	performanceEmoji := "⚡"
	if searchDuration > 2*time.Second {
		performanceEmoji = "🐌"
	}

	msg.WriteString(fmt.Sprintf("🎯 *Found %d relevant conversation%s*\n", len(results), pluralize(len(results))))
	msg.WriteString(fmt.Sprintf("📝 *Search:* \"%s\" | %s *Speed:* %v\n\n", query, performanceEmoji, formatDuration(searchDuration)))

	for _, result := range results {
		// Format timestamp in a more readable way
		timeStr := result.Message.Timestamp.Format("Jan 2 at 15:04")

		// Truncate long messages with smart cutoff
		text := result.Message.Text
		if len(text) > 180 {
			// Try to cut at sentence end
			cutoff := 180
			for i := 150; i < min(len(text), 180); i++ {
				if text[i] == '.' || text[i] == '!' || text[i] == '?' {
					cutoff = i + 1
					break
				}
			}
			text = text[:cutoff] + "..."
		}

		// Format similarity with emoji indicators
		similarityPercent := result.Similarity * 100
		var similarityEmoji string
		if similarityPercent >= 70 {
			similarityEmoji = "🎯"
		} else if similarityPercent >= 50 {
			similarityEmoji = "✅"
		} else {
			similarityEmoji = "📝"
		}

		msg.WriteString(fmt.Sprintf("*%d.* %s *%.0f%% match*\n",
			result.Rank, similarityEmoji, similarityPercent))
		msg.WriteString(fmt.Sprintf("👤 **%s** • 📅 %s\n",
			getDisplayName(result.Message.Username), timeStr))
		msg.WriteString(fmt.Sprintf("💬 %s\n\n", text))
	}

	// Footer with helpful tips
	msg.WriteString("💡 *Tips:* Results ranked by relevance • Try different keywords for more results")

	return msg.String()
}

func getDisplayName(username string) string {
	if username == "" {
		return "Anonymous"
	}
	return username
}

func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (b *Bot) storeMessage(message *tgbotapi.Message) {
	// Clean the message text (basic preprocessing)
	cleanText := b.cleanText(message.Text)

	// Skip very short messages
	if len(strings.TrimSpace(cleanText)) < 3 {
		return
	}

	// Create message object
	msg := database.Message{
		ChatID:    message.Chat.ID,
		UserID:    message.From.ID,
		Username:  message.From.UserName,
		Text:      cleanText,
		Timestamp: time.Unix(int64(message.Date), 0),
	}

	// Generate embedding asynchronously
	go func() {
		startTime := time.Now()

		embedding, err := b.embedding.GetEmbedding(cleanText)

		// Record embedding performance
		embeddingDuration := time.Since(startTime)
		b.perf.RecordEmbeddingTime(embeddingDuration)

		if err != nil {
			log.Printf("Failed to generate embedding for message: %v", err)
			// Save message without embedding
			if err := b.db.SaveMessage(msg); err != nil {
				log.Printf("Error saving message without embedding: %v", err)
			}
			return
		}

		// Add embedding to message
		msg.Embedding = embedding

		// Save message with embedding
		if err := b.db.SaveMessage(msg); err != nil {
			log.Printf("Error saving message with embedding: %v", err)
		} else {
			log.Printf("✅ Saved message with embedding (%d dims, %v) from %s",
				len(embedding), embeddingDuration, msg.Username)
		}
	}()
}

func (b *Bot) cleanText(text string) string {
	// Basic text cleaning
	// Remove multiple spaces
	text = strings.Join(strings.Fields(text), " ")

	// Trim whitespace
	text = strings.TrimSpace(text)

	return text
}

func (b *Bot) sendReply(message *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func getPerformanceStatus(searchAvg time.Duration) string {
	if searchAvg == 0 {
		return "🟡 No searches yet"
	} else if searchAvg < 1*time.Second {
		return "🚀 Lightning fast"
	} else if searchAvg < 2*time.Second {
		return "🟢 Excellent"
	} else if searchAvg < 5*time.Second {
		return "🟡 Good"
	} else {
		return "🔴 Could be faster"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
