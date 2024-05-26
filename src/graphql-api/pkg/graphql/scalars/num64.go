package scalars

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// Define the Int64 custom scalar type
var Int64Type = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64",
	Description: "The `Int64` scalar type represents a signed 64-bit numeric non-fractional value.",
	Serialize: func(value interface{}) interface{} {
		switch v := value.(type) {
		case int64:
			return v
		case *int64:
			return *v
		case int:
			return int64(v)
		case *int:
			return int64(*v)
		case float64:
			return int64(v)
		case *float64:
			return int64(*v)
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch v := value.(type) {
		case int64:
			return v
		case *int64:
			return *v
		case int:
			return int64(v)
		case *int:
			return int64(*v)
		case float64:
			return int64(v)
		case *float64:
			return int64(*v)
		case string:
			i, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				return i
			}
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch v := valueAST.(type) {
		case *ast.IntValue:
			i, err := strconv.ParseInt(v.Value, 10, 64)
			if err == nil {
				return i
			}
		}
		return nil
	},
})


// Define Float64 scalar
var Float64Type = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Float64",
	Description: "The `Float64` scalar type represents signed double-precision fractional values.",
	Serialize: func(value interface{}) interface{} {
		switch v := value.(type) {
		case float64:
			return v
		case *float64:
			return *v
		case float32:
			return float64(v)
		case *float32:
			return float64(*v)
		case int:
			return float64(v)
		case *int:
			return float64(*v)
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch v := value.(type) {
		case float64:
			return v
		case *float64:
			return *v
		case float32:
			return float64(v)
		case *float32:
			return float64(*v)
		case int:
			return float64(v)
		case *int:
			return float64(*v)
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err == nil {
				return f
			}
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch v := valueAST.(type) {
		case *ast.FloatValue:
			f, err := strconv.ParseFloat(v.Value, 64)
			if err == nil {
				return f
			}
		case *ast.IntValue:
			i, err := strconv.ParseInt(v.Value, 10, 64)
			if err == nil {
				return float64(i)
			}
		}
		return nil
	},
})