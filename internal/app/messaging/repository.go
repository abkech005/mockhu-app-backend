package messaging

import "context"

// ConversationRepository defines data access methods for conversations
type ConversationRepository interface {
	// CreateOrGetConversation creates a new conversation or returns existing one
	// Uses user1_id < user2_id ordering to prevent duplicates
	CreateOrGetConversation(ctx context.Context, user1ID, user2ID string) (*Conversation, error)

	// GetConversationByID retrieves a conversation by ID
	GetConversationByID(ctx context.Context, conversationID string) (*Conversation, error)

	// GetUserConversations retrieves all conversations for a user with pagination
	// Returns conversations ordered by updated_at DESC
	GetUserConversations(ctx context.Context, userID string, page, limit int) ([]Conversation, int, error)

	// UpdateLastMessage updates the denormalized last message fields
	UpdateLastMessage(ctx context.Context, conversationID string, messageID, messageText, senderID string) error

	// DeleteConversation soft deletes a conversation for a user
	DeleteConversation(ctx context.Context, conversationID, userID string) error

	// GetConversationByParticipants finds conversation between two users
	GetConversationByParticipants(ctx context.Context, user1ID, user2ID string) (*Conversation, error)
}

// MessageRepository defines data access methods for messages
type MessageRepository interface {
	// CreateMessage creates a new message in a conversation
	CreateMessage(ctx context.Context, message *Message) error

	// GetMessageByID retrieves a message by ID
	GetMessageByID(ctx context.Context, messageID string) (*Message, error)

	// GetConversationMessages retrieves messages for a conversation with pagination
	// Returns messages ordered by created_at DESC (newest first)
	GetConversationMessages(ctx context.Context, conversationID string, page, limit int) ([]Message, int, error)

	// UpdateMessage updates a message (for editing)
	UpdateMessage(ctx context.Context, messageID string, content string) error

	// DeleteMessage soft deletes a message
	DeleteMessage(ctx context.Context, messageID, userID string) error

	// MarkMessageAsRead marks a single message as read
	MarkMessageAsRead(ctx context.Context, messageID, userID string) error

	// MarkConversationAsRead marks all unread messages in conversation as read
	MarkConversationAsRead(ctx context.Context, conversationID, userID string) error

	// GetUnreadCount returns total unread messages count for a user
	GetUnreadCount(ctx context.Context, userID string) (int, error)

	// GetConversationUnreadCount returns unread count for specific conversation
	GetConversationUnreadCount(ctx context.Context, conversationID, userID string) (int, error)

	// GetUnreadConversationsCount returns count of conversations with unread messages
	GetUnreadConversationsCount(ctx context.Context, userID string) (int, error)
}

// BlockRepository defines data access methods for user blocking
type BlockRepository interface {
	// BlockUser creates a block relationship
	BlockUser(ctx context.Context, blockerID, blockedID, reason string) error

	// UnblockUser removes a block relationship
	UnblockUser(ctx context.Context, blockerID, blockedID string) error

	// IsBlocked checks if user1 has blocked user2 OR user2 has blocked user1
	IsBlocked(ctx context.Context, user1ID, user2ID string) (bool, error)

	// IsUserBlocked checks if blockerID has blocked blockedID (one direction only)
	IsUserBlocked(ctx context.Context, blockerID, blockedID string) (bool, error)

	// GetBlockedUsers retrieves all users blocked by a user
	GetBlockedUsers(ctx context.Context, blockerID string) ([]BlockedUser, error)

	// GetBlockedUserIDs retrieves just the IDs of blocked users
	GetBlockedUserIDs(ctx context.Context, blockerID string) ([]string, error)
}

