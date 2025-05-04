package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dtoReq "github.com/098765432m/internal/dto/request"
	"github.com/098765432m/internal/service"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler (service *service.UserService) *UserHandler{
	return &UserHandler{
		Service: service,
	}
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq dtoReq.LoginDtoRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := uh.Service.Login(&loginReq)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	res := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq dtoReq.CreatedUserDtoRequest
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = uh.Service.GuestRegister(&registerReq)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (uh *UserHandler) DashboardUpdateUser(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)	// Convert string to int
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var updateReq dtoReq.DashBoardUpdateUserDtoRequest
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = uh.Service.DashboardUpdateUser(id, &updateReq)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully!"})
}

func (uh *UserHandler) DeteleUser(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return 
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err = uh.Service.DeleteUser(id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	returnMessage := fmt.Sprintf("User Id %d has been deleted!", id)
	json.NewEncoder(w).Encode(map[string]string{"message": returnMessage})
}

func (uh *UserHandler) ChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	if idParam == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user Id", http.StatusBadRequest)
		return
	}

	err = uh.Service.ChangeUserPassword(id, r.FormValue("new_password"))
	if err != nil {
		http.Error(w, "Failed to change password", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "applcation/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed"})
}