package main

import (
	"fmt"
	"regexp"
	// "strings"
)

func main() {
	input := " { contact { gets() { name @substring(from:1,to:2) first_name @upper } }} "

	// This will store the final key-value pairs
	resultMap := make(map[string]string)

	// Regular expression to find the pattern "field @transformation"
	re := regexp.MustCompile(`(\w+)\s+(@\w+\([^)]+\)|@\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	// Populate the map with extracted key-value pairs
	for _, match := range matches {
		if len(match) == 3 {
			key := match[1]
			value := match[2]
			resultMap[key] = value
		}
	}

	// Print the result
	for key, value := range resultMap {
		fmt.Printf("map[\"%s\"] = \"%s\"\n", key, value)
	}
}
