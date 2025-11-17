package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/gorm"
    
    "qa-service/internal/models"
)

// MockGormDB создает mock для gorm.DB
type MockGormDB struct {
    mock.Mock
}

// Эти методы должны имитировать методы gorm.DB и возвращать *gorm.DB

func (m *MockGormDB) Create(value interface{}) *gorm.DB {
    args := m.Called(value)
    // Возвращаем фиктивный *gorm.DB с Error полем
    db := &gorm.DB{}
    if args.Get(0) != nil {
        // Можно установить ошибку если нужно
        // db.Error = args.Error(0)
    }
    return db
}

func (m *MockGormDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
    args := m.Called(dest, conds)
    db := &gorm.DB{}
    if args.Get(0) != nil {
        // Симулируем заполнение данных
        if question, ok := dest.(*models.Question); ok {
            question.ID = 1
            question.Text = "Test question"
        }
    }
    return db
}

func (m *MockGormDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
    args := m.Called(dest, conds)
    db := &gorm.DB{}
    
    // Симулируем заполнение списка вопросов
    if questions, ok := dest.(*[]models.Question); ok {
        *questions = []models.Question{
            {ID: 1, Text: "Question 1"},
            {ID: 2, Text: "Question 2"},
        }
    }
    
    return db
}

func (m *MockGormDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
    args := m.Called(value, conds)
    return &gorm.DB{}
}

func (m *MockGormDB) Preload(query string, args ...interface{}) *gorm.DB {
    m.Called(query, args)
    return m // Теперь возвращаем *MockGormDB, который должен имитировать *gorm.DB
}

// Реализуем минимальный набор методов gorm.DB для MockGormDB
func (m *MockGormDB) Error() error {
    args := m.Called()
    if err, ok := args.Get(0).(error); ok {
        return err
    }
    return nil
}

// Добавляем другие необходимые методы gorm.DB
func (m *MockGormDB) Model(value interface{}) *gorm.DB {
    return m
}

func (m *MockGormDB) Where(query interface{}, args ...interface{}) *gorm.DB {
    return m
}

func (m *MockGormDB) Scan(dest interface{}) *gorm.DB {
    return m
}

func TestCreateQuestion(t *testing.T) {
    // Создаем mock базы данных
    mockDB := new(MockGormDB)
    
    // Создаем обработчик с mock БД
    handler := &QuestionHandler{DB: mockDB}
    
    // Создаем тестовый вопрос
    question := models.Question{Text: "Test question"}
    body, _ := json.Marshal(question)
    
    // Настраиваем ожидания для mock
    mockDB.On("Create", &question).Return(&gorm.DB{})
    
    // Создаем запрос
    req, err := http.NewRequest("POST", "/questions/", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")
    
    // Создаем ResponseRecorder
    rr := httptest.NewRecorder()
    
    // Вызываем обработчик
    handler.CreateQuestion(rr, req)
    
    // Проверяем статус код
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
    }
    
    var response models.Question
    if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
        t.Fatal(err)
    }
    
    if response.Text != question.Text {
        t.Errorf("handler returned unexpected body: got %v want %v", response.Text, question.Text)
    }
    
    mockDB.AssertCalled(t, "Create", &question)
}