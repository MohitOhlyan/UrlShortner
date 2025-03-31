package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	utils "urlShortner/Utils"
	"urlShortner/config"
	"urlShortner/db"
	"urlShortner/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateShortURLHandler handles requests to create a short URL
func CreateShortURLHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.CreateURLRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate URL
		if !utils.ValidateURL(request.URL) {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		// Generate short code
		shortCode, err := utils.GenerateShortCode()
		if err != nil {
			http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
			return
		}

		// Create URL document
		now := time.Now()
		expiresAt := now.Add(cfg.URLExpiration)
		urlDoc := models.URL{
			OriginalURL: request.URL,
			ShortCode:   shortCode,
			CreatedAt:   now,
			ExpiresAt:   expiresAt,
			Clicks:      0,
		}

		// Insert into MongoDB
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = db.URLCollection.InsertOne(ctx, urlDoc)
		if err != nil {
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}

		// Create response
		response := models.CreateURLResponse{
			ShortURL:    cfg.BaseURL + "/" + shortCode,
			OriginalURL: request.URL,
			ExpiresAt:   expiresAt,
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// RedirectHandler handles the redirection from short URLs to original URLs
func RedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := mux.Vars(r)["shortCode"]

		// Find the URL in the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var urlDoc models.URL
		filter := bson.M{"short_code": shortCode}
		err := db.URLCollection.FindOne(ctx, filter).Decode(&urlDoc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "URL not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Check if URL has expired
		if time.Now().After(urlDoc.ExpiresAt) {
			http.Error(w, "URL has expired", http.StatusGone)
			return
		}

		// Increment click count
		update := bson.M{"$inc": bson.M{"clicks": 1}}
		_, err = db.URLCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			// Log the error but continue with redirection
			// In a production environment, you might want a more robust error handling
			// strategy
		}

		// Redirect to the original URL
		http.Redirect(w, r, urlDoc.OriginalURL, http.StatusFound)
	}
}

// URLStatsHandler handles requests for URL statistics
func URLStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := mux.Vars(r)["shortCode"]

		// Find the URL in the database
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var urlDoc models.URL
		filter := bson.M{"short_code": shortCode}
		err := db.URLCollection.FindOne(ctx, filter).Decode(&urlDoc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "URL not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// Create response
		response := models.URLStatsResponse{
			ShortCode:   urlDoc.ShortCode,
			OriginalURL: urlDoc.OriginalURL,
			CreatedAt:   urlDoc.CreatedAt,
			ExpiresAt:   urlDoc.ExpiresAt,
			Clicks:      urlDoc.Clicks,
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status": "healthy"}
		json.NewEncoder(w).Encode(response)
	}
}
