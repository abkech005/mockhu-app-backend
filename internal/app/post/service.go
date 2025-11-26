package post

import (
	"context"
	"errors"
	"mockhu-app-backend/internal/app/auth"
)

// Common errors returned by the service layer
var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnauthorized = errors.New("unauthorized to perform this action")
	ErrInvalidInput = errors.New("invalid input")
	ErrUserNotFound = errors.New("user not found")
)

// PostService defines the business logic for post operations
// This service layer sits between HTTP handlers and the repository
// It handles validation, permission checking, and data transformation
type PostService interface {
	// Post operations
	CreatePost(ctx context.Context, userID string, req CreatePostRequest) (*PostResponse, error)
	GetPost(ctx context.Context, postID, currentUserID string) (*PostResponse, error)
	GetUserPosts(ctx context.Context, userID, currentUserID string, page, limit int) ([]*PostResponse, error)
	DeletePost(ctx context.Context, postID, userID string) error

	// Feed operations
	GetFeed(ctx context.Context, userID string, page, limit int) (*FeedResponse, error)

	// Reactions
	ToggleReaction(ctx context.Context, postID, userID string) (*ReactionResponse, error)
}

// postService implements the PostService interface
type postService struct {
	postRepo PostRepository      // Repository for post operations
	userRepo auth.UserRepository // Repository for user operations (to get author info)
}

