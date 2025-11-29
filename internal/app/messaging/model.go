package messaging

import (
	"encoding/json"
	"time"
)

// Conversation represents a 1-on-1 conversation between two users
type Conversation struct {
	ID                  string     `json:"id"`
	User1ID             string     `json:"user1_id"`
	User2ID             string     `json:"user2_id"`
	LastMessageID       *string    `json:"last_message_id,omitempty"`
	LastMessageText     *string    `json:"last_message_text,omitempty"`
	LastMessageSenderID *string    `json:"last_message_sender_id,omitempty"`
	LastMessageAt       *time.Time `json:"last_message_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// Message represents a message within a conversation
type Message struct {
	ID             string               `json:"id"`
	ConversationID string               `json:"conversation_id"`
	SenderID       string               `json:"sender_id"`
	MessageType    string               `json:"message_type"` // "text", "image", "file"
	Content        *string              `json:"content,omitempty"`
	Attachments    []AttachmentMetadata `json:"attachments,omitempty"`
	Status         string               `json:"status"` // "sent", "delivered", "read"
	IsRead         bool                 `json:"is_read"`
	ReadAt         *time.Time           `json:"read_at,omitempty"`
	IsDeleted      bool                 `json:"is_deleted"`
	DeletedAt      *time.Time           `json:"deleted_at,omitempty"`
	DeletedBy      *string              `json:"deleted_by,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

// AttachmentMetadata represents metadata for file attachments (images, documents, etc.)
type AttachmentMetadata struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"` // "image" or "file"
	URL              string    `json:"url"`
	ThumbnailURL     string    `json:"thumbnail_url,omitempty"` // For images only
	Filename         string    `json:"filename"`
	OriginalFilename string    `json:"original_filename"`
	Size             int64     `json:"size"`
	MimeType         string    `json:"mime_type"`
	Width            int       `json:"width,omitempty"`  // For images only
	Height           int       `json:"height,omitempty"` // For images only
	UploadedAt       time.Time `json:"uploaded_at"`
}

// BlockedUser represents a user blocking relationship
type BlockedUser struct {
	ID        string    `json:"id"`
	BlockerID string    `json:"blocker_id"`
	BlockedID string    `json:"blocked_id"`
	Reason    *string   `json:"reason,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Helper methods for Message

// AttachmentsJSON converts attachments slice to JSON string for database storage
func (m *Message) AttachmentsJSON() (string, error) {
	if len(m.Attachments) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(m.Attachments)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseAttachments parses JSON string to attachments slice
func (m *Message) ParseAttachments(jsonStr string) error {
	if jsonStr == "" {
		m.Attachments = []AttachmentMetadata{}
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &m.Attachments)
}

// GetRecipientID returns the ID of the recipient (the other user in conversation)
func (c *Conversation) GetRecipientID(currentUserID string) string {
	if c.User1ID == currentUserID {
		return c.User2ID
	}
	return c.User1ID
}

// IsParticipant checks if a user is a participant in the conversation
func (c *Conversation) IsParticipant(userID string) bool {
	return c.User1ID == userID || c.User2ID == userID
}

// OrderUserIDs returns user IDs in correct order (smaller first) for conversation creation
func OrderUserIDs(userID1, userID2 string) (string, string) {
	if userID1 < userID2 {
		return userID1, userID2
	}
	return userID2, userID1
}
