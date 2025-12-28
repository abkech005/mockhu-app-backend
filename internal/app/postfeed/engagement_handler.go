package postfeed

import (
	"mockhu-app-backend/internal/app/auth"

	"github.com/gofiber/fiber/v2"
)

// EngagementHandler handles like, comment, share endpoints
type EngagementHandler struct {
	likeRepo    LikeRepository
	commentRepo CommentRepository
	shareRepo   ShareRepository
	authRepo    auth.UserRepository
}

// NewEngagementHandler creates a new engagement handler
func NewEngagementHandler(likeRepo LikeRepository, commentRepo CommentRepository, shareRepo ShareRepository, authRepo auth.UserRepository) *EngagementHandler {
	return &EngagementHandler{
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		shareRepo:   shareRepo,
		authRepo:    authRepo,
	}
}

// --- Like Handlers ---

// LikePostfeed handles POST /v1/postfeeds/:id/like
func (h *EngagementHandler) LikePostfeed(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	postfeedID := c.Params("id")
	if postfeedID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postfeed id is required"})
	}

	if err := h.likeRepo.Like(c.Context(), postfeedID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	count, _ := h.likeRepo.GetLikeCount(c.Context(), postfeedID)
	return c.JSON(LikeResponse{Liked: true, LikeCount: count})
}

// UnlikePostfeed handles DELETE /v1/postfeeds/:id/like
func (h *EngagementHandler) UnlikePostfeed(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	postfeedID := c.Params("id")
	if postfeedID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postfeed id is required"})
	}

	if err := h.likeRepo.Unlike(c.Context(), postfeedID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	count, _ := h.likeRepo.GetLikeCount(c.Context(), postfeedID)
	return c.JSON(LikeResponse{Liked: false, LikeCount: count})
}

// GetLikeStatus handles GET /v1/postfeeds/:id/like
func (h *EngagementHandler) GetLikeStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	postfeedID := c.Params("id")

	var liked bool
	if userIDStr, ok := userID.(string); ok && userIDStr != "" {
		liked, _ = h.likeRepo.IsLiked(c.Context(), postfeedID, userIDStr)
	}

	count, _ := h.likeRepo.GetLikeCount(c.Context(), postfeedID)
	return c.JSON(LikeResponse{Liked: liked, LikeCount: count})
}

// --- Comment Handlers ---

// CreateComment handles POST /v1/postfeeds/:id/comments
func (h *EngagementHandler) CreateComment(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	postfeedID := c.Params("id")
	if postfeedID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postfeed id is required"})
	}

	var req CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content is required"})
	}

	comment := &Comment{
		PostfeedID: postfeedID,
		UserID:     userID,
		ParentID:   req.ParentID,
		Content:    req.Content,
	}

	if err := h.commentRepo.Create(c.Context(), comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp := h.commentToResponse(c, comment)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListComments handles GET /v1/postfeeds/:id/comments
func (h *EngagementHandler) ListComments(c *fiber.Ctx) error {
	postfeedID := c.Params("id")
	if postfeedID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postfeed id is required"})
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	comments, total, err := h.commentRepo.ListByPostfeed(c.Context(), postfeedID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		resp := h.commentToResponse(c, &comment)
		// Get replies for each comment
		replies, _ := h.commentRepo.GetReplies(c.Context(), comment.ID)
		if len(replies) > 0 {
			replyResponses := make([]CommentResponse, len(replies))
			for j, reply := range replies {
				replyResponses[j] = h.commentToResponse(c, &reply)
			}
			resp.Replies = replyResponses
		}
		responses[i] = resp
	}

	totalPages := (total + limit - 1) / limit
	return c.JSON(ListCommentsResponse{
		Comments:   responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// UpdateComment handles PUT /v1/postfeeds/:id/comments/:comment_id
func (h *EngagementHandler) UpdateComment(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	commentID := c.Params("comment_id")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "comment id is required"})
	}

	// Check ownership
	comment, err := h.commentRepo.GetByID(c.Context(), commentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	if comment.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can only edit your own comments"})
	}

	var req UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content is required"})
	}

	if err := h.commentRepo.Update(c.Context(), commentID, req.Content); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updated, _ := h.commentRepo.GetByID(c.Context(), commentID)
	return c.JSON(h.commentToResponse(c, updated))
}

// DeleteComment handles DELETE /v1/postfeeds/:id/comments/:comment_id
func (h *EngagementHandler) DeleteComment(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	commentID := c.Params("comment_id")
	if commentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "comment id is required"})
	}

	// Check ownership
	comment, err := h.commentRepo.GetByID(c.Context(), commentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	if comment.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you can only delete your own comments"})
	}

	if err := h.commentRepo.Delete(c.Context(), commentID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "comment deleted successfully"})
}

// --- Share Handlers ---

// SharePostfeed handles POST /v1/postfeeds/:id/share
func (h *EngagementHandler) SharePostfeed(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	postfeedID := c.Params("id")
	if postfeedID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postfeed id is required"})
	}

	var req CreateShareRequest
	_ = c.BodyParser(&req)

	share := &Share{
		PostfeedID: postfeedID,
		UserID:     userID,
		Message:    req.Message,
	}

	if err := h.shareRepo.Create(c.Context(), share); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp := h.shareToResponse(c, share)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GetShareCount handles GET /v1/postfeeds/:id/shares
func (h *EngagementHandler) GetShareCount(c *fiber.Ctx) error {
	postfeedID := c.Params("id")
	if postfeedID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "postfeed id is required"})
	}

	count, _ := h.shareRepo.GetShareCount(c.Context(), postfeedID)
	return c.JSON(fiber.Map{"share_count": count})
}

// --- Helper Methods ---

func (h *EngagementHandler) commentToResponse(c *fiber.Ctx, comment *Comment) CommentResponse {
	resp := CommentResponse{
		ID:         comment.ID,
		PostfeedID: comment.PostfeedID,
		UserID:     comment.UserID,
		ParentID:   comment.ParentID,
		Content:    comment.Content,
		LikeCount:  comment.LikeCount,
		CreatedAt:  comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Get author info
	if h.authRepo != nil {
		user, err := h.authRepo.FindByID(c.Context(), comment.UserID)
		if err == nil && user != nil {
			resp.Author = &AuthorInfo{
				ID:        user.ID,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				AvatarURL: user.AvatarURL,
			}
		}
	}

	return resp
}

func (h *EngagementHandler) shareToResponse(c *fiber.Ctx, share *Share) ShareResponse {
	resp := ShareResponse{
		ID:         share.ID,
		PostfeedID: share.PostfeedID,
		UserID:     share.UserID,
		Message:    share.Message,
		CreatedAt:  share.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Get author info
	if h.authRepo != nil {
		user, err := h.authRepo.FindByID(c.Context(), share.UserID)
		if err == nil && user != nil {
			resp.Author = &AuthorInfo{
				ID:        user.ID,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				AvatarURL: user.AvatarURL,
			}
		}
	}

	return resp
}
