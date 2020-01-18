// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

callbacks = {
	filterNotTranslated: function () {
		setTimeout(function () {
			filterNonTranslated = !!findOne('[data-toggle=filterNotTranslated]').checked;
			renderState();
		}, 20);
	},
	setKey: function (input) {
		tainted = true;
		var row = findParent(input, ["TR"]);
		findAll(".is-label", row).forEach(function (label) {
			label.dataset["key"] = input.value;
		});
	},
	setContent: function () {
		serializeState();
		tainted = true;
	},
	moveUp: function (el) {
		tainted = true;
		var row = findParent(el, ["TR"]);
		row.parentNode.insertBefore(row, row.previousSibling);
		setButtonVisiblity();
		serializeState();
	},
	moveDown: function (el) {
		tainted = true;
		var row = findParent(el, ["TR"]);
		row.parentNode.insertBefore(row, row.nextSibling.nextSibling);
		setButtonVisiblity();
		serializeState();
	},
	sortLabels: function () {
		data.labels.sort(function (a, b) {
			if (a.id === b.id) {
				return 0;
			} else {
				return a.id > b.id ? 1 : -1
			}
		});
		renderState();
		setButtonVisiblity();
		serializeState();
		tainted = true;
	},
	add: function (btn) {
		var row = findParent(btn, "TR"),
			rowData = {
				id: "new." + (counter++),
				trans: displayedLanguages.map(function (lang) {
					return {
						lng: lang,
						content: ""
					};
				})
			};
		var newRow = tree(makeRow(rowData));

		if (row.nextSibling) {
			row.parentNode.insertBefore(newRow, row.nextSibling);
		} else {
			row.parentNode.appendChild(newRow);
		}
		tainted = true;
		serializeState();
		setButtonVisiblity();
	},
	remove: function (btn) {
		var row = findParent(btn, "TR");
		delete data.labels[findOne(".is-key", row).value];
		row.parentNode.removeChild(row);
		tainted = true;
		setButtonVisiblity();
		serializeState();
	},
	save: function () {
		doSave();
	},
	convert: function () {
		findOne("#ToXliff").value = "1";
		doSave();
		findOne("#ToXliffMessage").style.display = "none";
	},
	isDirty: function () {
		return tainted;
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
	addLanguage: function (select) {
		var langCode = select.value;

		if (data.languages.indexOf(langCode) > -1) {
			return;
		}

		data.languages.push(langCode);
		displayedLanguages.push(langCode);
		serializeState();

		renderState();
	},
	uploadReplace: function (el) {
		doUpload(el.files[0], true)
	},
	uploadMerge: function (el) {
		doUpload(el.files[0], false)
	},
	initMerge: function () {
		findOne("#FileSelectMerge").click();
	},
	initReplace: function () {
		findOne("#FileSelectReplace").click();
	},
	csvDropdown: function () {
		findOne("#CsvDropdown").classList.toggle('active');
	},
	csvHide: function () {
		findOne("#CsvDropdown").classList.remove("active");
	}
};

function doUpload(file, replace) {
	callbacks.csvHide();
	showMessage("Uploading");
	var reader = new FileReader();

	reader.onload = function (evt) {
		var url = "csv";

		if (replace === true) {
			url += "?mode=replace"
		}

		xhr(url, evt.target.result, function () {
			location.reload();
		});
	};

	reader.onerror = function (err) {
		showMessage(err, true)
	};

	reader.readAsBinaryString(file);
}

function doSave() {
	if (!tainted) {
		return showMessage("Data not changed");
	}

	serializeState();

	xhr("save", getFormData(), function (err, data) {
		showMessage(err || data.error || data.message, err || data.error);
		tainted = false;
	});
}

function runCallbacks(event, type) {
	var el = event.target;

	if (el && type == "click") {
		while (el.parentNode) {
			if (el.dataset && el.dataset["toggle"]) {
				break;
			}

			el = el.parentNode;
		}
	}

	if (el && (el.classList + "").indexOf("disabled") === -1 && el.dataset && el.dataset["toggle"]) {
		if (!el.dataset["event"] || el.dataset["event"] == type) {
			var cb = el.dataset["toggle"];

			if (callbacks[cb]) {
				callbacks[cb](el, event);
			}
		}
	}
}

function getFormData() {
	var d = new FormData();
	d.append("data", JSON.stringify(data));

	if (findOne("#ToXliff").value === "1") {
		d.append("format", "xlif");
	}

	return d;
}

function serializeState() {
	data.labels = [];

	findAll("#dataTable tbody tr").forEach(function (row) {
		var key, labels = [];

		findAll("textarea", row).forEach(function (el) {
			key = el.dataset["key"];
			labels.push({
				lng: el.dataset["lang"],
				content: el.value
			});
		});

		data.labels.push({
			id: key,
			trans: labels
		});
	});
}
