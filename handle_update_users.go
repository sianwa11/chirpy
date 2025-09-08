package main

import (
	"encoding/json"
	"net/http"

	"github.com/sianwa11/chirpy/internal/auth"
	"github.com/sianwa11/chirpy/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}


	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to get bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to validate jwt", err)
		return
	}

	_, err = cfg.db.FindUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to find user", err)
		return
	}

	var params Parameters
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode request body", err)
		return
	}


	if params.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "unathorized to perform this action", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to hash password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to update user", err)
		return
	}


	respondWithJSON(w, http.StatusOK, User {
		ID: updatedUser.ID,
		Email: updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
		CreatedAt: updatedUser.UpdatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	})

}

