// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

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

	data.languages.forEach(function (code) {
		headers.push({
			name: "th",
			hidden: displayedLanguages.indexOf(code) === -1,
			sub: [
				{
					name: "div",
					cls: "move-container",
					sub: [
						{
							name: "span",
							txt: knownLanguages[code]
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

	if (filterNonTranslated) {
		rows = rows.map(function (row) {
			row.hidden = row.sub.slice(1, -1).reduce(function (prev, cell) {
				if (cell.hidden) {
					return prev;
				} else {
					return ('' + cell.sub[0].txt).trim() + prev;
				}
			}, '') === '';
			return row;
		})
	}

	table.appendChild(tree({
		name: "tbody",
		sub: rows
	}));

	setButtonVisiblity();
}

function showMessage(msg, isError, persistent) {
	var c = empty(findOne("#messages")),
		m = tree({
			name: "span",
			txt: msg,
			cls: (isError ? "error" : "success")
		});

	c.appendChild(m);
	m.classList.add("show");

	if (!persistent) {
		setTimeout(function () {
			m.classList.remove("show");

			setTimeout(function () {
				c.removeChild(m);
			}, 210)
		}, 1000);
	}
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
