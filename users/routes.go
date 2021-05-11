package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/marianodsr/nura-api/authentication"
	"github.com/marianodsr/nura-api/middlewares"
)

//RegisterRoutes func
func RegisterRoutes(r chi.Router) {
	r.Post("/login", loginHandler)
	r.Post("/signup", signupHandler)

	r.Get("/refresh", refreshTokenHandler)

	r.With(middlewares.RequireAuth).Get("/{id}", getUserHandler)

	r.With(middlewares.RequireAuth).Put("/{id}", updateUserHandler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	type requestShape struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	loginAttempt := &requestShape{}
	json.NewDecoder(r.Body).Decode(loginAttempt)
	pair, err := attemptLogin(loginAttempt.Email, loginAttempt.Password)
	fmt.Printf("%+v", err)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	cookie := http.Cookie{
		Name:     "refresh-token",
		Value:    pair[1],
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(map[string]string{"access-token": pair[0]})
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	json.NewDecoder(r.Body).Decode(user)
	exists, _ := GetUserByEmail(user.Email)
	if exists != nil {
		http.Error(w, "Email already in use", 400)
		return
	}
	user, err := signUp(user)
	if err != nil {
		http.Error(w, "Internal error", 500)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)
}

func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := r.Cookie("refresh-token")
	fmt.Printf("%v", refreshToken)
	if err != nil {
		fmt.Println("error line 60")
		http.Error(w, "Authentication error...", 401)
		return
	}
	decoded, err := authentication.ValidateToken(refreshToken.Value)
	if err != nil {
		fmt.Println("error line 66")
		http.Error(w, "Authentication error...", 401)
		return
	}
	pair, err := authentication.GenerateTokenPair(uint(decoded["UserID"].(float64)))
	if err != nil {
		fmt.Println("error line 72")
		http.Error(w, "Authentication error...", 401)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"access-token": pair[0]})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "no id provided", 400)
		return
	}
	intID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "error parsing id", 400)
		return
	}
	user, err := GetUserByID(uint(intID))
	if err != nil {
		http.Error(w, "Invalid user id", 400)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "No id provided", http.StatusBadRequest)
		return
	}
	intID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Provided payload is not of the correct shape", http.StatusBadRequest)
		return
	}
	if _, err := GetUserByID(uint(intID)); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err := UpdateUser(user); err != nil {
		http.Error(w, "Unable to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
