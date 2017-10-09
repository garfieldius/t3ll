
/*
 * Copyright 2016 Georg Großberger
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
        x.open("POST", u);
        xhr.overrideMimeType('text/plain; charset=x-user-defined-binary');
        x.send(vals);
    } else {
        x.open("GET", u);
        x.send(null);
    }
}
