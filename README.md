## Wascra

### Description

Wascra is a tool written in Golang, which lets you extract all relevant Datasheet info from a Warhammer40K (9th edition) faction from [Wahapedia.ru](https://wahapedia.ru) and is powered by the [colly](https://github.com/gocolly/colly) framework.

Wascra searches for faction(s) via command-line arguments and supports terminal output as well as output to file.

For dependencies/ module information please see [go.mod](go.mod)

*Note that this module is [WIP](#future-improvements--planned-fixes).*

---

### Usage

For a full list of provided factions from the wiki, please refer to the [official site](https://wahapedia.ru/wh40k9ed/the-rules/playing-this-game/)


#### Synopsis

To scrape the desired faction(s), just execute the following command in a terminal:

```sh
go run wascra.go [{-flags}] {faction1} [faction2]
```


#### Flags

Supported flags are:

`-v` - for verbose output while scraping (shows URLs of visited models)

`-w` - writes data to an auto-generated file to `./factions/` (filename is a concatenation of given arguments/factions)

`-wv` - write verbose

`-j` - serializes scraped data with JSON

`-clear` - clears the cache

`-h` - lists additional help

***Note that `-clear` exits the program, so there will be no scraping even if other flags and arguments are provided!***


#### Examples

Scraping Space Marines, verbose output to console: 

```sh
go run wascra.go -v space-marines
```

Scraping Orks and Genestealer-Cults, write JSON to file: [`factions/orks_genestealer-cults.json`](/factions/orks_genestealer-cults.json):

```sh
go run wascra.go -w -j orks genestealer-cults
```

***Note:*** Faction arguments can be lower or uppercase, while whitespaces must be replaced with hyphens `-`.


---

### Future improvements / planned fixes

- [ ] - fix selector strings

- [ ] - correctly parse multiple instances of datasheet values

- [x] - serialize Model struct for JSON export

- [x] - add JSON support

- [ ] - concurrent scraping for multiple factions

- [x] - add option to clear cache with flag/ provide .sh