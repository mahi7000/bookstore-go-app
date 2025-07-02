package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mahi7000/bookstore-go-app/handlers"
	"github.com/mahi7000/bookstore-go-app/internal/hasura"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	hasuraURL := os.Getenv("HASURA_GRAPHQL_URL")
	hasuraAdminSecret := os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")

	// Initialize Hasura client
	hasuraClient := hasura.NewClient(hasuraURL, hasuraAdminSecret)

	// Initialize handlers with Hasura dependency
	handlerConfig := handlers.NewHandlerConfig(hasuraClient)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	
	// User management endpoints
	v1Router.Post("/register", handlerConfig.HandleRegister)
	v1Router.Get("/me", handlerConfig.AuthMiddleware(handlerConfig.HandleGetUser))
	v1Router.Post("/login", handlerConfig.HandleLogin)
	v1Router.Post("/logout", handlerConfig.AuthMiddleware(handlerConfig.HandleLogout))

	// Book management endpoints
	v1Router.Post("/book", handlerConfig.AuthMiddleware(handlerConfig.HandleAddBook))
	v1Router.Get("/books", handlerConfig.HandleGetBooks)
	v1Router.Get("/book/{id}", handlerConfig.HandleGetBook)
	v1Router.Put("/book/{id}", handlerConfig.AuthMiddleware(handlerConfig.HandleUpdateBook))
	v1Router.Delete("/book/{id}", handlerConfig.AuthMiddleware(handlerConfig.HandleDeleteBook))

	router.Mount("/v1", v1Router)

	server := http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}