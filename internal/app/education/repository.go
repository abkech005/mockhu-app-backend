package education

import "context"

// EducationRepository defines methods for education data access
type EducationRepository interface {
	// Education CRUD
	Create(ctx context.Context, education *Education) error
	FindByID(ctx context.Context, id string) (*Education, error)
	FindByUserID(ctx context.Context, userID string) ([]Education, error)
	Update(ctx context.Context, education *Education) error
	Delete(ctx context.Context, id string) error

	// Queries
	FindCurrentEducation(ctx context.Context, userID string) (*Education, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
}
