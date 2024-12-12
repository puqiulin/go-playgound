package main

import (
	"fmt"
	"go-playground/jsonparser"
)

func main() {
	// Example usage
	jsonStr := `{"name": "John Doe", "age": 30, "city": "New York", "hobbies": ["reading", "swimming"]}`
	parser := jsonparser.NewParser(jsonStr)
	result, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	// Use the parsed JSON data (result)
	fmt.Printf("%v", result)
}
