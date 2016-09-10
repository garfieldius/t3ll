# TYPO3 Local Lang

This is a small utility program to edit localization files in Xlif or the legacy format for [TYPO3 CMS](https://www.typo3.org/).

## Installation

There are no installer or package manager files, but installation is still easy:

#### From Archives

Go to the [releases page](https://github.com/garfieldius/t3ll/releases) and download the right archive for your system. Every archive contains exactly one file: the executable needed. It already contains all dependencies required. Put this executable binary in an accessible location, eg.: a directory inside `$PATH` (or `%PATH%` on windows).

#### From source

t3ll is written in go and uses node.js (npm) modules for building its frontend, so you need both tools installed and properly configured before proceeding.

Then simply clone the repository ...

```bash
# Manually

mkdir -p ${GOPATH}/src/github.com/garfieldius
cd ${GOPATH}/src/github.com/garfieldius
git clone https://github.com/garfieldius/t3ll.git
cd t3ll


# or using go get, the -d flag is important because
# building will fail without the frontend assets
# which are not included in the repository but must
# be built before compiling t3ll

go get -d github.com/garfieldius/t3ll
cd ${GOPATH}/src/github.com/garfieldius/t3ll
```

... and use `make` to build it:

```bash
# This will create a binary in the current directory
make

# Install the binary into ${GOPATH}/bin
make install

# This will create a binary with readable 
# frontend sources, in case you need them
make debug
```

## Usage

t3ll must be called from the command line. It takes exactly one argument: the XML or Xlif file to edit. 

```bash
# Xlif file
t3ll fr.locallang.xlf

# Legacy XML
t3ll locallang.xml
```

In the former case, the file can have a language prefix, or not. t3ll will automatically load all available translations within the same folder, but only that have the same *base name*. eg.: loading the file `fr.locallang.xlf` will also load `locallang.xlf` and `it.locallang.xlf`, but not `fr.locallang_be.xlf`.

If a file does not exist, it will be created.

Once the file is read, the editing mask will open in your default browser. It's interface should be self explanatory as it is very simple and reduced to the absolute minimum.

There are several shortcuts in the browser window:

1. <kbd>Tab</kbd> will focus the first input, jumping to the next if one already is. If the last input, or textarea is active, the first will be focused again.
2. <kbd>Ctrl</kbd>  / <kbd>⌘</kbd> + <kbd>s</kbd>  will save the file
3. <kbd>Ctrl</kbd> / <kbd>⌘</kbd> + <kbd>q</kbd> will save the file and close the window.
4. <kbd>Ctrl</kbd> / <kbd>⌘</kbd> + <kbd>←</kbd> / <kbd>↑</kbd> / <kbd>↓</kbd> / <kbd>→</kbd> will move the focus accordingly if an input is selected.
5. <kbd>Shift</kbd> + <kbd>←</kbd> / <kbd>↑</kbd> / <kbd>↓</kbd> / <kbd>→</kbd> will move the row of the entry or the language if an input is focused.
5. <kbd>Shift</kbd> + <kbd>+</kbd> will add a new row below the current if an input is selected.

This are the same actions as the buttons on screen provide.

## Credits & Notices

t3ll uses the following go packages:

* github.com/hydrogen18/stoppableListener
* github.com/kr/pretty
* github.com/jteeuwen/go-bindata

It also uses CSS Styles and HTML from, as well as (naming) references of, the [TYPO3 CMS](https://www.typo3.org) project.

## License

(c) 2016 Georg Großberger

Released under the Apache License 2.0

See the file [LICENSE](LICENSE) for further information
