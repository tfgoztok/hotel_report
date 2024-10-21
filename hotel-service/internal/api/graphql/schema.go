package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/tfgoztok/hotel-service/internal/service"
)

// GraphQLService is a struct that holds the hotel service.
type GraphQLService struct {
	hotelService *service.HotelService
}

// NewGraphQLService initializes a new GraphQLService with the provided hotel service.
func NewGraphQLService(hotelService *service.HotelService) *GraphQLService {
	return &GraphQLService{hotelService: hotelService}
}

// Schema defines the GraphQL schema for the service, including types and queries.
func (s *GraphQLService) Schema() (graphql.Schema, error) {
	// Define the Hotel type with its fields.
	hotelType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Hotel",
		Fields: graphql.Fields{
			"id":              &graphql.Field{Type: graphql.String},
			"officialName":    &graphql.Field{Type: graphql.String},
			"officialSurname": &graphql.Field{Type: graphql.String},
			"companyTitle":    &graphql.Field{Type: graphql.String},
			"location":        &graphql.Field{Type: graphql.String},
		},
	})

	// Define the Query type with its fields.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hotelsByLocation": &graphql.Field{
				Type: graphql.NewList(hotelType),
				Args: graphql.FieldConfigArgument{
					"location": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: s.resolveHotelsByLocation, // Resolver function for this query.
			},
		},
	})

	// Return the complete schema with the defined query type.
	return graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
}

// resolveHotelsByLocation is the resolver function for the hotelsByLocation query.
func (s *GraphQLService) resolveHotelsByLocation(p graphql.ResolveParams) (interface{}, error) {
	location, ok := p.Args["location"].(string) // Extract the location argument.
	if !ok {
		return nil, nil // Return nil if the location is not valid.
	}
	return s.hotelService.GetHotelsByLocation(p.Context, location) // Call the hotel service to get hotels by location.
}
