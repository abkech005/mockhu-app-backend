package auth

import "context"

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}
