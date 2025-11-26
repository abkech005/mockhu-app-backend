package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mockhu-app-backend/internal/app/auth"
	"mockhu-app-backend/internal/app/post"
)

// Errors
var (
	ErrCommentNotFound    = errors.New("comment not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidContent     = errors.New("invalid content")
	ErrPostNotFound       = errors.New("post not found")
	ErrParentNotFound     = errors.New("parent comment not found")
	ErrCannotReplyToReply = errors.New("cannot reply to a reply - only top-level comments can have replies")
)

// CommentService defines the business logic for comment operations
type CommentService interface {
	CreateComment(ctx context.Context, postID, userID string, req *CreateCommentRequest) (*CommentResponse, error)
	GetComment(ctx context.Context, commentID, currentUserID string) (*CommentResponse, error)
	GetPostComments(ctx context.Context, postID, currentUserID string, page, limit int) (*CommentListResponse, error)
	UpdateComment(ctx context.Context, commentID, userID string, content string) (*CommentResponse, error)
	DeleteComment(ctx context.Context, commentID, userID string) error
}

// commentService implements CommentService
type commentService struct {
	commentRepo CommentRepository
	userRepo    auth.UserRepository
	postRepo    post.PostRepository
}

// NewService creates a new comment service
func NewService(commentRepo CommentRepository, userRepo auth.UserRepository, postRepo post.PostRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postRepo:    postRepo,
	}
}

// CreateComment creates a new comment on a post
func (s *commentService) CreateComment(ctx context.Context, postID, userID string, req *CreateCommentRequest) (*CommentResponse, error) {
	// Validate content
	if len(req.Content) < 1 || len(req.Content) > 2000 {
		return nil, ErrInvalidContent
	}

	// Validate post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	// If parent_comment_id is provided, validate it exists and belongs to same post
	// ENFORCE SINGLE-LEVEL: Only allow replies to top-level comments (not replies to replies)
	if req.ParentCommentID != nil && *req.ParentCommentID != "" {
		parent, err := s.commentRepo.GetByID(ctx, *req.ParentCommentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent comment: %w", err)
		}
		if parent == nil {
			return nil, ErrParentNotFound
		}
		if parent.PostID != postID {
			return nil, errors.New("parent comment does not belong to this post")
		}
		// ENFORCE SINGLE-LEVEL: Reject replies to replies
		if parent.ParentCommentID != nil && *parent.ParentCommentID != "" {
			return nil, errors.New("cannot reply to a reply - only top-level comments can have replies")
		}
	}

	// Create comment
	comment := &Comment{
		PostID:          postID,
		UserID:          userID,
		Content:         req.Content,
		ParentCommentID: req.ParentCommentID,
		IsAnonymous:     req.IsAnonymous,
		IsActive:        true,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Get author info
	author, err := s.getAuthorInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author info: %w", err)
	}

	// Build response
	response := &CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		Author:          *author,
		Content:         comment.Content,
		ParentCommentID: comment.ParentCommentID,
		Replies:         []*CommentResponse{},
		ReplyCount:       0,
		CreatedAt:       comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       comment.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// GetComment retrieves a single comment by ID
func (s *commentService) GetComment(ctx context.Context, commentID, currentUserID string) (*CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, ErrCommentNotFound
	}

	return s.convertToResponse(ctx, comment, currentUserID)
}

// GetPostComments retrieves all comments for a post
func (s *commentService) GetPostComments(ctx context.Context, postID, currentUserID string, page, limit int) (*CommentListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get top-level comments
	comments, err := s.commentRepo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	// Convert to response with replies
	commentResponses := make([]*CommentResponse, 0, len(comments))
	for _, comment := range comments {
		response, err := s.convertToResponse(ctx, comment, currentUserID)
		if err != nil {
			continue // Skip comments with errors
		}
		commentResponses = append(commentResponses, response)
	}

	// Calculate total pages (simplified)
	totalPages := 1
	if len(commentResponses) == limit {
		totalPages = page + 1
	}

	return &CommentListResponse{
		Comments: commentResponses,
		Pagination: PaginationInfo{
			Page:       page,
			TotalPages: totalPages,
			TotalItems: len(commentResponses),
			Limit:      limit,
		},
	}, nil
}

// UpdateComment updates a comment
func (s *commentService) UpdateComment(ctx context.Context, commentID, userID string, content string) (*CommentResponse, error) {
	// Validate content
	if len(content) < 1 || len(content) > 2000 {
		return nil, ErrInvalidContent
	}

	// Get comment to verify ownership
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return nil, ErrCommentNotFound
	}

	// Check ownership
	if comment.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update comment
	comment.Content = content
	err = s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return s.convertToResponse(ctx, comment, userID)
}

// DeleteComment deletes a comment (soft delete)
func (s *commentService) DeleteComment(ctx context.Context, commentID, userID string) error {
	// Get comment to verify ownership
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}
	if comment == nil {
		return ErrCommentNotFound
	}

	// Check ownership
	if comment.UserID != userID {
		return ErrUnauthorized
	}

	// Delete comment
	err = s.commentRepo.Delete(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

// Helper methods

// getAuthorInfo retrieves author information for a comment
func (s *commentService) getAuthorInfo(ctx context.Context, userID string) (*AuthorInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &AuthorInfo{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		AvatarURL: user.AvatarURL,
	}, nil
}

// convertToResponse converts a comment to a response with replies
func (s *commentService) convertToResponse(ctx context.Context, comment *Comment, currentUserID string) (*CommentResponse, error) {
	// Get author info
	author, err := s.getAuthorInfo(ctx, comment.UserID)
	if err != nil {
		return nil, err
	}

	// Get reply count
	replyCount, _ := s.commentRepo.GetReplyCount(ctx, comment.ID)

	// Get replies (limit to 5 for preview)
	replies, _ := s.commentRepo.GetReplies(ctx, comment.ID, 5, 0)
	replyResponses := make([]*CommentResponse, 0, len(replies))
	for _, reply := range replies {
		replyAuthor, err := s.getAuthorInfo(ctx, reply.UserID)
		if err != nil {
			continue
		}
		replyResponses = append(replyResponses, &CommentResponse{
			ID:              reply.ID,
			PostID:          reply.PostID,
			Author:          *replyAuthor,
			Content:         reply.Content,
			ParentCommentID: reply.ParentCommentID,
			Replies:         []*CommentResponse{},
			ReplyCount:       0,
			CreatedAt:       reply.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       reply.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		Author:          *author,
		Content:         comment.Content,
		ParentCommentID: comment.ParentCommentID,
		Replies:         replyResponses,
		ReplyCount:       replyCount,
		CreatedAt:       comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       comment.UpdatedAt.Format(time.RFC3339),
	}, nil
}

