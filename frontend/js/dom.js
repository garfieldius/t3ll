
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

function iconHtml(name) {
    return "<svg viewBox=\"0 0 16 16\"><use xlink:href=\"#" + name + "\"></use></svg>"
}

function findAll(selector, parent) {
    var list = (parent || document).querySelectorAll(selector),
        arr = [];

    for (var i = list.length; i--; arr.unshift(list[i])) {
    }

    return arr;
}

function saveAndRestoreFocus(td, cb) {
    var currentCell = td.cellIndex,
        currentRow = td.parentNode.rowIndex;

    cb(td);

    activeElement = findOne(
        "input,textarea",
        findOne("#dataTable").rows[currentRow].cells[currentCell]
    );

    activeElement.focus();
}

function tree(start) {
    var parent = document.createElement(start.name);

    if (start.cls) {
        parent.className = start.cls;
    }

    if (start.attr) {
        Object.keys(start.attr).forEach(function (name) {
            if (name && start.attr[name]) {
                parent.setAttribute(name, start.attr[name]);
            }
        });
    }

    if (start.data) {
        Object.keys(start.data).forEach(function (name) {
            if (name && start.data[name]) {
                parent.dataset[name] = start.data[name];
            }
        });
    }

    if (start.events) {
        Object.keys(start.events).forEach(function (eventName) {
            var eventFunc = start.events[eventName];
            if (eventName && eventFunc) {
                parent.addEventListener(eventName, eventFunc);
            }
        });
    }

    if (start.txt) {
        parent.innerHTML = start.txt;
    }

    if (start.sub) {
        start.sub.forEach(function (sub) {
            parent.appendChild(tree(sub));
        });
    }

    return parent;
}

function makeRow(key) {
    var cells = [{
        name: "td",
        sub: [
            {
                name: "input",
                cls: "form-control is-key",
                attr: {
                    type: "text",
                    value: key
                },
                data: {
                    toggle: "setKey",
                    event: "input"
                },
                events: {
                    focus: function () {
                        activeElement = this;
                    },
                    blur: function () {
                        activeElement = null;
                    }
                }
            }
        ]
    }];

    displayedLanguages.forEach(function (lang) {
        cells.push({
            name: "td",
            sub: [
                {
                    name: "textarea",
                    cls: "form-control is-label",
                    txt: data.labels[key] && data.labels[key][lang] ? data.labels[key][lang] : "",
                    data: {
                        lang: lang,
                        key: key,
                        toggle: "setContent",
                        event: "input"
                    },
                    events: {
                        focus: function () {
                            activeElement = this;
                        },
                        blur: function () {
                            activeElement = null;
                        }
                    }
                }
            ]
        });
    });

    cells.push({
        name: "td",
        attr: {
            "width": "1%"
        },
        sub: [
            {
                name: "div",
                cls: "btn-group nowrap",
                sub: [
                    {
                        name: "button",
                        cls: "btn btn-sm btn-default",
                        txt: iconHtml("add"),
                        data: {
                            toggle: "add",
                            event: "click"
                        }
                    },
                    {
                        name: "button",
                        cls: "btn btn-sm btn-default",
                        txt: iconHtml("move-up"),
                        data: {
                            toggle: "moveUp",
                            event: "click"
                        }
                    },
                    {
                        name: "button",
                        cls: "btn btn-sm btn-default",
                        txt: iconHtml("move-down"),
                        data: {
                            toggle: "moveDown",
                            event: "click"
                        }
                    },
                    {
                        name: "button",
                        cls: "btn btn-sm btn-default",
                        txt: iconHtml("delete"),
                        data: {
                            toggle: "remove",
                            event: "click"
                        }
                    }
                ]
            }
        ]
    });

    return {
        name: "tr",
        sub: cells
    }
}

function findOne(selector, parent) {
    return (parent || document).querySelector(selector);
}

function empty(el) {
    while (el.childNodes.length) {
        el.removeChild(el.firstChild);
    }
    return el;
}

function findParent(el, tags) {
    while (el && tags.indexOf(el.tagName) === -1) {
        el = el.parentNode;
    }
    return el;
}
