package messaging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mockhu-app-backend/internal/app/auth"
)

// MessagingService defines the business logic for messaging operations
type MessagingService interface {
	// Conversation operations
	CreateOrGetConversation(ctx context.Context, currentUserID, recipientID string) (*ConversationResponse, error)
	GetConversations(ctx context.Context, currentUserID string, page, limit int, unreadOnly bool) (*ConversationListResponse, error)
	GetConversation(ctx context.Context, conversationID, currentUserID string) (*ConversationResponse, error)
	DeleteConversation(ctx context.Context, conversationID, currentUserID string) error

	// Message operations
	SendMessage(ctx context.Context, conversationID, senderID string, req *SendMessageRequest) (*MessageResponse, error)
	GetMessages(ctx context.Context, conversationID, currentUserID string, page, limit int) (*MessageListResponse, error)
	DeleteMessage(ctx context.Context, messageID, currentUserID string) error

	// Unread operations
	MarkConversationAsRead(ctx context.Context, conversationID, currentUserID string) error
	MarkMessageAsRead(ctx context.Context, messageID, currentUserID string) error
	GetUnreadCount(ctx context.Context, currentUserID string) (*UnreadCountResponse, error)

	// Privacy operations
	CanMessage(ctx context.Context, senderID, recipientID string) (*CanMessageResponse, error)
	BlockUser(ctx context.Context, blockerID, blockedID, reason string) error
	UnblockUser(ctx context.Context, blockerID, blockedID string) error
	GetBlockedUsers(ctx context.Context, blockerID string) (*BlockedUsersListResponse, error)
}

// messagingService implements MessagingService
type messagingService struct {
	convRepo       ConversationRepository
	msgRepo        MessageRepository
	blockRepo      BlockRepository
	userRepo       auth.UserRepository
	privacyChecker *PrivacyChecker
}

// NewService creates a new messaging service
func NewService(
	convRepo ConversationRepository,
	msgRepo MessageRepository,
	blockRepo BlockRepository,
	userRepo auth.UserRepository,
	privacyChecker *PrivacyChecker,
) MessagingService {
	return &messagingService{
		convRepo:       convRepo,
		msgRepo:        msgRepo,
		blockRepo:      blockRepo,
		userRepo:       userRepo,
		privacyChecker: privacyChecker,
	}
}

// ============================================================================
// CONVERSATION OPERATIONS
// ============================================================================

