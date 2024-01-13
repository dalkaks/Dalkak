package main

import (
	"dalkak/domain/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_app_routes_healthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

  mockUserService := user.NewUserService()

	router := app.NewRouter(mockUserService)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "ok" {
		t.Errorf("want %s; got %s", "ok", rr.Body.String())
	}
}
