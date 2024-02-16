package interfaces

import "github.com/go-chi/chi/v5"

type BoardHandler interface {
	Routes() chi.Router
}
