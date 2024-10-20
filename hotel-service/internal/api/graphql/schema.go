package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/tfgoztok/hotel-service/internal/service"
)

// GraphQLService struct holds the hotel service instance
type GraphQLService struct {
	hotelService *service.HotelService
}

// NewGraphQLService initializes a new GraphQLService with the provided hotel service
func NewGraphQLService(hotelService *service.HotelService) *GraphQLService {
	return &GraphQLService{hotelService: hotelService}
}

// Schema defines the GraphQL schema for the service
func (s *GraphQLService) Schema() (graphql.Schema, error) {
	// Define the Hotel type with its fields
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

	// Define the Contact type with its fields
	contactType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Contact",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.String},
				"hotelId": &graphql.Field{Type: graphql.String},
				"type":    &graphql.Field{Type: graphql.String},
				"content": &graphql.Field{Type: graphql.String},
		},
	})

	// Define the Query type with its fields and resolvers
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hotelsByLocation": &graphql.Field{
				Type: graphql.NewList(hotelType),
				Args: graphql.FieldConfigArgument{
					"location": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: s.resolveHotelsByLocation, // Resolver for fetching hotels by location
			},
			"contactsByLocation": &graphql.Field{
				Type: graphql.NewList(contactType),
				Args: graphql.FieldConfigArgument{
					"location": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: s.resolveContactsByLocation, // Resolver for fetching contacts by location
			},
		},
	})

	// Return the constructed schema
	return graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
}

// resolveHotelsByLocation fetches hotels based on the provided location argument
func (s *GraphQLService) resolveHotelsByLocation(p graphql.ResolveParams) (interface{}, error) {
	location, ok := p.Args["location"].(string)
	if !ok {
		return nil, nil // Return nil if location argument is not valid
	}
	return s.hotelService.GetHotelsByLocation(p.Context, location) // Call the service to get hotels
}

// resolveContactsByLocation fetches contacts based on the provided location argument
func (s *GraphQLService) resolveContactsByLocation(p graphql.ResolveParams) (interface{}, error) {
	location, ok := p.Args["location"].(string)
	if !ok {
		return nil, nil // Return nil if location argument is not valid
	}
	return s.hotelService.GetContactsByLocation(p.Context, location) // Call the service to get contacts
}