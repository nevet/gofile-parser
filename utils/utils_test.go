package utils

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LineTestCase struct {
	name           string
	data           string
	expectedResult []string
}

func TestGetNextLine(t *testing.T) {
	var testData = []LineTestCase{
		{
			name:           "Get Next Line Test 1 -- Simple",
			data:           "(canyou getme)",
			expectedResult: []string{"(", "canyou", "getme", ")"},
		},
		{
			name: "Get Next Line Test 2 -- Two Lines",
			data: `(canyou getme,
					another line)`,
			expectedResult: []string{"(", "canyou", "getme", ",", "another", "line", ")"},
		},
	}

	for _, test := range testData {
		buf := bufio.NewScanner(strings.NewReader(test.data))
		var actual []string

		for tokens, hasNext := GetNextLine(buf); hasNext; tokens, hasNext = GetNextLine(buf) {
			actual = append(actual, tokens...)
		}

		assert.Equal(t, test.expectedResult, actual, test.name)
	}
}
