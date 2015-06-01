package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func requestConfigValue(value *string, description string) {
	log.Warning(description)
	log.Info("[%v]", *value)
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _, err := reader.ReadLine()
	if err != nil {
		log.Critical("Error: %v", err)
	}
	log.Debug("Value provided: %v", string(input))
	log.Debug("Value provided: %v", input)
	log.Debug("Value length: %d", len(input))
	if len(input) != 0 {
		*value = string(input)
	}
	log.Notice("Using: %v", *value)
}

// general make http request

func executeCmd(cmdName string, cmdArgs ...string) {
	cmd := exec.Command(cmdName, cmdArgs...)

	log.Notice("> Executing: %s", strings.Join(cmd.Args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("%s\nThis ended abruptly.", err)
	}
}

/*
func handle(e error) {
    if e != nil {
        panic(e)
    }
}
*/
