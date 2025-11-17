package models

import (
    "time"
)

type Question struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Text      string    `gorm:"type:text;not null" json:"text"`
    CreatedAt time.Time `json:"created_at"`
    Answers   []Answer  `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}