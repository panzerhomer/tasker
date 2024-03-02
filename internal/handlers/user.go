package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"rest/internal/domain"
	myjwt "rest/internal/lib/jwt"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt")
		switch err {
		case nil:
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if ok := myjwt.ValidateJwtToken(c.Value); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idCtx := context.WithValue(r.Context(), "email", myjwt.DecodeJwtToken(c.Value))

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
}

type UserService interface {
	Login(ctx context.Context, user *domain.User) (string, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func response(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

const (
	success = "success"
	failed  = "failed"
)

func (uh *UserHandler) Routes() chi.Router {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(middleware.RequestID)
	root.Post("/login", uh.Login)
	root.Post("/signin", uh.CreateUser)

	r := chi.NewRouter()
	r.Use(Auth)
	r.Get("/hello", GetHello)

	root.Mount("/api", r)

	return root
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("email").(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("hello, " + id))
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response(w, "invalid data", http.StatusBadRequest)
		return
	}

	ctx := context.TODO()

	_, err := uh.service.CreateUser(ctx, &user)
	if err != nil {
		response(w, "unable to create user", http.StatusInternalServerError)
		return
	}

	response(w, success, http.StatusOK)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response(w, "unable to parse", http.StatusBadRequest)
		return
	}

	ctx := context.TODO()

	cookieAuth, err := uh.service.Login(ctx, &user)
	if err != nil {
		response(w, failed, http.StatusInternalServerError)
		return
	}

	c := &http.Cookie{
		Name:    "jwt",
		Value:   cookieAuth,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(time.Minute * 2),
	}

	http.SetCookie(w, c)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    "jwt",
		Value:   "",
		Path:    "/",
		Domain:  "",
		Expires: time.Now().AddDate(0, 0, -7),
	}
	http.SetCookie(w, c)
	response(w, success, http.StatusOK)
}
