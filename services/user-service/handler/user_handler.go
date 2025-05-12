package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/pollpulse/pkg/common/errors"
	"github.com/pollpulse/pkg/common/logging"
	"github.com/pollpulse/pkg/common/middleware"
	"github.com/pollpulse/services/user-service/models"
	"github.com/pollpulse/services/user-service/repository"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo      *repository.UserRepository
	validate  *validator.Validate
	logger    *logging.Logger
	jwtSecret string
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.UserRepository, logger *logging.Logger, jwtSecret string) *UserHandler {
	return &UserHandler{
		repo:      repo,
		validate:  validator.New(),
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

// RegisterRoutes registers the routes for the user handler
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtSecret))
		r.Get("/users", h.ListUsers)
		r.Get("/users/{id}", h.GetUser)
		r.Put("/users/{id}", h.UpdateUser)
		r.Delete("/users/{id}", h.DeleteUser)
		r.Get("/users/me", h.GetCurrentUser)
		r.Put("/users/me/password", h.UpdatePassword)
		r.Post("/users/{id}/roles", h.AddRole)
		r.Delete("/users/{id}/roles/{role}", h.RemoveRole)
		r.Get("/roles", h.GetRoles)
		r.Post("/roles", h.CreateRole)
	})
}

// RegisterUser registers a new user
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	// Check if user exists
	existingUser, _ := h.repo.GetUserByUsername(r.Context(), req.Username)
	if existingUser != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Username already exists")
		return
	}

	existingEmail, _ := h.repo.GetUserByEmail(r.Context(), req.Email)
	if existingEmail != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Email already exists")
		return
	}

	// Hash the password
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	// Create the user
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.repo.CreateUser(r.Context(), user); err != nil {
		h.logger.Error("Failed to create user", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	// Add default role
	if err := h.repo.AddRole(r.Context(), user.ID, "user"); err != nil {
		h.logger.Error("Failed to add default role", "error", err)
		// Continue anyway, user was created
	}

	// Add any additional roles
	for _, role := range req.Roles {
		if err := h.repo.AddRole(r.Context(), user.ID, role); err != nil {
			h.logger.Error("Failed to add role", "role", role, "error", err)
			// Continue anyway
		}
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username, user.Email, []string{"user"}, h.jwtSecret, 24*time.Hour)
	if err != nil {
		h.logger.Error("Failed to generate token", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	// Get user with roles
	user, err = h.repo.GetUserByID(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to get user with roles", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login logs in a user
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	// Get the user
	user, err := h.repo.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		errors.HandleError(w, errors.ErrUnauthorized, "Invalid username or password")
		return
	}

	// Check the password
	if !models.CheckPasswordHash(req.Password, user.Password) {
		errors.HandleError(w, errors.ErrUnauthorized, "Invalid username or password")
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username, user.Email, user.Roles, h.jwtSecret, 24*time.Hour)
	if err != nil {
		h.logger.Error("Failed to generate token", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListUsers lists all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	if !middleware.CheckRole(r.Context(), "admin") {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Parse pagination parameters
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil {
			limit = limitVal
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offsetVal, err := strconv.Atoi(offsetStr); err == nil {
			offset = offsetVal
		}
	}

	// Get users
	users, err := h.repo.ListUsers(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list users", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	// Convert to response objects
	response := make([]models.UserResponse, len(users))
	for i, user := range users {
		response[i] = user.ToResponse()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, errors.ErrUnauthorized, "")
		return
	}

	// Get the user ID from the URL
	userID := chi.URLParam(r, "id")

	// Check if user has admin role or is the same user
	if !middleware.CheckRole(r.Context(), "admin") && currentUser.UserID != userID {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Get the user
	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		errors.HandleError(w, errors.ErrNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

// GetCurrentUser gets the current user
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userClaims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, errors.ErrUnauthorized, "")
		return
	}

	// Get the user from the database
	user, err := h.repo.GetUserByID(r.Context(), userClaims.UserID)
	if err != nil {
		errors.HandleError(w, errors.ErrNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUser, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, errors.ErrUnauthorized, "")
		return
	}

	// Get the user ID from the URL
	userID := chi.URLParam(r, "id")

	// Check if user has admin role or is the same user
	if !middleware.CheckRole(r.Context(), "admin") && currentUser.UserID != userID {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Parse the request body
	var req struct {
		Username  string `json:"username" validate:"required,min=3,max=30"`
		Email     string `json:"email" validate:"required,email"`
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	// Get the user
	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		errors.HandleError(w, errors.ErrNotFound, "User not found")
		return
	}

	// Check if username or email is changing
	if user.Username != req.Username {
		// Check if new username exists
		existingUser, _ := h.repo.GetUserByUsername(r.Context(), req.Username)
		if existingUser != nil {
			errors.HandleError(w, errors.ErrBadRequest, "Username already exists")
			return
		}
	}

	if user.Email != req.Email {
		// Check if new email exists
		existingUser, _ := h.repo.GetUserByEmail(r.Context(), req.Email)
		if existingUser != nil {
			errors.HandleError(w, errors.ErrBadRequest, "Email already exists")
			return
		}
	}

	// Update the user
	user.Username = req.Username
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := h.repo.UpdateUser(r.Context(), user); err != nil {
		h.logger.Error("Failed to update user", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

// UpdatePassword updates a user's password
func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userClaims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		errors.HandleError(w, errors.ErrUnauthorized, "")
		return
	}

	// Parse the request body
	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	// Get the user
	user, err := h.repo.GetUserByID(r.Context(), userClaims.UserID)
	if err != nil {
		errors.HandleError(w, errors.ErrNotFound, "User not found")
		return
	}

	// Check the current password
	if !models.CheckPasswordHash(req.CurrentPassword, user.Password) {
		errors.HandleError(w, errors.ErrBadRequest, "Current password is incorrect")
		return
	}

	// Hash the new password
	hashedPassword, err := models.HashPassword(req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to hash password", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	// Update the password
	if err := h.repo.UpdatePassword(r.Context(), user.ID, hashedPassword); err != nil {
		h.logger.Error("Failed to update password", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	if !middleware.CheckRole(r.Context(), "admin") {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Get the user ID from the URL
	userID := chi.URLParam(r, "id")

	// Delete the user
	if err := h.repo.DeleteUser(r.Context(), userID); err != nil {
		h.logger.Error("Failed to delete user", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddRole adds a role to a user
func (h *UserHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	if !middleware.CheckRole(r.Context(), "admin") {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Get the user ID from the URL
	userID := chi.URLParam(r, "id")

	// Parse the request body
	var req struct {
		Role string `json:"role" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	// Add the role
	if err := h.repo.AddRole(r.Context(), userID, req.Role); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveRole removes a role from a user
func (h *UserHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	if !middleware.CheckRole(r.Context(), "admin") {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Get the user ID and role from the URL
	userID := chi.URLParam(r, "id")
	role := chi.URLParam(r, "role")

	// Remove the role
	if err := h.repo.RemoveRole(r.Context(), userID, role); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetRoles gets all roles
func (h *UserHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	if !middleware.CheckRole(r.Context(), "admin") {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Get all roles
	roles, err := h.repo.GetRoles(r.Context())
	if err != nil {
		h.logger.Error("Failed to get roles", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// CreateRole creates a new role
func (h *UserHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Check if user has admin role
	if !middleware.CheckRole(r.Context(), "admin") {
		errors.HandleError(w, errors.ErrForbidden, "")
		return
	}

	// Parse the request body
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	if err := h.validate.Struct(req); err != nil {
		errors.HandleError(w, errors.ErrBadRequest, err.Error())
		return
	}

	// Create the role
	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.repo.CreateRole(r.Context(), role); err != nil {
		h.logger.Error("Failed to create role", "error", err)
		errors.HandleError(w, errors.ErrInternalServer, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}
