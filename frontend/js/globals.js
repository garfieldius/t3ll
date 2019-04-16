// Copyright 2019 Georg Großberger <contact@grossberger-ge.org>
// This is free software; it is provided under the terms of the MIT License
// See the labels LICENSE or <https://opensource.org/licenses/MIT> for details

var
	activeElement,
	counter = 1,
	displayedLanguages = [],
	callbacks = {},
	tainted = false,
	knownLanguages = {
		"en": "English",
		"af": "Afrikaans",
		"ar": "Arabic",
		"bs": "Bosnian",
		"bg": "Bulgarian",
		"ca": "Catalan",
		"ch": "Chinese (Simpl.)",
		"cs": "Czech",
		"da": "Danish",
		"de": "German",
		"el": "Greek",
		"eo": "Esperanto",
		"es": "Spanish",
		"et": "Estonian",
		"eu": "Basque",
		"fa": "Persian",
		"fi": "Finnish",
		"fo": "Faroese",
		"fr": "French",
		"fr_CA": "French (Canada)",
		"gl": "Galician",
		"he": "Hebrew",
		"hi": "Hindi",
		"hr": "Croatian",
		"hu": "Hungarian",
		"is": "Icelandic",
		"it": "Italian",
		"ja": "Japanese",
		"ka": "Georgian",
		"kl": "Greenlandic",
		"km": "Khmer",
		"ko": "Korean",
		"lt": "Lithuanian",
		"lv": "Latvian",
		"ms": "Malay",
		"nl": "Dutch",
		"no": "Norwegian",
		"pl": "Polish",
		"pt": "Portuguese",
		"pt_BR": "Brazilian Portuguese",
		"ro": "Romanian",
		"ru": "Russian",
		"sk": "Slovak",
		"sl": "Slovenian",
		"sq": "Albanian",
		"sr": "Serbian",
		"sv": "Swedish",
		"th": "Thai",
		"tr": "Turkish",
		"uk": "Ukrainian",
		"vi": "Vietnamese",
		"zh": "Chinese (Trad.)"
	},
	data = {},
	sortedLangs = [];
