package messaging

import (
	"mockhu-app-backend/internal/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all messaging-related routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	v1 := app.Group("/v1")

	// All messaging routes require authentication
	auth := middleware.AuthMiddleware()

	// ========================================================================
	// CONVERSATION ROUTES
	// ========================================================================

	// Create or get conversation
	v1.Post("/conversations", auth, handler.CreateConversation)

	// List user's conversations
	v1.Get("/conversations", auth, handler.GetConversations)

	// Get specific conversation
	v1.Get("/conversations/:conversationId", auth, handler.GetConversation)

	// Delete conversation
	v1.Delete("/conversations/:conversationId", auth, handler.DeleteConversation)

	// ========================================================================
	// MESSAGE ROUTES
	// ========================================================================

	// Send message in conversation
	v1.Post("/conversations/:conversationId/messages", auth, handler.SendMessage)

	// Get messages in conversation
	v1.Get("/conversations/:conversationId/messages", auth, handler.GetMessages)

	// Delete message
	v1.Delete("/messages/:messageId", auth, handler.DeleteMessage)

	// ========================================================================
	// UNREAD ROUTES
	// ========================================================================

	// Mark conversation as read
	v1.Post("/conversations/:conversationId/read", auth, handler.MarkConversationAsRead)

	// Mark single message as read
	v1.Post("/messages/:messageId/read", auth, handler.MarkMessageAsRead)

	// Get unread count
	v1.Get("/conversations/unread-count", auth, handler.GetUnreadCount)

	// ========================================================================
	// PRIVACY & BLOCKING ROUTES
	// ========================================================================

	// Check if can message user
	v1.Get("/users/:userId/can-message", auth, handler.CanMessage)

	// Block user
	v1.Post("/users/:userId/block", auth, handler.BlockUser)

	// Unblock user
	v1.Delete("/users/:userId/block", auth, handler.UnblockUser)

	// Get blocked users list
	v1.Get("/users/blocked", auth, handler.GetBlockedUsers)
}
