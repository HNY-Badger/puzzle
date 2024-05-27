package handler

import (
	"encoding/json"
	"fmt"
	"godb/model"
	"godb/repository/user"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type User struct {
	Repo *user.SQLRepo
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := model.User{
		Email:     body.Email,
		Password:  body.Password,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	}
	err := u.Repo.CheckIfReg(user)
	if err != nil {
		fmt.Println("user exists:", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	err = u.Repo.Insert(user)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *User) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	const base = 10
	const bitSize = 64

	userID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := o.Repo.FindById(uint(userID))
	if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.UserID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *User) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	const base = 10
	const bitSize = 64

	userID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(uint(userID))

	if err != nil {
		fmt.Println("failed to deleted by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
