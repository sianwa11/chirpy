package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")

	if author_id != "" {
		userID, err := uuid.Parse(author_id)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to get user", err)
			return
		}

		chirps, err := cfg.db.GetAllChirpsByAuthor(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to get chirps", err)
			return
		}

		chirpsResponse := []Chirp{}
		for _, chirp := range chirps {
			chirpsResponse = append(chirpsResponse, Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserId: chirp.UserID,
			})
		}

		if sortBy == "desc" {
			sort.Slice(chirpsResponse, func(i, j int) bool {
				return chirpsResponse[i].CreatedAt.After(chirpsResponse[j].CreatedAt)
			})
		}

		respondWithJSON(w, http.StatusOK, chirpsResponse)
		return
	}

	chirpsArr, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch chirps", err)
		return
	}

	chirpsResponse := []Chirp{}
	for _, chirp := range chirpsArr {
		chirpsResponse = append(chirpsResponse, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		})
	}

	if sortBy == "desc" {
		sort.Slice(chirpsResponse, func(i, j int) bool {
			return chirpsResponse[i].CreatedAt.After(chirpsResponse[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsResponse)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")	
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse chirpID", err)
		return
	}
	
	
	chirpDB, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "failed to get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp {
		ID: chirpDB.ID,
		CreatedAt: chirpDB.CreatedAt,
		UpdatedAt: chirpDB.UpdatedAt,
		Body: chirpDB.Body,
		UserId: chirpDB.UserID,
	})
}