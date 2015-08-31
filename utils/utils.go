package utils

import (
	"errors"
	"fmt"
	"os/exec"
)

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
