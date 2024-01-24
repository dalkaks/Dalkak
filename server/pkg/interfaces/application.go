package interfaces

import (
	"net/http"
)

type Application interface {
	StartServer(port int, userService UserService) error
	NewRouter(userService UserService) http.Handler
	enableCORS(next http.Handler) http.Handler
}
