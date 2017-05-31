package split

import (
	"bufio"
	"strings"
)

// README PLEASE
// This part is different than what you may expect.
// I didn't know what is `Concatened SMS`` so I've made a `graceful` splitting.
// Splitter & RecursiveSplitter is the same but implemented in a different way.
// Now, I know it. Thanks to wikipedia.
// https://en.wikipedia.org/wiki/Concatenated_SMS

const smsMaxLength = 160
const whitespace = " "
const runeWhitespace = ' '

// Splitter explodes longer than 160 chars text to multiple parts.
func Splitter(message string) []string {
	messageParts := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(message))
	scanner.Split(gracefulSplit)

	for scanner.Scan() {
		part := strings.TrimLeft(scanner.Text(), whitespace)
		messageParts = append(messageParts, part)
	}
	return messageParts
}

// RecursiveSplitter same like Splitter but I wanted to measure which one is better
func RecursiveSplitter(text string, parts []string) []string {
	message := TrimMessage(text, smsMaxLength)
	parts = append(parts, message)
	if len(text) <= smsMaxLength {
		return parts
	}

	remainingText := text[len(message)+1:]
	return RecursiveSplitter(remainingText, parts)
}

// TrimMessage bite the first n chars from the given text, where n = maxLength
func TrimMessage(text string, maxLength int) string {
	if len(text) < maxLength {
		return text
	}

	if text[maxLength-1:maxLength] == whitespace {
		// TODO a következő msg text[maxLength:]-től kezdődik
		return text[:maxLength-1]
	}

	return TrimMessage(text, maxLength-1)
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
