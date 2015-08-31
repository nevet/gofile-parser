package utils

import (
	"bufio"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	text = `type readerLineTestCase struct {
	name           string
	data           string
	forwardNumber, something  int
	expectedResult []string
}`
)

type readerLineTestCase struct {
	name           string
	forwardNumber  int
	expectedResult []string
}

type readerTokenTestCase struct {
	name           string
	forwardNumber  int
	expectedResult string
}

func getMockReader() *TokenReader {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Scan()

	tokReader := &TokenReader{BufScanner: scanner}
	tokReader.MoveLineForward(0)

	return tokReader
}

func TestMoveLineForward(t *testing.T) {
	var testData = []readerLineTestCase{
		{
			name:           "Move Line Forward Test 1 -- Simple",
			forwardNumber:  1,
			expectedResult: []string{"name", "string"},
		},
		{
			name:           "Move Line Forward Test 2 -- Zero Forward",
			forwardNumber:  0,
			expectedResult: []string{"type", "readerLineTestCase", "struct", "{"},
		},
		{
			name:           "Move Line Forward Test 3 -- More Than One Forward",
			forwardNumber:  3,
			expectedResult: []string{"forwardNumber", ",", "something", "int"},
		},
		{
			name:           "Move Line Forward Test 4 -- More Than One Forward 2",
			forwardNumber:  4,
			expectedResult: []string{"expectedResult", "[", "]", "string"},
		},
		{
			name:           "Move Line Forward Test 5 -- Out Of Bound",
			forwardNumber:  12,
			expectedResult: nil,
		},
	}

	for _, test := range testData {
		tokReader := getMockReader()
		err := tokReader.MoveLineForward(test.forwardNumber)

		if err != nil {
			assert.Equal(t, errors.New("End of file"), err, test.name)
		}

		assert.Equal(t, test.expectedResult, tokReader.CurrentLine(), test.name)
		assert.Equal(t, 0, tokReader.CurrentPos(), test.name)
	}
}

func TestMoveTokenForward(t *testing.T) {
	var testData = []readerTokenTestCase{
		{
			name:           "Move Token Forward Test 1 -- Simple",
			forwardNumber:  1,
			expectedResult: "readerLineTestCase",
		},
		{
			name:           "Move Token Forward Test 2 -- Zero Forward",
			forwardNumber:  0,
			expectedResult: "type",
		},
		{
			name:           "Move Token Forward Test 3 -- More Than One Small Forward",
			forwardNumber:  3,
			expectedResult: "{",
		},
		{
			name:           "Move Token Forward Test 4 -- More Than One Large Forward ",
			forwardNumber:  10,
			expectedResult: "something",
		},
		{
			name:           "Move Token Forward Test 5 -- Out Of Bound",
			forwardNumber:  100,
			expectedResult: "",
		},
	}

	for _, test := range testData {
		tokReader := getMockReader()
		err := tokReader.MoveTokenForward(test.forwardNumber)

		if err != nil {
			assert.Equal(t, errors.New("End of file"), err, test.name)
		}

		assert.Equal(t, test.expectedResult, tokReader.CurrentToken(), test.name)
	}
}
