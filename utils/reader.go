package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"text/scanner"
)

type TokenReader struct {
	BufScanner *bufio.Scanner

	lineTokens []string
	cursor     int
}

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

func (tokReader *TokenReader) CurrentLine() []string {
	return tokReader.lineTokens
}

func (tokReader *TokenReader) CurrentPos() int {
	return tokReader.cursor
}

func (tokReader *TokenReader) CurrentToken() string {
	if tokReader.lineTokens == nil {
		return ""
	}

	return tokReader.lineTokens[tokReader.cursor]
}

func (tokReader *TokenReader) MoveLineForward(forward int) error {
	for forward > 0 && tokReader.BufScanner.Scan() {
		forward--
	}

	if forward == 0 {
		tokReader.lineTokens = getLineToken(tokReader.BufScanner.Text())
		tokReader.cursor = 0

		return nil
	}

	tokReader.lineTokens = nil
	tokReader.cursor = 0

	return errors.New("End of file")
}

func (tokReader *TokenReader) MoveLineNext() error {
	return tokReader.MoveLineForward(1)
}

func (tokReader *TokenReader) MoveTokenForward(forward int) error {
	lineLength := len(tokReader.lineTokens)

	if tokReader.cursor+forward < lineLength {
		tokReader.cursor += forward
		return nil
	}

	forward -= lineLength - tokReader.cursor

	err := tokReader.MoveLineNext()

	if err != nil {
		return err
	}

	return tokReader.MoveTokenForward(forward)
}

func (tokReader *TokenReader) MoveTokenNext() error {
	return tokReader.MoveTokenForward(1)
}

func NewTokenReader(filePath string) (*TokenReader, error) {
	stream, err := os.Open(filePath)
	defer func() {
		if err != nil {
			err = stream.Close()
		}
	}()

	if err != nil {
		return nil, err
	}

	token := &TokenReader{BufScanner: bufio.NewScanner(stream)}
	err = token.MoveLineNext()

	if err != nil {
		return nil, err
	}

	return token, nil
}
