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

type QuestionHandler struct {
    DB *gorm.DB
}

// GET /questions/ - получить все вопросы
func (h *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
    var questions []models.Question
    if result := h.DB.Find(&questions); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(questions)
}

// POST /questions/ - создать новый вопрос
func (h *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
    var question models.Question
    if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Устанавливаем время создания
    question.CreatedAt = time.Now()
    
    if result := h.DB.Create(&question); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(question)
}

// GET /questions/{id} - получить вопрос по ID с ответами
func (h *QuestionHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid question ID", http.StatusBadRequest)
        return
    }
    
    var question models.Question
    // Загружаем вопрос вместе с ответами
    if result := h.DB.Preload("Answers").First(&question, id); result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            http.Error(w, "Question not found", http.StatusNotFound)
        } else {
            http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        }
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(question)
}

// DELETE /questions/{id} - удалить вопрос
func (h *QuestionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid question ID", http.StatusBadRequest)
        return
    }
    
    if result := h.DB.Delete(&models.Question{}, id); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}