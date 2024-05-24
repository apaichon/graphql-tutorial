package directives

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
)

var SubstringDirective = graphql.NewDirective(graphql.DirectiveConfig{
	Name:        "substring",
	Description: "Truncates a string to a specified length",
	Locations:   []string{graphql.DirectiveLocationField},
	Args: graphql.FieldConfigArgument{
		"from": &graphql.ArgumentConfig{
			Type: graphql.Int, // Enforces non-null integer argument
		},
		"to": &graphql.ArgumentConfig{
			Type: graphql.Int, // Optional integer argument (defaults to string length)
		},
	},
	
})


func TruncateString(ctx context.Context, obj interface{}, next graphql.FieldResolveFn, params map[string]interface{}) (interface{}, error) {
	val, ok := obj.(string)
	if !ok {
		return nil, nil // Skip non-string fields
	}

	from, err := params["from"].(int)
	if err {
		return nil, fmt.Errorf("invalid argument 'from': %v", err)
	}

	to, ok := params["to"].(int)
	if !ok {
		to = len(val) // Default to full string length
	}

	if from < 0 || from >= len(val) || to < from {
		return nil, fmt.Errorf("invalid substring range")
	}

	if to > len(val) {
		to = len(val) // Clamp to string length
	}

	return val[from:to], nil
}
