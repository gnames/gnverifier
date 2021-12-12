# Global Names Verifier

[![DOI](https://zenodo.org/badge/297323648.svg)](https://zenodo.org/badge/latestdoi/297323648)

Try `GNverifier` [online][web-service].

Takes a scientific name or a list of scientific names and verifies them against
a variety of biodiversity [Data Sources][data_source_ids]. Includes an advanced
search feature.

<!-- vim-markdown-toc GFM -->

* [Citing](#citing)
* [Features](#features)
* [Installation](#installation)
  * [Using Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2])](#using-homebrew-on-mac-os-x-linux-and-linux-on-windows-wsl2)
  * [MS Windows](#ms-windows)
  * [Linux and Mac](#linux-and-mac)
  * [Compile from source](#compile-from-source)
* [Usage](#usage)
  * [As a web service](#as-a-web-service)
  * [One name-string](#one-name-string)
  * [Many name-strings in a file](#many-name-strings-in-a-file)
  * [Advanced search](#advanced-search)
  * [Options and flags](#options-and-flags)
    * [help](#help)
    * [version](#version)
    * [port](#port)
    * [capitalize](#capitalize)
    * [format](#format)
    * [sources](#sources)
    * [only_preferred](#only_preferred)
    * [all_matches](#all_matches)
    * [jobs](#jobs)
  * [Configuration file](#configuration-file)
  * [Advanced Search Query Language](#advanced-search-query-language)
    * [Examples of searches](#examples-of-searches)
* [Copyright](#copyright)

<!-- vim-markdown-toc -->

## Citing

If you want to cite GNverifier, use [DOI generated by Zenodo][Zenodo DOI]:

## Features

* Small and fast app to verify scientific names against many biodiversity
  databases. The app is a client to a [verifier API].
* It provides 6 different match levels:
  * Exact: complete match with a canonical form or a full name-string from a
     data source.
  * Fuzzy: if exact match did not happen, it tries to match name-strings
     assuming spelling errors.
  * Partial: strips  middle or last epithets from bi- or multi-nomial names
              and tries to match what is left.
  * PartialFuzzy: the same as Partial but assuming spelling mistakes.
  * FacetedSearch: marks [advanced-search](#advanced-search) queries.
* Taxonomic resolution. If a database contains taxonomic information, it
  returns the currently accepted name for the provided name-string.
* Best match is returned according to the match score. Data sources with some
  manual curation have priority over auto-curated and uncurated datasets. For
  example [Catalogue of Life] or [WoRMS] are considered curated,
  [GBIF] auto-curated, [uBio] not curated.
* Fine-tunng  the match score by matching authors, years, ranks etc.
* It is possible to map any name-strings checklist to any of registered
  Data Sources.
* If a Data Source provides a classification for a name, it will be returned in
  the output.
* The app works for checking just one name-string, or multiple ones written in
  a file.
* [Advanced search](#advanced-search) uses simple but powerful
  [query language](#advanced-search-query-language)
  to find abbreviated names, search by author, year etc.
* Supports feeding data via pipes of an operating system. This feature allows
  to chain the program together with other tools.
* `GNverifier` includes a web-based graphical user interface identical to its
  "official" [web-service].

## Installation

### Using Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2])

Homebrew is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on Windows 10, if Windows Subsystem for Linux (WSL) is
[installed][WSL install].

To use `GNverifier` with Homebrew:

1. Install [Homebrew]

2. Open terminal and run the following commands:

```bash
brew tap gnames/gn
brew install gnverifier
```

### MS Windows

Download the latest release from [github], unzip.

One possible way would be to create a default folder for executables and place
``GNverifier`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\Users\your_username\bin
copy path_to\gnverifier.exe C:\Users\your_username\bin
```

[Add ``C:\Users\your_username\bin`` directory to your ``PATH``][winpath]
environment variable.

Another, simpler way, would be to use ``cd C:\Users\your_username\bin`` command
in ``cmd`` terminal window. The ``GNverifier`` program then will be automatically
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

``GNverifier`` takes one name-string or a text file with one name-string per
line as an argument, sends a query with these data to [remote ``gnames``
server][gnames] to match the name-strigs against many different biodiversity
databases and returns results to STDOUT either in JSON, CSV or TSV format.

The app can alto take a query string like `g:M. sp:galloprovincialis
au:Olivier` to perform advanced searching, if the full scientific name is
undetermined.

### As a web service

```bash
gnverifier -p 8080
```

After running this command, you should be able to access web-based user
interface via a browser at ``http://localhost:8080``

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

### Advanced search

Advanced search allows to use a simple but powerful query language to find names
by abbreviated genus, a year or a range of years. See detailed description
in [Advanced Search Query Language](#advanced-search-query-language) section.

```bash
gnverifier "g:B. sp:bubo au:Linn. y:1700-"
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

Starts GNverifier as a web service using entered port

```bash
gnverifier -p 8080
```

This command will run user-interface accessible by a browser
at ``http://localhost:8080``

#### capitalize

If your names are co not have uninomials or genera capitalized according to
rules on nomenclature, you can still verify them using this option. If
`capitalize` flag is set, the first character of every name-string will be
capitalized (when appropriate). This flag is ignores by advanced search.

```bash
gnverifier -c "bubo bubo"
# or
gnverifier --capitalize "bubo bubo"
```

#### format

Allows to pick a format for output. Supported formats are

* compact: one-liner JSON.
* pretty: prettified JSON with new lines and tabs for easier reading.
* tsv: returns tab-separated values representation.
* csv: (DEFAULT) returns comma-separated values representation.

```bash
# short form for compact JSON format
gnverifier -f compact file.txt
# or long form for "pretty" JSON format
gnverifier --format="pretty" file.csv
# tsv format
gnverifier -f tsv file.csv
```

Note that a separate JSON "document" is returned for each separate record,
instead of returning one big JSON document for all records. For large lists it
significantly speeds up parsin of the JSON on the user side.

#### sources

By default ``GNverifier`` returns only one "best" result of a match. If a user
has a particular interest in a data set, s/he can set it with this option, and
all matches that exist for this source will be returned as well. You need to
provide a data source id for a dataset. Ids can be found at the following
[URL][data_source_ids]. Some of them are provided in the ``GNverifier`` help
output as well.

Data from such sources will be returned in preferred_results section of JSON
output, or with CSV/TSV rows that start with "PreferredMatch" string.

```bash
gnverifier file.csv -s "1,11,172"
# or
gnverifier file.tsv --sources="12"
# or
cat file.txt | gnverifier -s '1,12'
```

If all matched sources need to be returned, set the flag to "0".

WARNING: the result might be excessively large.

```bash
gnverifier "Bubo bubo" -s 0
# potentially even more results get returned by adding --all_matches flag
gnverifier "Bubo bubo" -s 0 -M
```
The `sources` option would overwrite `ds:` settings in case of advanced search.

#### only_preferred

Sometimes a users wants to map a list of names to a DataSource. They
are not interested if name matched anywhere else. In such case you can use
the ``only_preferred`` flag.

```bash
gnverifier -o -s '12' file.txt
# or
gnverifier --only_preferred --sources='1,12' file.tsv
```

In case of advanced search use `all:t` together with this flag.

#### all_matches

Sometimes data sources have more than one match to a name. To see all matches
instead of the best one per source use --all_matches flag.

WARNING: for some names the result will be excessively large.

```bash
gnverifier -s '1,12' -M file.txt
```

This flag is ignored by advanced search.

#### jobs

If the list of names if very large, it is possible to tell GNverifier to
run requests in parallel. In this example GNverifier will run 8 processes
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

This option is ignored by advanced search.

### Configuration file

If you find yourself using the same flags over and over again, it makes sense
to edit configuration file instead. It is located at
`$HOME/.config/gnverifier.yaml`. After that you do not need to use command line
options and flags.

```bash
gnverifier file.txt
```

### Advanced Search Query Language

Example: `g:M. sp:gallop. au:Oliv. y:1750-1799`

Query language allows searching for scientific names using name components
like genus name, specific epithet, infraspecific epithet, author, year.
It includes following operators:

`g:`
: Genus name, can be abbreviated (for example `g:Bubo`, `g:B.`).

`sp:`
: specific epithet, can be abbreviated (for example `sp:galloprovincialis`,
  `sp:gallop.`).

`isp:`
: Infraspecific epithet, can be abbreviated (for example `isp:auspicalis`,
  `isp:ausp.`).

`asp:`
: Either specific, or infraspecific epithet (for example `asp:bubo`).

`au:`
: One of the authors of a name, can be abbreviated (for example `au:Linn.`,
  `au:Linnaeus`).

`y:`
: Year. Can be one year, or a year range (for example `y:1888`, `y:1800-1802`,
`y:1756-`, `y:-1880`)

`ds:`
: Limit result to one or more data-sources. Note that command line `sources`
  option, if given, will overwrite this setting (`ds:1,2,172`).

`tx:`
: Parent taxon. Limits results to names that contain a particular clade in
  their classification. If `ds:` is given, uses the classification of the
  first data-source in the setting. If `ds:` is not given, uses managerial
  classification of the Catalogue of Life (`tx:Hemiptera`, `tx:Animalia`,
  `tx:Magnoliopsida`).

`all:`
: If true, [gnverifier] will show all results, not only the best ones.
  The setting can be `true` or `false` (`all:t`, `all:f`). This setting
  will become true if `sources` command line option is set to `0`.

`n:`
: A "name" setting, that allows to combine several query components together
  for convenience. Note that it is not a 'real' scientific name, but a shortcut
  to enter several settings at once (`n:B. bubo Linn. 1758`). This setting
  is limited to [gnparser] functionality, so, for example, it does not
  support abbreviated specific eptithets.

The query language is in `Beta` stage, and might change to some degree, to
improve its functionality.

Often there are errors in species eptithets gender. Because of that search
will try to detect names in any gender that correspond to the epithet.

The search requires to have either `sp:`, `isp:` or `asp:` setting,
or provide their analogs in `n:` setting.

#### Examples of searches

```text
gnverifier "n:Pom. saltator tx:Animalia y:1750-"

gnverifier "g:Plantago asp:major au:Linn."

gnverifier "g:Cara. isp:daurica ds:1,12"
```

## Copyright

Authors: [Dmitry Mozzherin][dimus]

Copyright © 2020-2021 Dmitry Mozzherin. See [LICENSE] for further
details.

[WSL2]: https://docs.microsoft.com/en-us/windows/wsl/install
[verifier API]: https://apidoc.globalnames.org/gnames-beta
[Catalogue of Life]: https://catalogueoflife.org/
[GBIF]: https://www.gbif.org/
[Homebrew]: https://brew.sh/
[WSL install]: https://docs.microsoft.com/en-us/windows/wsl/install-win10
[WoRMS]: https://marinespecies.org/
[Zenodo DOI]: https://zenodo.org/badge/latestdoi/297323648
[data_source_ids]: https://verifier.globalnames.org/data_sources
[dimus]: https://github.com/dimus
[github]: https://github.com/gnames/gnverifier/releases/latest
[gnames]: https://hub.apitree.com/dimus/gnames/
[go-install]: https://golang.org/doc/install
[LICENSE]: https://github.com/gnames/gnverifier/blob/master/LICENSE
[test directory]: https://github.com/gnames/gnverifier/tree/master/testdata
[uBio]: https://ubio.org/
[web-service]: https://verifier.globalnames.org
[win-pdf]: https://github.com/gnames/gnverifier/blob/master/use-gnverifier-windows.pdf
[winpath]: https://www.computerhope.com/issues/ch000549.htm
