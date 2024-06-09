package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/musab-olurode/lis_backend/database"
	"github.com/musab-olurode/lis_backend/routes"
	"github.com/musab-olurode/lis_backend/utils"
)

func main() {
	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	database.ConfigureDB()
	utils.ConfigureCloudinary()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithErr(w, http.StatusNotFound, "requested resource not found")
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		method := r.Method
		utils.RespondWithErr(w, http.StatusMethodNotAllowed, fmt.Sprintf("Cannot %s %v", method, url.Path))
	})

	v1Router := routes.ConfigureRouter()
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server is listening on port %s", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
