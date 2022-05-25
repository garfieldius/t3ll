# TYPO3 Local Lang

This is a small utility program to edit localization files in Xlif or the legacy
format for [TYPO3 CMS](https://www.typo3.org/).

## Installation

t3ll opens its editor inside a Google Chrome or Chromium window. One of this must
be available on your system.

Installation can be done with homebrew (MacOS and Linux) or simple binary downloads:

#### Homebrew

Tap into `garfieldius/taps` and install the package `t3ll`:

```bash
brew tap garfieldius/taps
brew install t3ll
```

#### Binary downloads

Go to the [releases page](https://github.com/garfieldius/t3ll/releases) and download the
right file for your system. Rename it to `t3ll` / `t3ll.exe` and put it into `$PATH`
/ `%PATH%`.

All binaries can be checked via gpg. The .sig files contain signatures created with the
key `0D1F16703AB055AA`. It is available on most common keyservers or on
<https://grossberger-ge.org/gpg.asc>

Example:

```bash
# Import key
gpg --keyserver hkp://pool.sks-keyservers.net --recv-keys 0D1F16703AB055AA

VERSION=$(curl -sSL 'https://api.github.com/repos/garfieldius/t3ll/releases/latest' | jq -r '.tag_name')
ARCH=linux_x64

# Download binary and signature
curl -sSLo t3ll https://github.com/garfieldius/t3ll/releases/download/${VERSION}/t3ll_${ARCH}
curl -sSLo t3ll.sig https://github.com/garfieldius/t3ll/releases/download/${VERSION}/t3ll_${ARCH}.sig

# Verify
gpg --verify t3ll.sig t3ll

# Install
install -m 0755 t3ll /usr/local/bin/
```

#### Building from source

t3ll is written in go and uses node.js and yarn modules for building its frontend,
so these tools need to be installed and properly configured before proceeding.

Then simply clone the repository and use `make` to build it:

```bash
git clone https://github.com/garfieldius/t3ll.git
cd t3ll

# This will create a production program in the current directory
make

# Build a debug program. Has the same functions but VERY verbose logging to stdout
# and readable frontend sources
make debug

# Install the program into /usr/local/bin
make install
```

## Usage

t3ll is called from the command line. It takes exactly one argument: the XML or
Xliff file to edit.

```bash
t3ll fr.locallang.xlf

# Legacy XML
t3ll locallang.xml
```

In the former case, the file can have a language prefix, or not. t3ll will
automatically load all available translations within the same folder, but only
those having the same **base name**. That is the name of the file without any
language prefix. eg.: loading the file `fr.locallang.xlf` will also load `
locallang.xlf` and `it.locallang.xlf`, but not `fr.locallang_backend.xlf`.

If a file does not exist, it will be created.

Only file names with a suffix / extension of .xml, .xlf, .xlif or .xliff are
supported. Others will be treated as an unknown file type.

Once the file is read, the editing mask will open in a chromium or google chrome
window. It's interface should be self explanatory as it is very simple and
reduced to the absolute minimum.

There are several shortcuts in the browser window:

<kbd>Meta</kbd> is <kbd>Ctrl</kbd>+<kbd>Shift</kbd> on MacOS, <kbd>Alt</kbd>+<kbd>Shift</kbd> on
other systems.

* <kbd>Tab</kbd> will focus the first input, jumping to the next if one is
  already focused. If the last input is active, the first will be focused again.
* <kbd>Shift</kbd>+<kbd>Tab</kbd> will focus the last input, jumping to the
  previous if one is already focused. If the first input is active, the last
  will be focused again.
* <kbd>Cmd</kbd> / <kbd>Ctrl</kbd> / <kbd>Alt</kbd> + <kbd>s</kbd>  will save
  the file
* <kbd>Meta</kbd> + <kbd>←</kbd> / <kbd>↑</kbd> / <kbd>↓</kbd> / <kbd>→</kbd>
  will move the focus accordingly if a field is focused.
* <kbd>Meta</kbd> + <kbd>Backspace</kbd> / <kbd>Del</kbd> will delete the focused
  row if a field is focused. The row below will be the new focused.
* <kbd>Meta</kbd> + <kbd>+</kbd> will add a new row below the current if a field
  is focused

When converting from XML to XLIF, the old .xml file will not be deleted, this must
be done manually.

## License

[MIT](https://opensource.org/licenses/MIT)
