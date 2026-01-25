package api

import (
	"net/http"
	"time"

	apperrors "github.com/waynenilsen/power-pro-v3/internal/errors"
	"github.com/waynenilsen/power-pro-v3/internal/middleware"
	"github.com/waynenilsen/power-pro-v3/internal/profile"
)

// ProfileHandler handles HTTP requests for user profile operations.
type ProfileHandler struct {
	profileService *profile.Service
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(profileService *profile.Service) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// UpdateProfileRequest represents the request body for updating a profile.
type UpdateProfileRequest struct {
	Name       *string `json:"name,omitempty"`
	WeightUnit *string `json:"weightUnit,omitempty"`
}

// ProfileResponse represents the response for profile operations.
type ProfileResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Name       *string   `json:"name"`
	WeightUnit string    `json:"weightUnit"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Get handles GET /users/{userId}/profile
func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: only the user themselves or an admin can view profile
	authUserID := middleware.GetUserID(r)
	isAdmin := middleware.IsAdmin(r)
	if authUserID != userID && !isAdmin {
		writeDomainError(w, apperrors.NewForbidden("you can only access your own profile"))
		return
	}

	// Get profile from service
	p, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	// Build and return response
	response := h.buildProfileResponse(p)
	writeData(w, http.StatusOK, response)
}

// Update handles PUT /users/{userId}/profile
func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		writeDomainError(w, apperrors.NewBadRequest("missing user ID"))
		return
	}

	// Authorization check: only the profile owner can update (not even admins)
	authUserID := middleware.GetUserID(r)
	if authUserID != userID {
		writeDomainError(w, apperrors.NewForbidden("profile updates are owner-only"))
		return
	}

	// Parse request body
	var req UpdateProfileRequest
	if err := readJSON(r, &req); err != nil {
		writeDomainError(w, apperrors.NewBadRequest("invalid request body"))
		return
	}

	// Convert to service request
	serviceReq := profile.UpdateProfileRequest{
		Name:       req.Name,
		WeightUnit: req.WeightUnit,
	}

	// Update profile via service
	p, err := h.profileService.UpdateProfile(r.Context(), userID, serviceReq)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	// Build and return response
	response := h.buildProfileResponse(p)
	writeData(w, http.StatusOK, response)
}

// buildProfileResponse builds the ProfileResponse from a profile.
func (h *ProfileHandler) buildProfileResponse(p *profile.Profile) ProfileResponse {
	return ProfileResponse{
		ID:         p.ID,
		Email:      p.Email,
		Name:       p.Name,
		WeightUnit: p.WeightUnit,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}
