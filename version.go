// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the file LICENSE or <https://opensource.org/licenses/MIT> for details

package main

// Version defines the current version of t3ll
var Version string

const helpText = `t3ll - TYPO3 Locallang Editor %s

Usage: t3ll <filename>|help|version

<filename> must be a text file in UTF-8 encoding, containing
a xliff or TYPO3 locallang XML structure.

If the filename argument is valid, but does not point to an
existing file, it is created.

There are two additional commands:

help:    Print this help text
version: Print version and license
`

const versionText = `t3ll - TYPO3 Locallang Editor %s

Copyright 2019 Georg Großberger <contact@grossberger-ge.org>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the Software
is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR
THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`
