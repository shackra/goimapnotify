package imap

import (
	"bytes"
	"os/exec"
)

func prepareCommand(cmd string, debug bool) *exec.Cmd {
	list := append([]string{"sh", "-c"}, cmd)

	if debug {
		// TODO: print command here
	}

	// #nosec
	command := exec.Command(list[0], list[1:]...)
	command.Stdout = nil

	return command
}

func runCommand(cmd *exec.Cmd) ([]byte, error) {
	buf, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return bytes.Trim(buf, "\n"), nil
}

func fetchThing(cmd string, debug bool) (string, error) {
	command := prepareCommand(cmd, debug)

	r, err := runCommand(command)
	if err != nil {
		return "", err
	}

	return string(r), nil
}

func fetchPassword(cmd string, debug bool) (string, error) {
	return fetchThing(cmd, debug)
}

func fetchToken(cmd string, debug bool) (string, error) {
	return fetchThing(cmd, debug)
}

func fetchUsername(cmd string, debug bool) (string, error) {
	return fetchThing(cmd, debug)
}
