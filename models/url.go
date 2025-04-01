package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URL represents a shortened URL in our system
type URL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OriginalURL string             `bson:"original_url" json:"original_url"`
	ShortCode   string             `bson:"short_code" json:"short_code"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
	Clicks      int                `bson:"clicks" json:"clicks"`
}

// CreateURLRequest represents the request to create a short URL
type CreateURLRequest struct {
	URL string `json:"url"`
}

// CreateURLResponse represents the response after creating a short URL
type CreateURLResponse struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// URLStatsResponse represents the statistics for a URL
type URLStatsResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Clicks      int       `json:"clicks"`
}
