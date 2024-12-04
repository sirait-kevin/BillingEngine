package restful

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/usecases"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserUseCase *usecases.UserUseCase
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.UserUseCase.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.UserUseCase.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/users/"+strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusCreated)
}
