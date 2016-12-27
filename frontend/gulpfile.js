var gulp = require("gulp");
var handlebars = require("gulp-compile-handlebars");
var sass = require("gulp-sass");
var del = require("del");
var bs = require("browser-sync").create();
var pump = require("pump");
var uncss = require("gulp-uncss");
var csso = require("gulp-csso");
var uglify = require("gulp-uglify");
var concat = require("gulp-concat");
var rename = require("gulp-rename");
var inline = require("./scripts/inline");
var chtml = require("./scripts/compress-html");
var wrap = require("./scripts/wrap");

gulp.task("css-live", ["css-dev", "html-dev"], function (cb) {
    pump([
        gulp.src("./_dev/assets/styles.css"),
        uncss({
            html: ["./_dev/*.html"]
        }),
        csso(),
        rename("styles.css"),
        gulp.dest("./_live/assets/")
    ], cb);
});

gulp.task("css-dev", function (cb) {
    pump([
        gulp.src("scss/styles.scss"),
        sass(),
        rename("styles.css"),
        gulp.dest("./_dev/assets/"),
        bs.stream()
    ], cb);
});

gulp.task("js-live", function (cb) {
    pump([
        gulp.src([
            "js/globals.js",
            "js/livedata.js",
            "js/dom.js",
            "js/xhr.js",
            "js/render.js",
            "js/callbacks.js",
            "js/_init.js"
        ]),
        concat("scripts.js"),
        wrap(),
        uglify({
            preserveComments: function () {
                return false;
            }
        }),
        gulp.dest("./_live/assets/")
    ], cb);
});

gulp.task("js-dev", function (cb) {
    pump([
        gulp.src([
            "js/globals.js",
            "js/testdata.js",
            "js/dom.js",
            "js/xhr.js",
            "js/render.js",
            "js/callbacks.js",
            "js/_init.js"
        ]),
        concat("scripts.js"),
        wrap(),
        gulp.dest("./_dev/assets/"),
        bs.stream()
    ], cb);
});

gulp.task("js-debug", function (cb) {
    pump([
        gulp.src([
            "js/globals.js",
            "js/livedata.js",
            "js/dom.js",
            "js/xhr.js",
            "js/render.js",
            "js/callbacks.js",
            "js/_init.js"
        ]),
        concat("scripts.js"),
        gulp.dest("./_dev/assets/"),
        bs.stream()
    ], cb);
});

gulp.task("html-dev", function (cb) {
    pump([
        gulp.src(["templates/*.hbs"]),
        handlebars({dev: true}, {}),
        rename("editor_tmp.html"),
        gulp.dest("./_dev/")
    ], cb);
});

gulp.task("html-debug", ["css-dev", "js-debug"], function (cb) {
    pump([
        gulp.src(["templates/*.hbs"]),
        handlebars({dev: true}, {}),
        inline(true),
        rename("editor.html"),
        gulp.dest("./_dev/")
    ], cb);
});

gulp.task("html-live", ["js-live", "css-live"], function (cb) {
    pump([
        gulp.src(["templates/*.hbs"]),
        handlebars({dev: false}, {}),
        inline(),
        chtml(),
        rename("editor.html"),
        gulp.dest("./_live/")
    ], cb);
});

gulp.task("clean", function (cb) {
    del(["./_dev/**.*", "./_live/**.*"], {force: true}).then(function () {
        cb();
    });
});

gulp.task("debug", ["html-debug"]);
gulp.task("dev", ["html-dev", "css-dev", "js-dev"]);
gulp.task("default", ["dev"]);

gulp.task("watch", ["dev"], function () {
    bs.init({
        server: {
            baseDir: "./_dev/"
        },
        startPath: "/editor_tmp.html",
        middleware: [
            {
                route: "/save",
                handle: function (req, res) {
                    res.setHeader('Content-Type', 'application/json;charset=UTF-8');
                    if (Math.random() * 100 < 10) {
                        res.statusCode = 503;
                        res.end('{"success":false,"message":"Error during save"}')
                    } else {
                        res.end('{"success":true,"message":"File saved successfully"}');
                    }
                }
            },
            {
                route: "/quit",
                handle: function (req, res) {
                    res.setHeader('Content-Type', 'application/json;charset=UTF-8');
                    res.end('{"success":true,"message":""}');
                }
            }
        ]
    });

    gulp.watch("scss/*.scss", ["css-dev"]);
    gulp.watch("js/*.js", ["js-dev"]);
    gulp.watch(["templates/*.hbs"], ["html-dev"]).on("change", bs.reload);
});

gulp.task("live", ["html-live"]);

