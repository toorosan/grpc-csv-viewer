package rest

import (
	"net/http"
)

// ListenAndServeREST initializes REST server handling UI requests.
func ListenAndServeREST(restListenAddr string) error {
	registerHandlers()

	return http.ListenAndServe(restListenAddr, nil)
}
