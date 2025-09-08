package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sianwa11/chirpy/internal/auth"
	"github.com/sianwa11/chirpy/internal/database"
)



func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding response body", err)
		return
	}

	if params.Password == "" || params.Email == "" {
		respondWithError(w, http.StatusForbidden, "password and email are required", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to hash password", err)
		return
	}

	usr, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating user", err)
		return
	}


	respondWithJSON(w, http.StatusCreated,User{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email: usr.Email,
		IsChirpyRed: usr.IsChirpyRed.Bool,
	})
}

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
	}

	ctx := r.Context()

	var params Parameters
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode request body", err)
		return
	}

	usr, err := cfg.db.FindByEmail(ctx, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error finding user by email", err)
		return
	}


	err = auth.CheckPasswordHash(params.Password, usr.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	expirationTime := 1 * time.Hour

	jwt, err := auth.MakeJWT(usr.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error making jwt", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to make refresh token", err)
		return
	} 

	refreshTknDB, err := cfg.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: usr.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email: usr.Email,
		IsChirpyRed: usr.IsChirpyRed.Bool,
		Token: jwt,
		RefreshToken: refreshTknDB.Token,
	})
}