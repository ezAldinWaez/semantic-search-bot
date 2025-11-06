package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return db, nil
}

func (db *DB) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT,
		text TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		embedding TEXT -- JSON array of floats
	);

	CREATE INDEX IF NOT EXISTS idx_chat_id ON messages(chat_id);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON messages(timestamp);
	`

	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) SaveMessage(msg Message) error {
	// Convert embedding to JSON string
	var embeddingJSON string
	if msg.Embedding != nil {
		embeddingBytes, err := json.Marshal(msg.Embedding)
		if err != nil {
			return fmt.Errorf("failed to marshal embedding: %w", err)
		}
		embeddingJSON = string(embeddingBytes)
	}

	query := `
	INSERT INTO messages (chat_id, user_id, username, text, timestamp, embedding)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, msg.ChatID, msg.UserID, msg.Username, msg.Text, msg.Timestamp, embeddingJSON)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	log.Printf("Saved message from %s in chat %d: %s", msg.Username, msg.ChatID, msg.Text[:min(50, len(msg.Text))])
	return nil
}

func (db *DB) GetMessages(chatID int64) ([]Message, error) {
	query := `
	SELECT id, chat_id, user_id, username, text, timestamp, embedding
	FROM messages
	WHERE chat_id = ?
	ORDER BY timestamp DESC
	`

	rows, err := db.conn.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var embeddingJSON sql.NullString

		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Username, &msg.Text, &msg.Timestamp, &embeddingJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse embedding JSON
		if embeddingJSON.Valid && embeddingJSON.String != "" {
			if err := json.Unmarshal([]byte(embeddingJSON.String), &msg.Embedding); err != nil {
				log.Printf("Failed to unmarshal embedding for message %d: %v", msg.ID, err)
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (db *DB) GetMessagesByIDs(ids []int64) ([]Message, error) {
	if len(ids) == 0 {
		return []Message{}, nil
	}

	// Build placeholders for IN clause
	placeholders := make([]interface{}, len(ids))
	queryPlaceholders := ""
	for i, id := range ids {
		placeholders[i] = id
		if i > 0 {
			queryPlaceholders += ","
		}
		queryPlaceholders += "?"
	}

	query := fmt.Sprintf(`
	SELECT id, chat_id, user_id, username, text, timestamp, embedding
	FROM messages
	WHERE id IN (%s)
	ORDER BY timestamp DESC
	`, queryPlaceholders)

	rows, err := db.conn.Query(query, placeholders...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages by IDs: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var embeddingJSON sql.NullString

		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Username, &msg.Text, &msg.Timestamp, &embeddingJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse embedding JSON
		if embeddingJSON.Valid && embeddingJSON.String != "" {
			if err := json.Unmarshal([]byte(embeddingJSON.String), &msg.Embedding); err != nil {
				log.Printf("Failed to unmarshal embedding for message %d: %v", msg.ID, err)
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetStats(chatID int64) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE chat_id = ?`
	var count int
	err := db.conn.QueryRow(query, chatID).Scan(&count)
	return count, err
}

func (db *DB) GetStatsWithEmbeddings(chatID int64) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE chat_id = ? AND embedding IS NOT NULL AND embedding != ''`
	var count int
	err := db.conn.QueryRow(query, chatID).Scan(&count)
	return count, err
}

func (db *DB) GetMessagesWithEmbeddings(chatID int64) ([]Message, error) {
	query := `
	SELECT id, chat_id, user_id, username, text, timestamp, embedding
	FROM messages
	WHERE chat_id = ? AND embedding IS NOT NULL AND embedding != ''
	ORDER BY timestamp DESC
	`

	rows, err := db.conn.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages with embeddings: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var embeddingJSON string

		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Username, &msg.Text, &msg.Timestamp, &embeddingJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse embedding JSON
		if err := json.Unmarshal([]byte(embeddingJSON), &msg.Embedding); err != nil {
			log.Printf("Failed to unmarshal embedding for message %d: %v", msg.ID, err)
			continue // Skip messages with invalid embeddings
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
