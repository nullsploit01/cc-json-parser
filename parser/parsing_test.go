package parser_test

import (
	"reflect"
	"testing"

	"github.com/nullsploit01/cc-json-parser/parser"
)

func TestParseObject(t *testing.T) {
	input := `{"key": "value", "number": 123, "boolean": true, "nullValue": null}`
	parser, err := parser.NewParser(input)
	if err != nil {
		t.Fatalf("Error initializing parser: %v", err)
	}

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error parsing object: %v", err)
	}

	expected := map[string]interface{}{
		"key":       "value",
		"number":    float64(123),
		"boolean":   true,
		"nullValue": nil,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseArray(t *testing.T) {
	input := `[1, "string", true, null, {"key": "value"}, [1, 2, 3]]`
	parser, err := parser.NewParser(input)
	if err != nil {
		t.Fatalf("Error initializing parser: %v", err)
	}

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error parsing array: %v", err)
	}

	expected := []interface{}{
		float64(1),
		"string",
		true,
		nil,
		map[string]interface{}{"key": "value"},
		[]interface{}{float64(1), float64(2), float64(3)},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestParseExtraData(t *testing.T) {
	input := `{"key": "value"} extra`
	parser, err := parser.NewParser(input)
	if err != nil {
		t.Fatalf("Error initializing parser: %v", err)
	}

	_, err = parser.Parse()
	if err == nil {
		t.Errorf("Expected error for extra data after valid JSON, but got none")
	}
}

func TestParseTrailingComma(t *testing.T) {
	input := `{"key": "value",}`
	parser, err := parser.NewParser(input)
	if err != nil {
		t.Fatalf("Error initializing parser: %v", err)
	}

	_, err = parser.Parse()
	if err == nil {
		t.Errorf("Expected error for trailing comma, but got none")
	}
}
