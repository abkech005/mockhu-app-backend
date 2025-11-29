package messaging

import "time"

// ===============================
// REQUEST DTOs
// ===============================

// CreateConversationRequest for creating/getting a conversation
type CreateConversationRequest struct {
	RecipientID string `json:"recipient_id" validate:"required,uuid"`
}

// SendMessageRequest for sending a message
type SendMessageRequest struct {
	MessageType string               `json:"message_type" validate:"required,oneof=text image file"`
	Content     string               `json:"content"`
	Attachments []AttachmentMetadata `json:"attachments"`
}

// UploadFilesRequest for uploading files
type UploadFilesRequest struct {
	Type string `json:"type" validate:"required,oneof=image file"`
}

// BlockUserRequest for blocking a user
type BlockUserRequest struct {
	Reason string `json:"reason"`
}

// ===============================
// RESPONSE DTOs
// ===============================

// UserBasicInfo represents basic user information for responses
type UserBasicInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
	IsOnline  bool   `json:"is_online,omitempty"`
	LastSeen  string `json:"last_seen,omitempty"`
}

// LastMessageInfo represents last message information for conversation list
type LastMessageInfo struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	SenderID    string    `json:"sender_id"`
	MessageType string    `json:"message_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// ConversationResponse for single conversation details
type ConversationResponse struct {
	ID            string          `json:"id"`
	Participant   UserBasicInfo   `json:"participant"`
	LastMessage   *LastMessageInfo `json:"last_message,omitempty"`
	UnreadCount   int             `json:"unread_count"`
	IsBlocked     bool            `json:"is_blocked"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// ConversationListItem for conversation list
type ConversationListItem struct {
	ID          string           `json:"id"`
	Participant UserBasicInfo    `json:"participant"`
	LastMessage *LastMessageInfo `json:"last_message,omitempty"`
	UnreadCount int              `json:"unread_count"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ConversationListResponse for list of conversations
type ConversationListResponse struct {
	Conversations []ConversationListItem `json:"conversations"`
	Pagination    PaginationMetadata     `json:"pagination"`
}

// MessageResponse for single message with sender info
type MessageResponse struct {
	ID             string               `json:"id"`
	ConversationID string               `json:"conversation_id"`
	Sender         UserBasicInfo        `json:"sender"`
	MessageType    string               `json:"message_type"`
	Content        string               `json:"content,omitempty"`
	Attachments    []AttachmentMetadata `json:"attachments,omitempty"`
	Status         string               `json:"status"`
	IsRead         bool                 `json:"is_read"`
	ReadAt         *time.Time           `json:"read_at,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
}

// MessageListResponse for list of messages
type MessageListResponse struct {
	Messages   []MessageResponse  `json:"messages"`
	Pagination PaginationMetadata `json:"pagination"`
}

// UnreadCountResponse for unread message counts
type UnreadCountResponse struct {
	TotalUnread          int `json:"total_unread"`
	UnreadConversations  int `json:"unread_conversations"`
}

// CanMessageResponse for checking if user can message another user
type CanMessageResponse struct {
	CanMessage bool   `json:"can_message"`
	Reason     string `json:"reason,omitempty"`
}

// FileUploadResponse for file upload result
type FileUploadResponse struct {
	Attachments []AttachmentMetadata `json:"attachments"`
}

// BlockedUserResponse for blocked user info
type BlockedUserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	BlockedAt time.Time `json:"blocked_at"`
}

// BlockedUsersListResponse for list of blocked users
type BlockedUsersListResponse struct {
	BlockedUsers []BlockedUserResponse `json:"blocked_users"`
}

// PaginationMetadata for pagination information
type PaginationMetadata struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	Total       int  `json:"total"`
	TotalPages  int  `json:"total_pages"`
	HasMore     bool `json:"has_more"`
	NextCursor  string `json:"next_cursor,omitempty"`
	TotalUnread int  `json:"total_unread,omitempty"` // For conversation list
}

// ===============================
// WEBSOCKET MESSAGE DTOs
// ===============================

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// WSNewMessageData for new message WebSocket event
type WSNewMessageData struct {
	ID             string               `json:"id"`
	ConversationID string               `json:"conversation_id"`
	Sender         UserBasicInfo        `json:"sender"`
	MessageType    string               `json:"message_type"`
	Content        string               `json:"content,omitempty"`
	Attachments    []AttachmentMetadata `json:"attachments,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
}

// WSMessageReadData for message read event
type WSMessageReadData struct {
	MessageID      string    `json:"message_id"`
	ConversationID string    `json:"conversation_id"`
	ReadBy         string    `json:"read_by"`
	ReadAt         time.Time `json:"read_at"`
}

// WSTypingData for typing indicator event
type WSTypingData struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}

// WSUserStatusData for user online/offline status
type WSUserStatusData struct {
	UserID   string     `json:"user_id"`
	IsOnline bool       `json:"is_online"`
	LastSeen *time.Time `json:"last_seen,omitempty"`
}

// ===============================
// CLIENT WebSocket MESSAGES
// ===============================

// WSClientMessage represents messages sent from client to server
type WSClientMessage struct {
	Type           string               `json:"type"` // "message", "typing", "read"
	ConversationID string               `json:"conversation_id,omitempty"`
	MessageType    string               `json:"message_type,omitempty"`
	Content        string               `json:"content,omitempty"`
	Attachments    []AttachmentMetadata `json:"attachments,omitempty"`
	MessageID      string               `json:"message_id,omitempty"`
	IsTyping       bool                 `json:"is_typing,omitempty"`
}

