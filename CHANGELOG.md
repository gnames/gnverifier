# Changelog

## [v0.6.0]

- Add [#66]: migrate old functionality to use gnames-beta (v0.5.6+)

## [v0.5.2]

- Add [#58]: add all matches/all sources to web interface.

## [v0.5.1]

- Fix: command line and configuration options

## [v0.5.0]

- Add [#57]: optionally, return all matches per source.
- Add [#48]: optionally, return results from all sources.

## [v0.4.1]

- Add [#56]: add AlgaeBase to web UI

## [v0.4.0]

- Add [#54]: add TSV format
- Fix [#55]: remove empty outputs from web-interface.

## [v0.3.3]

- Add [#52]: add DOI via Zenodo for citations.

## [v0.3.2]

- Add [#51]: show import date on web-pages for the verification records.

## [v0.3.1]

- Add [#49]: avoid possible race conditions by reimplementing GNverifier
             interface.

## [v0.3.0]

- Add [#45]: an option to capitalize names with low-case uninomials and genera.

## [v0.2.5]

- rename GNVerify to GNverifier to be consistent with GNparser.

## [v0.2.4]

- Add [#42]: rename project to gnverifier.

## [v0.2.3]

- Add [#40]: redirect to GET if names number is small.

## [v0.2.2]

- Add [#39]: add Wikispecies and Plazi to web-interface.

## [v0.2.1]

- Add [#38]: move web UI to POST from GET.
- Add [#34]: improve tests.
- Fix [#36]: correct output for preferred_only CSV files.

## [v0.2.0]

- Add [#34]: public service at globalnames.org.
- Add [#33]: use `embed` from Go 1.16 and remove `rice` dependency.
- Add [#32]: documenation for web UI.
- Add [#31]: install tools via Makefile.
- Add [#30]: display data-sources meta-information via web UI.
- Add [#29]: CSV, JSON, HTML formats for web UI.
- Add [#28]: verify names using preferred data sources in web UI.
- Add [#27]: verify names via web GUI.
- Add [#22]: cancel timed-out requests locally and on the server.
- Add [#25]: 'Error' fields in the csv output.

## [v0.1.0]

- Add [#19]: Add Homebrew package.
- Add [#15]: add jobs option.
- Add [#14]: add config file.
- Add [#8]: substitute zeroes to empty fields in CSV where it makes sense.
- Add [#5]: improve documentation.
- Add [#4]: improve architecture.
- Add [#3]: verify names from a file or STDIN.
- Add [#2]: remote verification with output in JSON and CSV formats.
- Add [#1]: command line interface.
- Fix [#9]: verification exits too early.

## Footnotes

This document follows [changelog guidelines]

[v0.6.0]: https://github.com/gnames/gnverifier/compare/v0.5.2...v1.6.0
[v0.5.2]: https://github.com/gnames/gnverifier/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/gnames/gnverifier/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/gnames/gnverifier/compare/v0.4.1...v0.5.0
[v0.4.1]: https://github.com/gnames/gnverifier/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/gnames/gnverifier/compare/v0.3.3...v0.4.0
[v0.3.3]: https://github.com/gnames/gnverifier/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/gnames/gnverifier/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/gnames/gnverifier/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/gnames/gnverifier/compare/v0.2.5...v0.3.0
[v0.2.5]: https://github.com/gnames/gnverifier/compare/v0.2.4...v0.2.5
[v0.2.4]: https://github.com/gnames/gnverifier/compare/v0.2.3...v0.2.4
[v0.2.3]: https://github.com/gnames/gnverifier/compare/v0.2.2...v0.2.3
[v0.2.2]: https://github.com/gnames/gnverifier/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/gnames/gnverifier/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/gnames/gnverifier/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/gnames/gnverifier/tree/v0.1.0

[#70]: https://github.com/gnames/gnverifier/issues/70
[#69]: https://github.com/gnames/gnverifier/issues/69
[#68]: https://github.com/gnames/gnverifier/issues/68
[#67]: https://github.com/gnames/gnverifier/issues/67
[#66]: https://github.com/gnames/gnverifier/issues/66
[#65]: https://github.com/gnames/gnverifier/issues/65
[#64]: https://github.com/gnames/gnverifier/issues/64
[#63]: https://github.com/gnames/gnverifier/issues/63
[#62]: https://github.com/gnames/gnverifier/issues/62
[#61]: https://github.com/gnames/gnverifier/issues/61
[#60]: https://github.com/gnames/gnverifier/issues/60
[#59]: https://github.com/gnames/gnverifier/issues/59
[#58]: https://github.com/gnames/gnverifier/issues/58
[#57]: https://github.com/gnames/gnverifier/issues/57
[#56]: https://github.com/gnames/gnverifier/issues/56
[#55]: https://github.com/gnames/gnverifier/issues/55
[#54]: https://github.com/gnames/gnverifier/issues/54
[#53]: https://github.com/gnames/gnverifier/issues/53
[#52]: https://github.com/gnames/gnverifier/issues/52
[#51]: https://github.com/gnames/gnverifier/issues/51
[#50]: https://github.com/gnames/gnverifier/issues/50
[#49]: https://github.com/gnames/gnverifier/issues/49
[#48]: https://github.com/gnames/gnverifier/issues/48
[#47]: https://github.com/gnames/gnverifier/issues/47
[#46]: https://github.com/gnames/gnverifier/issues/46
[#45]: https://github.com/gnames/gnverifier/issues/45
[#44]: https://github.com/gnames/gnverifier/issues/44
[#43]: https://github.com/gnames/gnverifier/issues/43
[#42]: https://github.com/gnames/gnverifier/issues/42
[#41]: https://github.com/gnames/gnverifier/issues/41
[#40]: https://github.com/gnames/gnverifier/issues/40
[#39]: https://github.com/gnames/gnverifier/issues/39
[#38]: https://github.com/gnames/gnverifier/issues/38
[#37]: https://github.com/gnames/gnverifier/issues/37
[#36]: https://github.com/gnames/gnverifier/issues/36
[#35]: https://github.com/gnames/gnverifier/issues/35
[#34]: https://github.com/gnames/gnverifier/issues/34
[#33]: https://github.com/gnames/gnverifier/issues/33
[#32]: https://github.com/gnames/gnverifier/issues/32
[#31]: https://github.com/gnames/gnverifier/issues/31
[#30]: https://github.com/gnames/gnverifier/issues/30
[#29]: https://github.com/gnames/gnverifier/issues/29
[#28]: https://github.com/gnames/gnverifier/issues/28
[#27]: https://github.com/gnames/gnverifier/issues/27
[#26]: https://github.com/gnames/gnverifier/issues/26
[#25]: https://github.com/gnames/gnverifier/issues/25
[#24]: https://github.com/gnames/gnverifier/issues/24
[#23]: https://github.com/gnames/gnverifier/issues/23
[#22]: https://github.com/gnames/gnverifier/issues/22
[#21]: https://github.com/gnames/gnverifier/issues/21
[#20]: https://github.com/gnames/gnverifier/issues/20
[#19]: https://github.com/gnames/gnverifier/issues/19
[#18]: https://github.com/gnames/gnverifier/issues/18
[#17]: https://github.com/gnames/gnverifier/issues/17
[#16]: https://github.com/gnames/gnverifier/issues/16
[#15]: https://github.com/gnames/gnverifier/issues/15
[#14]: https://github.com/gnames/gnverifier/issues/14
[#13]: https://github.com/gnames/gnverifier/issues/13
[#12]: https://github.com/gnames/gnverifier/issues/12
[#11]: https://github.com/gnames/gnverifier/issues/11
[#10]: https://github.com/gnames/gnverifier/issues/10
[#9]: https://github.com/gnames/gnverifier/issues/9
[#8]: https://github.com/gnames/gnverifier/issues/8
[#7]: https://github.com/gnames/gnverifier/issues/7
[#6]: https://github.com/gnames/gnverifier/issues/6
[#5]: https://github.com/gnames/gnverifier/issues/5
[#4]: https://github.com/gnames/gnverifier/issues/4
[#3]: https://github.com/gnames/gnverifier/issues/3
[#2]: https://github.com/gnames/gnverifier/issues/2
[#1]: https://github.com/gnames/gnverifier/issues/1

[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
