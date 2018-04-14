// Copyright 2018 Jacob Dufault
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
	"unicode"
)

// Strings which a comment can start with.
type commentMarker string

var commentMarkers = []commentMarker{
	"#",
	"*",
	"--",
	"//",
}

var (
	tabSize      = flag.Int("tab-size", 2, "Size of a tab (\\t) character")
	targetLength = flag.Int("width", 80, "Width to reflow text to")
)

// If word starts with a comment, removes the comment. If commentType is non-nil
// then only the comment continuation marker for commentTYpe will be checked.
func stripCommentMarker(word string, commentType commentMarker) (string, commentMarker) {
	if len(commentType) > 0 {
		return strings.TrimPrefix(word, string(commentType)), commentType
	}

	for _, marker := range commentMarkers {
		if strings.HasPrefix(word, string(marker)) {
			return strings.TrimPrefix(word, string(marker)), marker
		}
	}

	return word, ""
}

// Removes all comment markers from `txt`.
func getWordsLessCommentMarkers(lines []string) ([]string, commentMarker) {
	words := make([]string, 0)

	var commentType commentMarker
	for _, line := range lines {
		for n, word := range strings.Fields(line) {
			if n == 0 {
				word, commentType = stripCommentMarker(word, commentType)
			}

			if len(word) > 0 {
				words = append(words, word)
			}
		}
	}

	return words, commentType
}

// Returns the indent for |txt|, ie, the amount of space before the first
// non-whitespace character.
func getIndent(txt string) string {
	idx := strings.IndexFunc(txt, func(c rune) bool {
		return !unicode.IsSpace(c)
	})
	return txt[0:idx]
}

// Returns the "length" of a string, taking into account that a tab may consume
// >1 width.
func stringLength(s string) int {
	var len int
	for _, c := range s {
		if c == '\t' {
			len += *tabSize
		} else {
			len++
		}
	}
	return len
}

func reflow(inputLines []string) string {
	indent := getIndent(inputLines[0])

	words, commentType := getWordsLessCommentMarkers(inputLines)
	if len(commentType) > 0 {
		// Add the comment type plus any space after the comment to indent.
		indent += string(commentType) + getIndent(inputLines[0][len(indent)+len(commentType):])
	}

	lines := make([]string, 0)
	lines = append(lines, indent)
	var hasWordInLine bool

	for _, word := range words {
		wordLength := stringLength(word)
		lineLength := stringLength(lines[len(lines)-1])

		// If we have already added a word to this line and adding the next word
		// will cause it to overflow the target length, start a new line.
		//
		// add 1 for space
		if hasWordInLine &&
			lineLength+1+wordLength > *targetLength {

			lines = append(lines, indent)
			hasWordInLine = false
		}

		// Add the world to the line
		space := ""
		if hasWordInLine {
			space = " "
		}
		lines[len(lines)-1] += space + word
		hasWordInLine = true
	}

	return strings.Join(lines, "\n")
}

func main() {
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		panic(err)
	}

	println(reflow(lines))
}
