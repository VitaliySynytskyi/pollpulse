package models

import (
	"time"

	"github.com/google/uuid"
)

// SurveyStatus represents the status of a survey
type SurveyStatus string

// Survey statuses
const (
	SurveyStatusDraft     SurveyStatus = "draft"
	SurveyStatusPublished SurveyStatus = "published"
	SurveyStatusClosed    SurveyStatus = "closed"
)

// QuestionType represents the type of a question
type QuestionType string

// Question types
const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeSingleChoice   QuestionType = "single_choice"
	QuestionTypeText           QuestionType = "text"
	QuestionTypeRating         QuestionType = "rating"
	QuestionTypeDate           QuestionType = "date"
)

// Survey represents a survey in the system
type Survey struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	CreatedBy   uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	Questions   []Question `json:"questions,omitempty" db:"-"`
}

// SurveyQuestion represents a question in a survey
type SurveyQuestion struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SurveyID  uuid.UUID `json:"survey_id" db:"survey_id"`
	Question  string    `json:"question" db:"question"`
	Type      string    `json:"type" db:"type"` // e.g., "multiple_choice", "text", "rating"
	Required  bool      `json:"required" db:"required"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SurveyOption represents an option for a multiple choice question
type SurveyOption struct {
	ID         uuid.UUID `json:"id" db:"id"`
	QuestionID uuid.UUID `json:"question_id" db:"question_id"`
	OptionText string    `json:"option_text" db:"option_text"`
	Order      int       `json:"order" db:"order"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Question represents a question in a survey
type Question struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SurveyID  uuid.UUID `json:"survey_id" db:"survey_id"`
	Text      string    `json:"text" db:"question"`
	Type      string    `json:"type" db:"type"`
	Required  bool      `json:"required" db:"required"`
	Order     int       `json:"order" db:"order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Options   []Option  `json:"options,omitempty" db:"-"`
}

// Option represents an option for a multiple choice question
type Option struct {
	ID         uuid.UUID `json:"id" db:"id"`
	QuestionID uuid.UUID `json:"question_id" db:"question_id"`
	Text       string    `json:"text" db:"option_text"`
	Order      int       `json:"order" db:"order"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateSurveyRequest represents the request to create a new survey
type CreateSurveyRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Questions   []Question `json:"questions" validate:"required,min=1,dive"`
}

// UpdateSurveyRequest represents the request to update an existing survey
type UpdateSurveyRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Questions   []Question `json:"questions" validate:"required,min=1,dive"`
	IsActive    bool       `json:"is_active"`
}

// SurveyResponse represents the response to a survey request
type SurveyResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedBy   string     `json:"created_by"`
	IsActive    bool       `json:"is_active"`
	Questions   []Question `json:"questions"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SurveySummary represents a summary of a survey
type SurveySummary struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CreatedBy     string    `json:"created_by"`
	IsActive      bool      `json:"is_active"`
	QuestionCount int       `json:"question_count"`
	ResponseCount int       `json:"response_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse converts a Survey to a SurveyResponse
func (s *Survey) ToResponse() SurveyResponse {
	return SurveyResponse{
		ID:          s.ID.String(),
		Title:       s.Title,
		Description: s.Description,
		CreatedBy:   s.CreatedBy.String(),
		IsActive:    s.IsActive,
		Questions:   s.Questions,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// ToSummary converts a Survey to a SurveySummary
func (s *Survey) ToSummary(responseCount int) SurveySummary {
	return SurveySummary{
		ID:            s.ID.String(),
		Title:         s.Title,
		Description:   s.Description,
		CreatedBy:     s.CreatedBy.String(),
		IsActive:      s.IsActive,
		QuestionCount: len(s.Questions),
		ResponseCount: responseCount,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
