package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/scanner"

	"structs/stack"
)

const (
	Filename = "test.go"

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

func tokenize(s scanner.Scanner, line string) {
	s.Init(strings.NewReader(line))

	var (
		tok    rune
		tokens []string
	)

	for tok != scanner.EOF {
		tok = s.Scan()
		append(tokens, s.TokenText())
	}

	return tokens
}

func Parse(file string) (*GoFile, error) {
	stream, _ := os.Open(file)

	defer stream.Close()

	liner := bufio.NewScanner(stream)
	tokenizer := scanner.Scanner{}

	state := StateIdle

	goFile := &GoFile{}

	for liner.Scan() {
		line := liner.Text()
		tokens := tokenize(tokenizer, line)

		switch state {
		case StateIdle:
			goFile.Package = tokens[1]
			state = StateInPackage
		case StateInPackage:
			switch tokens[0] {
			case "import":

			case "const":
			case "var":
			case "func":
			case "type":
			}

		case StateInFunc:
		}
	}
}

func main() {
	// fmt, lint and vet the file before parsing
	if err := Exam(Filename); err != nil {
		log.Fatal(err)
	}

	res, err := Parse(Filename)

	log.Fatal(err)

	for _, v := range res {
		fmt.Println(v)
	}
}
