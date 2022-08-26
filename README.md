# fix-mp3-tag

## About

A lot of Cyrillic mp3's are broken by f%&amp;king windows encoding.
This binary will try to fix this.

The original fix-mp3-tag was written by [Dmitry Bukin](https://github.com/bukind).
As id3lib cannot properly write Unicode tags, this fork switched to Mutagen.


## Dependency

For reading and writing tags, you need mid3v2, part of [Mutagen](https://mutagen.readthedocs.io/en/latest/index.html).

To install you need a git client, and a go compiler to build.


## Installation

To be updated.


## Usage

```
$GOPATH/bin/fix-mp3-tag <mp3file>...
```

The program will try to decode the id3 tags of the mp3 using the
combination of the cp1251 and iso8859-1 encodings.
