package main

import (
	"fmt"
	"go-playground/jsonparser"
)

func main() {
	jsonStr := `{"name": "John Doe", "age": 30, "city": "New York", "hobbies": ["reading", "swimming"]}`
	parser := jsonparser.NewParser(jsonStr)
	result, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", result)
}
