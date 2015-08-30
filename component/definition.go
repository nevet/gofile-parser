package component

import (
	"bufio"

	"github.com/nevet/parser/utils"
)

type Definition map[string][]string

func parseDefinitionItem(lineTokens *[]string) *Definition {
	definition := make(Definition)

	paraNames := []string{(*lineTokens)[0]}
	curRune := 1

	for (*lineTokens)[curRune] == "," {
		paraNames = append(paraNames, (*lineTokens)[curRune+1])
		curRune += 2
	}

	definition[(*lineTokens)[curRune]] = paraNames
	*lineTokens = (*lineTokens)[curRune+1:]

	return &definition
}

func (def *Definition) Append(another *Definition) {
	for k, v := range *another {
		(*def)[k] = append((*def)[k], v...)
	}
}

func (def *Definition) Dump() {
	// print definition
}

func (definition *Definition) Parse(buf *bufio.Scanner, lineTokens *[]string) {
	for (*lineTokens)[0] != ")" {
		if len(*lineTokens) == 1 {
			*lineTokens, _ = utils.GetNextLine(buf)
		} else {
			definition.Append(parseDefinitionItem(lineTokens))
		}
	}

	*lineTokens = (*lineTokens)[1:]
}
