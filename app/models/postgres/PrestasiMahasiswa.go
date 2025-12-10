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

type AchievementGabung struct {
	Achievement Achievement            `json:"achievement"`
	Details     map[string]interface{} `json:"details"`
}

type CompetitionDetails struct {
	ID               uuid.UUID `json:"id"`
	AchievementID    uuid.UUID `json:"achievement_id"`
	CompetitionName  string    `json:"competition_name"`
	CompetitionLevel string    `json:"competition_level"`
	Rank             int       `json:"rank"`
	MedalType        string    `json:"medal_type"`
}

type PublicationDetails struct {
	ID               uuid.UUID `json:"id"`
	AchievementID    uuid.UUID `json:"achievement_id"`
	PublicationType  string    `json:"publication_type"`
	PublicationTitle string    `json:"publication_title"`
	Authors          []string  `json:"authors"` // TEXT[]
	Publisher        string    `json:"publisher"`
	ISSN             string    `json:"issn"`
}

type OrganizationDetails struct {
	ID               uuid.UUID `json:"id"`
	AchievementID    uuid.UUID `json:"achievement_id"`
	OrganizationName string    `json:"organization_name"`
	Position         string    `json:"position"`
	Period           string    `json:"period"` // or DATERANGE
}

type CertificationDetails struct {
	ID                  uuid.UUID `json:"id"`
	AchievementID       uuid.UUID `json:"achievement_id"`
	CertificationName   string    `json:"certification_name"`
	IssuedBy            string    `json:"issued_by"`
	CertificationNumber string    `json:"certification_number"`
	ValidUntil          time.Time `json:"valid_until"`
}
