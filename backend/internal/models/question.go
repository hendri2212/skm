package models

type Question struct {
	ID           uint     `json:"id" gorm:"primaryKey"`
	QuestionText string   `json:"question_text" gorm:"type:text;not null"`
	Choices      []Choice `json:"choices" gorm:"foreignKey:QuestionID"`
}
