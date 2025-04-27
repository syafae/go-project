package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/syafae/femProject/internal/store"
	"github.com/syafae/femProject/internal/tokens"
	"github.com/syafae/femProject/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("Error: decoding request body:%v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request body"})
		return
	}

	user, err := h.userStore.GetUserByName(req.UserName)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByName: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}

	isValid, err := user.PasswordHash.Matches(req.Password)
	if err != nil || !isValid {
		h.logger.Printf("ERROR PasswordHash.Matches:%v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid password"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR CreateNewToken: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"token": token})
}
