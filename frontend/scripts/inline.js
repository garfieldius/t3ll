
"use strict";

module.exports = require("./change-buffer")(function (data, file, debug) {
    var fs = require("fs");
    var path = require("path");
    var base = path.dirname(file.path);

    function abs(p) {
        return path.join(base, "..", (debug ? "_dev" : "_live"), p)
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
