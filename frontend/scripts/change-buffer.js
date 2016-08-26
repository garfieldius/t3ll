
"use strict";

module.exports = (converter) => {
    return () => {
        return require('through2').obj(function(file, enc, cb) {

            if (file.isNull()) {
                this.push(file);
                return cb();
            }

            if (file.isStream()) {
                let gutil = require('gulp-util');
                this.emit('error', new gutil.PluginError('change-buffer', 'Streaming not supported'));
                return cb();
            }

            let data = file.contents;
            let isBuffer = Buffer.isBuffer(data);

            if (isBuffer) {
                data = data.toString('utf-8');
            }

            let result = converter(data, file);

            file.contents = isBuffer ? new Buffer(result) : result;

            this.push(file);
            cb();
        });
    }
};
