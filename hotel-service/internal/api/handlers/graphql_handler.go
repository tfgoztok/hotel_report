package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/tfgoztok/hotel-service/internal/api/graphql"
)

// GraphQLHandler struct holds the schema for handling GraphQL requests
type GraphQLHandler struct {
	schema graphql.Schema // The GraphQL schema
}

// NewGraphQLHandler initializes a new GraphQLHandler with the provided GraphQL service
func NewGraphQLHandler(graphqlService *graphql.GraphQLService) (*GraphQLHandler, error) {
	schema, err := graphqlService.Schema() // Retrieve the schema from the service
	if err != nil {
		return nil, err // Return an error if schema retrieval fails
	}
	return &GraphQLHandler{schema: schema}, nil // Return a new GraphQLHandler instance
}

// ServeHTTP handles incoming HTTP requests for GraphQL
func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`         // The GraphQL query string
		OperationName string                 `json:"operationName"` // The name of the operation
		Variables     map[string]interface{} `json:"variables"`     // Variables for the query
	}

	// Decode the JSON request body into the params struct
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Return a bad request error if decoding fails
		return
	}

	// Execute the GraphQL query and get the result
	result := graphql.Do(graphql.Params{
		Schema:         h.schema,          // The schema to use for execution
		RequestString:  params.Query,      // The query string
		OperationName:  params.OperationName, // The operation name
		VariableValues: params.Variables,  // The variable values
		Context:        r.Context(),       // The request context
	})

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Encode the result as JSON and send it in the response
	json.NewEncoder(w).Encode(result)
}
