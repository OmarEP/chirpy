package main

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/OmarEP/chirpy/internal/auth"
	"github.com/OmarEP/chirpy/internal/database"
)

type User struct {
	ID 		int `json:"id"`
	Email 	string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return 
	}

	HashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
	}
	
	user, err := cfg.DB.CreateUser(params.Email, HashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExits) {
			respondWithJSON(w, http.StatusConflict, "User already exists")
			return 
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:		user.ID,
			Email: 	user.Email,
		},
	})
}

