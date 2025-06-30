package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mahi7000/bookstore-go-app/handlers"
	"github.com/mahi7000/bookstore-go-app/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Unable to connect with database:", err)
	}

	apiCfg := handlers.ApiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300, 
	}))
	
	v1Router := chi.NewRouter()
	v1Router.Post("/register", apiCfg.HandlerCreateUser)
	v1Router.Get("/me", apiCfg.AuthMiddleware(apiCfg.HandlerGetUser))
	v1Router.Post("/login", apiCfg.HandlerLoginUser)
	v1Router.Post("/logout", apiCfg.AuthMiddleware(apiCfg.HandlerLogoutUser))

	v1Router.Post("/book", apiCfg.AuthMiddleware(apiCfg.HandlerAddBook))

	router.Mount("/v1", v1Router)

	server := http.Server {
		Handler: router,
		Addr: ":" + portString,
	}

	server.ListenAndServe()
	log.Printf("Server starting on port %v", portString)
}