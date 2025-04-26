package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/syafae/femProject/internal/store"
	"github.com/syafae/femProject/internal/utils"
)

type registeredUserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateregisterRequest(reg *registeredUserRequest) error {
	if reg.UserName == "" {
		return errors.New("username is required")
	}
	if len(reg.UserName) > 50 {
		return errors.New("username cannot be greater than 50 characters")
	}
	emailregex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(?:\.[a-zA-Z]{2,})?$`)
	if !emailregex.MatchString(reg.Email) {
		return errors.New("invalid email format")
	}
	if len(reg.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(reg.Password) {
		return errors.New("password must include at least one lowercase letter")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(reg.Password) {
		return errors.New("password must include at least one uppercase letter")
	}
	if !regexp.MustCompile(`\d`).MatchString(reg.Password) {
		return errors.New("password must include at least one number")
	}
	if !regexp.MustCompile(`[@$!%*?&]`).MatchString(reg.Password) {
		return errors.New("password must include at least one special character")
	}

	return nil
}

func (uh *UserHandler) HandleRegiserUserRequest(w http.ResponseWriter, r *http.Request) {
	var req registeredUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Error decoding request body:%v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	err = uh.validateregisterRequest(&req)
	if err != nil {
		uh.logger.Printf("Error: validateregisterRequest %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		UserName: req.UserName,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	_, err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: CreateUser %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username is required"})
		return
	}

	var req registeredUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Error decoding request body:%v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	err = uh.validateregisterRequest(&req)
	if err != nil {
		uh.logger.Printf("Error: validateregisterRequest %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		UserName: username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if req.Password != "" {
		_, err = user.PasswordHash.Set(req.Password)
		if err != nil {
			uh.logger.Printf("ERROR: hashing password %v", err)
			utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "internal server error"})
			return
		}
	}

	err = uh.userStore.UpdateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: UpdateUser %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})

}

func (uh *UserHandler) HandleGetUserByName(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	if username == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username is required"})
		return
	}

	user, err := uh.userStore.GetUserByName(username)
	if err != nil {
		uh.logger.Printf("ERROR: GetUserByName %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if user == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}
