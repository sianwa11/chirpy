package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) middlewareValidateChirp(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string      `json:"body"`
			UserId uuid.UUID `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to decode request body", err)
			return
		}

		const maxChirpLength = 140
		if len(params.Body) > maxChirpLength {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
			return 
		}

		profane := []string{"kerfuffle", "sharbert", "fornax"}
		splitBody := strings.Split(params.Body, " ")
		for i, word := range splitBody {
			if slices.Contains(profane, strings.ToLower(word)) {
				splitBody[i] = "****"
			}
		}

		params.Body = strings.Join(splitBody, " ")

		newBody, err := json.Marshal(params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to process request", err)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(newBody))

		next.ServeHTTP(w, r)
	})
}

// func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Body string `json:"body"`
// 	}

// 	type returnVals struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
// 		return
// 	}

// 	const maxChirpLength = 140
// 	if len(params.Body) > maxChirpLength {
// 		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
// 		return
// 	}

// 	profane := []string{"kerfuffle", "sharbert", "fornax"}
// 	splitBody := strings.Split(params.Body, " ")
// 	for i, word := range splitBody {
// 		if slices.Contains(profane, strings.ToLower(word)) {
// 			splitBody[i] = "****"
// 		}
// 	}


// 	respondWithJSON(w, http.StatusOK, returnVals{
// 		CleanedBody: strings.Join(splitBody, " "),
// 	})
// }