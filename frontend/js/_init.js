// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

["input", "change", "click"].forEach(function (eventName) {
	document.addEventListener(eventName, function (event) {
		runCallbacks(event, eventName);
	});
});

window.addEventListener("keydown", function (event) {
	var
		key = event.which,
		char = String.fromCharCode(key).toLowerCase(),
		isCtrl = event.ctrlKey,
		isMeta = event.metaKey,
		isAlt = event.altKey,
		el = activeElement,
		hasInput = false,
		cell, row, table,
		tdNum = 0,
		trNum = 0,
		newActive,
		isQuit = char == "w" || char == "q",
		tries = 10000;

	if (key == 9 || isAlt) {
		hasInput = el && ["INPUT", "TEXTAREA"].indexOf(el.tagName) > -1;

		if (hasInput) {
			cell = findParent(el, ["TD"]);
			row = findParent(cell, ["TR"]);
			table = findParent(row, ["TABLE"]);

			tdNum = cell.cellIndex;
			trNum = row.rowIndex;
		}
	}

	function moveToNextCell() {
		tdNum++;
		if (tdNum >= row.cells.length - 1) {
			tdNum = 0;
			trNum++;
		}

		if (trNum >= table.rows.length) {
			trNum = 1;
		}
	}

	function isVisible(el) {
		var td = findParent(el, ["TD"]), tr;

		if (td && td.style.display !== "none") {
			tr = findParent(td, ["TR"]);

			if (tr && tr.style.display !== "none") {
				return true;
			}
		}

		return false;
	}

	if (key == 9) {
		if (hasInput) {
			moveToNextCell();
		} else {
			tdNum = 0;
			trNum = 1;
		}

		while (!newActive) {
			newActive = findOne(
				"input,textarea",
				findOne("#dataTable").rows[trNum].cells[tdNum]
			);

			if (!isVisible(newActive)) {
				newActive = null;
				moveToNextCell();
			}

			if (!tries--) {
				return;
			}
		}

		activeElement = findOne(
			"input,textarea",
			findOne("#dataTable").rows[trNum].cells[tdNum]
		);

		activeElement = newActive;
		activeElement.focus();
		event.preventDefault();
	} else if (isMeta || isCtrl || isAlt) {
		hasInput = el && ["INPUT", "TEXTAREA"].indexOf(el.tagName) > -1;

		if (hasInput) {
			cell = findParent(el, ["TD"]);
			row = findParent(cell, ["TR"]);
			table = findParent(row, ["TABLE"]);

			tdNum = cell.cellIndex;
			trNum = row.rowIndex;
		}

		switch (true) {
			case char == 's' && (isMeta || isCtrl || isAlt):
				callbacks.save();
				event.preventDefault();
				break;

			// Remove row
			case (key == 8 || key == 46) && hasInput && isAlt:
				callbacks.remove(el);
				event.preventDefault();
				break;

			// Add row
			case (key == 107 || key == 187) && hasInput && isAlt:
				callbacks.add(el);
				event.preventDefault();
				break;

			// Move left
			case key == 37 && hasInput && isAlt && cell && tdNum > 0:
				while (true) {
					tdNum--;

					if (tdNum < 0) {
						return;
					}

					activeElement = findOne(
						"input,textarea",
						row.cells[tdNum]
					);

					if (isVisible(activeElement)) {
						break;
					}
				}

				activeElement.focus();
				event.preventDefault();
				break;

			// Move right
			case key == 39 && hasInput && isAlt && cell && tdNum < row.cells.length - 1:
				while (true) {
					tdNum++;

					// -1 because there is an additional column on the right!
					if (tdNum >= row.cells.length - 1) {
						return;
					}

					activeElement = findOne(
						"input,textarea",
						row.cells[tdNum]
					);

					if (isVisible(activeElement)) {
						break;
					}
				}

				activeElement.focus();
				event.preventDefault();
				break;

			// Move up
			case key == 38 && hasInput && isAlt && row && row.rowIndex > 1:
				while (true) {
					trNum--;

					if (trNum < 0) {
						return;
					}

					activeElement = findOne(
						"input,textarea",
						table.rows[trNum].cells[tdNum]
					);

					if (isVisible(activeElement)) {
						break;
					}
				}

				activeElement.focus();
				event.preventDefault();
				break;

			// Move down
			case key == 40 && hasInput && isAlt && row && row.rowIndex < table.rows.length - 1:
				while (true) {
					trNum++;

					if (trNum >= table.rows.length) {
						return;
					}

					activeElement = findOne(
						"input,textarea",
						table.rows[trNum].cells[tdNum]
					);

					if (isVisible(activeElement)) {
						break;
					}
				}

				activeElement.focus();
				event.preventDefault();
		}
	}
});

xhr("data", function (_, resp) {
	data = resp;

	window.addEventListener("unload", function () {
		xhr("quit");
	});

	sortedLangs = data.languages.filter(function (lang) {
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

	function heartbeat() {
		xhr("hb", function (err, resp) {
			if (err || !resp || !resp.success) {
				showMessage("t3ll does not seem to be running. Close and reopen this window", true, true);
			} else {
				setTimeout(heartbeat, 800);
			}
		})
	}

	heartbeat();
})
