package postgres

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID            uuid.UUID `json:"id"`
	StudentID     uuid.UUID `json:"student_id"`
	MongoId       string    `json:"mongo_achievement_id"`
	Status        string    `json:"status"`
	SubmittedAt   time.Time `json:"submitted_at"`
	RejectionNote *string   `json:"rejection_note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AchievmentRejectRequest struct {
	RejectionNote string `json:"rejection_note"`
}

type AchievementMongo struct {
	Details     map[string]interface{} `json:"details"`
}

type AchievementGabung struct {
	Achievement Achievement            `json:"achievement"`
	Details     map[string]interface{} `json:"details"`
}
