package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rest/internal/domain"
	"rest/internal/services"
	"strconv"
	"time"
)

var ctx context.Context = context.Background()

type UserHandler struct {
	userServo services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userServo: userService}
}

func (h *UserHandler) GetHello(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("hello, " + strconv.Itoa(int(userID))))
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response(w, "invalid data", http.StatusBadRequest)
		return
	}

	if err := (&user).Validate(); err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.userServo.CreateUser(ctx, &user)
	if err != nil {
		response(w, "unable to create user", http.StatusInternalServerError)
		return
	}

	response(w, success, http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response(w, "unable to parse", http.StatusBadRequest)
		return
	}
	// log.Println("[handler1]", user)
	if err := user.Validate(); err != nil {
		response(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("[handler]", user)
	cookieAuth, err := h.userServo.Login(ctx, &user)
	if err != nil {
		response(w, failed, http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:     "jwt",
		HttpOnly: true,
		Value:    cookieAuth,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
	}

	http.SetCookie(w, c)
	response(w, success, http.StatusOK)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:     "jwt",
		HttpOnly: true,
		Value:    "",
		Path:     "/",
		Domain:   "",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().AddDate(0, 0, -7),
	}

	http.SetCookie(w, c)
	response(w, success, http.StatusOK)
}
