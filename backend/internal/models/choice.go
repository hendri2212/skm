package models

type Choice struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	QuestionID uint     `json:"question_id"`
	Question   Question `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ChoiceText string   `json:"choice_text" gorm:"size:255;not null"`
	Points     int      `json:"points" gorm:"not null"`
}
