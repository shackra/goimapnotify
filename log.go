package main

import (
	"bufio"
	"io"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	emailRegexp           = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	detectPasswordInLOGIN = regexp.MustCompile(`^(.*LOGIN\s+\S+\s+)"[^"]+"(.*)$`)
)

func censorCredentials(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		censoredLine := censorEmailAddress(censorPasswordInLogin(line))

		_, err := out.Write([]byte(censoredLine + "\n"))
		if err != nil {
			logrus.WithError(err).Error("unable to write censored line")
		}
	}
}

func censorPasswordInLogin(in string) string {
	matches := detectPasswordInLOGIN.FindStringSubmatch(in)

	if len(matches) == 0 {
		return in
	}

	return matches[1] + `"****"` + matches[2]
}

func censorEmailAddress(in string) string {
	matches := emailRegexp.FindAllString(in, -1)

	if len(matches) == 0 {
		return in
	}

	return emailRegexp.ReplaceAllString(in, "*******@*****.***")
}
