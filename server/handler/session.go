package handler

import (
	"encoding/json"
	"fmt"
	"godb/model"
	"godb/repository/user"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type Session struct {
	URepo *user.SQLRepo
	Key   []byte
}

func (u *Session) Create(w http.ResponseWriter, r *http.Request) {
	emailParam := r.URL.Query().Get("email")
	passParam := r.URL.Query().Get("password")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	user_exists, err := u.URepo.UserExists(emailParam)
	if err != nil || !user_exists {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := u.URepo.FindIdByCredentials(emailParam, passParam)
	if err != nil {
		fmt.Println("failed to find user:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expires := 604800000

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.UserID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(u.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session := model.Session{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Value:     tokenString,
		Expires:   expires,
	}

	json.NewEncoder(w).Encode(session)
}
