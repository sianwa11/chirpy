package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sianwa11/chirpy/internal/auth"
)

func (cfg *apiConfig) handleUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to get api key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req Req
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode request body", err)
		return
	}

	if req.Event != "user.upgraded" || req.Data.UserID == "" {
		respondWithError(w,http.StatusNoContent, "No content", nil)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse uuid", err)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed tp upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}