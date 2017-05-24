package main

import (
	"bufio"
	"strings"
)

const smsMaxLength = 160
const whitespace = " "
const runeWhitespace = ' '

// Splitter explodes longer than 160 chars text to multiple parts.
func Splitter(message string) []string {
	if len(message) <= smsMaxLength {
		s := make([]string, 1)
		s[0] = message
		return s
	}

	messageParts := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(message))
	scanner.Split(gracefulSplit)

	for scanner.Scan() {
		part := strings.TrimLeft(scanner.Text(), whitespace)
		messageParts = append(messageParts, part)
	}
	return messageParts
}

// Graceful scanner, implements bufio.SplitFunc
// Originally, I would to use bufio.ScanWords but it skips multiple whitespaces.
func gracefulSplit(data []byte, atEOF bool) (int, []byte, error) {
	i := smsMaxLength
	for {
		if i > len(data) || i < 1 {
			// I can imagine a msg without any whitespace. This would fix that case.
			return 0, data, bufio.ErrFinalToken
		}

		if data[i] == runeWhitespace {
			return i, data[:i], nil
		}
		i--
	}

}
