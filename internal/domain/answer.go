package domain

import "time"

type Answer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	QuestionID uint      `gorm:"not null;index" json:"question_id"`
	UserID     string    `gorm:"type:varchar(255);not null;index" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Answer) TableName() string {
	return "answers"
}
