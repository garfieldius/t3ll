
"use strict";

module.exports = require("./change-buffer")(function (data) {
    return "(function(){\n\n" + data + "\n\n})();\n";
});
