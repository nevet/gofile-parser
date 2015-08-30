package main

import (
	"bufio"
	"errors"
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

	boundMapper = map[string]string{
		"(": ")",
		"{": "}",
	}
)

type GoFile struct {
	Package string
	Import  []string
	Const   []string
	Var     []string
	Type    []string
	Func    []Function
}

type Function struct {
	Name       string
	Parameters map[string]string
	Return     map[string]string
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
	if len(goFile.Var) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Var {
			fmt.Println(v)
		}
	}

	fmt.Println("\nType:")
	if len(goFile.Type) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Type {
			fmt.Println(v)
		}
	}

	fmt.Println("\nFunc:")
	if len(goFile.Func) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Func {
			fmt.Println(v)
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

func parseDefinitionPair(lineTokens []string) (mapper map[string][]string, leftOver []string) {
	paraNames := []string{lineTokens[0]}
	curRune := 1

	for lineTokens[curRune] == "," {
		paraNames = append(paraNames, lineTokens[curRune+1])
		curRune += 2
	}

	mapper[lineTokens[curRune]] = paraNames

	return mapper, lineTokens[curRune+1:]
}

func parseDefinition(buf *bufio.Scanner, lineTokens []string, bound string) (mapper map[string]string, leftOver []string) {
	// parse function parameters
	for lineTokens[0] != boundMapper[bound] {
		if len(lineTokens) == 1 {
			buf.Scan()
			lineTokens = getLineToken(buf.Text())
		} else {

		}

		for curRune < len(lineTokens) && lineTokens[curRune] != boundMapper[bound] {
			if lineTokens[curRune] == bound {
				curRune++
			} else if isMathcingName {
				curName = lineTokens[curRune]
				isMathcingName = false
				curRune++
			} else {
				curType += lineTokens[curRune]
				curRune++

				if lineTokens[curRune] == "," || lineTokens[curRune] == boundMapper[bound] {
					mapper[curName] = curType
					isMathcingName = true

					curName = ""
					curType = ""
				}
			}
		}

		if curRune == len(lineTokens) {

			curRune = 0
		}
	}

	return mapper, lineTokens[curRune+1:]
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

func parseFunc(buf *bufio.Scanner, lineTokens []string) (memberOf string, funcName *Function) {
	// parse function membership
	if lineTokens[0] == "(" {
		memberOf = lineTokens[1]

		if lineTokens[1] == "*" {
			memberOf = lineTokens[2]
			_, funcName = parseFunc(buf, lineTokens[4:])
		}

		_, funcName = parseFunc(buf, lineTokens[3:])

		return
	}

	// parse function name
	function = &Function{Name: lineTokens[0]}

	// parse parameter
	function.Parameters = parseDefinition(buf, lineTokens[1:], "(")

	if lineTokens[0] == "(" {
		returnMapper
	}
}

func Exam(file string) error {
	output, err := exec.Command("goimports", "-d", file).Output()

	if err != nil {
		return err
	} else {
		if len(output) != 0 {
			fmt.Println("Format check failed.")
			return errors.New(string(output[:]))
		}

		fmt.Println("Format check passed.")
	}

	output, err = exec.Command("errcheck", "-blank=false", "-asserts=true", "-ignore=Walk", file).Output()

	if err != nil {
		return err
	} else {
		if len(output) != 0 {
			fmt.Println("Error check failed.")
			return errors.New(string(output[:]))
		}

		fmt.Println("Error check passed.")
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
			// case "var":
			// 	append(goFile.Var, parseOneOrMoreLine(buf, tokens))
			case "func":
				goFile.Func = append(goFile.Func, parseFunc(buf, tokens[1:])...)
				// case "type":
				// 	append(goFile.Type, parseOneOrMoreLine(buf, tokens))
				// 	}

				// case StateInFunc:
			}
		}
	}

	return goFile, nil
}

func main() {
	// fmt, lint and vet the file before parsing
	if err := Exam(Filename); err != nil {
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
