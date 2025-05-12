package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pollpulse/services/user-service/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, username, email, password_hash, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user roles
	roles, err := r.GetUserRoles(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user roles
	roles, err := r.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user roles
	roles, err := r.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

// ListUsers retrieves a list of users with pagination
func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, first_name, last_name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var users []*models.User
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get roles for each user
	for _, user := range users {
		roles, err := r.GetUserRoles(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		user.Roles = roles
	}

	return users, nil
}

// UpdateUser updates a user in the database
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		passwordHash,
		time.Now().UTC(),
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	// First delete from user_roles
	_, err := r.db.ExecContext(ctx, "DELETE FROM user_roles WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user roles: %w", err)
	}

	// Then delete the user
	_, err = r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// AddRole adds a role to a user
func (r *UserRepository) AddRole(ctx context.Context, userID, roleName string) error {
	// Get role ID by name
	var roleID string
	err := r.db.GetContext(ctx, &roleID, "SELECT id FROM roles WHERE name = $1", roleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("role not found: %s", roleName)
		}
		return fmt.Errorf("failed to get role ID: %w", err)
	}

	// Check if user already has this role
	var exists bool
	err = r.db.GetContext(
		ctx,
		&exists,
		"SELECT EXISTS(SELECT 1 FROM user_roles WHERE user_id = $1 AND role_id = $2)",
		userID,
		roleID,
	)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if exists {
		return nil // User already has this role
	}

	// Add role to user
	_, err = r.db.ExecContext(
		ctx,
		"INSERT INTO user_roles (user_id, role_id, created_at) VALUES ($1, $2, $3)",
		userID,
		roleID,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("failed to add role to user: %w", err)
	}

	return nil
}

// RemoveRole removes a role from a user
func (r *UserRepository) RemoveRole(ctx context.Context, userID, roleName string) error {
	// Get role ID by name
	var roleID string
	err := r.db.GetContext(ctx, &roleID, "SELECT id FROM roles WHERE name = $1", roleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("role not found: %s", roleName)
		}
		return fmt.Errorf("failed to get role ID: %w", err)
	}

	// Remove role from user
	_, err = r.db.ExecContext(
		ctx,
		"DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2",
		userID,
		roleID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// GetUserRoles gets the roles for a user
func (r *UserRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	var roles []string
	err := r.db.SelectContext(ctx, &roles, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}

// CreateRole creates a new role in the database
func (r *UserRepository) CreateRole(ctx context.Context, role *models.Role) error {
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	role.CreatedAt = now
	role.UpdatedAt = now

	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetRoles gets all roles
func (r *UserRepository) GetRoles(ctx context.Context) ([]*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
	`

	var roles []*models.Role
	err := r.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, nil
} 