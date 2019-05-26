// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

const gulp = require("gulp");
const sass = require("gulp-sass");
const bs = require("browser-sync").create();
const csso = require("gulp-csso");
const uglify = require("gulp-uglify");
const concat = require("gulp-concat");
const rename = require("gulp-rename");
const inline = require("./scripts/inline");
const chtml = require("./scripts/compress-html");
const wrap = require("./scripts/wrap");
const when = require("gulp-if");
const isProd = process.env.NODE_ENV === "production";

gulp.task("css", () => {
	return gulp.src("scss/styles.scss")
		.pipe(sass())
		.pipe(rename("styles.css"))
		.pipe(when(isProd, csso()))
		.pipe(gulp.dest("./build/assets/"))
});

gulp.task("js", () => {
	return gulp.src([
			"js/globals.js",
			"js/dom.js",
			"js/xhr.js",
			"js/render.js",
			"js/callbacks.js",
			"js/_init.js"
		])
		.pipe(concat("scripts.js"))
		.pipe(wrap())
		.pipe(when(isProd, uglify({
			mangle: {
				toplevel: true
			},
			compress: {
				drop_console: true,
				keep_fargs: false,
				toplevel: true,
				global_defs: {
					NODE_ENV: "production"
				}
			},
			output: {
				comments: false
			}
		})))
		.pipe(gulp.dest("./build/assets/"))
});

gulp.task("html", gulp.series(gulp.parallel("css", "js"), () => {
	return gulp.src(["templates/*.html"])
		.pipe(inline())
		.pipe(when(isProd, chtml()))
		.pipe(gulp.dest("./build/"))
}));

gulp.task("default", gulp.parallel("html"));

gulp.task("serve", (cb) => {
	bs.init({
		server: {
			baseDir: "./build/"
		},
		startPath: "/index.html",
		middleware: [
			{
				route: "/data",
				handle: function (req, res) {
					res.setHeader('Content-Type', 'application/json;charset=UTF-8');
					res.end(JSON.stringify({
						format: "xlif",
						languages: ["en", "fr", "de"],
						labels: [
							{
								id: "hello.world",
								trans: [
									{
										content: "Hello World",
										lng: "en"
									},
									{
										content: "Bonjour le monde",
										lng: "fr"
									},
									{
										content: "Hallo Welt",
										lng: "de"
									}
								]
							}
						]
					}));
				}
			},
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
				route: "/data",
				handle: function (req, res) {
					res.setHeader('Content-Type', 'application/json;charset=UTF-8');
					res.end('{"success":true,"message":""}');
				}
			}
		]
	}, cb);
});

gulp.task("watch", gulp.series("serve", () => {
	gulp.watch("scss/*.scss", gulp.series("css", "html"));
	gulp.watch("js/*.js", gulp.series("js", "html"));
	gulp.watch(["templates/*.html"], gulp.series("html"));
}));
