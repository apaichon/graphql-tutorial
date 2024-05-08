package main

import (
    "fmt"
    "log"
    // "strings"

    "github.com/graphql-go/graphql/language/ast"
    "github.com/graphql-go/graphql/language/parser"
    "github.com/graphql-go/graphql/language/source"
)

// Sample GraphQL query
var query = `
  mutation CreateContact($input: CreateContactInput, $input2: CreateContactInput) {
    contactMutations {
      contact1: createContact(input: $input) {
        contact_id
      }
      contact2: createContact(input: $input2) {
        contact_id
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

// Recursively generate dot notation for "resolve" fields (fields with arguments)
func generateResolveDotNotation(selectionSet *ast.SelectionSet, parentPath string, result *[]string) {
    for _, selection := range selectionSet.Selections {
        if field, ok := selection.(*ast.Field); ok {
            // Determine the full field path
            var currentPath string
            if parentPath == "" {
                currentPath = field.Name.Value
            } else {
                currentPath = fmt.Sprintf("%s.%s", parentPath, field.Name.Value)
            }

            // If the field has arguments, consider it a "resolve" field
            if len(field.Arguments) > 0 {
                *result = append(*result, currentPath)
            }

            // If there's a nested selection set, recurse into it
            if field.SelectionSet != nil {
                generateResolveDotNotation(field.SelectionSet, currentPath, result)
            }
        }
    }
}

func main() {
    // Parse the GraphQL query into an AST
    document, err := parseGraphQLQuery(query)
    if err != nil {
        log.Fatalf("Failed to parse GraphQL query: %v", err)
    }

    // List to store the dot notation for "resolve" fields
    var resolveDotNotation []string

    // Find the first operation and generate dot notation for resolve fields
    for _, definition := range document.Definitions {
        if operation, ok := definition.(*ast.OperationDefinition); ok {
            generateResolveDotNotation(operation.SelectionSet, "", &resolveDotNotation)
        }
    }

    // Display the dot notation representation of "resolve" fields
    fmt.Println("Resolve fields in dot notation:")
    for _, path := range resolveDotNotation {
        fmt.Println(path)
    }
}
