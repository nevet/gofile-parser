package component

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nevet/parser/utils"
)

type ParseTestData struct {
	name           string
	data           string
	expectedResult Definition
}

func TestDefinitionParse(t *testing.T) {
	var testData = []ParseTestData{
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
			expectedResult: Definition{"string": []string{"name", "work"}, "bool": []string{"another"}},
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

		fmt.Println(definition)

		assert.Equal(t, test.expectedResult, definition, test.name)

		// tokens should possess only 1 item now
		assert.Equal(t, 1, len(tokens), test.name)
		assert.Equal(t, ")", tokens[0], test.name)

		// next buf.Scan should return false
		assert.Equal(t, false, buf.Scan(), test.name)
	}
}