// NewPostService creates a new post service instance
// Parameters:
//   - postRepo: Repository for post database operations
//   - userRepo: Repository for user database operations
//
// Returns: PostService interface implementation
func NewPostService(postRepo PostRepository, userRepo auth.UserRepository) PostService {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// CreatePost creates a new post
// This method:
//   - Validates input (content length, image count)
//   - Creates post in database
//   - Returns formatted response with author info
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - userID: ID of user creating the post
//   - req: Request data (content, images, is_anonymous)
//
// Returns:
//   - *PostResponse: Formatted post response with author info
//   - error: Validation or database errors
func (s *postService) CreatePost(ctx context.Context, userID string, req CreatePostRequest) (*PostResponse, error) {
	// Validate input
	if req.Content == "" {
		return nil, ErrInvalidInput
	}

	if len(req.Content) > 5000 {
		return nil, errors.New("content exceeds maximum length of 5000 characters")
	}

	if len(req.Images) > 10 {
		return nil, errors.New("maximum 10 images allowed")
	}

	// Create post entity
	post := &Post{
		UserID:      userID,
		Content:     req.Content,
		Images:      req.Images,
		IsAnonymous: req.IsAnonymous,
	}

	// Save to database
	err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	// Get user info for response
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build response with author info
	response := &PostResponse{
		ID:      post.ID,
		Content: post.Content,
		Images:  post.Images,
		Author: AuthorInfo{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			AvatarURL: user.AvatarURL,
		},
		Reactions: ReactionInfo{
			FireCount:   0,
			IsFiredByMe: false,
			RecentUsers: []AuthorInfo{},
		},
		CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Hide author if anonymous
	if post.IsAnonymous {
		response.Author = AuthorInfo{
			ID:        "",
			Username:  "Anonymous Student",
			FirstName: "Anonymous",
			AvatarURL: "",
		}
	}

	return response, nil
}

// GetPost retrieves a single post by ID with complete information
// This method:
//   - Gets post from database
//   - Enriches with author info
//   - Adds reaction count and user's reaction status
//   - Adds recent reactions (up to 3)
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - postID: ID of post to retrieve
//   - currentUserID: ID of user requesting (to check if they reacted)
//
// Returns:
//   - *PostResponse: Complete post with author and reaction info
//   - error: ErrPostNotFound if not found, other errors on failure
func (s *postService) GetPost(ctx context.Context, postID, currentUserID string) (*PostResponse, error) {
	// Get post from database
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, ErrPostNotFound
	}

	// Get author info
	author, err := s.userRepo.FindByID(ctx, post.UserID)
	if err != nil {
		return nil, err
	}

	// Get reaction count
	reactionCount, err := s.postRepo.GetReactionCount(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Check if current user reacted
	hasReacted, err := s.postRepo.HasUserReacted(ctx, postID, currentUserID)
	if err != nil {
		return nil, err
	}

	// Get recent reactions (up to 3 for display)
	reactions, err := s.postRepo.GetReactions(ctx, postID, 3, 0)
	if err != nil {
		return nil, err
	}

	// Get user info for recent reactions
	recentUsers := []AuthorInfo{}
	for _, reaction := range reactions {
		reactUser, err := s.userRepo.FindByID(ctx, reaction.UserID)
		if err != nil {
			continue // Skip if user not found
		}
		recentUsers = append(recentUsers, AuthorInfo{
			ID:        reactUser.ID,
			Username:  reactUser.Username,
			FirstName: reactUser.FirstName,
			AvatarURL: reactUser.AvatarURL,
		})
	}

	// Build complete response
	response := &PostResponse{
		ID:      post.ID,
		Content: post.Content,
		Images:  post.Images,
		Author: AuthorInfo{
			ID:        author.ID,
			Username:  author.Username,
			FirstName: author.FirstName,
			AvatarURL: author.AvatarURL,
		},
		Reactions: ReactionInfo{
			FireCount:   reactionCount,
			IsFiredByMe: hasReacted,
			RecentUsers: recentUsers,
		},
		CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Hide author if anonymous
	if post.IsAnonymous {
		response.Author = AuthorInfo{
			ID:        "",
			Username:  "Anonymous Student",
			FirstName: "Anonymous",
			AvatarURL: "",
		}
	}

	return response, nil
}

// GetUserPosts retrieves all posts by a specific user
// This method:
//   - Gets posts from database with pagination
//   - Enriches each post with reaction info
//   - Checks if current user reacted to each post
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - userID: ID of user whose posts to retrieve
//   - currentUserID: ID of user requesting (to check reactions)
//   - page: Page number (1-indexed)
//   - limit: Number of posts per page
//
// Returns:
//   - []*PostResponse: List of posts with complete info
//   - error: Errors from database or user lookup
func (s *postService) GetUserPosts(ctx context.Context, userID, currentUserID string, page, limit int) ([]*PostResponse, error) {
	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Get posts from database
	posts, err := s.postRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get author info (same user for all posts)
	author, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build responses
	responses := make([]*PostResponse, 0, len(posts))
	for _, post := range posts {
		// Get reaction count for this post
		reactionCount, _ := s.postRepo.GetReactionCount(ctx, post.ID)

		// Check if current user reacted to this post
		hasReacted, _ := s.postRepo.HasUserReacted(ctx, post.ID, currentUserID)

		response := &PostResponse{
			ID:      post.ID,
			Content: post.Content,
			Images:  post.Images,
			Author: AuthorInfo{
				ID:        author.ID,
				Username:  author.Username,
				FirstName: author.FirstName,
				AvatarURL: author.AvatarURL,
			},
			Reactions: ReactionInfo{
				FireCount:   reactionCount,
				IsFiredByMe: hasReacted,
				RecentUsers: []AuthorInfo{}, // Don't load for lists (performance)
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Hide author if anonymous
		if post.IsAnonymous {
			response.Author = AuthorInfo{
				ID:        "",
				Username:  "Anonymous Student",
				FirstName: "Anonymous",
				AvatarURL: "",
			}
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// DeletePost soft deletes a post
// This method:
//   - Verifies post exists
//   - Checks if user owns the post (authorization)
//   - Performs soft delete
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - postID: ID of post to delete
//   - userID: ID of user requesting deletion
//
// Returns:
//   - error: ErrPostNotFound, ErrUnauthorized, or database errors
func (s *postService) DeletePost(ctx context.Context, postID, userID string) error {
	// Get post to verify ownership
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	if post == nil {
		return ErrPostNotFound
	}

	// Check if user owns the post (authorization check)
	if post.UserID != userID {
		return ErrUnauthorized
	}

	// Delete the post (soft delete)
	err = s.postRepo.Delete(ctx, postID)
	if err != nil {
		return err
	}

	return nil
}

// GetFeed retrieves the personalized feed for a user
// This method:
//   - Gets posts from users that current user follows
//   - Enriches each post with author and reaction info
//   - Returns paginated results
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - userID: ID of user requesting feed
//   - page: Page number (1-indexed)
//   - limit: Number of posts per page
//
// Returns:
//   - *FeedResponse: Feed with posts and pagination info
//   - error: Database or user lookup errors
func (s *postService) GetFeed(ctx context.Context, userID string, page, limit int) (*FeedResponse, error) {
	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Get posts from followed users
	posts, err := s.postRepo.GetFeed(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Build responses with author info
	postResponses := make([]*PostResponse, 0, len(posts))
	for _, post := range posts {
		// Get author info for this post
		author, err := s.userRepo.FindByID(ctx, post.UserID)
		if err != nil {
			continue // Skip if user not found
		}

		// Get reaction count
		reactionCount, _ := s.postRepo.GetReactionCount(ctx, post.ID)

		// Check if current user reacted
		hasReacted, _ := s.postRepo.HasUserReacted(ctx, post.ID, userID)

		response := &PostResponse{
			ID:      post.ID,
			Content: post.Content,
			Images:  post.Images,
			Author: AuthorInfo{
				ID:        author.ID,
				Username:  author.Username,
				FirstName: author.FirstName,
				AvatarURL: author.AvatarURL,
			},
			Reactions: ReactionInfo{
				FireCount:   reactionCount,
				IsFiredByMe: hasReacted,
				RecentUsers: []AuthorInfo{}, // Don't load for feed (performance)
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Hide author if anonymous
		if post.IsAnonymous {
			response.Author = AuthorInfo{
				ID:        "",
				Username:  "Anonymous Student",
				FirstName: "Anonymous",
				AvatarURL: "",
			}
		}

		postResponses = append(postResponses, response)
	}

	// Build pagination info (simplified for MVP)
	pagination := PaginationInfo{
		Page:       page,
		TotalPages: -1, // TODO: Calculate total pages (requires COUNT query)
		TotalItems: -1, // TODO: Count total items (requires COUNT query)
		Limit:      limit,
	}

	return &FeedResponse{
		Posts:      postResponses,
		Pagination: pagination,
	}, nil
}

// ToggleReaction adds or removes a reaction (fire) on a post
// This method:
//   - Checks if user already reacted
//   - If yes: Removes reaction
//   - If no: Adds reaction
//   - Returns updated reaction count and status
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - postID: ID of post to react to
//   - userID: ID of user reacting
//
// Returns:
//   - *ReactionResponse: Updated reaction count and user's reaction status
//   - error: ErrPostNotFound or database errors
func (s *postService) ToggleReaction(ctx context.Context, postID, userID string) (*ReactionResponse, error) {
	// Check if post exists
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post == nil {
		return nil, ErrPostNotFound
	}

	// Check if user already reacted
	hasReacted, err := s.postRepo.HasUserReacted(ctx, postID, userID)
	if err != nil {
		return nil, err
	}

	if hasReacted {
		// Remove reaction (unlike)
		err = s.postRepo.RemoveReaction(ctx, postID, userID)
		if err != nil {
			return nil, err
		}
	} else {
		// Add reaction (like)
		reaction := &Reaction{
			PostID:       postID,
			UserID:       userID,
			ReactionType: "fire", // Currently only fire reaction
		}
		err = s.postRepo.AddReaction(ctx, reaction)
		if err != nil {
			return nil, err
		}
	}

	// Get updated reaction count
	reactionCount, err := s.postRepo.GetReactionCount(ctx, postID)
	if err != nil {
		return nil, err
	}

	return &ReactionResponse{
		PostID:      postID,
		FireCount:   reactionCount,
		IsFiredByMe: !hasReacted, // Toggled state
	}, nil
}

// ReactionResponse is the response for toggling a reaction
type ReactionResponse struct {
	PostID      string `json:"post_id"`
	FireCount   int    `json:"fire_count"`
	IsFiredByMe bool   `json:"is_fired_by_me"`
}
