package models

type Answer struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	QuestionID uint     `json:"question_id"`
	Question   Question `json:"question" gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ChoiceID uint   `json:"choice_id"`
	Choice   Choice `json:"choice" gorm:"foreignKey:ChoiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
