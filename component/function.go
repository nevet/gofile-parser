package component

import (
	"bufio"

	"github.com/nevet/parser/utils"
)

type FuncBody struct {
	Raw []string
}

type Function struct {
	MemberOf   string
	Name       string
	Parameters *Definition
	Return     *Definition
	Body       *FuncBody
}

func parseBody(buf *bufio.Scanner, lineTokens *[]string) (funcBody []string) {
	// since the function is well formatted, there are only 2 possible cases:
	// 1. { function_body }
	// 2.
	// { <--- must be on a new line with itself
	//     function_body
	// } <--- must be on a new line with itself

	// case 1
	funcBody = append(funcBody, *lineTokens...)

	// case 2
	if len(*lineTokens) == 1 {
		for *lineTokens, _ = utils.GetNextLine(buf); len(*lineTokens) != 1 || (*lineTokens)[0] != "}"; *lineTokens, _ = utils.GetNextLine(buf) {
			funcBody = append(funcBody, *lineTokens...)
		}
	}

	return
}

func (function *Function) Dump() {
	// print function
}

func (function *Function) Parse(buf *bufio.Scanner, lineTokens *[]string) {
	// parse function membership
	if (*lineTokens)[0] == "(" {
		function.MemberOf = (*lineTokens)[1]

		if (*lineTokens)[1] == "*" {
			function.MemberOf = (*lineTokens)[2]
			*lineTokens = (*lineTokens)[4:]
		}

		*lineTokens = (*lineTokens)[3:]
	}

	// parse function name
	function.Name = (*lineTokens)[0]

	// parse parameter
	*lineTokens = (*lineTokens)[1:]
	function.Parameters = &Definition{}
	function.Parameters.Parse(buf, lineTokens)

	// parse return
	if (*lineTokens)[0] == "(" {
		*lineTokens = (*lineTokens)[1:]
		function.Return = &Definition{}
		function.Return.Parse(buf, lineTokens)
	} else {
		function.Return = &Definition{(*lineTokens)[0]: nil}
	}

	// parse body
	*lineTokens = (*lineTokens)[1:]
	function.Body = &FuncBody{Raw: parseBody(buf, lineTokens)}

	return
}
