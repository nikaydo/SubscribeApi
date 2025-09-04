package main

import (
	"fmt"
	"log"
	_ "main/docs"
	"main/internal/config"
	"main/internal/database"
	"main/internal/handlers"
	"net/http"
)

// @title Subscriptions API
// @version 1.0
// @description CRUDL API for managing subscriptions
// @host localhost:8080
// @BasePath /
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
