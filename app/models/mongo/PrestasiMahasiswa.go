package mongo

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `json:"student_id"`
	AchievementType string             `json:"achievement_type"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Details         any                `json:"details"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type AchievementAttachementRequest struct {
	AchievementReferencesID        string
	StudentID  string             `json:"student_id"`
	Attachment *multipart.FileHeader
}

type CompetitionDetails struct {
	CompetitionName  string `json:"competition_name"`
	CompetitionLevel string `json:"competition_level"`
	Rank             int    `json:"rank"`
	MedalType        string `json:"medal_type"`
}

type PublicationDetails struct {
	PublicationType  string   `json:"publication_type"`
	PublicationTitle string   `json:"publication_title"`
	Authors          []string `json:"authors"`
	Publisher        string   `json:"publisher"`
	ISSN             string   `json:"issn"`
}

type OrganizationDetails struct {
	OrganizationName string `json:"organization_name"`
	Position         string `json:"position"`
	Period           string `json:"period"`
}

type CertificationDetails struct {
	CertificationName   string    `json:"certification_name"`
	IssuedBy            string    `json:"issued_by"`
	CertificationNumber string    `json:"certification_number"`
	ValidUntil          time.Time `json:"valid_until"`
}
