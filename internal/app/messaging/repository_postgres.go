package messaging

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresConversationRepository implements ConversationRepository using PostgreSQL
type PostgresConversationRepository struct {
	db *pgxpool.Pool
}

// NewPostgresConversationRepository creates a new PostgreSQL conversation repository
func NewPostgresConversationRepository(db *pgxpool.Pool) ConversationRepository {
	return &PostgresConversationRepository{db: db}
}

// CreateOrGetConversation creates a new conversation or returns existing one
func (r *PostgresConversationRepository) CreateOrGetConversation(ctx context.Context, user1ID, user2ID string) (*Conversation, error) {
	// Order user IDs to ensure uniqueness
	orderedUser1, orderedUser2 := OrderUserIDs(user1ID, user2ID)

	query := `
		INSERT INTO conversations (user1_id, user2_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (user1_id, user2_id) 
		DO UPDATE SET updated_at = NOW()
		RETURNING id, user1_id, user2_id, last_message_id, last_message_text, 
		          last_message_sender_id, last_message_at, created_at, updated_at
	`

	var conv Conversation
	err := r.db.QueryRow(ctx, query, orderedUser1, orderedUser2).Scan(
		&conv.ID,
		&conv.User1ID,
		&conv.User2ID,
		&conv.LastMessageID,
		&conv.LastMessageText,
		&conv.LastMessageSenderID,
		&conv.LastMessageAt,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create or get conversation: %w", err)
	}

	return &conv, nil
}

