package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/gorilla/mux"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "qa-service/internal/handlers"
    "qa-service/internal/models"
)

func main() {
    // Подключение к БД с повторными попытками
    var db *gorm.DB
    var err error
    
    dsn := "host=db user=postgres password=password dbname=qa_service port=5432 sslmode=disable"
    
    for i := 0; i < 10; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Printf("Failed to connect to database (attempt %d/10): %v", i+1, err)
            time.Sleep(2 * time.Second)
            continue
        }
        break
    }
    
    if err != nil {
        log.Fatal("Failed to connect to database after 10 attempts:", err)
    }
    
    log.Println("Successfully connected to database")
    
    // Автоматическое создание таблиц
    log.Println("Running database migrations...")
    err = db.AutoMigrate(&models.Question{}, &models.Answer{})
    if err != nil {
        log.Fatal("Failed to run migrations:", err)
    }
    log.Println("Migrations completed successfully")
    
    // Проверяем, что таблицы созданы
    checkTables(db)
    
    // Инициализация обработчиков
    questionHandler := &handlers.QuestionHandler{DB: db}
    answerHandler := &handlers.AnswerHandler{DB: db}
    
    r := mux.NewRouter()
    
    // Вопросы
    r.HandleFunc("/questions", questionHandler.GetQuestions).Methods("GET")
    r.HandleFunc("/questions", questionHandler.CreateQuestion).Methods("POST")
    r.HandleFunc("/questions/{id}", questionHandler.GetQuestion).Methods("GET")
    r.HandleFunc("/questions/{id}", questionHandler.DeleteQuestion).Methods("DELETE")
    
    // Ответы
    r.HandleFunc("/questions/{id}/answers", answerHandler.CreateAnswer).Methods("POST")
    r.HandleFunc("/answers/{id}", answerHandler.GetAnswer).Methods("GET")
    r.HandleFunc("/answers/{id}", answerHandler.DeleteAnswer).Methods("DELETE")
    
    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        // Проверяем подключение к БД
        sqlDB, err := db.DB()
        dbStatus := "connected"
        if err != nil {
            dbStatus = "error"
        } else if err = sqlDB.Ping(); err != nil {
            dbStatus = "disconnected"
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status":   "ok",
            "database": dbStatus,
        })
    }).Methods("GET")
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func checkTables(db *gorm.DB) {
    var tableExists bool
    
    // Проверяем таблицу questions
    db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'questions')").Scan(&tableExists)
    if tableExists {
        log.Println("✓ Table 'questions' exists")
        
        // Проверяем количество записей
        var count int64
        db.Model(&models.Question{}).Count(&count)
        log.Printf("  - Contains %d records", count)
    } else {
        log.Println("✗ Table 'questions' does not exist")
    }
    
    // Проверяем таблицу answers
    db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'answers')").Scan(&tableExists)
    if tableExists {
        log.Println("✓ Table 'answers' exists")
        
        // Проверяем количество записей
        var count int64
        db.Model(&models.Answer{}).Count(&count)
        log.Printf("  - Contains %d records", count)
    } else {
        log.Println("✗ Table 'answers' does not exist")
    }
}