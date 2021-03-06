package component

import (
	"bufio"
	"fmt"

	"github.com/nevet/parser/utils"
)

type Definition map[string][]string

func parseDefinitionItem(lineTokens *[]string) *Definition {
	// we have 2 format here:
	// 1. with bracket: (def
	// 2. no bracket: def

	// we may have 4 cases for definition:
	// 1. name type[,)]
	// 2. n1, n2 type[,)]
	// 3. type)
	// 4. t1, t2)
	definition := make(Definition)

	// skip starting open bracket
	if (*lineTokens)[0] == "(" || (*lineTokens)[0] == "," {
		*lineTokens = (*lineTokens)[1:]
	}

	temp := []string{(*lineTokens)[0]}
	curRune := 1

	for (*lineTokens)[curRune] == "," {
		temp = append(temp, (*lineTokens)[curRune+1])
		curRune += 2
	}

	// if current token is close bracket, then all names should be type;
	if (*lineTokens)[curRune] == ")" {
		for _, v := range temp {
			definition[v] = nil
		}
	} else
	// construct type name until we hit ")", all names should be parameter name.
	{
		typeName := ""

		for (*lineTokens)[curRune] != ")" && (*lineTokens)[curRune] != "," {
			typeName += (*lineTokens)[curRune]
			curRune++
		}

		definition[typeName] = temp
	}

	*lineTokens = (*lineTokens)[curRune:]

	return &definition
}

func (def *Definition) Append(another *Definition) {
	for k, v := range *another {
		(*def)[k] = append((*def)[k], v...)
	}
}

func (def *Definition) Dump() {
	for key, names := range *def {
		fmt.Printf("%v: ", key)
		for _, name := range names {
			fmt.Printf("%v ", name)
		}

		fmt.Println("")
	}
}

func (definition *Definition) Parse(buf *bufio.Scanner, lineTokens *[]string) {
	for (*lineTokens)[0] != ")" {
		if len(*lineTokens) == 1 {
			*lineTokens, _ = utils.GetNextLine(buf)
		} else if (*lineTokens)[0] == "," {
			*lineTokens = (*lineTokens)[1:]
		} else {
			definition.Append(parseDefinitionItem(lineTokens))
		}
	}

	*lineTokens = (*lineTokens)[1:]
}
