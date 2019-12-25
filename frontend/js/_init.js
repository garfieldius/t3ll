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
		isQuit = char == "w" || char == "q";

	if (key == 9 || isMeta) {
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

			case isQuit && (isMeta || isCtrl || isAlt):
				callbacks.close();
				event.preventDefault();
				break;

			case (key == 8 || key == 46) && hasInput && (isMeta || isCtrl || isAlt):
				callbacks.remove(el);
				event.preventDefault();
				break;

			case (key == 107 || key == 187) && hasInput && isCtrl:
				callbacks.add(el);
				event.preventDefault();
				break;

			case key == 37 && hasInput && isCtrl:
				if (cell && tdNum > 0 && isMeta) {
					activeElement = findOne(
						"input,textarea",
						table.rows[trNum].cells[tdNum - 1]
					);
					activeElement.focus();
					event.preventDefault();
				}
				break;

			case key == 39 && hasInput && isCtrl:
				if (isMeta && cell && tdNum < row.cells.length - 1) {
					activeElement = findOne(
						"input,textarea",
						table.rows[trNum].cells[tdNum + 1]
					);
					activeElement.focus();
					event.preventDefault();
				}
				break;

			case key == 38 && hasInput && isCtrl:
				if (isMeta && row && row.rowIndex > 1) {
					activeElement = findOne(
						"input,textarea",
						table.rows[trNum - 1].cells[tdNum]
					);
					activeElement.focus();
					event.preventDefault();
				}
				break;

			case key == 40 && hasInput && isCtrl:
				if (isMeta && row && row.rowIndex < table.rows.length - 1) {
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
