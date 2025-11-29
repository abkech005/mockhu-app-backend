package messaging

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for messaging operations
type Handler struct {
	service MessagingService
}

// NewHandler creates a new messaging handler
func NewHandler(service MessagingService) *Handler {
	return &Handler{
		service: service,
	}
}

// ============================================================================
// CONVERSATION HANDLERS
// ============================================================================

// CreateConversation handles POST /v1/conversations
func (h *Handler) CreateConversation(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Parse request body
	var req CreateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	// Validate recipient ID
	if req.RecipientID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "recipient_id is required",
		})
	}

	// Create or get conversation
	conversation, err := h.service.CreateOrGetConversation(c.Context(), currentUserID, req.RecipientID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    conversation,
	})
}

// GetConversations handles GET /v1/conversations
func (h *Handler) GetConversations(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Parse pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	unreadOnly := c.QueryBool("unread_only", false)

	// Get conversations
	response, err := h.service.GetConversations(c.Context(), currentUserID, page, limit, unreadOnly)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// GetConversation handles GET /v1/conversations/:conversationId
func (h *Handler) GetConversation(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get conversation ID from params
	conversationID := c.Params("conversationId")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "conversation_id is required",
		})
	}

	// Get conversation
	conversation, err := h.service.GetConversation(c.Context(), conversationID, currentUserID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    conversation,
	})
}

// DeleteConversation handles DELETE /v1/conversations/:conversationId
func (h *Handler) DeleteConversation(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get conversation ID from params
	conversationID := c.Params("conversationId")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "conversation_id is required",
		})
	}

	// Delete conversation
	if err := h.service.DeleteConversation(c.Context(), conversationID, currentUserID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "conversation deleted successfully",
	})
}

// ============================================================================
// MESSAGE HANDLERS
// ============================================================================

// SendMessage handles POST /v1/conversations/:conversationId/messages
func (h *Handler) SendMessage(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get conversation ID from params
	conversationID := c.Params("conversationId")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "conversation_id is required",
		})
	}

	// Parse request body
	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	// Send message
	message, err := h.service.SendMessage(c.Context(), conversationID, currentUserID, &req)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    message,
	})
}

// GetMessages handles GET /v1/conversations/:conversationId/messages
func (h *Handler) GetMessages(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get conversation ID from params
	conversationID := c.Params("conversationId")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "conversation_id is required",
		})
	}

	// Parse pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)

	// Get messages
	response, err := h.service.GetMessages(c.Context(), conversationID, currentUserID, page, limit)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// DeleteMessage handles DELETE /v1/messages/:messageId
func (h *Handler) DeleteMessage(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get message ID from params
	messageID := c.Params("messageId")
	if messageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "message_id is required",
		})
	}

	// Delete message
	if err := h.service.DeleteMessage(c.Context(), messageID, currentUserID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "message deleted successfully",
	})
}

// ============================================================================
// UNREAD HANDLERS
// ============================================================================

// MarkConversationAsRead handles POST /v1/conversations/:conversationId/read
func (h *Handler) MarkConversationAsRead(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get conversation ID from params
	conversationID := c.Params("conversationId")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "conversation_id is required",
		})
	}

	// Mark as read
	if err := h.service.MarkConversationAsRead(c.Context(), conversationID, currentUserID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "conversation marked as read",
	})
}

// MarkMessageAsRead handles POST /v1/messages/:messageId/read
func (h *Handler) MarkMessageAsRead(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get message ID from params
	messageID := c.Params("messageId")
	if messageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "message_id is required",
		})
	}

	// Mark as read
	if err := h.service.MarkMessageAsRead(c.Context(), messageID, currentUserID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "message marked as read",
	})
}

// GetUnreadCount handles GET /v1/conversations/unread-count
func (h *Handler) GetUnreadCount(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get unread count
	response, err := h.service.GetUnreadCount(c.Context(), currentUserID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ============================================================================
// PRIVACY HANDLERS
// ============================================================================

// CanMessage handles GET /v1/users/:userId/can-message
func (h *Handler) CanMessage(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get user ID from params
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user_id is required",
		})
	}

	// Check if can message
	response, err := h.service.CanMessage(c.Context(), currentUserID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// BlockUser handles POST /v1/users/:userId/block
func (h *Handler) BlockUser(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get user ID from params
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user_id is required",
		})
	}

	// Parse request body (optional reason)
	var req BlockUserRequest
	_ = c.BodyParser(&req) // Ignore error, reason is optional

	// Block user
	if err := h.service.BlockUser(c.Context(), currentUserID, userID, req.Reason); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "user blocked successfully",
	})
}

// UnblockUser handles DELETE /v1/users/:userId/block
func (h *Handler) UnblockUser(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get user ID from params
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "user_id is required",
		})
	}

	// Unblock user
	if err := h.service.UnblockUser(c.Context(), currentUserID, userID); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "user unblocked successfully",
	})
}

// GetBlockedUsers handles GET /v1/users/blocked
func (h *Handler) GetBlockedUsers(c *fiber.Ctx) error {
	// Get current user ID from context
	currentUserID, ok := c.Locals("user_id").(string)
	if !ok || currentUserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get blocked users
	response, err := h.service.GetBlockedUsers(c.Context(), currentUserID)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// handleError handles service errors and returns appropriate HTTP responses
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	// Map common errors to HTTP status codes
	switch {
	case strings.Contains(errMsg, "not found"):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   errMsg,
		})
	case strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "not a participant"):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   errMsg,
		})
	case strings.Contains(errMsg, "blocked") || strings.Contains(errMsg, "cannot message"):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   errMsg,
		})
	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "required") ||
		strings.Contains(errMsg, "too long") || strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "maximum"):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   errMsg,
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "internal server error",
		})
	}
}

// Helper function to parse int with default value
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}
