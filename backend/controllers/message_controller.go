package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frallan97/react-go-app-backend/database"
	"github.com/frallan97/react-go-app-backend/models"
)

// MessageController handles message-related endpoints
type MessageController struct {
	dbManager *database.DBManager
}

// NewMessageController creates a new message controller
func NewMessageController(dbManager *database.DBManager) *MessageController {
	return &MessageController{
		dbManager: dbManager,
	}
}

// MessagesHandler lists or creates messages
// @Summary     List messages
// @Description Get all messages
// @Tags        messages
// @Produce     json
// @Success     200  {array}   models.Message
// @Router      /api/messages [get]
//
// @Summary     Create message
// @Description Insert a new message
// @Tags        messages
// @Accept      json
// @Produce     json
// @Param       msg  body   models.MessageInput  true  "message payload"
// @Success     201   {object}  map[string]int
// @Router      /api/messages [post]
func (mc *MessageController) MessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !mc.dbManager.IsConnected() {
			http.Error(w, `{"status":"db unavailable"}`, http.StatusServiceUnavailable)
			return
		}

		switch r.Method {
		case http.MethodGet:
			mc.handleGetMessages(w, r)
		case http.MethodPost:
			mc.handleCreateMessage(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", 405)
		}
	}
}

// handleGetMessages handles GET requests to retrieve all messages
func (mc *MessageController) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := mc.dbManager.DB.Query(`SELECT id, content, created_at FROM messages ORDER BY id`)
	if err != nil {
		http.Error(w, "db query failed", 500)
		return
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.CreatedAt); err != nil {
			http.Error(w, "scan failed", 500)
			return
		}
		msgs = append(msgs, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

// handleCreateMessage handles POST requests to create a new message
func (mc *MessageController) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var in models.MessageInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid payload", 400)
		return
	}

	var id int
	err := mc.dbManager.DB.QueryRow(
		`INSERT INTO messages(content) VALUES($1) RETURNING id`, in.Content,
	).Scan(&id)
	if err != nil {
		http.Error(w, "insert failed", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}
