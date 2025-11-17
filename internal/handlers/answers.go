package handlers

import (
    "encoding/json"
    "net/http"
    "qa-service/internal/models"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

type AnswerHandler struct {
    DB *gorm.DB
}

// POST /questions/{id}/answers/
func (h *AnswerHandler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    questionID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid question ID", http.StatusBadRequest)
        return
    }
    // Проверяем, существует ли вопрос
    var question models.Question
    if result := h.DB.First(&question, questionID); result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Question not found", http.StatusNotFound)
        } else {
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }

    var answer models.Answer
    if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Устанавливаем question_id и время создания
    answer.QuestionID = uint(questionID)
    answer.CreatedAt = time.Now()

    if result := h.DB.Create(&answer); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(answer)
}

// GET /answers/{id}
func (h *AnswerHandler) GetAnswer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid answer ID", http.StatusBadRequest)
        return
    }
    var answer models.Answer
    if result := h.DB.First(&answer, id); result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Answer not found", http.StatusNotFound)
        } else {
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(answer)
}

// DELETE /answers/{id}
func (h *AnswerHandler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid answer ID", http.StatusBadRequest)
        return
    }
    if result := h.DB.Delete(&models.Answer{}, id); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}