package handler

import (
	"encoding/json"
	"godb/model"
	"godb/repository/level"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type Level struct {
	LRepo *level.SQLRepo
	Key   []byte
}

func (l *Level) Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	req := model.Request{}
	json.NewDecoder(r.Body).Decode(&req)
	userJWT := req.Cookie

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(userJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.Key), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var userId uint
	for key, val := range claims {
		if key == "id" {
			userId = uint(val.(float64))
		}
	}

	level, err := l.LRepo.GetLevel(req.Id, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(level)
	w.WriteHeader(http.StatusOK)
}

func (l *Level) Completed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	req := model.Request{}
	json.NewDecoder(r.Body).Decode(&req)
	userJWT := req.Cookie

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(userJWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.Key), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var userId uint
	for key, val := range claims {
		if key == "id" {
			userId = uint(val.(float64))
		}
	}

	err = l.LRepo.CompleteRound(req.Id, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
