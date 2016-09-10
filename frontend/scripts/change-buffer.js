
"use strict";

module.exports = (converter) => {
    return (flags) => {
        var args = [flags];

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
            let arg = [data, file];
            arg.push.apply(arg, args);

            let result = converter.apply(null, arg);

            file.contents = isBuffer ? new Buffer(result) : result;

            this.push(file);
            cb();
        });
    }
};
