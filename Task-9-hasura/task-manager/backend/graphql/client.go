package graphql

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const graphqlEndpoint = "http://localhost:8080/v1/graphql"

// ExecuteQuery sends a GraphQL query or mutation request to Hasura.
func ExecuteQuery(query string, variables map[string]interface{}) (*http.Response, error) {
	_ = godotenv.Load()

	requestBody, _ := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	req, _ := http.NewRequest("POST", graphqlEndpoint, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))

	client := &http.Client{}
	return client.Do(req)
}
