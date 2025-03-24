package main

import (
	"fmt"
	"graphql-crud/config"
	"graphql-crud/graph"
	"graphql-crud/middleware"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

const defaultPort = "8000"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	database := config.ConnectDatabase()

	publicSrv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Client: database}}))
	publicSrv.AddTransport(transport.Options{})
	publicSrv.AddTransport(transport.GET{})
	publicSrv.AddTransport(transport.POST{})
	publicSrv.Use(extension.Introspection{})

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Client: database}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.Use(extension.Introspection{})

	http.Handle("/public", corsMiddleware(publicSrv))
	http.Handle("/query", corsMiddleware(middleware.AuthMiddleware(srv)))
	http.Handle("/", corsMiddleware(playground.Handler("GraphQL playground", "/public")))

	fmt.Printf("Server running at http://localhost:%s/", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error serving locally")
	}
}
