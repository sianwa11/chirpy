package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sianwa11/chirpy/internal/auth"
	"github.com/sianwa11/chirpy/internal/database"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		respondWithError(w, 400, "missing chirp ID", nil)
		return
	}

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, 400, "invalid chirp ID format", err)
		return
	}
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "failed to get bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to validate jwt", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), uuid.UUID(chirpID))
	if err != nil {
		respondWithError(w, 403, "failed to get chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, 403, "failed to get chrip", nil)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID: chirp.ID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}