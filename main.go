package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/nevet/parser/component"
	"github.com/nevet/parser/utils"
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

func parseFunc(buf *bufio.Scanner, lineTokens []string) (function *Function) {

}

func parseVar(buf *bufio.Scanner, lineTokens []string) *Definition {
	definition := make(Definition)

	return &definition
}

func Parse(buf *bufio.Scanner) (*GoFile, error) {
	state := StateIdle

	goFile := &GoFile{}

	for tokens, hasNext := utils.GetNextLine(buf); hasNext; tokens, hasNext := utils.GetNextLine(buf) {
		switch state {
		case StateIdle:
			goFile.Package = tokens[1]
			state = StateInPackage

		case StateInPackage:
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
	if err := utils.PreCheck(Filename, *noFormatCheck, *noErrorCheck); err != nil {
		log.Fatal(err)
	}

	buf, err := utils.GetScanner(Filename)

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
