
"use strict";

module.exports = require("./change-buffer")(function (data) {
    return data
        .toString()
        .replace(/>\s+</g, "><")
        .replace(/\s+/g, " ")
        .replace(/>\s+\{/g, ">{")
        .replace(/\}\s+</g, "}<")
        .replace(/\}\s+\{</g, "} {")
        .trim();
});
