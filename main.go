package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/scanner"

	// "github.com/nevet/parser/structs"
)

const (
	Filename = "test-file1.go"

	StateIdle      = 0 // starting state, waiting for package name
	StateInPackage = 1 // after package name, before EOF, outside any funcs
)

var (
	noFormatCheck, noErrorCheck *bool

	boundMapper = map[string]string{
		"(": ")",
		"{": "}",
	}
)

type Definition map[string][]string

func (def Definition) Append(another *Definition) {
	for k, v := range *another {
		def[k] = append(def[k], v...)
	}
}

func (def *Definition) Dump() {
	// print definition
}

type GoFile struct {
	Package   string
	Import    []string
	Const     []string
	Var       *Definition
	Type      *Definition
	Functions []*Function
}

type Function struct {
	MemberOf   string
	Name       string
	Parameters *Definition
	Return     *Definition
	Body       *FuncBody
}

func (function *Function) Dump() {
	// print function
}

type FuncBody struct {
	Raw []string
}

func (goFile *GoFile) Dump() {
	fmt.Println("File Info:")

	fmt.Println("\nPacakge Name:\n" + goFile.Package)

	fmt.Println("\nImport:")
	if len(goFile.Import) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Import {
			fmt.Println(v)
		}
	}

	fmt.Println("\nConst:")
	if len(goFile.Const) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Const {
			fmt.Println(v)
		}
	}

	fmt.Println("\nVar:")
	if len(*goFile.Var) == 0 {
		fmt.Println("Not Found")
	} else {
		goFile.Var.Dump()
	}

	fmt.Println("\nType:")
	if len(*goFile.Type) == 0 {
		fmt.Println("Not Found")
	} else {
		goFile.Type.Dump()
	}

	fmt.Println("\nFunc:")
	if len(goFile.Functions) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Functions {
			v.Dump()
		}
	}
}

func readFileToBuffer(file string) (*bufio.Scanner, error) {
	stream, err := os.Open(file)
	defer func() {
		if err != nil {
			err = stream.Close()
		}
	}()

	return bufio.NewScanner(stream), nil
}

func tokenize(s *scanner.Scanner) []string {
	var (
		tok    rune
		tokens []string
	)

	for tok != scanner.EOF {
		tok = s.Scan()
		tokens = append(tokens, s.TokenText())
	}

	return tokens
}

func getLineToken(line string) []string {
	var tokenScanner scanner.Scanner

	tokenScanner.Init(strings.NewReader(line))

	return tokenize(&tokenScanner)
}

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

func parseDefinition(buf *bufio.Scanner, lineTokens *[]string, bound string) *Definition {
	definition := make(Definition)

	for (*lineTokens)[0] != boundMapper[bound] {
		if len(*lineTokens) == 1 {
			buf.Scan()
			*lineTokens = getLineToken(buf.Text())
		} else {
			definition.Append(parseDefinitionItem(lineTokens))
		}
	}

	*lineTokens = (*lineTokens)[1:]

	return &definition
}

func parseFlag() {
	noFormatCheck = flag.Bool("nofmt", false, "Skip code format checking")
	noErrorCheck = flag.Bool("noerr", false, "Skip code error checking")

	flag.Parse()
}

func parseImport(buf *bufio.Scanner, lineTokens []string) (packageNames []string) {
	if lineTokens[0] != "(" {
		packageNames = append(packageNames, lineTokens[0])
	} else {
		for buf.Scan(); buf.Text() != ")"; buf.Scan() {
			lineTokens = getLineToken(buf.Text())

			if len(lineTokens[0]) > 0 {
				packageNames = append(packageNames, lineTokens[0])
			}
		}
	}

	return
}

