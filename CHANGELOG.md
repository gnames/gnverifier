# Changelog

## Unreleased

## [v1.1.3]

- Add: Wikidata link for the web interface.
- Add: update modules and Go.

## [v1.1.2] - 2023-09-27 Wed

- Add: ent and io are in pkg, so they can be used by other projects.

## [v1.1.1] - 2023-06-22 Thu

- Add [#105]: option for uninomials fuzzy-matching.

## [v1.1.0] - 2023-05-12 Fri

- Add [#104]: name-string widget.
- Add [#103]: name-string endpoint.

## [v1.0.3] - 2023-05-01 Mon

- Add: Nomenclator Zoologicus data-source.

## [v1.0.2] - 2023-03-09 Thu

- Add [#101]: more data-sources for web-UI.
- Add [#100]: refactor to a better file organization.

## [v1.0.1] - 2022-09-30 Fri

- Add: update all modules.

## [v1.0.0] - 2022-08-24 Wed

- Add: prepare for v1.0.0.

## [v1.0.0-RC1] - 2022-05-10 Tue

- Add: update gnmatcher, gnames to v1.0.0-RC1.

## [v0.9.5] - 2022-05-02 Mon

- Add: update gnlib to v0.13.7, gnames to v0.13.3.
- Add: species group option.
- Add: cardinality score

## [v0.9.4] - 2022-04-28 Thu

- Add: update gnlib to v0.13.2

## [v0.9.3] - 2022-04-09 Sat

- Add: update gnlib to v0.13.0

## [v0.9.2] - 2022-04-08 Fri

- Add: update gnlib to v0.12.0
- Add: IRMNG data-source to web UI

## [v0.9.1] - 2022-03-22 Tue

- Add: update Go (v1.18), modules

## [v0.9.0] - 2022-03-13 Sun

- Add [#91]: provide context to verification and searh methods.
- Add: PaleoBioDB to Web GUI

## [v0.8.2] - 2022-02-25 Fri

- Fix: go.mod bug

## [v0.8.1] - 2022-02-25 Fri

- Fix: bug in go.mod

## [v0.8.0] - 2022-02-24 Thu

- Add[#89]: compatibility with gnames v0.8.0

## [v0.7.3] - 2022-02-14 Mon

- Add: make gnverifier compatible with gnames v0.7.1

## [v0.7.2] - 2022-02-10 Thu

- Add [#84]: set filters for NSQ logs via config file.

## [v0.7.1] - 2022-02-08 Tue

- Add [#79]: better NSQ logs, switch to zerolog library.

## [v0.7.0]

- Add [#77]: optional log aggregation by an NSQ-messaging service.
- Add [#76]: verification of viruses.

## [v0.6.6]

- Fix [#75]: faceted search with one letter author without punctuation.

## [v0.6.5]

- Add [#73]: World Checklist of Vascular Plants for web UI.

## [v0.6.4]

- Add [#72]: World Flora Online for the web UI.

## [v0.6.3]

- Add [#70]: add a `-q` flag to suppress log output.
- Add [#71]: Advanced Search `n:` field support abbreviated specific epithets,
  year ranges.

## [v0.6.2]

- Add [#69]: add faceted search to web UI.
- Add [#61]: add faceted search to CLI app.

## [v0.6.1]

- Add [#67]: score details for web-UI.

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

<!-- VERSIONS -->

[v1.1.1]: https://github.com/gnames/gnverifier/compare/v1.1.0...v1.1.1
[v1.1.0]: https://github.com/gnames/gnverifier/compare/v1.0.3...v1.1.0
[v1.0.3]: https://github.com/gnames/gnverifier/compare/v1.0.2...v1.0.3
[v1.0.2]: https://github.com/gnames/gnverifier/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/gnames/gnverifier/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/gnames/gnverifier/compare/v1.0.0-RC1...v1.0.0
[v1.0.0-rc1]: https://github.com/gnames/gnverifier/compare/v0.9.5...v1.0.0-RC1
[v0.9.5]: https://github.com/gnames/gnverifier/compare/v0.9.4...v0.9.5
[v0.9.4]: https://github.com/gnames/gnverifier/compare/v0.9.3...v0.9.4
[v0.9.3]: https://github.com/gnames/gnverifier/compare/v0.9.2...v0.9.3
[v0.9.2]: https://github.com/gnames/gnverifier/compare/v0.9.1...v0.9.2
[v0.9.1]: https://github.com/gnames/gnverifier/compare/v0.9.0...v0.9.1
[v0.9.0]: https://github.com/gnames/gnverifier/compare/v0.8.2...v0.9.0
[v0.8.2]: https://github.com/gnames/gnverifier/compare/v0.8.1...v0.8.2
[v0.8.1]: https://github.com/gnames/gnverifier/compare/v0.8.0...v0.8.1
[v0.8.0]: https://github.com/gnames/gnverifier/compare/v0.7.3...v0.8.0
[v0.7.3]: https://github.com/gnames/gnverifier/compare/v0.7.2...v0.7.3
[v0.7.2]: https://github.com/gnames/gnverifier/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/gnames/gnverifier/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/gnames/gnverifier/compare/v0.6.6...v0.7.0
[v0.6.6]: https://github.com/gnames/gnverifier/compare/v0.6.5...v0.6.6
[v0.6.5]: https://github.com/gnames/gnverifier/compare/v0.6.4...v0.6.5
[v0.6.4]: https://github.com/gnames/gnverifier/compare/v0.6.3...v0.6.4
[v0.6.3]: https://github.com/gnames/gnverifier/compare/v0.6.2...v0.6.3
[v0.6.2]: https://github.com/gnames/gnverifier/compare/v0.6.1...v0.6.2
[v0.6.1]: https://github.com/gnames/gnverifier/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/gnames/gnverifier/compare/v0.5.2...v0.6.0
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

<!-- Issues -->

[#110]: https://github.com/gnames/gnverifier/issues/110
[#109]: https://github.com/gnames/gnverifier/issues/109
[#108]: https://github.com/gnames/gnverifier/issues/108
[#107]: https://github.com/gnames/gnverifier/issues/107
[#106]: https://github.com/gnames/gnverifier/issues/106
[#105]: https://github.com/gnames/gnverifier/issues/105
[#104]: https://github.com/gnames/gnverifier/issues/104
[#103]: https://github.com/gnames/gnverifier/issues/103
[#102]: https://github.com/gnames/gnverifier/issues/102
[#101]: https://github.com/gnames/gnverifier/issues/101
[#100]: https://github.com/gnames/gnverifier/issues/100
[#99]: https://github.com/gnames/gnverifier/issues/99
[#98]: https://github.com/gnames/gnverifier/issues/98
[#97]: https://github.com/gnames/gnverifier/issues/97
[#96]: https://github.com/gnames/gnverifier/issues/96
[#95]: https://github.com/gnames/gnverifier/issues/95
[#94]: https://github.com/gnames/gnverifier/issues/94
[#93]: https://github.com/gnames/gnverifier/issues/93
[#92]: https://github.com/gnames/gnverifier/issues/92
[#91]: https://github.com/gnames/gnverifier/issues/91
[#90]: https://github.com/gnames/gnverifier/issues/90
[#89]: https://github.com/gnames/gnverifier/issues/89
[#88]: https://github.com/gnames/gnverifier/issues/88
[#87]: https://github.com/gnames/gnverifier/issues/87
[#86]: https://github.com/gnames/gnverifier/issues/86
[#85]: https://github.com/gnames/gnverifier/issues/85
[#84]: https://github.com/gnames/gnverifier/issues/84
[#83]: https://github.com/gnames/gnverifier/issues/83
[#82]: https://github.com/gnames/gnverifier/issues/82
[#81]: https://github.com/gnames/gnverifier/issues/81
[#80]: https://github.com/gnames/gnverifier/issues/80
[#79]: https://github.com/gnames/gnverifier/issues/79
[#78]: https://github.com/gnames/gnverifier/issues/78
[#77]: https://github.com/gnames/gnverifier/issues/77
[#76]: https://github.com/gnames/gnverifier/issues/76
[#75]: https://github.com/gnames/gnverifier/issues/75
[#74]: https://github.com/gnames/gnverifier/issues/74
[#73]: https://github.com/gnames/gnverifier/issues/73
[#72]: https://github.com/gnames/gnverifier/issues/72
[#71]: https://github.com/gnames/gnverifier/issues/71
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
