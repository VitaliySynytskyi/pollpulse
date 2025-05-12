package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitaliySynytskyi/pollpulse/services/survey-service/models"
	"github.com/VitaliySynytskyi/pollpulse/services/survey-service/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SurveyHandler struct {
	repo *repository.SurveyRepository
}

func NewSurveyHandler(repo *repository.SurveyRepository) *SurveyHandler {
	return &SurveyHandler{repo: repo}
}

// CreateSurvey handles the creation of a new survey
func (h *SurveyHandler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var survey models.Survey
	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	survey.CreatedBy = userID

	if err := h.repo.CreateSurvey(r.Context(), &survey); err != nil {
		http.Error(w, "Failed to create survey", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(survey)
}

// GetSurvey handles retrieving a survey by ID
func (h *SurveyHandler) GetSurvey(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	survey, err := h.repo.GetSurvey(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get survey", http.StatusInternalServerError)
		return
	}
	if survey == nil {
		http.Error(w, "Survey not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey)
}

// UpdateSurvey handles updating an existing survey
func (h *SurveyHandler) UpdateSurvey(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	var survey models.Survey
	if err := json.NewDecoder(r.Body).Decode(&survey); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	survey.ID = id

	if err := h.repo.UpdateSurvey(r.Context(), &survey); err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "Survey not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update survey", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSurvey handles deleting a survey
func (h *SurveyHandler) DeleteSurvey(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteSurvey(r.Context(), id); err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, "Survey not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete survey", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListSurveys handles retrieving a list of surveys with pagination
func (h *SurveyHandler) ListSurveys(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	surveys, err := h.repo.ListSurveys(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, "Failed to list surveys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveys)
}
