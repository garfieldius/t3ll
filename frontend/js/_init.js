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

["input", "change", "click"].forEach(function (eventName) {
    document.addEventListener(eventName, function (event) {
        runCallbacks(event, eventName);
    });
});

window.addEventListener("keydown", function (event) {
    var
        key = event.which,
        char = String.fromCharCode(key).toLowerCase(),
        isShift = event.shiftKey,
        isMeta = event.ctrlKey || event.metaKey,
        el = activeElement,
        hasInput = false,
        cell, row, table,
        tdNum = 0,
        trNum = 0;

    if (key == 9 || isMeta || isShift) {
        hasInput = el && ["INPUT", "TEXTAREA"].indexOf(el.tagName) > -1;

        if (hasInput) {
            cell = findParent(el, ["TD"]);
            row = findParent(cell, ["TR"]);
            table = findParent(row, ["TABLE"]);

            tdNum = cell.cellIndex;
            trNum = row.rowIndex;
        }
    }

    if (key == 9) {
        if (hasInput) {
            tdNum++;
            if (tdNum >= row.cells.length - 1) {
                tdNum = 0;
                trNum++;
            }

            if (trNum >= table.rows.length) {
                trNum = 1;
            }
        } else {
            tdNum = 0;
            trNum = 1;
        }

        activeElement = findOne(
            "input,textarea",
            findOne("#dataTable").rows[trNum].cells[tdNum]
        );
        activeElement.focus();
        event.preventDefault();
    } else if (isMeta || isShift) {
        hasInput = el && ["INPUT", "TEXTAREA"].indexOf(el.tagName) > -1;

        if (hasInput) {
            cell = findParent(el, ["TD"]);
            row = findParent(cell, ["TR"]);
            table = findParent(row, ["TABLE"]);

            tdNum = cell.cellIndex;
            trNum = row.rowIndex;
        }

        switch (true) {
            case char == 's' && !isShift:
                callbacks.save();
                event.preventDefault();
                break;

            case char == 'q' && !isShift:
                callbacks.close();
                event.preventDefault();
                break;

            case (key == 107 || key == 187) && hasInput && !isShift:
                callbacks.add(el);
                event.preventDefault();
                break;

            case key == 37 && hasInput:
                // Move lang left
                if (cell && tdNum > 0) {
                    if (isShift && tdNum > 1) {
                        callbacks.moveLeft(el)
                    }

                    // Move focus left
                    activeElement = findOne(
                        "input,textarea",
                        table.rows[trNum].cells[tdNum - 1]
                    );
                    activeElement.focus();
                    event.preventDefault();
                }

                break;

            case key == 39 && hasInput:
                if (cell && tdNum < row.cells.length - 1) {
                    // Move lang right
                    if (isShift && tdNum > 0) {
                        callbacks.moveRight(el);
                    }

                    // Move focus right
                    activeElement = findOne(
                        "input,textarea",
                        table.rows[trNum].cells[tdNum + 1]
                    );
                    activeElement.focus();
                    event.preventDefault();
                }
                break;

            case key == 38 && hasInput:
                if (row && row.rowIndex > 1) {

                    // Move entry up
                    if (isShift) {
                        callbacks.moveUp(activeElement);
                    }

                    // Move focus up
                    activeElement = findOne(
                        "input,textarea",
                        table.rows[trNum - 1].cells[tdNum]
                    );
                    activeElement.focus();
                    event.preventDefault();
                }
                break;

            case key == 40 && hasInput:
                if (row && row.rowIndex < table.rows.length - 1) {
                    // Move entry down
                    if (isShift) {
                        callbacks.moveDown(activeElement);
                    }

                    // Move cursor down
                    activeElement = findOne(
                        "input,textarea",
                        table.rows[trNum + 1].cells[tdNum]
                    );
                    activeElement.focus();
                    event.preventDefault();
                }
        }
    }
});

var sortedLangs = data.languages.filter(function (lang) {
    return (lang != "en");
});
sortedLangs.sort(function (a, b) {
    return knownLanguages[a] < knownLanguages[b];
});
sortedLangs.unshift("en");

data.languages = sortedLangs;
[].push.apply(displayedLanguages, data.languages);

if (data.format == "xml") {
    findOne("#ToXliffMessage").style.display = "";
}

findOne("#messages").className = "flash-message";
findOne("#messages").innerHTML = "";

renderState(data);
