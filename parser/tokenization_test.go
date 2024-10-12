package parser_test

import (
	"reflect"
	"testing"

	"github.com/nullsploit01/cc-json-parser/parser"
)

func TestNextToken(t *testing.T) {
	input := `{"key": "value", "number": 123, "boolean": true, "nullValue": null}`
	expectedTokens := []parser.Token{
		{parser.TknLeftBrace, "{"},
		{parser.TknString, "key"},
		{parser.TknColon, ":"},
		{parser.TknString, "value"},
		{parser.TknComma, ","},
		{parser.TknString, "number"},
		{parser.TknColon, ":"},
		{parser.TknNumber, "123"},
		{parser.TknComma, ","},
		{parser.TknString, "boolean"},
		{parser.TknColon, ":"},
		{parser.TknBoolean, "true"},
		{parser.TknComma, ","},
		{parser.TknString, "nullValue"},
		{parser.TknColon, ":"},
		{parser.TknNull, "null"},
		{parser.TknRightBrace, "}"},
		{parser.TknEOF, ""},
	}

	lexer := parser.NewLexer(input)

	for i, expected := range expectedTokens {
		token, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("Error on token %d: %v", i, err)
		}
		if !reflect.DeepEqual(token, expected) {
			t.Errorf("Token %d - Expected: %v, Got: %v", i, expected, token)
		}
	}
}

func TestReadString(t *testing.T) {
	input := `"hello"`
	lexer := parser.NewLexer(input)

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatalf("Error reading string: %v", err)
	}

	expected := parser.Token{parser.TknString, "hello"}
	if !reflect.DeepEqual(token, expected) {
		t.Errorf("Expected %v, got %v", expected, token)
	}
}

func TestReadStringWithEscape(t *testing.T) {
	input := `"escaped\"quote\""`
	lexer := parser.NewLexer(input)

	token, err := lexer.NextToken()
	if err != nil {
		t.Fatalf("Error reading string with escape: %v", err)
	}

	expected := parser.Token{parser.TknString, `escaped"quote"`}
	if !reflect.DeepEqual(token, expected) {
		t.Errorf("Expected %v, got %v", expected, token)
	}
}

func TestReadNumber(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"123", "123"},
		{"-123", "-123"},
		{"3.14159", "3.14159"},
		{"1e10", "1e10"},
		{"1.23e-10", "1.23e-10"},
	}

	for _, testCase := range testCases {
		lexer := parser.NewLexer(testCase.input)
		token, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("Error reading number from input %s: %v", testCase.input, err)
		}

		expectedToken := parser.Token{parser.TknNumber, testCase.expected}
		if !reflect.DeepEqual(token, expectedToken) {
			t.Errorf("For input %s: expected %v, got %v", testCase.input, expectedToken, token)
		}
	}
}

func TestInvalidNumberLeadingZero(t *testing.T) {
	input := "0123"
	lexer := parser.NewLexer(input)

	_, err := lexer.NextToken()
	if err == nil {
		t.Errorf("Expected error for number with leading zero, but got none")
	}
}

func TestInvalidEscapeSequence(t *testing.T) {
	input := `"invalid\escape"`
	lexer := parser.NewLexer(input)

	_, err := lexer.NextToken()
	if err == nil {
		t.Errorf("Expected error for invalid escape sequence, but got none")
	}
}

func TestBooleanAndNull(t *testing.T) {
	testCases := []struct {
		input    string
		expected parser.Token
	}{
		{"true", parser.Token{parser.TknBoolean, "true"}},
		{"false", parser.Token{parser.TknBoolean, "false"}},
		{"null", parser.Token{parser.TknNull, "null"}},
	}

	for _, testCase := range testCases {
		lexer := parser.NewLexer(testCase.input)
		token, err := lexer.NextToken()
		if err != nil {
			t.Fatalf("Error reading value from input %s: %v", testCase.input, err)
		}

		if !reflect.DeepEqual(token, testCase.expected) {
			t.Errorf("For input %s: expected %v, got %v", testCase.input, testCase.expected, token)
		}
	}
}
