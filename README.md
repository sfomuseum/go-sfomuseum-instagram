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
    "body": "In 1994, Gilbert Baker, the original creator of the rainbow flag and a team of volunteers created a mile-long rainbow flag for the 25th Anniversary of the 1969 Stonewall riots. The flag was carried by 5,000 people on First Avenue in New York City. Baker worked tirelessly to ensure the rainbow flag would become a powerful and enduring symbol of pride and inclusion that transcends languages and borders, gender and race, and now, four decades after its creation, generations.  Courtesy of Mick Hicks. See \"#ALegacyOfPride: #GilbertBaker and the 40th Anniversary of the #RainbowFlag\" on view pre-security in the International Terminal. http://bit.ly/RainbowFlagSFO",
    "excerpt": "In 1994, Gilbert Baker, the original creator of the rainbow flag and a team of volunteers created a mile-long rainbow flag for the 25th Anniversary of the 1969 Stonewall riots.",
    "hashtags": [
      "GilbertBaker", 
      "RainbowFlag", 
      "ALegacyofPride", 
      "gaypride", 
      "lgbtpride", 
      "pride", 
      "LGBT", 
      "🏳️‍🌈", 
      "👭", 
      "👬", 
      "👩‍❤️‍👩", 
      "👩‍❤️‍💋‍👩", 
      "👨‍❤️‍👨", 
      "👨‍❤️‍💋‍👨"
    ],
    "users": []
  },
  "taken_at": "2018-09-20T03:40:04+00:00",
  "location": "San Francisco International Airport (SFO)",
  "path": "photos/201809/0ebfa6dda7247127fb67475768299db2.jpg",
  "taken": 1537414804,
  "media_id": "0ebfa6dda7247127fb67475768299db2"
}
  ... and so on
```