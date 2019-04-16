// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

"use strict";

module.exports = require("./change-buffer")(function (data) {
	return "(function(){\n\n" + data + "\n\n})();\n";
});
