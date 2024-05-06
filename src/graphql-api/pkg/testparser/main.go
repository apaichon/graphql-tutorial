package main

import (
	"fmt"
	"log"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// Sample GraphQL query
var query = `
{
	contacts {
	  data: getPagination(searchText: "", page: 1, pageSize: 50) {
		contacts {
		  contact_id
		  name
		  first_name
		  last_name
		  gender_id
		}
		pagination {
		  page
		  totalItems
		  totalPages
		  hasNext
		}
	  }
	}
  }
`

// Parse the GraphQL query into an AST
func parseGraphQLQuery(queryString string) (*ast.Document, error) {
	src := source.NewSource(&source.Source{
		Body: []byte(queryString),
	})

	document, err := parser.Parse(parser.ParseParams{
		Source: src,
	})
	if err != nil {
		return nil, err
	}

	return document, nil
}

// Recursively generate the dot notation for field hierarchy
func generateDotNotation(selectionSet *ast.SelectionSet, parent string, result *[]string) {
	// Iterate over selections in the selection set
	for _, selection := range selectionSet.Selections {
		if field, ok := selection.(*ast.Field); ok {
			// Create the full field path
			fullPath := parent + "." + field.Name.Value

			// If it's a top-level field, remove the leading dot
			if parent == "" {
				fullPath = field.Name.Value
			}

			// Add the full path to the result list
			*result = append(*result, fullPath)

			// If there's a nested selection set, recurse into it
			if field.SelectionSet != nil {
				generateDotNotation(field.SelectionSet, fullPath, result)
			}
		}
	}
}

func main() {
	// Parse the query into an AST
	document, err := parseGraphQLQuery(query)
	if err != nil {
		log.Fatalf("Failed to parse GraphQL query: %v", err)
	}

	// Initialize a list to store the dot notation representation
	var dotNotation []string

	// Find the first operation and traverse its selection set
	for _, definition := range document.Definitions {
		if operation, ok := definition.(*ast.OperationDefinition); ok {
			// Recursively generate the dot notation
			generateDotNotation(operation.SelectionSet, "", &dotNotation)
		}
	}

	// Print the dot notation representation
	fmt.Println("Field hierarchy in dot notation:")
	for _, path := range dotNotation {
		fmt.Println(path)
	}

}
