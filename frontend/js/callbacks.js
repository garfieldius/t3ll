
/*
 * Copyright 2016 Georg Großberger
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

var callbacks = {
    setKey: function (input) {
        tainted = true;
        var row = findParent(input, "TR");
        findAll(".is-label", row).forEach(function (label) {
            label.dataset["key"] = input.value;
        });
    },
    setContent: function (input) {
        tainted = true;
        data.labels[input.dataset["key"]][input.dataset["lang"]] = input.value;
    },
    moveUp: function (el) {
        tainted = true;
        var row = findParent(el, ["TR"]);
        row.parentNode.insertBefore(row, row.previousSibling);
        setButtonVisiblity();
    },
    moveDown: function (el) {
        tainted = true;
        var row = findParent(el, ["TR"]);
        row.parentNode.insertBefore(row, row.nextSibling.nextSibling);
        setButtonVisiblity();
    },
    add: function (btn) {
        var row = findParent(btn, "TR");
        var newRow = tree(makeRow("new." + (counter++)));

        if (row.nextSibling) {
            row.parentNode.insertBefore(newRow, row.nextSibling);
        } else {
            row.parentNode.appendChild(newRow);
        }
        setButtonVisiblity();
        serializeState();
        tainted = true;
    },
    remove: function (btn) {
        var row = findParent(btn, "TR");
        delete data.labels[findOne(".is-key", row).value];
        row.parentNode.removeChild(row);
        tainted = true;
        setButtonVisiblity();
    },
    save: function () {
        if (!tainted) {
            return showMessage("Data not changed");
        }

        serializeState();

        xhr("save", getFormData(), function (err, data) {
            showMessage(err || data.error || data.message, err || data.error);
            tainted = false;
        });
    },
    close: function () {
        if (!tainted) {
            return quit();
        }

        serializeState();

        xhr("save", getFormData(), function () {
            quit();
        });
    },
    setDisplayLanguages: function () {
        displayedLanguages = [];

        findAll("input", findOne("#visibleLanguages")).forEach(function (input) {
            if (input.checked) {
                displayedLanguages.push(input.value);
            }
        });

        setTimeout(renderState, 10);
    },
    moveLeft: function (el) {
        var lang = el.dataset["language"],
            pos = displayedLanguages.indexOf(lang);

        displayedLanguages[pos] = displayedLanguages[pos - 1];
        displayedLanguages[pos - 1] = lang;
        renderState();
    },
    moveRight: function (el) {
        var lang = el.dataset["language"],
            pos = displayedLanguages.indexOf(lang);

        displayedLanguages[pos] = displayedLanguages[pos + 1];
        displayedLanguages[pos + 1] = lang;
        renderState();
    },
    addLanguage: function (select) {
        var langCode = select.value;

        if (data.languages.indexOf(langCode) > -1) {
            return;
        }

        data.languages.push(langCode);
        displayedLanguages.push(langCode);

        renderState();
    }
};

function quit() {
    xhr("quit", function () {
        window.close();
    });
}

function runCallbacks(event, type) {
    var el = event.target;

    if (el && type == "click") {
        el = findParent(el, ["A", "BUTTON"]);
    }

    if (el && !el.classList.contains("disabled") && el.dataset && el.dataset["toggle"]) {
        if (!el.dataset["event"] || el.dataset["event"] == type) {
            var cb = el.dataset["toggle"];
            callbacks[cb](event.target, event);
        }
    }
}

function getFormData() {
    var d = {
        data: JSON.stringify(data)
    };

    if (findOne("#ToXliff").checked) {
        d.format = "xlif";
    }

    return d;
}

function serializeState() {
    data.labels = {};

    findAll("#dataTable textarea").forEach(function (el) {
        var lang = el.dataset["lang"],
            key  = el.dataset["key"],
            content = el.value;

        if (!data.labels[key]) {
            data.labels[key] = {};
        }

        data.labels[key][lang] = content;
    });
}
