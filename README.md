# go-sfomuseum-instagram

Go package for working with Instagram archives.

## Important

Work in progress. Documentation to follow.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/emit cmd/emit/main.go
```

### emit

A command-line tool for parsing and emitting individual records from an Instagram `media.json` file.

```
$> ./bin/emit -h
Usage of ./bin/emit:
  -append-all -append-
    	Enable all the -append- flags.
  -append-id media_id
    	Append a media_id property derived from the `path` property.
  -append-timestamp taken
    	Append a taken property containing a Unix timestamp derived from the `taken_at` property.
  -expand-caption caption
    	Parse and replace the string caption property with a `document.Caption` struct.
  -format-json
    	Format JSON output for each record.
  -json
    	Emit a JSON list.
  -media-uri media.json
    	A valid gocloud.dev/blob URI to your Instagram media.json file.
  -null
    	Emit to /dev/null
  -stdout
    	Emit to STDOUT (default true)
```

For example:

```
$> ./bin/emit \
	-append-all  \
	-expand-caption \
	-json \
	-format-json \
| jq

{
    "caption": {
      "body": "When it went up and the wind finally took it out of my hands, it blew my mind…I saw immediately how everyone around me owned the flag. I thought, 'It's better than I ever dreamed.' - Gilbert Baker (1951-2017)  Did you know the #rainbowflag was created right here in San Francisco in 1978? To commemorate the 40th anniversary of the flag's creation, SFO Museum is presenting an exhibition about Gilbert Baker, the flag's creator featuring a flag created by Baker for the ABC television miniseries When We Rise in 2016. The sewing machine and table on display were also used by Gilbert Baker in the show's production.  See \"#ALegacyOfPride: #GilbertBaker and the 40th Anniversary of the #RainbowFlag\" on view pre-security in the International Terminal. http://bit.ly/RainbowFlagSFO",
      "hashtags": [
        "GilbertBaker",
        "RainbowFlag",
        "ALegacyofPride",
        "gaypride",
        "lgbtpride",
        "pride",
        "LGBT",
        "🏳️‍🌈"
      ],
      "users": []
    },
    "taken_at": "2018-08-09T16:24:20+00:00",
    "location": "SFO Museum",
    "path": "photos/201808/956afae7ece4253954366ae5973385c4.jpg",
    "taken": 1533831860,
    "media_id": "956afae7ece4253954366ae5973385c4"
  }
  ... and so on
```