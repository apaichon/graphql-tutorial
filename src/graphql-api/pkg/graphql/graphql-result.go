package graphql

type GraphQLError struct {
    Message   string   `json:"message"`
    Locations []struct {
        Line   int `json:"line"`
        Column int `json:"column"`
    } `json:"locations"`
    Path []interface{} `json:"path"`
}

// Struct to represent a GraphQL response
type GraphQLResponse struct {
    Data   interface{}   `json:"data"`  // `interface{}` allows for flexible data parsing
    Errors []GraphQLError `json:"errors"` // List of errors in the response
}

// GraphQLRequest represents the structure of a typical GraphQL request
type GraphQLRequest struct {
	Query         string                 `json:"query"`         // The GraphQL query
	OperationName string                 `json:"operationName"` // Optional operation name
	Variables     map[string]interface{} `json:"variables"`     // Optional variables
}