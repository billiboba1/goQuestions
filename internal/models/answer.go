package models

import (
    "time"
)

type Answer struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    QuestionID uint      `gorm:"not null;index" json:"question_id"`
    UserID     string    `gorm:"type:uuid;not null;index" json:"user_id"`
    Text       string    `gorm:"type:text;not null" json:"text"`
    CreatedAt  time.Time `json:"created_at"`
}