// CreateOrGetConversation creates or retrieves a conversation
func (s *messagingService) CreateOrGetConversation(ctx context.Context, currentUserID, recipientID string) (*ConversationResponse, error) {
	// Validate input
	if currentUserID == "" || recipientID == "" {
		return nil, fmt.Errorf("user IDs cannot be empty")
	}

	if currentUserID == recipientID {
		return nil, fmt.Errorf("cannot create conversation with yourself")
	}

	// Check if recipient exists
	recipient, err := s.userRepo.FindByID(ctx, recipientID)
	if err != nil {
		return nil, fmt.Errorf("recipient not found")
	}

	// Check if existing conversation
	existingConv, err := s.convRepo.GetConversationByParticipants(ctx, currentUserID, recipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}

	// Check privacy - can sender message recipient?
	canMsg, reason, err := s.privacyChecker.CanMessage(ctx, currentUserID, recipientID, existingConv != nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check messaging permission: %w", err)
	}
	if !canMsg {
		return nil, fmt.Errorf("%s", reason)
	}

	// Create or get conversation
	conv, err := s.convRepo.CreateOrGetConversation(ctx, currentUserID, recipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Get unread count
	unreadCount, err := s.msgRepo.GetConversationUnreadCount(ctx, conv.ID, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread count: %w", err)
	}

	// Check if blocked
	isBlocked, err := s.privacyChecker.IsBlocked(ctx, currentUserID, recipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check blocking status: %w", err)
	}

	// Build response
	participantInfo := s.buildUserBasicInfo(recipient)

	response := &ConversationResponse{
		ID:          conv.ID,
		Participant: participantInfo,
		UnreadCount: unreadCount,
		IsBlocked:   isBlocked,
		CreatedAt:   conv.CreatedAt,
		UpdatedAt:   conv.UpdatedAt,
	}

	// Add last message if exists
	if conv.LastMessageID != nil {
		response.LastMessage = &LastMessageInfo{
			ID:          *conv.LastMessageID,
			Content:     getStringValue(conv.LastMessageText),
			SenderID:    getStringValue(conv.LastMessageSenderID),
			CreatedAt:   getTimeValue(conv.LastMessageAt),
		}
	}

	return response, nil
}

// GetConversations retrieves all conversations for a user
func (s *messagingService) GetConversations(ctx context.Context, currentUserID string, page, limit int, unreadOnly bool) (*ConversationListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	// Get conversations
	conversations, total, err := s.convRepo.GetUserConversations(ctx, currentUserID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	// Get total unread count
	totalUnread, err := s.msgRepo.GetUnreadCount(ctx, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total unread: %w", err)
	}

	// Build response
	items := make([]ConversationListItem, 0, len(conversations))
	for _, conv := range conversations {
		// Get recipient ID (the other user)
		recipientID := conv.GetRecipientID(currentUserID)

		// Get recipient info
		recipient, err := s.userRepo.FindByID(ctx, recipientID)
		if err != nil {
			continue // Skip if user not found
		}

		// Get unread count for this conversation
		unreadCount, err := s.msgRepo.GetConversationUnreadCount(ctx, conv.ID, currentUserID)
		if err != nil {
			unreadCount = 0
		}

		// Filter by unread if requested
		if unreadOnly && unreadCount == 0 {
			continue
		}

		item := ConversationListItem{
			ID:          conv.ID,
			Participant: s.buildUserBasicInfo(recipient),
			UnreadCount: unreadCount,
			UpdatedAt:   conv.UpdatedAt,
		}

		// Add last message if exists
		if conv.LastMessageID != nil {
			item.LastMessage = &LastMessageInfo{
				ID:        *conv.LastMessageID,
				Content:   getStringValue(conv.LastMessageText),
				SenderID:  getStringValue(conv.LastMessageSenderID),
				CreatedAt: getTimeValue(conv.LastMessageAt),
			}
		}

		items = append(items, item)
	}

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	hasMore := page < totalPages

	response := &ConversationListResponse{
		Conversations: items,
		Pagination: PaginationMetadata{
			Page:        page,
			Limit:       limit,
			Total:       total,
			TotalPages:  totalPages,
			HasMore:     hasMore,
			TotalUnread: totalUnread,
		},
	}

	return response, nil
}

// GetConversation retrieves a single conversation
func (s *messagingService) GetConversation(ctx context.Context, conversationID, currentUserID string) (*ConversationResponse, error) {
	// Get conversation
	conv, err := s.convRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found")
	}

	// Verify user is a participant
	if !conv.IsParticipant(currentUserID) {
		return nil, fmt.Errorf("you are not a participant in this conversation")
	}

	// Get recipient ID
	recipientID := conv.GetRecipientID(currentUserID)

	// Get recipient info
	recipient, err := s.userRepo.FindByID(ctx, recipientID)
	if err != nil {
		return nil, fmt.Errorf("recipient not found")
	}

	// Get unread count
	unreadCount, err := s.msgRepo.GetConversationUnreadCount(ctx, conv.ID, currentUserID)
	if err != nil {
		unreadCount = 0
	}

	// Check if blocked
	isBlocked, err := s.privacyChecker.IsBlocked(ctx, currentUserID, recipientID)
	if err != nil {
		isBlocked = false
	}

	// Build response
	response := &ConversationResponse{
		ID:          conv.ID,
		Participant: s.buildUserBasicInfo(recipient),
		UnreadCount: unreadCount,
		IsBlocked:   isBlocked,
		CreatedAt:   conv.CreatedAt,
		UpdatedAt:   conv.UpdatedAt,
	}

	// Add last message if exists
	if conv.LastMessageID != nil {
		response.LastMessage = &LastMessageInfo{
			ID:        *conv.LastMessageID,
			Content:   getStringValue(conv.LastMessageText),
			SenderID:  getStringValue(conv.LastMessageSenderID),
			CreatedAt: getTimeValue(conv.LastMessageAt),
		}
	}

	return response, nil
}

// DeleteConversation soft deletes a conversation
func (s *messagingService) DeleteConversation(ctx context.Context, conversationID, currentUserID string) error {
	// Get conversation to verify participation
	conv, err := s.convRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("conversation not found")
	}

	// Verify user is a participant
	if !conv.IsParticipant(currentUserID) {
		return fmt.Errorf("you are not a participant in this conversation")
	}

	// Delete conversation
	if err := s.convRepo.DeleteConversation(ctx, conversationID, currentUserID); err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	return nil
}

