package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frallan97/hackaton-demo-backend/database"
	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/frallan97/hackaton-demo-backend/utils"
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

// CreateMessageResponse represents the response for creating a message
type CreateMessageResponse struct {
	ID int `json:"id"`
}

// MessagesHandler lists or creates messages
// @Summary     List messages
// @Description Get all messages
// @Tags        messages
// @Produce     json
// @Success     200  {object}  utils.APIResponse{data=[]models.Message}
// @Failure     503  {object}  utils.APIResponse
// @Failure     500  {object}  utils.APIResponse
// @Router      /api/messages [get]
//
// @Summary     Create message
// @Description Insert a new message
// @Tags        messages
// @Accept      json
// @Produce     json
// @Param       msg  body   models.MessageInput  true  "message payload"
// @Success     201   {object}  utils.APIResponse{data=CreateMessageResponse}
// @Failure     400   {object}  utils.APIResponse
// @Failure     500   {object}  utils.APIResponse
// @Router      /api/messages [post]
func (mc *MessageController) MessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !mc.dbManager.IsConnected() {
			utils.WriteError(w, http.StatusServiceUnavailable, "Database connection unavailable", nil)
			return
		}

		switch r.Method {
		case http.MethodGet:
			mc.handleGetMessages(w, r)
		case http.MethodPost:
			mc.handleCreateMessage(w, r)
		default:
			utils.WriteMethodNotAllowed(w, "GET, POST")
		}
	}
}

// handleGetMessages handles GET requests to retrieve all messages
func (mc *MessageController) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := mc.dbManager.DB.Query(`SELECT id, content, created_at FROM messages ORDER BY id`)
	if err != nil {
		utils.WriteInternalServerError(w, "Failed to query messages", err)
		return
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.Content, &m.CreatedAt); err != nil {
			utils.WriteInternalServerError(w, "Failed to scan message data", err)
			return
		}
		msgs = append(msgs, m)
	}

	utils.WriteOK(w, msgs, "Messages retrieved successfully")
}

// handleCreateMessage handles POST requests to create a new message
func (mc *MessageController) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var in models.MessageInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		utils.WriteBadRequest(w, "Invalid request payload", err)
		return
	}

	// Validate input
	if in.Content == "" {
		utils.WriteValidationError(w, map[string]string{
			"content": "Message content is required",
		})
		return
	}

	var id int
	err := mc.dbManager.DB.QueryRow(
		`INSERT INTO messages(content) VALUES($1) RETURNING id`, in.Content,
	).Scan(&id)
	if err != nil {
		utils.WriteInternalServerError(w, "Failed to create message", err)
		return
	}

	response := &CreateMessageResponse{ID: id}
	utils.WriteCreated(w, response, "Message created successfully")
}
