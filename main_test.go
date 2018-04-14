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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripCommentMarker(t *testing.T) {
	var doTest = func(inputString, inputCommentMarker, expectedStrippedString, expectedCommentMarker string) {
		actualString, actualCommentMarker := stripCommentMarker(inputString, commentMarker(inputCommentMarker))
		assert.Equal(t, expectedStrippedString, actualString, "string")
		assert.Equal(t, commentMarker(expectedCommentMarker), actualCommentMarker, "comment")
	}
	doTest("//", "", "", "//")
	doTest("*", "*", "", "*")
	doTest("*foobar", "*", "foobar", "*")
	doTest("--foobar", "", "foobar", "--")
	doTest("--foobar", "--", "foobar", "--")
	doTest("--", "--", "", "--")
	doTest("*", "--", "*", "--")
}

func TestGetWordsLessCommentMarkers(t *testing.T) {
	var doTest = func(inputLines []string, expectedWords []string, expectedCommentMarker string) {
		actualWords, actualCommentMarker := getWordsLessCommentMarkers(inputLines)
		assert.Equal(t, expectedWords, actualWords, "words")
		assert.Equal(t, commentMarker(expectedCommentMarker), actualCommentMarker, "comment")
	}

	doTest([]string{"-- foobar --"}, []string{"foobar", "--"}, "--")
	doTest([]string{"-- foobar --", "--baz"}, []string{"foobar", "--", "baz"}, "--")
	doTest([]string{"okay", " * more", "*less"}, []string{"okay", "more", "less"}, "*")
}

func TestIndent(t *testing.T) {
	assert.Equal(t, "", getIndent("foo"))
	assert.Equal(t, "", getIndent("fo  o"))
	assert.Equal(t, "", getIndent("foo  "))
	assert.Equal(t, "  ", getIndent("  foo"))
	assert.Equal(t, " \t ", getIndent(" \t foo"))
}

func TestStringLength(t *testing.T) {
	assert.Equal(t, 0, stringLength(""))
	assert.Equal(t, 1, stringLength(" "))
	assert.Equal(t, 3, stringLength("abc"))
	assert.Equal(t, 2+*tabSize, stringLength(" \t "))
}

func TestReflow(t *testing.T) {
	doTest := func(line string, width int, expected string) {
		*targetLength = width
		actual := reflow(strings.Split(line, "\n"))
		assert.Equal(t, expected, actual)
	}

	doTest("ab", 10, "ab")
	doTest("ab cd ef\nhi jk\r\nlm no\r\n", 100, "ab cd ef hi jk lm no")
	doTest("ab cd", 2, "ab\ncd")
	doTest("-- ab cd", 2, "-- ab\n-- cd")
	doTest("aa bb cc dd ee ff", 5, "aa bb\ncc dd\nee ff")
}
