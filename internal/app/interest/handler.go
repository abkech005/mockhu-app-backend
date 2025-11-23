package interest

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for interest endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new interest handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAllInterests handles GET /v1/interests
func (h *Handler) GetAllInterests(c *fiber.Ctx) error {
	// Check for category filter
	category := c.Query("category")
	
	var interests []Interest
	var err error
	
	if category != "" {
		interests, err = h.service.GetInterestsByCategory(c.Context(), category)
	} else {
		interests, err = h.service.GetAllInterests(c.Context())
	}
	
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(ListInterestsResponse{
		Interests: interests,
		Total:     len(interests),
	})
}

// GetCategories handles GET /v1/interests/categories
func (h *Handler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(ListCategoriesResponse{
		Categories: categories,
	})
}

// GetUserInterests handles GET /v1/users/:id/interests
func (h *Handler) GetUserInterests(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}
	
	interests, err := h.service.GetUserInterests(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(GetUserInterestsResponse{
		UserID:    userID,
		Interests: interests,
		Count:     len(interests),
	})
}

// AddUserInterests handles POST /v1/users/:id/interests
func (h *Handler) AddUserInterests(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}
	
	var req AddUserInterestsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	
	addedInterests, err := h.service.AddUserInterests(c.Context(), userID, req.InterestSlugs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	// Get total count after adding
	allInterests, _ := h.service.GetUserInterests(c.Context(), userID)
	
	return c.Status(fiber.StatusCreated).JSON(AddUserInterestsResponse{
		Message:        "interests added successfully",
		AddedInterests: addedInterests,
		TotalCount:     len(allInterests),
	})
}

// RemoveUserInterest handles DELETE /v1/users/:id/interests/:slug
func (h *Handler) RemoveUserInterest(c *fiber.Ctx) error {
	userID := c.Params("id")
	slug := c.Params("slug")
	
	if userID == "" || slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id and interest slug are required",
		})
	}
	
	if err := h.service.RemoveUserInterest(c.Context(), userID, slug); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(RemoveUserInterestResponse{
		Message: "interest removed successfully",
	})
}

// ReplaceUserInterests handles PUT /v1/users/:id/interests
func (h *Handler) ReplaceUserInterests(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}
	
	var req ReplaceUserInterestsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	
	interests, err := h.service.ReplaceUserInterests(c.Context(), userID, req.InterestSlugs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.JSON(ReplaceUserInterestsResponse{
		Message:   "interests updated successfully",
		Interests: interests,
		Count:     len(interests),
	})
}

// CreateInterest handles POST /v1/interests (admin only - add auth later)
func (h *Handler) CreateInterest(c *fiber.Ctx) error {
	var req CreateInterestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	
	interest, err := h.service.CreateInterest(c.Context(), req.Name, req.Slug, req.Category, req.Icon)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(CreateInterestResponse{
		Message:  "interest created successfully",
		Interest: *interest,
	})
}

