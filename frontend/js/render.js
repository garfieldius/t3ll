
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

function renderState() {
    var visibleLangs = empty(findOne("#visibleLanguages")),
        selectLangs = empty(findOne("#availableLanguages")),
        table = empty(findOne("#dataTable")),
        headers = [], rows = [];

    selectLangs.appendChild(tree({
        name: "option",
        txt: "Add language",
    }));

    Object.keys(knownLanguages).forEach(function (key) {
        if (data.languages.indexOf(key) == -1) {
            selectLangs.appendChild(tree({
                name: "option",
                txt: knownLanguages[key],
                attr: {
                    value: key
                }
            }));
        }
    });

    data.languages.forEach(function (langCode) {
        visibleLangs.appendChild(tree({
            name: "div",
            cls: "checkbox",
            sub: [
                {
                    name: "label",
                    sub: [
                        {
                            name: "input",
                            attr: {
                                value: langCode,
                                checked: (displayedLanguages.indexOf(langCode) > -1 ? "checked" : null),
                                type: "checkbox"
                            },
                            data: {
                                toggle: "setDisplayLanguages",
                                event: "change"
                            }
                        },
                        {
                            name: "span",
                            txt: knownLanguages[langCode]
                        }
                    ]
                }
            ]
        }));
    });

    headers.push({
        name: "th",
        txt: "Key"
    });

    displayedLanguages.forEach(function (code) {
        var btns = [];

            btns.push({
                name: "btn",
                cls: "btn btn-xs btn-default",
                txt: iconHtml("move-left"),
                data: {
                    toggle: "moveLeft",
                    language: code
                }
            });

            btns.push({
                name: "btn",
                cls: "btn btn-xs btn-default",
                txt: iconHtml("move-right"),
                data: {
                    toggle: "moveRight",
                    language: code
                }
            });

        headers.push({
            name: "th",
            sub: [
                {
                    name: "div",
                    cls: "move-container",
                    sub: [
                        {
                            name: "span",
                            txt: knownLanguages[code]
                        },
                        {
                            name: "div",
                            cls: "move-btns",
                            sub: [
                                {
                                    name: "span",
                                    cls: "btn-group nowrap",
                                    sub: btns
                                }
                            ]
                        }
                    ]
                }
            ]
        });
    });
    headers.push({
        name: "th",
        cls: "actions",
        txt: "Actions"
    });

    table.appendChild(tree({
        name: "thead",
        sub: [
            {
                name: "tr",
                sub: headers
            }
        ]
    }));

    data.labels.forEach(function (label) {
        rows.push(makeRow(label));
    });

    table.appendChild(tree({
        name: "tbody",
        sub: rows
    }));

    setButtonVisiblity();
}
function showMessage(msg, isError) {
    var c = empty(findOne("#messages")),
        m = tree({
            name: "span",
            txt: msg,
            cls: (isError ? "error" : "success")
        });

    c.appendChild(m);
    m.classList.add("show");

    setTimeout(function () {
        m.classList.remove("show");

        setTimeout(function() {
            c.removeChild(m);
        }, 210)
    }, 1000);
}

function setButtonVisiblity() {
    var rows = findAll("tr", findOne("#dataTable tbody")), max = rows.length - 1;

    rows.forEach(function (row, index) {
        findAll("[data-toggle]", row).forEach(function (btn) {
            btn.classList.remove("disabled");
        });

        if (index == 0) {
            findOne("[data-toggle=moveUp]", row).classList.add("disabled");
        } else if (index == max) {
            findOne("[data-toggle=moveDown]", row).classList.add("disabled");
        }
    });
}
