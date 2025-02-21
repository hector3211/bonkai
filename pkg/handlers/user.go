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
	ctx *context.Context
}

const UserDataPath = "/tmp/uer-data.json"

func NewUserRouter(ctx *context.Context) *UserRouter {
	return &UserRouter{
		ctx: ctx,
	}
}

func (R UserRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	db, err := pkg.NewDB()
	if err != nil {
		fmt.Println("[DB] No database connection")
		return
	}

	newUser := models.User{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("[DB] failed reading body")
		http.Error(w, "failed reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		fmt.Println("[router] failed parsing body to data")
		http.Error(w, "invalid json", http.StatusBadRequest)
	}

	newData := append(db.UserData, newUser)

	db.UpdateWithData(newData)
}