// ============================================================================
// MESSAGE OPERATIONS
// ============================================================================

// SendMessage sends a message in a conversation
func (s *messagingService) SendMessage(ctx context.Context, conversationID, senderID string, req *SendMessageRequest) (*MessageResponse, error) {
	// Validate input
	if err := s.validateSendMessageRequest(req); err != nil {
		return nil, err
	}

	// Get conversation
	conv, err := s.convRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found")
	}

	// Verify sender is a participant
	if !conv.IsParticipant(senderID) {
		return nil, fmt.Errorf("you are not a participant in this conversation")
	}

	// Get recipient ID
	recipientID := conv.GetRecipientID(senderID)

	// Check if sender can message recipient
	canMsg, reason, err := s.privacyChecker.CanMessage(ctx, senderID, recipientID, true) // existing conversation
	if err != nil {
		return nil, fmt.Errorf("failed to check messaging permission: %w", err)
	}
	if !canMsg {
		return nil, fmt.Errorf("%s", reason)
	}

	// Create message
	message := &Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		MessageType:    req.MessageType,
		Content:        &req.Content,
		Attachments:    req.Attachments,
		Status:         "sent",
		IsRead:         false,
	}

	if err := s.msgRepo.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Update conversation's last message
	lastMsgText := s.getLastMessageText(req)
	if err := s.convRepo.UpdateLastMessage(ctx, conversationID, message.ID, lastMsgText, senderID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to update last message: %v\n", err)
	}

	// Get sender info
	sender, err := s.userRepo.FindByID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender info: %w", err)
	}

	// Build response
	response := &MessageResponse{
		ID:             message.ID,
		ConversationID: message.ConversationID,
		Sender:         s.buildUserBasicInfo(sender),
		MessageType:    message.MessageType,
		Content:        req.Content,
		Attachments:    message.Attachments,
		Status:         message.Status,
		IsRead:         message.IsRead,
		CreatedAt:      message.CreatedAt,
	}

	return response, nil
}

// GetMessages retrieves messages for a conversation
func (s *messagingService) GetMessages(ctx context.Context, conversationID, currentUserID string, page, limit int) (*MessageListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Get conversation to verify participation
	conv, err := s.convRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found")
	}

	// Verify user is a participant
	if !conv.IsParticipant(currentUserID) {
		return nil, fmt.Errorf("you are not a participant in this conversation")
	}

	// Get messages
	messages, total, err := s.msgRepo.GetConversationMessages(ctx, conversationID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Build response
	items := make([]MessageResponse, 0, len(messages))
	for _, msg := range messages {
		// Get sender info
		sender, err := s.userRepo.FindByID(ctx, msg.SenderID)
		if err != nil {
			continue // Skip if sender not found
		}

		item := MessageResponse{
			ID:             msg.ID,
			ConversationID: msg.ConversationID,
			Sender:         s.buildUserBasicInfo(sender),
			MessageType:    msg.MessageType,
			Content:        getStringValue(msg.Content),
			Attachments:    msg.Attachments,
			Status:         msg.Status,
			IsRead:         msg.IsRead,
			ReadAt:         msg.ReadAt,
			CreatedAt:      msg.CreatedAt,
		}

		items = append(items, item)
	}

	// Auto-mark as read (optional - could be done explicitly by client)
	// Uncomment if you want automatic read marking when fetching messages
	// _ = s.msgRepo.MarkConversationAsRead(ctx, conversationID, currentUserID)

	// Calculate pagination
	totalPages := (total + limit - 1) / limit
	hasMore := page < totalPages

	response := &MessageListResponse{
		Messages: items,
		Pagination: PaginationMetadata{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasMore:    hasMore,
		},
	}

	return response, nil
}