func parseBody(buf *bufio.Scanner, lineTokens []string) (funcBody []string) {
	// since the function is well formatted, there are only 2 possible cases:
	// 1. { function_body }
	// 2.
	// { <--- must be on a new line with itself
	//     function_body
	// } <--- must be on a new line with itself

	// case 1
	funcBody = append(funcBody, lineTokens...)

	// case 2
	if len(lineTokens) == 1 {
		for buf.Scan(); len(lineTokens) != 1 || lineTokens[0] != "}"; buf.Scan() {
			lineTokens = getLineToken(buf.Text())
			funcBody = append(funcBody, lineTokens...)
		}
	}

	return
}

func parseFunc(buf *bufio.Scanner, lineTokens []string) (function *Function) {
	function = &Function{}

	// parse function membership
	if lineTokens[0] == "(" {
		function.MemberOf = lineTokens[1]

		if lineTokens[1] == "*" {
			function.MemberOf = lineTokens[2]
			lineTokens = lineTokens[4:]
		}

		lineTokens = lineTokens[3:]
	}

	// parse function name
	function.Name = lineTokens[0]

	// parse parameter
	lineTokens = lineTokens[1:]
	function.Parameters = parseDefinition(buf, &lineTokens, "(")

	// parse return
	if lineTokens[0] == "(" {
		lineTokens = lineTokens[1:]
		function.Return = parseDefinition(buf, &lineTokens, "(")
	} else {
		function.Return = &Definition{lineTokens[0]: nil}
	}

	// parse body
	lineTokens = lineTokens[1:]
	function.Body = &FuncBody{Raw: parseBody(buf, lineTokens)}

	return
}

func parseVar(buf *bufio.Scanner, lineTokens []string) *Definition {
	definition := make(Definition)

	return &definition
}

func preCheckErr(file string) error {
	output, err := exec.Command("errcheck", "-blank=false", "-asserts=true", "-ignore=Walk", file).Output()

	if err != nil {
		return err
	} else if len(output) != 0 {
		fmt.Println("Error check failed.")
		return errors.New(string(output[:]))
	}

	fmt.Println("Error check passed.")
	return nil
}

func preCheckFmt(file string) error {
	output, err := exec.Command("goimports", "-d", file).Output()

	if err != nil {
		return err
	} else if len(output) != 0 {
		fmt.Println("Format check failed.")
		return errors.New(string(output[:]))
	}

	fmt.Println("Format check passed.")
	return nil
}

func PreCheck(file string) (err error) {
	if !*noFormatCheck {
		err = preCheckFmt(file)
	} else {
		fmt.Println("Skip format check")
	}

	if err != nil {
		return err
	}

	if !*noErrorCheck {
		return preCheckErr(file)
	} else {
		fmt.Println("Skip error check")
	}

	return nil
}

func Parse(buf *bufio.Scanner) (*GoFile, error) {
	state := StateIdle

	goFile := &GoFile{}

	for buf.Scan() {
		switch state {
		case StateIdle:
			goFile.Package = getLineToken(buf.Text())[1]
			state = StateInPackage

		case StateInPackage:
			tokens := getLineToken(buf.Text())

			switch tokens[0] {
			case "import":
				goFile.Import = append(goFile.Import, parseImport(buf, tokens[1:])...)
			// case "const":
			// 	append(goFile.Const, parseOneOrMoreLine(buf, tokens))
			case "var":
				goFile.Var.Append(parseVar(buf, tokens[1:]))
			case "func":
				goFile.Functions = append(goFile.Functions, parseFunc(buf, tokens[1:]))
				// case "type":
				// 	append(goFile.Type, parseOneOrMoreLine(buf, tokens))
				// 	}
			}
		}
	}

	return goFile, nil
}

func main() {
	parseFlag()

	// fmt, lint and vet the file before parsing
	if err := PreCheck(Filename); err != nil {
		log.Fatal(err)
	}

	buf, err := readFileToBuffer(Filename)

	if err != nil {
		log.Fatal(err)
	}

	parsedGoFile, err := Parse(buf)

	if err != nil {
		log.Fatal(err)
	}

	parsedGoFile.Dump()
}

//
