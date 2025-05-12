package models

import (
	"time"
)

// Response represents a survey response from a respondent
type Response struct {
	ID           string     `json:"id" db:"id"`
	SurveyID     string     `json:"survey_id" db:"survey_id"`
	RespondentID *string    `json:"respondent_id,omitempty" db:"respondent_id"` // Can be NULL for anonymous responses
	StartedAt    time.Time  `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"` // NULL until completed
	IPAddress    string     `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    string     `json:"user_agent,omitempty" db:"user_agent"`
	Answers      []Answer   `json:"answers,omitempty" db:"-"` // Handled separately
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Answer represents an answer to a specific question
type Answer struct {
	ID         string    `json:"id" db:"id"`
	ResponseID string    `json:"response_id" db:"response_id"`
	SurveyID   string    `json:"survey_id" db:"survey_id"`
	QuestionID string    `json:"question_id" db:"question_id"`
	OptionID   *string   `json:"option_id,omitempty" db:"option_id"`     // NULL for text answers
	TextAnswer *string   `json:"text_answer,omitempty" db:"text_answer"` // NULL for choice answers
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// SurveyResult represents the aggregated results of a survey
type SurveyResult struct {
	SurveyID       string           `json:"survey_id"`
	ResponseCount  int              `json:"response_count"`
	CompletionRate float64          `json:"completion_rate"` // Percentage of completed responses
	Questions      []QuestionResult `json:"questions"`
}

// QuestionResult represents the results for a specific question
type QuestionResult struct {
	QuestionID    string                 `json:"question_id"`
	QuestionText  string                 `json:"question_text"`
	QuestionType  string                 `json:"question_type"`
	ResponseCount int                    `json:"response_count"`
	Options       []OptionResult         `json:"options,omitempty"`      // For choice questions
	TextAnswers   []string               `json:"text_answers,omitempty"` // For text questions
	Statistics    map[string]interface{} `json:"statistics,omitempty"`   // For rating questions
}

// OptionResult represents the results for a specific option
type OptionResult struct {
	OptionID   string  `json:"option_id"`
	OptionText string  `json:"option_text"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// SubmitResponseRequest represents the request to submit a response to a survey
type SubmitResponseRequest struct {
	SurveyID     string                `json:"survey_id" validate:"required"`
	RespondentID *string               `json:"respondent_id,omitempty"`
	Answers      []SubmitAnswerRequest `json:"answers" validate:"required,min=1,dive"`
	IPAddress    string                `json:"ip_address,omitempty"`
	UserAgent    string                `json:"user_agent,omitempty"`
}

// SubmitAnswerRequest represents the request to submit an answer to a question
type SubmitAnswerRequest struct {
	QuestionID string  `json:"question_id" validate:"required"`
	OptionID   *string `json:"option_id,omitempty"`
	TextAnswer *string `json:"text_answer,omitempty"`
}

// ResponseSummary represents a summary of a response
type ResponseSummary struct {
	ID           string     `json:"id"`
	SurveyID     string     `json:"survey_id"`
	RespondentID *string    `json:"respondent_id,omitempty"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	AnswerCount  int        `json:"answer_count"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ToSummary converts a Response to a ResponseSummary
func (r *Response) ToSummary() ResponseSummary {
	return ResponseSummary{
		ID:           r.ID,
		SurveyID:     r.SurveyID,
		RespondentID: r.RespondentID,
		StartedAt:    r.StartedAt,
		CompletedAt:  r.CompletedAt,
		AnswerCount:  len(r.Answers),
		CreatedAt:    r.CreatedAt,
	}
}

// ExportFormat represents the format for exporting survey results
type ExportFormat string

// Export formats
const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
	ExportFormatXLSX ExportFormat = "xlsx"
)

// ExportRequest represents the request to export survey results
type ExportRequest struct {
	SurveyID   string       `json:"survey_id" validate:"required"`
	Format     ExportFormat `json:"format" validate:"required,oneof=csv json xlsx"`
	IncludeRaw bool         `json:"include_raw"` // Whether to include raw responses
}

// TimePeriod represents a time period for analytics
type TimePeriod string

// Time periods
const (
	TimePeriodDay   TimePeriod = "day"
	TimePeriodWeek  TimePeriod = "week"
	TimePeriodMonth TimePeriod = "month"
	TimePeriodYear  TimePeriod = "year"
	TimePeriodAll   TimePeriod = "all"
)

// AnalyticsRequest represents the request for survey analytics
type AnalyticsRequest struct {
	SurveyID  string     `json:"survey_id" validate:"required"`
	Period    TimePeriod `json:"period" validate:"required,oneof=day week month year all"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

// ResponseTrend represents a trend in responses over time
type ResponseTrend struct {
	Period string  `json:"period"`
	Count  int     `json:"count"`
	Change float64 `json:"change"` // Percentage change from previous period
}

// Analytics represents analytics for a survey
type Analytics struct {
	SurveyID              string                 `json:"survey_id"`
	TotalResponses        int                    `json:"total_responses"`
	CompletionRate        float64                `json:"completion_rate"`
	AverageTimeToComplete float64                `json:"average_time_to_complete"` // In seconds
	ResponseTrends        []ResponseTrend        `json:"response_trends"`
	QuestionInsights      []QuestionInsight      `json:"question_insights"`
	Demographics          map[string]interface{} `json:"demographics,omitempty"` // If demographic data is available
}

// QuestionInsight represents insights for a specific question
type QuestionInsight struct {
	QuestionID   string                   `json:"question_id"`
	QuestionText string                   `json:"question_text"`
	SkipRate     float64                  `json:"skip_rate"`              // Percentage of respondents who skipped this question
	TopAnswers   []interface{}            `json:"top_answers,omitempty"`  // Most common answers
	Correlations []map[string]interface{} `json:"correlations,omitempty"` // Correlations with other questions
}
