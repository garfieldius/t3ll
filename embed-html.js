
// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

const { readFileSync, writeFileSync } = require("fs");
const { join } = require("path");

const content = readFileSync(join(__dirname, "frontend", "build", "index.html"));

let bytes = [];
const lines = [];

for (const pair of content.entries()) {
    bytes.push("0x" + pair.pop().toString(16))

    if (bytes.length >= 16) {
        lines.push("\t" + bytes.join(", "))
        bytes = []
    }
}

if (bytes.length) {
    lines.push("\t" + bytes.join(", "))
}

const data = lines.join(",\n")

writeFileSync(join(__dirname, "server", "html.go"), `// +build !debug

// Generated file, do not edit manually

package server

var html = []byte{
${data},
}
`);
