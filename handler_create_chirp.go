package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sianwa11/chirpy/internal/database"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	


	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode request", err)
		return
	}

	chirpBody, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams{UserID: params.UserId, Body: params.Body})
	if err != nil {
		respondWithError(w, 404, "failed to create chirp", err)
	}

	
	respondWithJSON(w, 201, Chirp{
		ID: chirpBody.ID,
		CreatedAt: chirpBody.CreatedAt,
		UpdatedAt: chirpBody.UpdatedAt,
		Body: chirpBody.Body,
		UserId: chirpBody.UserID,
	})

}