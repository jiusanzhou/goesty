package goesty

import (
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"
)

func echo(name string) string {
	fmt.Println("Name:", name)
	return name
}

func TestMushNew(t *testing.T) {
	req, _ := http.NewRequest("GET", "/echo/simple-0", nil)
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.Handle("/echo/simple-0", MustNew(echo))
	mux.ServeHTTP(rr, req)
}

func BenchmarkMustNew(b *testing.B) {
	req, _ := http.NewRequest("GET", "/echo/simple-0", nil)
	mux := http.NewServeMux()
	mux.Handle("/echo/simple-0", MustNew(echo))

	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(nil, req)
	}
}