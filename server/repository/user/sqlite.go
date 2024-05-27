package user

import (
	"database/sql"
	"fmt"
	"godb/model"
)

type SQLRepo struct {
	DB *sql.DB
}

func (r *SQLRepo) CheckIfReg(user model.User) error {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM `+"`user`"+` WHERE email = ?`, user.Email)
	if err != nil {
		return fmt.Errorf("failed to perfom select: %w", err)
	}
	for rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Email, &model_user.Password, &model_user.FirstName, &model_user.LastName)
	}
	if model_user.UserID != 0 {
		return fmt.Errorf("user registered")
	}
	return nil
}

func (r *SQLRepo) Insert(user model.User) error {
	statement, err := r.DB.Prepare(`INSERT INTO` + "`user`" + `(
	email, 
	password,
	first_name,
	last_name) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert user: %w", err)
	}
	_, err = statement.Exec(user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *SQLRepo) FindById(user_id uint) (model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM`+"`user`"+`WHERE user_id = ?`, user_id)
	if err != nil {
		return model_user, fmt.Errorf("failed to find by id user: %w", err)
	}

	for rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Email, &model_user.Password, &model_user.FirstName, &model_user.LastName)
	}

	return model_user, nil
}

func (r *SQLRepo) UserExists(email string) (bool, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT user_id FROM`+" `user` "+`WHERE email = ?`, email)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		rows.Scan(&model_user.UserID)
	}
	if model_user.UserID > 0 {
		return true, nil
	}
	return false, nil
}

func (r *SQLRepo) FindIdByCredentials(email string, password string) (model.User, error) {
	var user = model.User{}
	rows, err := r.DB.Query(`SELECT user_id, first_name, last_name FROM`+"`user`"+`WHERE email = ? AND password = ?`, email, password)
	if err != nil {
		return user, fmt.Errorf("failed to find user by credentials: %w", err)
	}

	for rows.Next() {
		rows.Scan(&user.UserID, &user.FirstName, &user.LastName)
	}
	if user.UserID == 0 {
		return user, fmt.Errorf("incorrect password: %w", err)
	}

	return user, nil
}

func (r *SQLRepo) DeleteByID(user_id uint) error {
	statement, err := r.DB.Prepare(`DELETE FROM ` + "`user`" + ` WHERE user_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete user: %w", err)
	}

	_, err = statement.Exec(user_id)
	if err != nil {
		return fmt.Errorf("failed to prepare delete EXEC user: %w", err)
	}

	return nil
}
