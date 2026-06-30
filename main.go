package main

import (
	"net/http"

	"url_shortner/database"
	"url_shortner/handlers"
	"url_shortner/middleware"
)

func main() {

	database.ConnectDB()

	http.HandleFunc(
		"/register",
		handlers.RegisterHandler,
	)

	http.HandleFunc(
		"/login",
		handlers.LoginHandler,
	)

	http.HandleFunc(
		"/shorten",
		middleware.JWTMiddleware(
			handlers.ShortenURLHandler,
		),
	)

	http.HandleFunc(
		"/",
		handlers.HandleRedirect,
	)

	http.ListenAndServe(
		":8080",
		nil,
	)
}