// GetConversationByID retrieves a conversation by ID
func (r *PostgresConversationRepository) GetConversationByID(ctx context.Context, conversationID string) (*Conversation, error) {
	query := `
		SELECT id, user1_id, user2_id, last_message_id, last_message_text,
		       last_message_sender_id, last_message_at, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`

	var conv Conversation
	err := r.db.QueryRow(ctx, query, conversationID).Scan(
		&conv.ID,
		&conv.User1ID,
		&conv.User2ID,
		&conv.LastMessageID,
		&conv.LastMessageText,
		&conv.LastMessageSenderID,
		&conv.LastMessageAt,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conv, nil
}

// GetUserConversations retrieves all conversations for a user with pagination
func (r *PostgresConversationRepository) GetUserConversations(ctx context.Context, userID string, page, limit int) ([]Conversation, int, error) {
	offset := (page - 1) * limit

	// Query for conversations
	query := `
		SELECT id, user1_id, user2_id, last_message_id, last_message_text,
		       last_message_sender_id, last_message_at, created_at, updated_at
		FROM conversations
		WHERE user1_id = $1 OR user2_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.LastMessageID,
			&conv.LastMessageText,
			&conv.LastMessageSenderID,
			&conv.LastMessageAt,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM conversations
		WHERE user1_id = $1 OR user2_id = $1
	`

	var total int
	err = r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	return conversations, total, nil
}

// UpdateLastMessage updates the denormalized last message fields
func (r *PostgresConversationRepository) UpdateLastMessage(ctx context.Context, conversationID string, messageID, messageText, senderID string) error {
	query := `
		UPDATE conversations
		SET last_message_id = $1,
		    last_message_text = $2,
		    last_message_sender_id = $3,
		    last_message_at = NOW(),
		    updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(ctx, query, messageID, messageText, senderID, conversationID)
	if err != nil {
		return fmt.Errorf("failed to update last message: %w", err)
	}

	return nil
}

// DeleteConversation soft deletes a conversation (future implementation)
func (r *PostgresConversationRepository) DeleteConversation(ctx context.Context, conversationID, userID string) error {
	// For MVP, we'll just verify the user is a participant
	// In the future, we can add a deleted_by_user1/deleted_by_user2 flag
	conv, err := r.GetConversationByID(ctx, conversationID)
	if err != nil {
		return err
	}

	if !conv.IsParticipant(userID) {
		return fmt.Errorf("user is not a participant in this conversation")
	}

	// For now, we don't actually delete anything
	// TODO: Implement soft delete with separate flags per user
	return nil
}

// GetConversationByParticipants finds conversation between two users
func (r *PostgresConversationRepository) GetConversationByParticipants(ctx context.Context, user1ID, user2ID string) (*Conversation, error) {
	orderedUser1, orderedUser2 := OrderUserIDs(user1ID, user2ID)

	query := `
		SELECT id, user1_id, user2_id, last_message_id, last_message_text,
		       last_message_sender_id, last_message_at, created_at, updated_at
		FROM conversations
		WHERE user1_id = $1 AND user2_id = $2
	`

	var conv Conversation
	err := r.db.QueryRow(ctx, query, orderedUser1, orderedUser2).Scan(
		&conv.ID,
		&conv.User1ID,
		&conv.User2ID,
		&conv.LastMessageID,
		&conv.LastMessageText,
		&conv.LastMessageSenderID,
		&conv.LastMessageAt,
		&conv.CreatedAt,
		&conv.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found, not an error
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conv, nil
}

// ============================================================================
// MESSAGE REPOSITORY
// ============================================================================

// PostgresMessageRepository implements MessageRepository using PostgreSQL
type PostgresMessageRepository struct {
	db *pgxpool.Pool
}

// NewPostgresMessageRepository creates a new PostgreSQL message repository
func NewPostgresMessageRepository(db *pgxpool.Pool) MessageRepository {
	return &PostgresMessageRepository{db: db}
}

// CreateMessage creates a new message in a conversation
func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, message *Message) error {
	// Convert attachments to JSON
	attachmentsJSON, err := message.AttachmentsJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal attachments: %w", err)
	}

	query := `
		INSERT INTO messages (
			conversation_id, sender_id, message_type, content, attachments,
			status, is_read, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var attachmentsParam interface{}
	if attachmentsJSON != "" {
		attachmentsParam = attachmentsJSON
	} else {
		attachmentsParam = nil
	}

	err = r.db.QueryRow(ctx, query,
		message.ConversationID,
		message.SenderID,
		message.MessageType,
		message.Content,
		attachmentsParam,
		message.Status,
		message.IsRead,
	).Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

// GetMessageByID retrieves a message by ID
func (r *PostgresMessageRepository) GetMessageByID(ctx context.Context, messageID string) (*Message, error) {
	query := `
		SELECT id, conversation_id, sender_id, message_type, content, attachments,
		       status, is_read, read_at, is_deleted, deleted_at, deleted_by,
		       created_at, updated_at
		FROM messages
		WHERE id = $1 AND is_deleted = FALSE
	`

	var msg Message
	var attachmentsJSON *string
	err := r.db.QueryRow(ctx, query, messageID).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.MessageType,
		&msg.Content,
		&attachmentsJSON,
		&msg.Status,
		&msg.IsRead,
		&msg.ReadAt,
		&msg.IsDeleted,
		&msg.DeletedAt,
		&msg.DeletedBy,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Parse attachments
	if attachmentsJSON != nil {
		if err := msg.ParseAttachments(*attachmentsJSON); err != nil {
			return nil, fmt.Errorf("failed to parse attachments: %w", err)
		}
	}

	return &msg, nil
}

// GetConversationMessages retrieves messages for a conversation with pagination
func (r *PostgresMessageRepository) GetConversationMessages(ctx context.Context, conversationID string, page, limit int) ([]Message, int, error) {
	offset := (page - 1) * limit

	// Query messages (ordered by created_at DESC for chat UX)
	query := `
		SELECT id, conversation_id, sender_id, message_type, content, attachments,
		       status, is_read, read_at, is_deleted, deleted_at, deleted_by,
		       created_at, updated_at
		FROM messages
		WHERE conversation_id = $1 AND is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, conversationID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var attachmentsJSON *string

		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.MessageType,
			&msg.Content,
			&attachmentsJSON,
			&msg.Status,
			&msg.IsRead,
			&msg.ReadAt,
			&msg.IsDeleted,
			&msg.DeletedAt,
			&msg.DeletedBy,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan message: %w", err)
		}

		// Parse attachments
		if attachmentsJSON != nil {
			if err := msg.ParseAttachments(*attachmentsJSON); err != nil {
				return nil, 0, fmt.Errorf("failed to parse attachments: %w", err)
			}
		}

		messages = append(messages, msg)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM messages
		WHERE conversation_id = $1 AND is_deleted = FALSE
	`

	var total int
	err = r.db.QueryRow(ctx, countQuery, conversationID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return messages, total, nil
}

// UpdateMessage updates a message content (for editing)
func (r *PostgresMessageRepository) UpdateMessage(ctx context.Context, messageID string, content string) error {
	query := `
		UPDATE messages
		SET content = $1, updated_at = NOW()
		WHERE id = $2 AND is_deleted = FALSE
	`

	result, err := r.db.Exec(ctx, query, content, messageID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("message not found or already deleted")
	}

	return nil
}

// DeleteMessage soft deletes a message
func (r *PostgresMessageRepository) DeleteMessage(ctx context.Context, messageID, userID string) error {
	query := `
		UPDATE messages
		SET is_deleted = TRUE,
		    deleted_at = NOW(),
		    deleted_by = $1,
		    updated_at = NOW()
		WHERE id = $2 AND sender_id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.Exec(ctx, query, userID, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("message not found or you don't have permission to delete it")
	}

	return nil
}

// MarkMessageAsRead marks a single message as read
func (r *PostgresMessageRepository) MarkMessageAsRead(ctx context.Context, messageID, userID string) error {
	// Only mark as read if the user is the recipient (not the sender)
	query := `
		UPDATE messages
		SET is_read = TRUE,
		    read_at = NOW(),
		    status = 'read',
		    updated_at = NOW()
		WHERE id = $1 
		  AND sender_id != $2
		  AND is_read = FALSE
		  AND is_deleted = FALSE
	`

	_, err := r.db.Exec(ctx, query, messageID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

// MarkConversationAsRead marks all unread messages in conversation as read
func (r *PostgresMessageRepository) MarkConversationAsRead(ctx context.Context, conversationID, userID string) error {
	query := `
		UPDATE messages
		SET is_read = TRUE,
		    read_at = NOW(),
		    status = 'read',
		    updated_at = NOW()
		WHERE conversation_id = $1
		  AND sender_id != $2
		  AND is_read = FALSE
		  AND is_deleted = FALSE
	`

	_, err := r.db.Exec(ctx, query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark conversation as read: %w", err)
	}

	return nil
}

// GetUnreadCount returns total unread messages count for a user
func (r *PostgresMessageRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages m
		JOIN conversations c ON m.conversation_id = c.id
		WHERE (c.user1_id = $1 OR c.user2_id = $1)
		  AND m.sender_id != $1
		  AND m.is_read = FALSE
		  AND m.is_deleted = FALSE
	`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// GetConversationUnreadCount returns unread count for specific conversation
func (r *PostgresMessageRepository) GetConversationUnreadCount(ctx context.Context, conversationID, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE conversation_id = $1
		  AND sender_id != $2
		  AND is_read = FALSE
		  AND is_deleted = FALSE
	`

	var count int
	err := r.db.QueryRow(ctx, query, conversationID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get conversation unread count: %w", err)
	}

	return count, nil
}

// GetUnreadConversationsCount returns count of conversations with unread messages
func (r *PostgresMessageRepository) GetUnreadConversationsCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT m.conversation_id)
		FROM messages m
		JOIN conversations c ON m.conversation_id = c.id
		WHERE (c.user1_id = $1 OR c.user2_id = $1)
		  AND m.sender_id != $1
		  AND m.is_read = FALSE
		  AND m.is_deleted = FALSE
	`

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread conversations count: %w", err)
	}

	return count, nil
}

// ============================================================================
// BLOCK REPOSITORY
// ============================================================================

// PostgresBlockRepository implements BlockRepository using PostgreSQL
type PostgresBlockRepository struct {
	db *pgxpool.Pool
}

// NewPostgresBlockRepository creates a new PostgreSQL block repository
func NewPostgresBlockRepository(db *pgxpool.Pool) BlockRepository {
	return &PostgresBlockRepository{db: db}
}

// BlockUser creates a block relationship
func (r *PostgresBlockRepository) BlockUser(ctx context.Context, blockerID, blockedID, reason string) error {
	query := `
		INSERT INTO blocked_users (blocker_id, blocked_id, reason, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (blocker_id, blocked_id) DO NOTHING
	`

	var reasonParam interface{}
	if reason != "" {
		reasonParam = reason
	} else {
		reasonParam = nil
	}

	_, err := r.db.Exec(ctx, query, blockerID, blockedID, reasonParam)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	return nil
}

// UnblockUser removes a block relationship
func (r *PostgresBlockRepository) UnblockUser(ctx context.Context, blockerID, blockedID string) error {
	query := `
		DELETE FROM blocked_users
		WHERE blocker_id = $1 AND blocked_id = $2
	`

	result, err := r.db.Exec(ctx, query, blockerID, blockedID)
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("block relationship not found")
	}

	return nil
}

// IsBlocked checks if either user has blocked the other (bidirectional check)
func (r *PostgresBlockRepository) IsBlocked(ctx context.Context, user1ID, user2ID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM blocked_users
			WHERE (blocker_id = $1 AND blocked_id = $2)
			   OR (blocker_id = $2 AND blocked_id = $1)
		)
	`

	var blocked bool
	err := r.db.QueryRow(ctx, query, user1ID, user2ID).Scan(&blocked)
	if err != nil {
		return false, fmt.Errorf("failed to check if blocked: %w", err)
	}

	return blocked, nil
}

// IsUserBlocked checks if blockerID has blocked blockedID (one direction only)
func (r *PostgresBlockRepository) IsUserBlocked(ctx context.Context, blockerID, blockedID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM blocked_users
			WHERE blocker_id = $1 AND blocked_id = $2
		)
	`

	var blocked bool
	err := r.db.QueryRow(ctx, query, blockerID, blockedID).Scan(&blocked)
	if err != nil {
		return false, fmt.Errorf("failed to check if user blocked: %w", err)
	}

	return blocked, nil
}

// GetBlockedUsers retrieves all users blocked by a user
func (r *PostgresBlockRepository) GetBlockedUsers(ctx context.Context, blockerID string) ([]BlockedUser, error) {
	query := `
		SELECT id, blocker_id, blocked_id, reason, created_at
		FROM blocked_users
		WHERE blocker_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, blockerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked users: %w", err)
	}
	defer rows.Close()

	var blockedUsers []BlockedUser
	for rows.Next() {
		var bu BlockedUser
		err := rows.Scan(
			&bu.ID,
			&bu.BlockerID,
			&bu.BlockedID,
			&bu.Reason,
			&bu.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan blocked user: %w", err)
		}
		blockedUsers = append(blockedUsers, bu)
	}

	return blockedUsers, nil
}

// GetBlockedUserIDs retrieves just the IDs of blocked users (for quick checks)
func (r *PostgresBlockRepository) GetBlockedUserIDs(ctx context.Context, blockerID string) ([]string, error) {
	query := `
		SELECT blocked_id
		FROM blocked_users
		WHERE blocker_id = $1
	`

	rows, err := r.db.Query(ctx, query, blockerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked user IDs: %w", err)
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan blocked user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

