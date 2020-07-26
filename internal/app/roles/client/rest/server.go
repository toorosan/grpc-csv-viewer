package rest

import (
	"net/http"
)

// ListenAndServeREST initializes REST server handling UI requests.
func ListenAndServeREST() error {
	registerHandlers()

	return http.ListenAndServe(":8081", nil)
}
