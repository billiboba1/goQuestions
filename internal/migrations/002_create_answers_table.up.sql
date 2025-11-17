CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Внешний ключ с каскадным удалением
    CONSTRAINT fk_answers_question 
        FOREIGN KEY (question_id) 
        REFERENCES questions(id) 
        ON DELETE CASCADE
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_answers_question_id ON answers(question_id);
CREATE INDEX idx_answers_user_id ON answers(user_id);
CREATE INDEX idx_answers_created_at ON answers(created_at DESC);