// DeleteMessage soft deletes a message
func (s *messagingService) DeleteMessage(ctx context.Context, messageID, currentUserID string) error {
	// Get message to verify ownership
	msg, err := s.msgRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found")
	}

	// Verify user is the sender
	if msg.SenderID != currentUserID {
		return fmt.Errorf("you can only delete your own messages")
	}

	// Delete message
	if err := s.msgRepo.DeleteMessage(ctx, messageID, currentUserID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// ============================================================================
// UNREAD OPERATIONS
// ============================================================================

// MarkConversationAsRead marks all messages in a conversation as read
func (s *messagingService) MarkConversationAsRead(ctx context.Context, conversationID, currentUserID string) error {
	// Get conversation to verify participation
	conv, err := s.convRepo.GetConversationByID(ctx, conversationID)
	if err != nil {
		return fmt.Errorf("conversation not found")
	}

	// Verify user is a participant
	if !conv.IsParticipant(currentUserID) {
		return fmt.Errorf("you are not a participant in this conversation")
	}

	// Mark as read
	if err := s.msgRepo.MarkConversationAsRead(ctx, conversationID, currentUserID); err != nil {
		return fmt.Errorf("failed to mark conversation as read: %w", err)
	}

	return nil
}

// MarkMessageAsRead marks a single message as read
func (s *messagingService) MarkMessageAsRead(ctx context.Context, messageID, currentUserID string) error {
	// Get message to verify it's not sender's own message
	msg, err := s.msgRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found")
	}

	// Cannot mark own message as read
	if msg.SenderID == currentUserID {
		return fmt.Errorf("cannot mark your own message as read")
	}

	// Verify user is in the conversation
	conv, err := s.convRepo.GetConversationByID(ctx, msg.ConversationID)
	if err != nil {
		return fmt.Errorf("conversation not found")
	}

	if !conv.IsParticipant(currentUserID) {
		return fmt.Errorf("you are not a participant in this conversation")
	}

	// Mark as read
	if err := s.msgRepo.MarkMessageAsRead(ctx, messageID, currentUserID); err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	return nil
}

// GetUnreadCount returns unread message counts
func (s *messagingService) GetUnreadCount(ctx context.Context, currentUserID string) (*UnreadCountResponse, error) {
	// Get total unread count
	totalUnread, err := s.msgRepo.GetUnreadCount(ctx, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread count: %w", err)
	}

	// Get unread conversations count
	unreadConversations, err := s.msgRepo.GetUnreadConversationsCount(ctx, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread conversations count: %w", err)
	}

	response := &UnreadCountResponse{
		TotalUnread:         totalUnread,
		UnreadConversations: unreadConversations,
	}

	return response, nil
}

// ============================================================================
// PRIVACY OPERATIONS
// ============================================================================

// CanMessage checks if a user can message another user
func (s *messagingService) CanMessage(ctx context.Context, senderID, recipientID string) (*CanMessageResponse, error) {
	// Check if existing conversation
	existingConv, err := s.convRepo.GetConversationByParticipants(ctx, senderID, recipientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}

	// Check privacy
	canMsg, reason, err := s.privacyChecker.CanMessage(ctx, senderID, recipientID, existingConv != nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check messaging permission: %w", err)
	}

	response := &CanMessageResponse{
		CanMessage: canMsg,
		Reason:     reason,
	}

	return response, nil
}

