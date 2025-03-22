package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/mickali02/templates/internal/validator"
)

type Feedback struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Fullname  string    `json:"fullname"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Email     string    `json:"email"`
}

func ValidateFeedback(v *validator.Validator, feedback *Feedback) {
	v.Check(validator.NotBlank(feedback.Fullname), "fullname", "must be provided")
	v.Check(validator.MaxLength(feedback.Fullname, 50), "fullname", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(feedback.Subject), "subject", "must be provided")
	v.Check(validator.MaxLength(feedback.Subject, 50), "subject", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(feedback.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(feedback.Email), "email", "invalid email address")
	v.Check(validator.MaxLength(feedback.Email, 100), "email", "must not be more than 100 bytes long")
	v.Check(validator.NotBlank(feedback.Message), "message", "must be provided")
	v.Check(validator.MaxLength(feedback.Message, 500), "message", "must not be more than 500 bytes long")
}

type FeedbackModel struct {
	DB *sql.DB
}

func (m *FeedbackModel) Insert(feedback *Feedback) error {
	query := `
		INSERT INTO feedback (fullname, subject, message, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		feedback.Fullname,
		feedback.Subject,
		feedback.Message,
		feedback.Email,
	).Scan(&feedback.ID, &feedback.CreatedAt)
}
