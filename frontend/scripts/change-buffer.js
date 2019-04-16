// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

"use strict";
const PluginError = require('plugin-error');
const { obj } = require("through2")

module.exports = (converter) => {
	return (flags) => {
		let args = [flags];

		return obj(function (file, enc, cb) {

			if (file.isNull()) {
				this.push(file);
				return cb();
			}

			if (file.isStream()) {
				this.emit('error', new PluginError('change-buffer', 'Streaming not supported'));
				return cb();
			}

			let data = file.contents;
			let isBuffer = Buffer.isBuffer(data);

			if (isBuffer) {
				data = data.toString('utf-8');
			}
			let arg = [data, file];
			arg.push.apply(arg, args);

			let result = converter.apply(null, arg);

			file.contents = isBuffer ? Buffer.from(result) : result;

			this.push(file);
			cb();
		});
	}
};
