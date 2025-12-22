package location

import (
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for location endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new location handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAllLocations handles GET /v1/locations
func (h *Handler) GetAllLocations(c *fiber.Ctx) error {
	country := c.Query("country")

	var locations []Location
	var err error

	if country != "" {
		locations, err = h.service.GetLocationsByCountry(c.Context(), country)
	} else {
		locations, err = h.service.GetAllLocations(c.Context())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ListLocationsResponse{
		Locations: locations,
		Total:     len(locations),
	})
}

// SearchLocations handles GET /v1/locations/search
func (h *Handler) SearchLocations(c *fiber.Ctx) error {
	query := c.Query("q")

	locations, err := h.service.SearchLocations(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(SearchLocationsResponse{
		Locations: locations,
		Total:     len(locations),
	})
}

// GetPopularLocations handles GET /v1/locations/popular
func (h *Handler) GetPopularLocations(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)

	locations, err := h.service.GetMostUsedLocations(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ListLocationsResponse{
		Locations: locations,
		Total:     len(locations),
	})
}

// GetLocation handles GET /v1/locations/:id
func (h *Handler) GetLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "location id is required",
		})
	}

	location, err := h.service.GetLocationByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(GetLocationResponse{
		Location: *location,
	})
}

// CreateLocation handles POST /v1/locations
func (h *Handler) CreateLocation(c *fiber.Ctx) error {
	var req CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	location, err := h.service.CreateLocation(c.Context(), req.City, req.Country)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(CreateLocationResponse{
		Message:  "location created successfully",
		Location: *location,
	})
}

// UpdateLocation handles PUT /v1/locations/:id
func (h *Handler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "location id is required",
		})
	}

	var req UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	location, err := h.service.UpdateLocation(c.Context(), id, req.City, req.Country)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(UpdateLocationResponse{
		Message:  "location updated successfully",
		Location: *location,
	})
}

// DeleteLocation handles DELETE /v1/locations/:id
func (h *Handler) DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "location id is required",
		})
	}

	if err := h.service.DeleteLocation(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(DeleteLocationResponse{
		Message: "location deleted successfully",
	})
}

// IncrementUsage handles POST /v1/locations/:id/increment
func (h *Handler) IncrementUsage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "location id is required",
		})
	}

	count, err := h.service.IncrementUsage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(IncrementUsageResponse{
		Message:     "usage incremented successfully",
		UsedByCount: count,
	})
}

// DecrementUsage handles POST /v1/locations/:id/decrement
func (h *Handler) DecrementUsage(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "location id is required",
		})
	}

	count, err := h.service.DecrementUsage(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(DecrementUsageResponse{
		Message:     "usage decremented successfully",
		UsedByCount: count,
	})
}
