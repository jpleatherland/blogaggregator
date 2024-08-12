package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/jpleatherland/blogaggregator/internal/database"
)

func (resources *Resources) authMiddleware(handler authHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.Background()

		// get apikey from header
		authToken := req.Header.Get("Authorization")
		if strings.HasPrefix(authToken, "ApiKey ") {
			authToken = strings.TrimPrefix(authToken, "ApiKey ")
		} else {
			respondWithError(rw, http.StatusUnauthorized, "invalid apikey")
			return
		}

		// get user by api key
		user, err := resources.DB.GetUserByApiKey(ctx, authToken)
		if err != nil {
			respondWithError(rw, http.StatusNotFound, "couldn't get user")
			return
		}

		// pass authenticated user to handler
		handler(rw, req, user)
	}
}

type authHandler func(rw http.ResponseWriter, req *http.Request, dbUser database.User)
