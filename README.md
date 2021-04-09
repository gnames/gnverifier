# gnverifier

Takes a name or a list of names and verifies them against a variety of
biodiversity [Data Sources][data_source_ids]

<!-- vim-markdown-toc GFM -->

* [Features](#features)
* [Installation](#installation)
  * [Using Homebrew on Mac OS, Linux, and Linux on Windows X (WSL2)](#using-homebrew-on-mac-os-linux-and-linux-on-windows-x-wsl2)
  * [MS Windows](#ms-windows)
  * [Linux and Mac](#linux-and-mac)
  * [Compile from source](#compile-from-source)
* [Usage](#usage)
  * [As a web service](#as-a-web-service)
  * [One name-string](#one-name-string)
  * [Many name-strings in a file](#many-name-strings-in-a-file)
  * [Options and flags](#options-and-flags)
    * [help](#help)
    * [version](#version)
    * [port](#port)
    * [capitalize](#capitalize)
    * [format](#format)
    * [sources](#sources)
    * [only_preferred](#only_preferred)
    * [jobs](#jobs)
  * [Configuration file](#configuration-file)
* [Copyright](#copyright)

<!-- vim-markdown-toc -->

## Features

* Small and fast app to verify scientific names against many biodiversity
  databases.
* Has 4 different match levels:
  * Exact: complete match with a canonical form or full name-string from a
     data source.
  * Fuzzy: if exact match did not happen, it tries to match name-strings
     assuming spelling errors.
  * Partial: strips  middle or last epithets from bi- or multi-nomial names
              and tries to match what is left.
  * PartialFuzzy: the same as Partial but assuming spelling mistakes.
* Taxonomic resolution. If a database contains taxonomic information, returns
  currently accepted name for a name-string, if it is different from the
  matched name.
* Best match is returned according to the match score. Data sources with some
  manual curation have priority over auto-curated and uncurated datasets. For
  example [Catalogue of Life] or [WoRMS] are considered curated,
  [GBIF] auto-curated, [uBio] not curated.
* It is possible to map any name-strings checklist to any of registered
  Data Sources.
* If a Data Source provides classification for a name, it will be returned in
  the output.
* Works for checking just one name-string, or multiple ones written in a file.
* Supports feeding data via pipes of an operating system. This feature allows
  to chain the program together with other tools.

## Installation

### Using Homebrew on Mac OS, Linux, and Linux on Windows X (WSL2)

Homebrew is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on Windows 10, if Windows Subsystem for Linux (WSL) is
[installed][WSL install].

To use `gnverifier` with Homebrew:

1. Install [Homebrew]

2. Open terminal and run the following commands:

```bash
brew tap gnames/gn
brew install gnverifier
```

### MS Windows

Download the latest release from [github], unzip.

One possible way would be to create a default folder for executables and place
``gnverifier`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\Users\your_username\bin
copy path_to\gnverifier.exe C:\Users\your_username\bin
```

[Add ``C:\Users\your_username\bin`` directory to your ``PATH``][winpath]
environment variable.

Another, simpler way, would be to use ``cd C:\Users\your_username\bin`` command
in ``cmd`` terminal window. The ``gnverifier`` program then will be automatically
found by Windows operating system when you run its commands from that
directory.

You can also read a more detailed guide for Windows users in
[a PDF document][win-pdf].

### Linux and Mac

Download the latest release from [github], untar, and install binary somewhere
in your path.

```bash
tar xvf gnverifier-linux-0.1.0.tar.xz
# or tar xvf gnverifier-mac-0.1.0.tar.gz
sudo mv gnverifier /usr/local/bin
```

### Compile from source

Install Go according to [installation instructions][go-install]

```bash
go get github.com/gnames/gnverifier/gnverifier
```

## Usage

``gnverifier`` takes one name-string or a text file with one name-string per line
as an argument, sends a query with these data to [remote
``gnames`` server][gnames] to match the name-strigs against many different
biodiversity databases and returns results to STDOUT either in JSON or CSV
format.

### As a web service

```bash
gnverifier -p 8080
```

You should be able to access web user interface via a browser at
``http://localhost:8080``

### One name-string

```bash
gnverifier "Monohamus galloprovincialis"
```

### Many name-strings in a file

```bash
gnverifier /path/to/names.txt
```

The app assumes that a file contains a simple list of names, one per line.

It is also possible to feed data via STDIN:

```bash
cat /path/to/names.txt | gnverifier
```

### Options and flags

According to POSIX standard flags and options can be given either before or
after name-string or file name.

#### help

```bash
gnverifier -h
# or
gnverifier --help
# or
gnverifier
```

#### version

```bash
gnverifier -V
# or
gnverifier --version
```

#### port

Starts gnverifier as a web service using entered port

```bash
gnverifier -p 8080
```

This command will run user-interface accessible by a browser
at ``http://localhost:8080``

#### capitalize

If your names are co not have uninomials or genera capitalized according to
rules on nomenclature, you can still verify them using this option. If
`capitalize` flag is set, the first character of every name-string will be
capitalized (when appropriate).

```bash
gnverifier -c "bubo bubo"
# or
gnverifier --capitalize "bubo bubo"
```

#### format

Allows to pick a format for output. Supported formats are

* compact: one-liner JSON.
* pretty: prettified JSON with new lines and tabs for easier reading.
* csv: (DEFAULT) returns CSV representation.

```bash
gnverifier -f compact file.txt
# or
gnverifier --format="pretty" file.csv
```

Note that a separate JSON "document" is returned for each separate record,
instead of returning one big JSON document for all records. For large lists it
significantly speeds up parsin of the JSON on the user side.

#### sources

By default ``gnverifier`` returns only one "best" result of a match. If a user
has a particular interest in a data set, s/he can set it with this option, and
all matches that exist for this source will be returned as well. You need to
provide a data source id for a dataset. Ids can be found at the following
[URL][data_source_ids]. Some of them are provided in the ``gnverifier`` help
output as well.

Data from such sources will be returned in preferred_results section of JSON
output, or with CSV rows that start with "PreferredMatch" string.

```bash
gnverifier file.csv -s "1,11,172"
# or
gnverifier file.tsv --sources="12"
# or
cat file.txt | gnverifier -s '1,12'
```

#### only_preferred

Sometimes a users wants to map a list of names to a DataSource. They
are not interested if name matched anywhere else. In such case you can use
the ``only_preferred`` flag.

```bash
gnverifier -o -s '12' file.txt
# or
gnverifier --only_preferred --sources='1,12' file.tsv
```

#### jobs

If the list of names if very large, it is possible to tell gnverifier to
run requests in parallel. In this example gnverifier will run 8 processes
simultaneously. The order of returned names will be somewhat randomized.

```bash
gnverifier -j 8 file.txt
# or
gnverifier --jobs=8 file.tsv
```

Sometimes it is important to return names in exactly same order. For such
cases set `jobs` flag to 1.

```bash
gnverifier -j 1 file.txt
```

### Configuration file

If you find yourself using the same flags over and over again, it makes sense
to edit configuration file instead. It is located at
`$HOME/.config/gnverifier.yaml`. After that you do not need to use command line
options and flags.

```bash
gnverifier file.txt
```

## Copyright

Authors: [Dmitry Mozzherin][dimus]

Copyright © 2020 Dmitry Mozzherin. See [LICENSE][license] for further details.

[github]: https://github.com/gnames/gnverifier/releases/latest
[gnames]: https://hub.apitree.com/dimus/gnames/
[Catalogue of Life]: https://catalogueoflife.org/
[WoRMS]: https://marinespecies.org/
[GBIF]: https://www.gbif.org/
[uBio]: https://ubio.org/
[test directory]: https://github.com/gnames/gnverifier/tree/master/testdata
[data_source_ids]: http://resolver.globalnames.org/data_sources
[dimus]: https://github.com/dimus
[license]: https://github.com/gnames/gnverifier/blob/master/LICENSE
[winpath]: https://www.computerhope.com/issues/ch000549.htm
[win-pdf]: https://github.com/gnames/gnverifier/blob/master/use-gnverifier-windows.pdf
[go-install]: https://golang.org/doc/install
[WSL install]: https://docs.microsoft.com/en-us/windows/wsl/install-win10
[Homebrew]: https://brew.sh/
