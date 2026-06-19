package middleware

import (
	"context"
	"net/http"
	"strings"
	"url_shortner/utils"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(
	next http.HandlerFunc,
) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		authHeader :=
			r.Header.Get(
				"Authorization",
			)

		if authHeader == "" {

			http.Error(
				w,
				"Token Missing",
				401,
			)

			return
		}

		tokenString :=
			strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

		token, err := jwt.Parse(
			tokenString,
			func(
				token *jwt.Token,
			) (interface{}, error) {
				return utils.SecretKey, nil
			},
		)

		if err != nil ||
			!token.Valid {

			http.Error(
				w,
				"Invalid Token",
				401,
			)

			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userID := int(
			claims["id"].(float64),
		)

		ctx := context.WithValue(
			r.Context(),
			"userID",
			userID,
		)

		next(
			w,
			r.WithContext(ctx),
		)
	}

}
