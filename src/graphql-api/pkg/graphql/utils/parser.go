package utils

import (
	"fmt"
	"graphql-api/pkg/graphql"
	"strings"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
)

// Parse the GraphQL query into an AST
func ParseGraphQLQuery(queryString string) (*ast.Document, error) {
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
func GenerateResolveDotNotation(selectionSet *ast.SelectionSet, parentPath string, result *[]string) {
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
                GenerateResolveDotNotation(field.SelectionSet, currentPath, result)
            }
        }
    }
}


// Convert GraphQL errors into a single string
func ErrorsToString(errors []graphql.GraphQLError) string {
    var parts []string
    for _, err := range errors {
        part := fmt.Sprintf("Message: %s", err.Message)

        if len(err.Locations) > 0 {
            locs := make([]string, len(err.Locations))
            for i, loc := range err.Locations {
                locs[i] = fmt.Sprintf("Line %d, Column %d", loc.Line, loc.Column)
            }
            part += fmt.Sprintf(" | Locations: %s", strings.Join(locs, "; "))
        }

        if len(err.Path) > 0 {
            pathParts := make([]string, len(err.Path))
            for i, p := range err.Path {
                pathParts[i] = fmt.Sprintf("%v", p)
            }
            part += fmt.Sprintf(" | Path: %s", strings.Join(pathParts, "."))
        }

        parts = append(parts, part)
    }

    return strings.Join(parts, "\n")
}
