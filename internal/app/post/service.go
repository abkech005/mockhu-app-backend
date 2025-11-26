package post

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mockhu-app-backend/internal/app/auth"
)

// Errors
var (
	ErrPostNotFound      = errors.New("post not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidContent    = errors.New("invalid content")
	ErrTooManyImages     = errors.New("too many images (max 10)")
	ErrPostAlreadyExists = errors.New("post already exists")
)

// PostService defines the business logic for post operations
type PostService interface {
	CreatePost(ctx context.Context, userID string, req *CreatePostRequest) (*PostResponse, error)
	GetPost(ctx context.Context, postID, currentUserID string) (*PostResponse, error)
	GetUserPosts(ctx context.Context, userID, currentUserID string, page, limit int) (*FeedResponse, error)
	DeletePost(ctx context.Context, postID, userID string) error
	ToggleReaction(ctx context.Context, postID, userID string) (*ReactionResponse, error)
	GetFeed(ctx context.Context, userID string, page, limit int) (*FeedResponse, error)
}

// postService implements PostService
type postService struct {
	postRepo PostRepository
	userRepo auth.UserRepository
}

// NewService creates a new post service
func NewService(postRepo PostRepository, userRepo auth.UserRepository) PostService {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// CreatePost creates a new post
func (s *postService) CreatePost(ctx context.Context, userID string, req *CreatePostRequest) (*PostResponse, error) {
	// Validate content
	if len(req.Content) < 1 || len(req.Content) > 5000 {
		return nil, ErrInvalidContent
	}

	// Validate images
	if len(req.Images) > 10 {
		return nil, ErrTooManyImages
	}

	// Create post
	post := &Post{
		UserID:      userID,
		Content:     req.Content,
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
		IsActive:    true,
		ViewCount:   0,
	}

	err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Get author info
	author, err := s.getAuthorInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author info: %w", err)
	}

	// Build response
	response := &PostResponse{
		ID:      post.ID,
		Author:  *author,
		Content: post.Content,
		Images:  post.Images,
		Reactions: ReactionInfo{
			FireCount:   0,
			IsFiredByMe: false,
			RecentUsers: []AuthorInfo{},
		},
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// GetPost retrieves a single post by ID
func (s *postService) GetPost(ctx context.Context, postID, currentUserID string) (*PostResponse, error) {
	// Get post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	// Increment view count (async, don't wait)
	go func() {
		// In production, you might want to batch these updates
		_ = s.postRepo.Update(ctx, post)
	}()

	// Get author info
	author, err := s.getAuthorInfo(ctx, post.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get author info: %w", err)
	}

	// Get reaction info
	reactionInfo, err := s.getReactionInfo(ctx, postID, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reaction info: %w", err)
	}

	// Build response
	response := &PostResponse{
		ID:        post.ID,
		Author:    *author,
		Content:   post.Content,
		Images:    post.Images,
		Reactions: *reactionInfo,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

// GetUserPosts retrieves all posts by a specific user
func (s *postService) GetUserPosts(ctx context.Context, userID, currentUserID string, page, limit int) (*FeedResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get posts
	posts, err := s.postRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	// Convert to response
	postResponses, err := s.convertPostsToResponse(ctx, posts, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert posts: %w", err)
	}

	// Calculate total pages (simplified - in production, get total count)
	totalPages := 1
	if len(postResponses) == limit {
		totalPages = page + 1 // Assume there might be more
	}

	return &FeedResponse{
		Posts: postResponses,
		Pagination: PaginationInfo{
			Page:       page,
			TotalPages: totalPages,
			TotalItems: len(postResponses),
			Limit:      limit,
		},
	}, nil
}

// DeletePost deletes a post (soft delete)
func (s *postService) DeletePost(ctx context.Context, postID, userID string) error {
	// Get post to verify ownership
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return ErrPostNotFound
	}

	// Check ownership
	if post.UserID != userID {
		return ErrUnauthorized
	}

	// Delete post
	err = s.postRepo.Delete(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// ToggleReaction toggles a fire reaction on a post
func (s *postService) ToggleReaction(ctx context.Context, postID, userID string) (*ReactionResponse, error) {
	// Check if post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	// Check if user already reacted
	hasReacted, err := s.postRepo.HasUserReacted(ctx, postID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check reaction: %w", err)
	}

	if hasReacted {
		// Remove reaction
		err = s.postRepo.RemoveReaction(ctx, postID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to remove reaction: %w", err)
		}
	} else {
		// Add reaction
		reaction := &Reaction{
			PostID:       postID,
			UserID:       userID,
			ReactionType: "fire",
		}
		err = s.postRepo.AddReaction(ctx, reaction)
		if err != nil {
			return nil, fmt.Errorf("failed to add reaction: %w", err)
		}
	}

	// Get updated reaction count
	count, err := s.postRepo.GetReactionCount(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reaction count: %w", err)
	}

	// Check if user reacted (after toggle)
	hasReacted, _ = s.postRepo.HasUserReacted(ctx, postID, userID)

	return &ReactionResponse{
		PostID:      postID,
		FireCount:   count,
		IsFiredByMe: hasReacted,
	}, nil
}

// GetFeed retrieves posts from users that the current user follows
func (s *postService) GetFeed(ctx context.Context, userID string, page, limit int) (*FeedResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get feed posts
	posts, err := s.postRepo.GetFeed(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}

	// Convert to response
	postResponses, err := s.convertPostsToResponse(ctx, posts, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert posts: %w", err)
	}

	// Calculate total pages (simplified)
	totalPages := 1
	if len(postResponses) == limit {
		totalPages = page + 1
	}

	return &FeedResponse{
		Posts: postResponses,
		Pagination: PaginationInfo{
			Page:       page,
			TotalPages: totalPages,
			TotalItems: len(postResponses),
			Limit:      limit,
		},
	}, nil
}

// Helper methods

// getAuthorInfo retrieves author information for a post
func (s *postService) getAuthorInfo(ctx context.Context, userID string) (*AuthorInfo, error) {
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

// getReactionInfo retrieves reaction information for a post
func (s *postService) getReactionInfo(ctx context.Context, postID, currentUserID string) (*ReactionInfo, error) {
	// Get reaction count
	count, err := s.postRepo.GetReactionCount(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check if current user reacted
	isFiredByMe := false
	if currentUserID != "" {
		hasReacted, err := s.postRepo.HasUserReacted(ctx, postID, currentUserID)
		if err == nil {
			isFiredByMe = hasReacted
		}
	}

	// Get recent users who reacted (limit 5)
	reactions, err := s.postRepo.GetReactions(ctx, postID, 5, 0)
	if err != nil {
		return nil, err
	}

	// Convert to author info
	recentUsers := make([]AuthorInfo, 0, len(reactions))
	for _, reaction := range reactions {
		author, err := s.getAuthorInfo(ctx, reaction.UserID)
		if err == nil && author != nil {
			recentUsers = append(recentUsers, *author)
		}
	}

	return &ReactionInfo{
		FireCount:   count,
		IsFiredByMe: isFiredByMe,
		RecentUsers: recentUsers,
	}, nil
}

// convertPostsToResponse converts a slice of posts to post responses
func (s *postService) convertPostsToResponse(ctx context.Context, posts []*Post, currentUserID string) ([]*PostResponse, error) {
	responses := make([]*PostResponse, 0, len(posts))

	for _, post := range posts {
		// Get author info
		author, err := s.getAuthorInfo(ctx, post.UserID)
		if err != nil {
			continue // Skip posts with invalid authors
		}

		// Get reaction info
		reactionInfo, err := s.getReactionInfo(ctx, post.ID, currentUserID)
		if err != nil {
			reactionInfo = &ReactionInfo{
				FireCount:   0,
				IsFiredByMe: false,
				RecentUsers: []AuthorInfo{},
			}
		}

		responses = append(responses, &PostResponse{
			ID:        post.ID,
			Author:    *author,
			Content:   post.Content,
			Images:    post.Images,
			Reactions: *reactionInfo,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		})
	}

	return responses, nil
}
