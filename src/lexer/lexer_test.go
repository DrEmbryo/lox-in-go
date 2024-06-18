package lexer

import (
	"fmt"
	"testing"
)

func TestParseDigit(t *testing.T) {
	var tests = []struct {
		name string
		arg rune
        expect bool
    }{
        {"min border value", '0', true},
		{"mid border value", '5', true},
		{"max border value", '9', true},
		{"out of bound value", 'a', false},
		{"out of bound value", '!', false},
    }

	for _, tc := range tests {
		testname := fmt.Sprintf("%s: %s,", tc.name, string(tc.arg))
		t.Run(testname, func(t *testing.T) {
            result := parseDigit(&tc.arg)
            if result != tc.expect {
                t.Errorf("got %v, want %v", result, tc.expect)
            }
        })
	} 
	
}

func TestParseChar(t *testing.T) {
	var tests = []struct {
		name string
		arg rune
        expect bool
    }{
        {"min border value", 'a', true},
		{"mid border value", 'k', true},
		{"max border value", 'z', true},
		{"out of bound value", '5', false},
		{"out of bound value", '!', false},
    }

	for _, tc := range tests {
        testname := fmt.Sprintf("%s: %s,", tc.name, string(tc.arg))
		t.Run(testname, func(t *testing.T) {
            result := parseChar(&tc.arg)
            if result != tc.expect {
                t.Errorf("got %v, want %v", result, tc.expect)
            }
        })
	} 
	
}

func TestParseSkippable(t *testing.T) {
	var tests = []struct {
		name string
		arg rune
        expect bool
    }{
        {"border value", '\t', true},
		{"border value", '\r', true},
		{"border value", ' ', true},
		{"out of bound value", '5', false},
		{"out of bound value", 'a', false},
		{"out of bound value", '!', false},
    }

	for _, tc := range tests {
		testname := fmt.Sprintf("%s: %s,", tc.name, string(tc.arg))
		t.Run(testname, func(t *testing.T) {
            result := parseSkippable(&tc.arg)
            if result != tc.expect {
                t.Errorf("got %v, want %v", result, tc.expect)
            }
        })
	} 
	
}