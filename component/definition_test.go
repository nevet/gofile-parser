package component

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nevet/parser/utils"
)

type ParseTestCase struct {
	name           string
	data           string
	expectedResult Definition
}

type ParseItemTestCase struct {
	name           string
	data           []string
	expectedResult Definition
	remaining      []string
}

func TestDefinitionItemParse(t *testing.T) {
	var testData = []ParseItemTestCase{
		{
			name:           "Definition Item Parse Test 1 -- Normal Defition",
			data:           []string{"name", "type", ")"},
			expectedResult: Definition{"type": []string{"name"}},
			remaining:      []string{")"},
		},
		{
			name:           "Definition Item Parse Test 2 -- Normal Defition End With Comma",
			data:           []string{"name", "type", ","},
			expectedResult: Definition{"type": []string{"name"}},
			remaining:      []string{","},
		},
		{
			name:           "Definition Item Parse Test 3 -- Multiple Name Single Type",
			data:           []string{"n1", ",", "n2", "type", ")"},
			expectedResult: Definition{"type": []string{"n1", "n2"}},
			remaining:      []string{")"},
		},
		{
			name:           "Definition Item Parse Test 4 -- Multiple Name Single Type End With Comma",
			data:           []string{"n1", ",", "n2", "type", ","},
			expectedResult: Definition{"type": []string{"n1", "n2"}},
			remaining:      []string{","},
		},
		{
			name:           "Definition Item Parse Test 5 -- Single Type",
			data:           []string{"type", ")"},
			expectedResult: Definition{"type": nil},
			remaining:      []string{")"},
		},
		{
			name:           "Definition Item Parse Test 6 -- Multiple Type",
			data:           []string{"t1", ",", "t2", ")"},
			expectedResult: Definition{"t1": nil, "t2": nil},
			remaining:      []string{")"},
		},
		{
			name:           "Definition Item Parse Test 7 -- Normal Defition Start With Bracket",
			data:           []string{"(", "name", "type", ")"},
			expectedResult: Definition{"type": []string{"name"}},
			remaining:      []string{")"},
		},
	}

	for _, test := range testData {
		definition := parseDefinitionItem(&test.data)

		assert.Equal(t, test.expectedResult, *definition, test.name)
		assert.Equal(t, test.data, test.remaining, test.name)
	}
}

func TestDefinitionParse(t *testing.T) {
	var testData = []ParseTestCase{
		{
			name:           "Definition Parse Test 1 -- One Pair",
			data:           "(name type)",
			expectedResult: Definition{"type": []string{"name"}},
		},
		{
			name:           "Definition Parse Test 2 -- Two Pairs",
			data:           "(name string, another bool)",
			expectedResult: Definition{"string": []string{"name"}, "bool": []string{"another"}},
		},
		{
			name:           "Definition Parse Test 3 -- Two Pairs, Multiple Per Pair",
			data:           "(name, work string, age int)",
			expectedResult: Definition{"string": []string{"name", "work"}, "int": []string{"age"}},
		},
		{
			name: "Definition Parse Test 4 -- One Pair, Multiple Line",
			data: `(
					name string,
					)`,
			expectedResult: Definition{"string": []string{"name"}},
		},
		{
			name: "Definition Parse Test 5 -- Two Pairs, Multiple Line",
			data: `(
					name string,
					another int)`,
			expectedResult: Definition{"string": []string{"name"}, "int": []string{"another"}},
		},
		{
			name:           "Definition Parse Test 6 -- Type Aggregation In One Line",
			data:           `(name string, work string)`,
			expectedResult: Definition{"string": []string{"name", "work"}},
		},
		{
			name: "Definition Parse Test 7 -- Type Aggregation In Multiple Lines",
			data: `(
					name string,
					work string)`,
			expectedResult: Definition{"string": []string{"name", "work"}},
		},
	}

	for _, test := range testData {
		buf := bufio.NewScanner(strings.NewReader(test.data))
		tokens, _ := utils.GetNextLine(buf)
		definition := Definition{}

		definition.Parse(buf, &tokens)

		assert.Equal(t, test.expectedResult, definition, test.name)

		// tokens should possess only 1 item now
		assert.Equal(t, 0, len(tokens), test.name)

		// next buf.Scan should return false
		assert.Equal(t, false, buf.Scan(), test.name)
	}
}
