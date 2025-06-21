package models

import "time"

// UserRole defines valid roles for a user.
type UserRole string

const (
	UserRoleSuperadmin UserRole = "superadmin"
	UserRoleAdmin      UserRole = "admin"
	UserRoleUser       UserRole = "user"
)

type User struct {
	ID               uint        `json:"id" gorm:"primaryKey"`
	FullName         string      `json:"full_name" gorm:"size:50"`
	PlaceOfBirth     string      `json:"place_of_birth" gorm:"size:100"`
	DateOfBirth      time.Time   `json:"date_of_birth"`
	IsMale           bool        `json:"is_male" gorm:"type:boolean"`
	LastEducationID  uint        `json:"last_education_id"`
	LastEducation    *Education  `json:"last_education" gorm:"foreignKey:LastEducationID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	MainOccupationID uint        `json:"main_occupation_id"`
	MainOccupation   *Occupation `json:"main_occupation" gorm:"foreignKey:MainOccupationID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SubmitAnswersPayload struct {
	User struct {
		Name         string `json:"name" binding:"required"`
		BirthPlace   string `json:"birthPlace" binding:"required"`
		BirthDate    string `json:"birthDate" binding:"required,datetime=2006-01-02"` // format YYYY-MM-DD
		IsMale       bool   `json:"is_male"`
		EducationID  uint   `json:"education" binding:"required"`
		OccupationID uint   `json:"occupation" binding:"required"`
	} `json:"user" binding:"required"`
	Answers []struct {
		QuestionID uint `json:"question_id" binding:"required"`
		ChoiceID   uint `json:"choice_id" binding:"required"`
	} `json:"answers" binding:"required,dive"`
}
