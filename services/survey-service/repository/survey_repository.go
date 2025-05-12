package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/VitaliySynytskyi/pollpulse/services/survey-service/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("not found")
)

// SurveyRepository handles database operations for surveys
type SurveyRepository struct {
	db *sqlx.DB
}

// NewSurveyRepository creates a new survey repository
func NewSurveyRepository(db *sqlx.DB) *SurveyRepository {
	return &SurveyRepository{
		db: db,
	}
}

// CreateSurvey creates a new survey in the database
func (r *SurveyRepository) CreateSurvey(ctx context.Context, survey *models.Survey) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Generate ID if not provided
	if survey.ID == uuid.Nil {
		survey.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now().UTC()
	survey.CreatedAt = now
	survey.UpdatedAt = now

	// Insert survey
	query := `
		INSERT INTO surveys (id, title, description, created_by, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		survey.ID,
		survey.Title,
		survey.Description,
		survey.CreatedBy,
		survey.CreatedAt,
		survey.UpdatedAt,
		survey.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create survey: %w", err)
	}

	// Insert questions
	for i, question := range survey.Questions {
		if question.ID == uuid.Nil {
			question.ID = uuid.New()
		}

		question.SurveyID = survey.ID
		question.CreatedAt = now
		question.UpdatedAt = now
		question.Order = i + 1 // Set order based on index

		query := `
			INSERT INTO survey_questions (id, survey_id, question, type, required, "order", created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err = tx.ExecContext(
			ctx,
			query,
			question.ID,
			question.SurveyID,
			question.Text,
			question.Type,
			question.Required,
			question.Order,
			question.CreatedAt,
			question.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create question: %w", err)
		}

		// Insert options for choice questions
		if question.Type == "multiple_choice" || question.Type == "single_choice" {
			for j, option := range question.Options {
				if option.ID == uuid.Nil {
					option.ID = uuid.New()
				}

				option.QuestionID = question.ID
				option.CreatedAt = now
				option.UpdatedAt = now
				option.Order = j + 1 // Set order based on index

				query := `
					INSERT INTO survey_options (id, question_id, option_text, "order", created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6)
				`

				_, err = tx.ExecContext(
					ctx,
					query,
					option.ID,
					option.QuestionID,
					option.Text,
					option.Order,
					option.CreatedAt,
					option.UpdatedAt,
				)

				if err != nil {
					return fmt.Errorf("failed to create option: %w", err)
				}
			}
		}

		// Update the question in the survey object
		survey.Questions[i] = question
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetSurvey retrieves a survey by ID
func (r *SurveyRepository) GetSurvey(ctx context.Context, id uuid.UUID) (*models.Survey, error) {
	// Get the survey
	query := `
		SELECT id, title, description, created_by, created_at, updated_at, is_active
		FROM surveys
		WHERE id = $1
	`

	var survey models.Survey
	err := r.db.GetContext(ctx, &survey, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get survey: %w", err)
	}

	// Get the questions
	questionsQuery := `
		SELECT id, survey_id, question as text, type, required, "order", created_at, updated_at
		FROM survey_questions
		WHERE survey_id = $1
		ORDER BY "order"
	`

	var questions []models.Question
	err = r.db.SelectContext(ctx, &questions, questionsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	// Get the options for each question
	for i, question := range questions {
		if question.Type == "multiple_choice" || question.Type == "single_choice" {
			optionsQuery := `
				SELECT id, question_id, option_text as text, "order", created_at, updated_at
				FROM survey_options
				WHERE question_id = $1
				ORDER BY "order"
			`

			var options []models.Option
			err = r.db.SelectContext(ctx, &options, optionsQuery, question.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get options: %w", err)
			}

			questions[i].Options = options
		}
	}

	survey.Questions = questions

	return &survey, nil
}

// GetSurveysByUserID retrieves all surveys created by a user
func (r *SurveyRepository) GetSurveysByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Survey, error) {
	query := `
		SELECT id, title, description, user_id, status, start_date, end_date, created_at, updated_at
		FROM surveys
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var surveys []*models.Survey
	err := r.db.SelectContext(ctx, &surveys, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get surveys: %w", err)
	}

	// For each survey, get the number of questions
	for i, survey := range surveys {
		countQuery := `
			SELECT COUNT(*) FROM questions WHERE survey_id = $1
		`

		var count int
		err = r.db.GetContext(ctx, &count, countQuery, survey.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get question count: %w", err)
		}

		// Create a slice with the capacity for the questions (they'll be loaded on demand)
		surveys[i].Questions = make([]models.Question, 0, count)
	}

	return surveys, nil
}

// UpdateSurvey updates a survey
func (r *SurveyRepository) UpdateSurvey(ctx context.Context, survey *models.Survey) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update timestamp
	survey.UpdatedAt = time.Now().UTC()

	// Update survey
	query := `
		UPDATE surveys
		SET title = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		survey.Title,
		survey.Description,
		survey.IsActive,
		survey.UpdatedAt,
		survey.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update survey: %w", err)
	}

	// Delete existing questions and options
	_, err = tx.ExecContext(ctx, "DELETE FROM survey_questions WHERE survey_id = $1", survey.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing questions: %w", err)
	}

	// Insert updated questions and options
	for i, question := range survey.Questions {
		if question.ID == uuid.Nil {
			question.ID = uuid.New()
		}

		question.SurveyID = survey.ID
		question.UpdatedAt = survey.UpdatedAt
		question.Order = i + 1 // Set order based on index

		query := `
			INSERT INTO survey_questions (id, survey_id, question, type, required, "order", created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err = tx.ExecContext(
			ctx,
			query,
			question.ID,
			question.SurveyID,
			question.Text,
			question.Type,
			question.Required,
			question.Order,
			question.CreatedAt,
			survey.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to create question: %w", err)
		}

		// Insert options for choice questions
		if question.Type == "multiple_choice" || question.Type == "single_choice" {
			for j, option := range question.Options {
				if option.ID == uuid.Nil {
					option.ID = uuid.New()
				}

				option.QuestionID = question.ID
				option.UpdatedAt = survey.UpdatedAt
				option.Order = j + 1 // Set order based on index

				query := `
					INSERT INTO survey_options (id, question_id, option_text, "order", created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6)
				`

				_, err = tx.ExecContext(
					ctx,
					query,
					option.ID,
					option.QuestionID,
					option.Text,
					option.Order,
					option.CreatedAt,
					survey.UpdatedAt,
				)

				if err != nil {
					return fmt.Errorf("failed to create option: %w", err)
				}
			}
		}

		// Update the question in the survey object
		survey.Questions[i] = question
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteSurvey deletes a survey
func (r *SurveyRepository) DeleteSurvey(ctx context.Context, id uuid.UUID) error {
	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete survey - questions and options will be deleted by foreign key cascade
	_, err = tx.ExecContext(ctx, "DELETE FROM surveys WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete survey: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateSurveyStatus updates a survey's status
func (r *SurveyRepository) UpdateSurveyStatus(ctx context.Context, id string, status models.SurveyStatus) error {
	query := `
		UPDATE surveys
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update survey status: %w", err)
	}

	return nil
}

// ListSurveys lists all surveys with pagination
func (r *SurveyRepository) ListSurveys(ctx context.Context, offset, limit int) ([]*models.Survey, error) {
	query := `
		SELECT id, title, description, created_by, created_at, updated_at, is_active
		FROM surveys
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var surveys []*models.Survey
	err := r.db.SelectContext(ctx, &surveys, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list surveys: %w", err)
	}

	// For each survey, get the number of questions
	for i, survey := range surveys {
		countQuery := `
			SELECT COUNT(*) FROM survey_questions WHERE survey_id = $1
		`

		var count int
		err = r.db.GetContext(ctx, &count, countQuery, survey.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get question count: %w", err)
		}

		// Create a slice with the capacity for the questions (they'll be loaded on demand)
		surveys[i].Questions = make([]models.Question, 0, count)
	}

	return surveys, nil
}

// GetPublicSurveys gets all published surveys
func (r *SurveyRepository) GetPublicSurveys(ctx context.Context, limit, offset int) ([]*models.Survey, error) {
	query := `
		SELECT id, title, description, user_id, status, start_date, end_date, created_at, updated_at
		FROM surveys
		WHERE status = $1
		AND (start_date IS NULL OR start_date <= NOW())
		AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var surveys []*models.Survey
	err := r.db.SelectContext(ctx, &surveys, query, models.SurveyStatusPublished, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get public surveys: %w", err)
	}

	// For each survey, get the number of questions
	for i, survey := range surveys {
		countQuery := `
			SELECT COUNT(*) FROM questions WHERE survey_id = $1
		`

		var count int
		err = r.db.GetContext(ctx, &count, countQuery, survey.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get question count: %w", err)
		}

		// Create a slice with the capacity for the questions (they'll be loaded on demand)
		surveys[i].Questions = make([]models.Question, 0, count)
	}

	return surveys, nil
}

// CountSurveyResponses counts the number of responses to a survey
func (r *SurveyRepository) CountSurveyResponses(ctx context.Context, surveyID string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT response_id) FROM responses WHERE survey_id = $1
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, surveyID)
	if err != nil {
		return 0, fmt.Errorf("failed to count survey responses: %w", err)
	}

	return count, nil
}
