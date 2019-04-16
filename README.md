# TYPO3 Local Lang

This is a small utility program to edit localization files in Xlif or the legacy format for [TYPO3 CMS](https://www.typo3.org/).

## Installation

t3ll opens its editor inside a Google Chrome or Chromium window. One of this must be available on your system.

There are no installer or package manager files, but installation is still easy:

#### Binary downloads

Go to the [releases page](https://github.com/garfieldius/t3ll/releases) and download the right file for your system. Rename it to `t3ll` / `t3ll.exe` and put it into `$PATH` / `%PATH%`.

All binaries can be checked via gpg. The .sig files contain signatures created with the key `0D1F16703AB055AA`. It is available on most common keyservers or on <https://grossberger-ge.org/gpg.asc>

Example:

```bash
# Import key
gpg --keyserver hkp://pool.sks-keyservers.net --recv-keys 0D1F16703AB055AA

VERSION=1.0.0
ARCH=linux_x64

# Download binary and signature
curl -sSLo t3ll https://github.com/lextoumbourou/goodhosts/releases/download/v${VERSION}/t3ll_${ARCH}
curl -sSLo t3ll.sig https://github.com/lextoumbourou/goodhosts/releases/download/v${VERSION}/t3ll_${ARCH}.sig

# Verify
gpg --verify t3ll.sig.sig t3ll.sig

# Install
mv t3ll /usr/local/bin
chmod +x /usr/local/bin/t3ll
```

#### Building from source

t3ll is written in go and uses node.js and yarn modules for building its frontend, so these tools need to be installed and properly configured before proceeding.

Then simply clone the repository:

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
# This will create a (debug) binary in the current directory
make

# Install the (debug) binary into ${GOPATH}/bin
make install

# This will create a release binary for Linux, Windows
# and MacOS in the folder dist/
make dist
```

## Usage

t3ll is called from the command line. It takes exactly one argument: the XML or Xlif file to edit.

```bash
labels
t3ll fr.locallang.xlf

# Legacy XML
t3ll locallang.xml
```

In the former case, the file can have a language prefix, or not. t3ll will automatically load all available translations within the same folder, but only that have the same *base name*. eg.: loading the file `fr.locallang.xlf` will also load `locallang.xlf` and `it.locallang.xlf`, but not `fr.locallang_be.xlf`.

If a file does not exist, it will be created.

Once the file is read, the editing mask will open in your default browser. It's interface should be self explanatory as it is very simple and reduced to the absolute minimum.

There are several shortcuts in the browser window (<kbd>Meta</kbd> means one of <kbd>Alt</kbd> or <kbd>Command</kbd> as in <kbd>⌘</kbd> or <kbd>win</kbd>):

* <kbd>Tab</kbd> will focus the first input, jumping to the next if one already is. If the last input is active, the first will be focused again.
* <kbd>Meta</kbd> + <kbd>s</kbd>  will save the file
* <kbd>Meta</kbd> + <kbd>q</kbd> / <kbd>w</kbd> will save the file and close the window.
* <kbd>Meta</kbd> + <kbd>←</kbd> / <kbd>↑</kbd> / <kbd>↓</kbd> / <kbd>→</kbd> will move the focus accordingly if an input is selected.
* <kbd>Ctrl</kbd> / <kbd>Meta</kbd> + <kbd>Backspace</kbd> / <kbd>Del</kbd> will delete a row if an input or textarea is focused
* <kbd>Ctrl</kbd> / <kbd>Meta</kbd> + <kbd>+</kbd> will add a bew row below the current if an input or textarea is focused

When converting from XML to XLIF, the old .xml file will not be deleted, this must be done manually.

## License

(c) 2019 Georg Großberger <contact@grossberger-ge.org>

Released under the MIT License; see the file [LICENSE](LICENSE) for further information
