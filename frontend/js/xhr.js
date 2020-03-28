// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

function xhr(url, vals, cb) {
	if (!cb && vals) {
		cb = vals;
		vals = null
	}

	var x = new XMLHttpRequest();
	var u = location.protocol + "//" + location.host + "/" + url;

	x.addEventListener("readystatechange", function () {
		if (this.readyState == 4) {
			try {
				var res = JSON.parse(this.responseText);
				cb(null, res);
			} catch (err) {
				return cb(err);
			}
		}
	});

	if (vals) {
		x.overrideMimeType('text/plain; charset=x-user-defined-binary');
		x.open("POST", u);
		x.send(vals);
	} else {
		x.open("GET", u);
		x.send(null);
	}
}