// BlockUser blocks a user
func (s *messagingService) BlockUser(ctx context.Context, blockerID, blockedID, reason string) error {
	// Validate input
	if blockerID == "" || blockedID == "" {
		return fmt.Errorf("user IDs cannot be empty")
	}

	if blockerID == blockedID {
		return fmt.Errorf("cannot block yourself")
	}

	// Check if blocked user exists
	_, err := s.userRepo.FindByID(ctx, blockedID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Block user
	if err := s.blockRepo.BlockUser(ctx, blockerID, blockedID, reason); err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	return nil
}

// UnblockUser unblocks a user
func (s *messagingService) UnblockUser(ctx context.Context, blockerID, blockedID string) error {
	// Validate input
	if blockerID == "" || blockedID == "" {
		return fmt.Errorf("user IDs cannot be empty")
	}

	// Unblock user
	if err := s.blockRepo.UnblockUser(ctx, blockerID, blockedID); err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	return nil
}

// GetBlockedUsers retrieves all blocked users
func (s *messagingService) GetBlockedUsers(ctx context.Context, blockerID string) (*BlockedUsersListResponse, error) {
	// Get blocked users
	blockedUsers, err := s.blockRepo.GetBlockedUsers(ctx, blockerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked users: %w", err)
	}

	// Build response
	items := make([]BlockedUserResponse, 0, len(blockedUsers))
	for _, bu := range blockedUsers {
		// Get user info
		user, err := s.userRepo.FindByID(ctx, bu.BlockedID)
		if err != nil {
			continue // Skip if user not found
		}

		item := BlockedUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FirstName + " " + user.LastName,
			AvatarURL: user.AvatarURL,
			BlockedAt: bu.CreatedAt,
		}

		items = append(items, item)
	}

	response := &BlockedUsersListResponse{
		BlockedUsers: items,
	}

	return response, nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// validateSendMessageRequest validates the send message request
func (s *messagingService) validateSendMessageRequest(req *SendMessageRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate message type
	if req.MessageType != "text" && req.MessageType != "image" && req.MessageType != "file" {
		return fmt.Errorf("invalid message type: must be text, image, or file")
	}

	// For text messages, content is required
	if req.MessageType == "text" {
		if strings.TrimSpace(req.Content) == "" {
			return fmt.Errorf("text message content cannot be empty")
		}
		if len(req.Content) > 10000 {
			return fmt.Errorf("message content too long (max 10,000 characters)")
		}
	}

	// For image/file messages, attachments are required
	if req.MessageType == "image" || req.MessageType == "file" {
		if len(req.Attachments) == 0 {
			return fmt.Errorf("%s message requires attachments", req.MessageType)
		}
		if len(req.Attachments) > 5 {
			return fmt.Errorf("maximum 5 attachments per message")
		}
	}

	return nil
}

// buildUserBasicInfo builds UserBasicInfo from User model
func (s *messagingService) buildUserBasicInfo(user *auth.User) UserBasicInfo {
	return UserBasicInfo{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FirstName + " " + user.LastName,
		AvatarURL: user.AvatarURL,
	}
}

// getLastMessageText gets a preview text for the last message
func (s *messagingService) getLastMessageText(req *SendMessageRequest) string {
	switch req.MessageType {
	case "text":
		if len(req.Content) > 100 {
			return req.Content[:100] + "..."
		}
		return req.Content
	case "image":
		if req.Content != "" {
			return "📷 " + req.Content
		}
		return "📷 Photo"
	case "file":
		if len(req.Attachments) > 0 {
			return "📎 " + req.Attachments[0].Filename
		}
		return "📎 File"
	default:
		return ""
	}
}

// Helper functions for pointer values
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getTimeValue(ptr *time.Time) time.Time {
	if ptr == nil {
		return time.Time{}
	}
	return *ptr
}

