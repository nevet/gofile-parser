package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/scanner"

	"github.com/nevet/parser/structs"
)

const (
	Filename = "test-file.go"

	StateIdle      = 0 // starting state, waiting for package name
	StateInPackage = 1 // after package name, before EOF, outside any funcs
	StateInFunc    = 2 // inside a func
)

var (
	operatorStack stack.Stack
)

type GoFile struct {
	Package string
	Import  []string
	Const   []string
	Var     []string
	Type    []string
	Func    []string
}

func (goFile *GoFile) Dump() {
	fmt.Println("File Info:")

	fmt.Println("\nPacakge Name: " + goFile.Package + "\n")

	fmt.Println("\nImport:")
	if len(goFile.Import) == 0 {
		fmt.Println("Not Found")
	} else {
		for _, v := range goFile.Import {
			fmt.Println(v)
		}
	}

	fmt.Println("\nConst:")
	if len(goFile.Import) == 0 {
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

/*
func parseOneOrMoreLine(liner *scanner.Scanner, curLineTokens []string) []string {
	return nil
}

func parseOneLine(liner *scanner.Scanner, curLineTokens []string) []string {
	return nil
}

func skipCurrentBlock(liner *scanner.Scanner) {
	return nil
}
*/
func Exam(file string) error {
	// execute "go fmt", "go lint" and "go vet" on the file
	err := exec.Command("goimports", "-d", file).Run()

	if err != nil {
		return err
	} else {
		fmt.Println("Format check passed.")
	}

	err = exec.Command("errcheck", "-blank=false", "-asserts=true", "-ignore=Walk", file).Run()

	if err != nil {
		return err
	} else {
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
			_ = getLineToken(buf.Text())

			// switch tokens[0] {
			// case "import":
			// 	append(goFile.Import, parseOneOrMoreLine(buf, tokens))
			// case "const":
			// 	append(goFile.Const, parseOneOrMoreLine(buf, tokens))
			// case "var":
			// 	append(goFile.Var, parseOneOrMoreLine(buf, tokens))
			// case "func":
			// 	append(goFile.Func, parseOneLine(buf, tokens))
			// 	skipCurrentBlock(buf)
			// case "type":
			// 	append(goFile.Type, parseOneOrMoreLine(buf, tokens))
			// }

		case StateInFunc:
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
