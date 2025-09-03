package main

import (
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"net/http"
)

func main() {
	e, err := config.ReadEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := database.InitBD(e)
	if err != nil {
		fmt.Println(err)
		return
	}
	mux := http.NewServeMux()
	handlers.HandlersInit(&db, e, mux)
	http.ListenAndServe(fmt.Sprintf("%s:%s", e.Host, e.Port), loggingMiddleware(mux))

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s  %s  %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
