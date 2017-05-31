package queue

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var middlewareOrder = ""

func TestMiddlewareShouldKeepTheOrder(t *testing.T) {
	clearStack()
	f := Create(middleware1, middleware2, middleware3).Then(finalHandler)
	dummyServer(t, f)
	expectMiddlewareOrder(t, "123Final")
}

func TestReusableMiddlewareShouldBeFine(t *testing.T) {
	clearStack()
	f := Create(middleware3, middleware1, middleware2, middleware3).Then(finalHandler)
	dummyServer(t, f)
	expectMiddlewareOrder(t, "3123Final")
}

func expectMiddlewareOrder(t *testing.T, expected string) {
	if middlewareOrder != expected {
		t.Error("Middleware ordering fails. Expected order is", expected, "got:", middlewareOrder)
	}
}

func clearStack() {
	middlewareOrder = ""
}

func dummyServer(t *testing.T, f http.HandlerFunc) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(rr, req)
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	reportMiddlewarePosition("Final")
}

func middleware1(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportMiddlewarePosition("1")
		next.ServeHTTP(w, r)
	})
}

func middleware2(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportMiddlewarePosition("2")
		next.ServeHTTP(w, r)
	})
}

func middleware3(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportMiddlewarePosition("3")
		next.ServeHTTP(w, r)
	})
}

func reportMiddlewarePosition(s string) {
	middlewareOrder += s
}
