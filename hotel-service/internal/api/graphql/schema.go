package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/tfgoztok/hotel-service/internal/service"
)

// GraphQLService is a struct that holds the hotel service.
type GraphQLService struct {
	hotelService *service.HotelService // Reference to the hotel service for data retrieval.
}

// NewGraphQLService initializes a new GraphQLService with the provided hotel service.
func NewGraphQLService(hotelService *service.HotelService) *GraphQLService {
	return &GraphQLService{hotelService: hotelService} // Return a new instance of GraphQLService.
}

// Schema defines the GraphQL schema for the service, including types and queries.
func (s *GraphQLService) Schema() (graphql.Schema, error) {
	// Define the Hotel type with its fields.
	hotelType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Hotel", // Name of the GraphQL type.
		Fields: graphql.Fields{
			"id":              &graphql.Field{Type: graphql.String}, // Field for hotel ID.
			"officialName":    &graphql.Field{Type: graphql.String}, // Field for official name.
			"officialSurname": &graphql.Field{Type: graphql.String}, // Field for official surname.
			"companyTitle":    &graphql.Field{Type: graphql.String}, // Field for company title.
			"location":        &graphql.Field{Type: graphql.String}, // Field for location.
		},
	})

	// Define the Contact type with its fields.
	contactType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Contact", // Name of the GraphQL type.
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.String}, // Field for contact ID.
			"HotelID": &graphql.Field{Type: graphql.String}, // Field for associated hotel ID.
			"Type":    &graphql.Field{Type: graphql.String}, // Field for contact type.
			"Content": &graphql.Field{Type: graphql.String}, // Field for contact content.
		},
	})

	// Define the Query type with its fields.
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query", // Name of the query type.
		Fields: graphql.Fields{
			"hotelsByLocation": &graphql.Field{
				Type: graphql.NewList(hotelType), // Return a list of hotels.
				Args: graphql.FieldConfigArgument{
					"location": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)}, // Required location argument.
				},
				Resolve: s.resolveHotelsByLocation, // Resolver function for this query.
			},
			"contactsByLocation": &graphql.Field{
				Type: graphql.NewList(contactType), // Return a list of contacts.
				Args: graphql.FieldConfigArgument{
					"location": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)}, // Required location argument.
				},
				Resolve: s.resolveContactsByLocation, // Resolver function for this query.
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

func (s *GraphQLService) resolveContactsByLocation(p graphql.ResolveParams) (interface{}, error) {
	location, ok := p.Args["location"].(string) // Extract the location argument.
	if !ok {
		return nil, nil // Return nil if the location is not valid.
	}
	return s.hotelService.GetContactsByLocation(p.Context, location) // Call the hotel service to get contact by location.
}
