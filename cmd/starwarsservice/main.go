package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/lonmarsDev/starwars-be/graph"
	"github.com/lonmarsDev/starwars-be/internal/service/swservice"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const defaultPort = "8080"

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		loggingLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Fatalf("invalid log level value %+w", err)
		}
		logrus.SetLevel(loggingLevel)

	}
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := chi.NewRouter()
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	dbConn := os.Getenv("DB_CONNECTION")
	if dbConn == "" {
		dbConn = "mongodb://root:root@localhost:27017"
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Service: swservice.NewService(dbConn),
	}}))
	router.Handle("/", playground.Handler("Starwars", "/query"))
	router.Handle("/query", srv)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logrus.Fatal(err)
	}
}
