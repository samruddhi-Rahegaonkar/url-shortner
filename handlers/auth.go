package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"url_shortner/database"
	"url_shortner/models"
	"url_shortner/utils"
)

func RegisterHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var user models.User

	err := json.NewDecoder(
		r.Body,
	).Decode(&user)

	if err != nil {
		http.Error(
			w,
			"Invalid Body",
			http.StatusBadRequest,
		)
		return
	}

	query := `
	INSERT INTO users(username,password)
	VALUES($1,$2)
	`

	_, err = database.DB.Exec(
		context.Background(),
		query,
		user.Username,
		user.Password,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			500,
		)
		return
	}

	w.Write([]byte("User Registered"))
}
func LoginHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var user models.User

	json.NewDecoder(
		r.Body,
	).Decode(&user)

	var (
		userID   int
		password string
	)

	query := `
	SELECT id,password
	FROM users
	WHERE username=$1
	`

	err := database.DB.QueryRow(
		context.Background(),
		query,
		user.Username,
	).Scan(&userID, &password)

	if err != nil {
		http.Error(
			w,
			"User Not Found",
			401,
		)
		return
	}

	if password != user.Password {

		http.Error(
			w,
			"Wrong Password",
			401,
		)
		return
	}

	token, err := utils.GenerateToken(
		userID,
		user.Username,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			500,
		)
		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"token": token,
		},
	)
}
