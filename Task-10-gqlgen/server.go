package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"taskManager/graph"
	"taskManager/internal/database"
)

func main() {
	database.Connect()
	database.Migrate()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Println("Server is running on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
