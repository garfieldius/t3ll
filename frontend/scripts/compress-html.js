// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

"use strict";

module.exports = require("./change-buffer")(function (data) {
	const parts = data.toString().split("<!--")
	const content = [parts.shift()];

	while (parts.length) {
		const chunk = parts.shift()
		content.push(chunk.substr(chunk.indexOf("-->")+3))
	}

	return content
		.join("")
		.replace(/>\s+</g, "><")
		.replace(/\s+/g, " ")
		.replace(/>\s+\{/g, ">{")
		.replace(/\}\s+</g, "}<")
		.replace(/\}\s+\{</g, "} {")
		.trim();
});
