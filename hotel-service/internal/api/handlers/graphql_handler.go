package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/graphql-go/graphql"
    localGraphQL "github.com/tfgoztok/hotel-service/internal/api/graphql"
)

// GraphQLHandler struct holds the schema for handling GraphQL requests
type GraphQLHandler struct {
    schema graphql.Schema // The GraphQL schema
}

// NewGraphQLHandler initializes a new GraphQLHandler with the provided GraphQL service
func NewGraphQLHandler(graphqlService *localGraphQL.GraphQLService) (*GraphQLHandler, error) {
    schema, err := graphqlService.Schema() // Retrieve the schema from the service
    if err != nil {
        return nil, err // Return an error if schema retrieval fails
    }
    return &GraphQLHandler{schema: schema}, nil // Return a new GraphQLHandler instance
}

// ServeHTTP handles incoming HTTP requests for GraphQL
func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var params struct {
        Query         string                 `json:"query"`
        OperationName string                 `json:"operationName"`
        Variables     map[string]interface{} `json:"variables"`
    }

    if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result := graphql.Do(graphql.Params{
        Schema:         h.schema,
        RequestString:  params.Query,
        OperationName:  params.OperationName,
        VariableValues: params.Variables,
        Context:        r.Context(),
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}