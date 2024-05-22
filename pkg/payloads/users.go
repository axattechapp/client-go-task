package payloads

import (
	db_sqlc "client_task/pkg/common/db/sqlc"

	// "github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUser struct {
	FullName     string           `json:"full_name" binding:"required`
	Email        string           `json:"email" binding:"required,email"`
	PasswordHash string           `json:"password_hash" binding:"required" `
	UserType     db_sqlc.UserRole `json:"user_type" binding:"required"`
	AvatarUrl    pgtype.Text      `json:"avatar_url" binding:"required`
}

type UpdateUser struct {
	FullName     string `json:"full_name" binding:"required`
	PasswordHash string `json:"password_hash" binding:"required" `
}

type CreateProfile struct {
	UserID      int32       `json:"user_id" binding:"required`
	Bio         pgtype.Text `json:"bio" binding:"required`
	Company     string      `json:"company" binding:"required`
	JobRole     string      `json:"job_role" binding:"required`
	Description pgtype.Text `json:"description" binding:"required`
	Skills      []string    `json:"skills"` // Added skills list

}

type UpdateProfile struct {
	ID          int32       `json:"id" binding:"required`
	Bio         pgtype.Text `json:"bio" binding:"required`
	Company     string      `json:"company" binding:"required`
	JobRole     string      `json:"job_role" binding:"required`
	Description pgtype.Text `json:"description" binding:"required`
	Skills      []string    `json:"skills" binding:"required` // Added skills list

}

type CreateSkill struct {
	Name string `json:"Name"`
}

type CreateCareer struct {
	UserID      int32  `json:"user_id" binding:"required`
	Title       string `json:"title" binding:"required`
	Company     string `json:"company" binding:"required`     // Optional
	Description string `json:"description" binding:"required` // Optional
	SkillID     int32  `json:"skill_id" binding:"required`
}

type UpdateCareer struct {
	Title       string `json:"title" binding:"required`
	Company     string `json:"company" binding:"required`     // Optional
	Description string `json:"description" binding:"required` // Optional
	SkillID     int32  `json:"skill_id" binding:"required`
}
