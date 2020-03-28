// Copyright 2020 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

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

	if (start.hidden) {
		parent.style.display = "none";
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

function makeRow(row) {
	var cells = [{
		name: "td",
		sub: [
			{
				name: "input",
				cls: "form-control is-key",
				attr: {
					type: "text",
					value: row.id
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

	data.languages.forEach(function (lang) {
		var cell = {
			name: "td",
			hidden: displayedLanguages.indexOf(lang) === -1,
			sub: [
				{
					name: "textarea",
					cls: "form-control is-label",
					data: {
						lang: lang,
						key: row.id,
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
		};

		row.trans.forEach(function (lbl) {
			if (lbl.lng == lang) {
				cell.sub[0].txt = lbl.content;
			}
		});

		cells.push(cell);
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

function contains(el, container) {
	while (el) {
		if (el.id && el.id == container) {
			return true;
		}
		el = el.parentNode;
	}

	return false;
}
