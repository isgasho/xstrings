// Copyright 2015 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package xstrings

import (
	"bytes"
	"unicode/utf8"
)

// ExpandTabs can expand tabs ('\t') rune in str to one or more spaces dpending on
// current column and tabSize.
// The column number is reset to zero after each newline ('\n') occurring in the str.
//
// ExpandTabs uses RuneWidth to decide rune's width.
// For example, CJK characters will be treated as two characters.
//
// If tabSize <= 0, ExpandTabs panics with error.
//
// Samples:
//     ExpandTabs("a\tbc\tdef\tghij\tk", 4) => "a   bc  def ghij    k"
//     ExpandTabs("abcdefg\thij\nk\tl", 4)  => "abcdefg hij\nk   l"
//     ExpandTabs("z中\t文\tw", 4)           => "z中 文  w"
func ExpandTabs(str string, tabSize int) string {
	if tabSize <= 0 {
		panic("tab size must be positive")
	}

	var r rune
	var i, size, column, expand int
	var output *bytes.Buffer

	orig := str

	for len(str) > 0 {
		r, size = utf8.DecodeRuneInString(str)

		if r == '\t' {
			expand = tabSize - column%tabSize

			if output == nil {
				output = allocBuffer(orig, str)
			}

			for i = 0; i < expand; i++ {
				output.WriteByte(byte(' '))
			}

			column += expand
		} else {
			if r == '\n' {
				column = 0
			} else {
				column += RuneWidth(r)
			}

			if output != nil {
				output.WriteRune(r)
			}
		}

		str = str[size:]
	}

	if output == nil {
		return orig
	}

	return output.String()
}

// LeftJustify returns a string with pad string at right side if str's rune length is smaller than length.
// If str's rune length is larger than length, str itself will be returned.
//
// If pad is an empty string, str will be returned.
//
// Samples:
//     LeftJustify("hello", 4, " ")    => "hello"
//     LeftJustify("hello", 10, " ")   => "hello     "
//     LeftJustify("hello", 10, "123") => "hello12312"
func LeftJustify(str string, length int, pad string) string {
	l := Len(str)

	if l >= length || pad == "" {
		return str
	}

	var r rune
	var size int

	paddingLen := Len(pad)
	remains := length - l
	repeats := remains / paddingLen

	output := &bytes.Buffer{}
	output.Grow(len(str) + (repeats+1)*len(pad))
	output.WriteString(str)

	for i := 0; i < repeats; i++ {
		output.WriteString(pad)
	}

	remains = remains % paddingLen

	if remains != 0 {
		for i := 0; i < remains; i++ {
			r, size = utf8.DecodeRuneInString(pad)
			output.WriteRune(r)
			pad = pad[size:]
		}
	}

	return output.String()
}
