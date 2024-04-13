# Global Names Verifier

[![DOI](https://zenodo.org/badge/297323648.svg)](https://zenodo.org/badge/latestdoi/297323648)

**Warning**: Version v1.2.0 introduces some backward incompatible features:
Some flags for command line application are changed, CSV output now returns
`TaxonomicStatus` instead of `IsSynonym`. The term `isSynonym` stays in JSON
format for backward compatibility, but is deprecated.

Try `GNverifier` [online][web-service].

[GNverifier with OpenRefine]

[GNverifier API]

[Feedback]

Takes a scientific name or a list of scientific names and verifies them against
a variety of biodiversity [Data Sources][data_source_ids]. Includes an advanced
search feature.

<!-- vim-markdown-toc GFM -->

* [Citing](#citing)
* [Features](#features)
* [Installation](#installation)
  * [Using Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2])](#using-homebrew-on-mac-os-x-linux-and-linux-on-windows-wsl2)
  * [MS Windows](#ms-windows)
  * [Linux and Mac (without Homebrew)](#linux-and-mac-without-homebrew)
  * [Compile from source](#compile-from-source)
* [Usage](#usage)
  * [As a web service](#as-a-web-service)
  * [As a RESTful API](#as-a-restful-api)
  * [One name-string](#one-name-string)
  * [Many name-strings in a file](#many-name-strings-in-a-file)
  * [Advanced search](#advanced-search)
  * [Options and flags](#options-and-flags)
    * [help](#help)
    * [version](#version)
    * [port](#port)
    * [all_matches](#all_matches)
    * [capitalize](#capitalize)
    * [species group](#species-group)
    * [relaxed fuzzy-match](#relaxed-fuzzy-match)
    * [fuzzy-match of uninomial names](#fuzzy-match-of-uninomial-names)
    * [format](#format)
    * [jobs](#jobs)
    * [quiet](#quiet)
    * [sources](#sources)
  * [Configuration file](#configuration-file)
  * [Advanced Search Query Language](#advanced-search-query-language)
    * [Examples of searches](#examples-of-searches)
* [Copyright](#copyright)

<!-- vim-markdown-toc -->

## Citing

If you want to cite GNverifier, use [DOI generated by Zenodo][zenodo doi]:

## Features

- Small and fast app to verify scientific names against many biodiversity
  databases. The app is a client to a [verifier API].
- It provides different match levels:
  - **Exact**: complete match with a canonical form or a full name-string
    from a data source.
  - **Fuzzy**: if exact match did not happen, it tries to match name-strings
    assuming spelling errors.
  - **FuzzyRelaxed**: if exact match did not happen, it tries to match
    name-strings using 'relaxed' fuzzy-matching rules.
  - **Partial**: strips middle or last epithets from bi- or multi-nomial names
    and tries to match what is left.
  - **PartialFuzzy**: the same as Partial but assuming spelling mistakes.
  - **PartialFuzzyRelaxed**: the same as PartialFuzzy but with relaxed
    fuzzy-matchng rules
  - **Virus**: verification of virus names.
  - **FacetedSearch**: marks [advanced-search](#advanced-search) queries.
- Fuzzy matching that tries to balance number of false positives and false
  negatives (more information on: [fuzzy-matching]).
- Taxonomic resolution. If a database contains taxonomic information, it
  returns the currently accepted name for the provided name-string.
- Best match is returned according to the match score. Data sources with some
  manual curation have priority over auto-curated and uncurated datasets. For
  example [Catalogue of Life] or [WoRMS] are considered curated,
  [GBIF] auto-curated, [uBio] not curated.
- Fine-tuning the match score by matching authors, years, ranks etc.
- It is possible to map any name-strings checklist to any of registered
  Data Sources.
- If a Data Source provides a classification for a name, it will be returned to
  the output.
- The app works for checking just one name-string, or multiple ones written in
  a file.
- [Advanced search](#advanced-search) uses simple but powerful
  [query language](#advanced-search-query-language)
  to find abbreviated names, search by author, year etc.
- Supports feeding data via pipes of an operating system. This feature allows
  to chain the program together with other tools.
- [GNverifier] includes a web-based graphical user interface identical to its
  "official" [web-service].

## Installation

### Using Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2])

Homebrew is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on Windows 10 or 11, if Windows Subsystem for Linux (WSL) is
[installed][wsl install].

To use [GNverifier] with Homebrew:

1. Install [Homebrew]

2. Open terminal and run the following commands:

```bash
brew tap gnames/gn
brew install gnverifier
```

### MS Windows

Download the [latest release] from GitHub, unzip.

One possible way would be to create a default folder for executables and place
`GNverifier` there.

Use `Windows+R` keys
combination and type "`cmd`". In the appeared terminal window type:

```cmd
mkdir C:\Users\your_username\bin
copy path_to\gnverifier.exe C:\Users\your_username\bin
```

[Add `C:\Users\your_username\bin` directory to your `PATH`][winpath] `user`
and/or `system` environment variable.

Another, simpler way, would be to use `cd C:\Users\your_username\bin` command
in `cmd` terminal window. The [GNverifier] program then will be automatically
found by Windows operating system when you run its commands from that
directory.

You can also read a more detailed guide for Windows users in
[a PDF document][win-pdf].

### Linux and Mac (without Homebrew)

If [Homebrew] is not installed, download the [latest release] from GitHub,
untar, and install binary somewhere in your path.

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

[GNverifier] takes one name-string or a text file with one name-string per
line as an argument, sends a query with these data to a [remote GNames
server][gnames] to match the name-strings against many biodiversity
databases and returns results to STDOUT either in JSON, CSV or TSV format.

The app can alto take a query string like
`g:M. sp:galloprovincialis au:Olivier` to perform advanced searching,
if the full scientific name is undetermined.

### As a web service

```bash
gnverifier -p 8080
```

After running this command, you should be able to access web-based user
interface via a browser at `http://localhost:8080`

### As a RESTful API

Refer to the [RESTful API docs][gnames] to learn how to use the same
functionality via scripts.

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
at `http://localhost:8080`

#### all_matches

To see all matches instead of the best one use --all_matches flag.

WARNING: for some names the result will be excessively large.

```bash
gnverifier -s '1,12' -M file.txt
gnverifier --all_matches "Pardosa moesta"
```

This flag is ignored by advanced search.

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

#### species group

If `species_group` flag is on, a search of `Aus bus` would also search for
`Aus bus bus` and vice versa. This flag expands search to a species group of
a name if applicable. It means it involves into search botanical autonyms and
coordinated names in zoology.

```bash
gnverifier -G "Bubo bubo"
gnverifier  --species_group "Bubo bubo"
```

#### relaxed fuzzy-match

Relaxes fuzzy-matching rules, allowing fuzzy match for words of any size, and
increasing maximum edit distance (for stems) to two. This creates many more
false positives, but increases recall. It is recommended to check results by
hand if this feature is enabled. The maximum number of names allowed when this
option is enabled is 50.

```bash
gnverifier -R "Bbo bbo"
gnverifier --fuzzy_relaxed "Bbo bbo"
```

#### fuzzy-match of uninomial names

When `fuzzy_uninomial` flag is on, uninomials are allowed to go through
fuzzy matching, if needed. Normally this flag is off because fuzzy-matched
uninomials create a significant amount of false positives.

```bash
gnverifier -U "Pomatmus"
gnverifier --fuzzy_uninomial "Pomatmus"
```

#### format

Allows to pick a format for output. Supported formats are

- compact: one-liner JSON.
- pretty: prettified JSON with new lines and tabs for easier reading.
- tsv: returns tab-separated values representation.
- csv: (DEFAULT) returns comma-separated values representation.

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
significantly speeds up parsing of the JSON on the user side.

#### jobs

If the list of names if very large, it is possible to tell [GNverifier] to
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

#### quiet

Removes log messages from the output. Note that results of verification go
to STDOUT, while log messages go to STDERR. So instead of using `-q` flag
STDERR can be redirected to `/dev/null`:

```bash
gnverifier "Puma concolor" -q >verif-results.csv

#or

gnverifier "Puma concolor 2>/dev/null >verif-results.csv
```

#### sources

By default [GNverifier] returns only one "best" result of a match. If a user
has a particular interest in a data set, s/he can set it with this option, and
all matches that exist for this source will be returned as well. You need to
provide a data source id for a dataset. Ids can be found at the following
[URL][data_source_ids]. Some of them are provided in the GNverifier help
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

### Configuration file

If you find yourself using the same flags over and over again, it makes sense
to edit configuration file instead. It is located at
`$HOME/.config/gnverifier.yaml`. After that you do not need to use command line
options and flags. Configuration file is self-documented, the [default
gnverifier.yaml] is located on GitHub

```bash
gnverifier file.txt
```

In case if [GNverifier] runs as a web-based user interface, it is also
possible to use environment variables for configuration.

| Env. Var.               | Configuration      |
| :---------------------- | :----------------- |
| GNV_FORMAT              | Format             |
| GNV_DATA_SOURCES        | DataSources        |
| GNV_WITH_ALL_MATCHES    | WithAllMatches     |
| GNV_WITH_CAPITALIZATION | WithCapitalization |
| GNV_VERIFIER_URL        | VerifierURL        |
| GNV_JOBS                | Jobs               |

### Advanced Search Query Language

Example: `g:M. sp:gallop. au:Oliv. y:1750-1799` or `n:M. gallop. Oliv. 1750-1799`

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
: Parent taxon. Limit results to names that contain a particular higher taxon
in their classification. If `ds:` is given, uses the classification of the
first data-source in the setting. If `ds:` is not given, uses managerial
classification of the Catalogue of Life (`tx:Hemiptera`, `tx:Animalia`,
`tx:Magnoliopsida`).

`all:`
: If true, [GNverifier] will show all results, not only the best ones.
The setting can be `true` or `false` (`all:t`, `all:f`). This setting
will also become true if `sources` command line option is set to `0`.

`n:`
: A "name" setting. It allows to combine several query components together
for convenience. Note that it is not a 'real' scientific name, but a shortcut
to enter several settings at once loosely following rules of nomenclature
(`n:B. bubo Linn. 1758`). For example, in contrast with GNparser results, it
is possible to have abbreviated specific epithets or range in
years: `n:Mono. gall. Oliv. 1750-1800`.

Often there are errors in species epithets gender. Because of that search
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

Copyright © 2020-2024 Dmitry Mozzherin. See [LICENSE] for further
details.

[Feedback]: https://github.com/gnames/gnverifier/issues
[GNverifier API]: https://apidoc.globalnames.org/gnames
[GNverifier with OpenRefine]: https://github.com/gnames/gnverifier/wiki/OpenRefine-readme
[catalogue of life]: https://catalogueoflife.org/
[data_source_ids]: https://verifier.globalnames.org/data_sources
[default gnverifier.yaml]: https://github.com/gnames/gnverifier/blob/master/gnverifier/cmd/gnverifier.yaml
[dimus]: https://github.com/dimus
[fuzzy-matching]: https://github.com/gnames/gnverifier/blob/master/fuzzy-matching.md
[gbif]: https://www.gbif.org/
[gnames]: https://apidoc.globalnames.org/gnames
[gnverifier]: https://github.com/gnames/gnverifier
[go-install]: https://golang.org/doc/install
[homebrew]: https://brew.sh/
[latest release]: https://github.com/gnames/gnverifier/releases/latest
[license]: https://github.com/gnames/gnverifier/blob/master/LICENSE
[test directory]: https://github.com/gnames/gnverifier/tree/master/testdata
[ubio]: https://ubio.org/
[verifier api]: https://apidoc.globalnames.org/gnames
[web-service]: https://verifier.globalnames.org
[win-pdf]: https://github.com/gnames/gnverifier/blob/master/use-gnverifier-windows.pdf
[winpath]: https://www.computerhope.com/issues/ch000549.htm
[worms]: https://marinespecies.org/
[wsl install]: https://docs.microsoft.com/en-us/windows/wsl/install-win10
[wsl2]: https://docs.microsoft.com/en-us/windows/wsl/install
[zenodo doi]: https://zenodo.org/badge/latestdoi/297323648

