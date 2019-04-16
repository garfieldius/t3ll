// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

"use strict";

const path = require("path");
const fs = require("fs");

module.exports = require("./change-buffer")(function (data, file) {
	const base = path.join(path.dirname(file.path), "..", "build");

	function abs(p) {
		return path.join(base, p)
	}

	data.match(/<link[^>]+href="([^"]+)">/ig).forEach(function (el) {
		let start = el.indexOf("href=\"") + 6;
		let p1 = el.substr(start);
		let p = p1.substr(0, p1.indexOf('"'));
		let d = fs.readFileSync(abs(p));

		data = data.replace(el, "<style>" + d + "</style>");
	});

	data.match(/<script[^>]+src="([^"]+)">\s*<\/script>/ig).forEach(function (el) {
		let start = el.indexOf("src=\"") + 5;
		let p1 = el.substr(start);
		let p = p1.substr(0, p1.indexOf('"'));
		let d = fs.readFileSync(abs(p))

		data = data.replace(el, "<script>" + d + "</script>");
	});

	return data;
});
