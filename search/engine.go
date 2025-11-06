package search

import (
	"fmt"
	"math"
	"semantic-search-bot/database"
	"semantic-search-bot/embedding"
	"sort"
	"strings"
)

type Engine struct {
	db         *database.DB
	embedding  *embedding.Client
	maxResults int
}

type SearchResult struct {
	Message    database.Message
	Similarity float64
	Rank       int
}

func NewEngine(db *database.DB, embeddingClient *embedding.Client, maxResults int) *Engine {
	return &Engine{
		db:         db,
		embedding:  embeddingClient,
		maxResults: maxResults,
	}
}

func (e *Engine) Search(query string, chatID int64) ([]SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Generate embedding for the search query
	queryEmbedding, err := e.embedding.GetEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Get all messages with embeddings from the chat
	messages, err := e.db.GetMessagesWithEmbeddings(chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	if len(messages) == 0 {
		return []SearchResult{}, nil
	}

	// Calculate similarities and create results
	var results []SearchResult
	for _, msg := range messages {
		if len(msg.Embedding) == 0 {
			continue // Skip messages without embeddings
		}

		similarity := cosineSimilarity(queryEmbedding, msg.Embedding)
		if similarity > 0.1 { // Filter out very low similarities
			results = append(results, SearchResult{
				Message:    msg,
				Similarity: similarity,
			})
		}
	}

	// Sort by similarity (highest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Limit results and add ranking
	if len(results) > e.maxResults {
		results = results[:e.maxResults]
	}

	for i := range results {
		results[i].Rank = i + 1
	}

	return results, nil
}

func (e *Engine) SearchStats(chatID int64) (int, int, error) {
	totalMessages, err := e.db.GetStats(chatID)
	if err != nil {
		return 0, 0, err
	}

	messagesWithEmbeddings, err := e.db.GetStatsWithEmbeddings(chatID)
	if err != nil {
		return totalMessages, 0, err
	}

	return totalMessages, messagesWithEmbeddings, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// GetSimilarMessages finds messages similar to a given message
func (e *Engine) GetSimilarMessages(messageID int64, chatID int64) ([]SearchResult, error) {
	// Get the source message
	sourceMessages, err := e.db.GetMessagesByIDs([]int64{messageID})
	if err != nil || len(sourceMessages) == 0 {
		return nil, fmt.Errorf("source message not found")
	}

	sourceMsg := sourceMessages[0]
	if len(sourceMsg.Embedding) == 0 {
		return nil, fmt.Errorf("source message has no embedding")
	}

	// Get all other messages with embeddings
	messages, err := e.db.GetMessagesWithEmbeddings(chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve messages: %w", err)
	}

	var results []SearchResult
	for _, msg := range messages {
		if msg.ID == messageID || len(msg.Embedding) == 0 {
			continue // Skip the source message and messages without embeddings
		}

		similarity := cosineSimilarity(sourceMsg.Embedding, msg.Embedding)
		if similarity > 0.3 { // Higher threshold for similar messages
			results = append(results, SearchResult{
				Message:    msg,
				Similarity: similarity,
			})
		}
	}

	// Sort by similarity
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Limit results
	maxSimilar := 3
	if len(results) > maxSimilar {
		results = results[:maxSimilar]
	}

	for i := range results {
		results[i].Rank = i + 1
	}

	return results, nil
}
