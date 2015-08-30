package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/scanner"
)

func tokenize(s *scanner.Scanner) []string {
	var (
		tok    rune
		tokens []string
		text   string
	)

	for tok != scanner.EOF {
		tok = s.Scan()
		text = s.TokenText()

		if text != "" {
			tokens = append(tokens, s.TokenText())
		}
	}

	return tokens
}

func getLineToken(line string) []string {
	var tokenScanner scanner.Scanner

	tokenScanner.Init(strings.NewReader(line))

	return tokenize(&tokenScanner)
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

func GetNextLine(buf *bufio.Scanner) ([]string, bool) {
	if buf.Scan() {
		return getLineToken(buf.Text()), true
	}

	return nil, false
}

func GetScanner(file string) (*bufio.Scanner, error) {
	stream, err := os.Open(file)
	defer func() {
		if err != nil {
			err = stream.Close()
		}
	}()

	return bufio.NewScanner(stream), nil
}

func PreCheck(file string, noFormatCheck, noErrorCheck bool) (err error) {
	if !noFormatCheck {
		err = preCheckFmt(file)
	} else {
		fmt.Println("Skip format check")
	}

	if err != nil {
		return err
	}

	if !noErrorCheck {
		return preCheckErr(file)
	} else {
		fmt.Println("Skip error check")
	}

	return nil
}
