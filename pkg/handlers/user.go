package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"my-chi/models"
	"my-chi/pkg"
	"net/http"
)

type UserRouter interface {
	NewUserRouter(ctx context.Context, db *pkg.DB) *UserHandler
	GetAllusers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	Ctx context.Context
	DB  *pkg.DB
}

func NewUserRouter(ctx context.Context, db *pkg.DB) *UserHandler {
	return &UserHandler{
		Ctx: ctx,
		DB:  db,
	}
}

func (R UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		fmt.Println("[User-Router] failed reading body")
		http.Error(w, "failed reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := R.DB.UpdateWithData(newUser); err != nil {
		fmt.Println("[UserHandler] Failed updating database:", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(R.Ctx, "latest-id", newUser.Id)

	r = r.WithContext(ctx)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func (R UserHandler) GetAllusers(w http.ResponseWriter, r *http.Request) {
	stringifyData, err := json.Marshal(R.DB.UserData)
	if err != nil {
		http.Error(w, "failed parsing json", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(stringifyData)
}
