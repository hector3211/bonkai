package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-chi/models"
	"my-chi/pkg"
	"net/http"
)

type UserRouter struct {
	Ctx context.Context
	DB  *pkg.DB
}

func NewUserRouter(ctx context.Context, db *pkg.DB) *UserRouter {
	return &UserRouter{
		Ctx: ctx,
		DB:  db,
	}
}

func (R UserRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := models.User{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("[User-Router] failed reading body")
		http.Error(w, "failed reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		fmt.Println("[User-Router] failed parsing body to data")
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	R.DB.UpdateWithData(newUser)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func (R UserRouter) GetAllusers(w http.ResponseWriter, r *http.Request) {
	stringifyData, err := json.Marshal(R.DB.UserData)
	if err != nil {
		http.Error(w, "failed parsing json", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(stringifyData)
}
