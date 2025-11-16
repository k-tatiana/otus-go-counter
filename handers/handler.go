package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	guuid "github.com/google/uuid"

	"go-server-counters/models"
)

type CounterService interface {
	IncrementTotalCounter(ctx context.Context, key string) error
	ReadMessages(ctx context.Context, key string) error
	GetMessageCounter(ctx context.Context, key string) (*models.MessageCounter, error)
}

type Handler struct {
	counterService CounterService
}

func NewHandler(counterService CounterService) *Handler {
	return &Handler{
		counterService: counterService,
	}
}

func validateBody(msg *models.MessageCounter) error {
	// Validate required fields
	if msg.FromUserID == "" || msg.ToUserID == "" {
		return fmt.Errorf("fromUserID and ToUserID are required")
	}
	// Validate UUID format
	if err := guuid.Validate(msg.FromUserID); err != nil {
		return fmt.Errorf("invalid FromUserID format")
	}
	if err := guuid.Validate(msg.ToUserID); err != nil {
		return fmt.Errorf("invalid ToUserID format")
	}

	return nil
}

// POST /increment
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var msg models.MessageCounter
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate body
	if err := validateBody(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a consistent key for the user pair
	md5Hash := md5.Sum([]byte(msg.FromUserID + msg.ToUserID))
	md5HashStr := fmt.Sprintf("%x", md5Hash)

	if err := h.counterService.IncrementTotalCounter(ctx, md5HashStr); err != nil {
		http.Error(w, "Failed to increment counter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("{\"success\": true}"))
}

// PATCH /read-messages
func (h *Handler) ReadMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var msg models.MessageCounter
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate body
	if err := validateBody(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a consistent key for the user pair
	md5Hash := md5.Sum([]byte(msg.FromUserID + msg.ToUserID))
	md5HashStr := fmt.Sprintf("%x", md5Hash)

	if err := h.counterService.ReadMessages(ctx, md5HashStr); err != nil {
		http.Error(w, "Failed to read messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"success\": true}"))
}

// GET /unread-count
func (h *Handler) GetMessageCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Here you would typically extract user IDs from query parameters or headers
	// For simplicity, let's assume they are provided as query parameters
	fromUserID := r.URL.Query().Get("from_user_id")
	toUserID := r.URL.Query().Get("to_user_id")

	if fromUserID == "" || toUserID == "" {
		http.Error(w, "from_user_id and to_user_id are required", http.StatusBadRequest)
		return
	}

	// Create a consistent key for the user pair
	md5Hash := md5.Sum([]byte(fromUserID + toUserID))
	md5HashStr := fmt.Sprintf("%x", md5Hash)

	// Here you would typically call a method to get the unread count
	// For simplicity, let's just return a placeholder response
	stat, err := h.counterService.GetMessageCounter(ctx, md5HashStr)
	if err != nil {
		http.Error(w, "Failed to get unread count", http.StatusInternalServerError)
		return
	}

	stat.FromUserID = fromUserID
	stat.ToUserID = toUserID

	responseData, err := json.Marshal(stat)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	// Placeholder response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
	w.WriteHeader(http.StatusOK)
}
