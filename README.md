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
> ./bin/emit -h
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
  -query value
    	One or more {PATH}={REGEXP} parameters for filtering records.
  -query-mode string
    	Specify how query filtering should be evaluated. Valid modes are: ALL, ANY (default "ALL")
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
	-media-uri file:///usr/local/instagram/media.json
	
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

#### Inline queries

You can also specify inline queries by passing a `-query` parameter which is a string in the format of:

```
{PATH}={REGULAR EXPRESSION}
```

Paths follow the dot notation syntax used by the [tidwall/gjson](https://github.com/tidwall/gjson) package and regular expressions are any valid [Go language regular expression](https://golang.org/pkg/regexp/). Successful path lookups will be treated as a list of candidates and each candidate's string value will be tested against the regular expression's [MatchString](https://golang.org/pkg/regexp/#Regexp.MatchString) method.

For example:

```
 $> ./bin/emit \
 	-json \
	-format-json \
	-query 'caption=🏳️‍🌈' \
	-media-uri file:///usr/local/instagram/sfomuseum_20201008_part_2/media.json

[{
  "caption": "In June 15, 2003, LGBTQ activist and labor organizer Cleve Jones (with hand raised) and colleagues unfurled a 1.25 mile-long “Sea-to-Sea” flag in Key West, Florida, for the 25th Anniversary of the rainbow flag stretching from the Gulf of Mexico to the Atlantic Ocean. The flag was created by Gilbert Baker, the original creator of the rainbow flag, and a team of volunteers. \nPhoto courtesy of Mick Hicks. See \"#ALegacyOfPride: #GilbertBaker and the 40th Anniversary of the #RainbowFlag\" on view pre-security in the International Terminal. http://bit.ly/RainbowFlagSFO\n.\n.\n.\n#RainbowFlagSFO #GilbertBaker #RainbowFlag #ALegacyofPride #CleveJones #gaypride #lgbtpride #pride #LGBT #🏳️‍🌈 #👭 #👬 #👩‍❤️‍👩 #👩‍❤️‍💋‍👩 #👨‍❤️‍👨 #👨‍❤️‍💋‍👨",
  "taken_at": "2018-12-01T01:57:11+00:00",
  "location": "San Francisco International Airport (SFO)",
  "path": "photos/201811/d5a2cc3513b40164f52a9607d13a69a2.jpg"
}

,{
  "caption": "National #ComingOutDay was yesterday! The creator of the rainbow flag, Gilbert Baker, worked tirelessly to ensure the rainbow flag would become a powerful and enduring symbol of pride and inclusion that transcends languages and borders, gender and race, and now, four decades after its creation, generations. After suffering a stroke in 2012, Baker retaught himself how to sew, continued to make art every day, and witnessed the full appreciation of the symbol he created. This included its recognition as a historically important design icon when acquired by the Museum of Modern Art’s design collection in 2015, and a White House ceremony in which Baker presented President Barack Obama with a framed copy of his original eight-stripe rainbow flag in 2016. Shown here is the White House illuminated with the colors of the flag in celebration of the Supreme Court ruling legalizing same-sex marriage in the United States.\nPhoto by Pete Souza. See \"#ALegacyOfPride: #GilbertBaker and the 40th Anniversary of the #RainbowFlag\" on view pre-security in the International Terminal. http://bit.ly/RainbowFlagSFO\n.\n.\n.\n#GilbertBaker #RainbowFlag #ALegacyofPride #gaypride #lgbtpride #pride #LGBT #🏳️‍🌈 #👭 #👬 #👩‍❤️‍👩 #👩‍❤️‍💋‍👩 #👨‍❤️‍👨 #👨‍❤️‍💋‍👨",
  "taken_at": "2018-10-12T15:16:03+00:00",
  "location": "SFO Museum",
  "path": "photos/201810/414c02a5c136cedec6abf90ad2eda4e2.jpg"
}
...and so on
]
```

You can pass multiple `-query` parameters:

```
 $> ./bin/emit \
 	-json \
	-format-json \
	-query 'caption=behindthescenes' \
	-query 'caption=conservation' \	 
	-media-uri file:///usr/local/instagram/sfomuseum_20201008_part_2/media.json

[{
  "caption": "Our conservators are working on a months-long project to restore this United Air Lines DC-8 cutaway model. The 1:10 model features dozens of people in various positions passing time on their flight, from knitting to reading magazines to eating meals. The model is from the late 1950s and has areas of loss and surface dirt. Several missing parts, such as feet, hands, and armrests, were cast in plaster from intact components and are carefully being shaped and finished, and will eventually be painted to completely blend in. .\n.\n.\n#conservation #museum #museummonday #aircraftmodel #airplanemodel #plastercast #behindthescenes",
  "taken_at": "2019-08-12T18:41:37+00:00",
  "location": "",
  "path": "photos/201908/7c9060f8784378e2f0a3b4bd44e776ed.jpg"
}

,{
  "caption": "Our conservators are working on a months-long project to restore this United Air Lines DC-8 cutaway model. The 1:10 model features dozens of people in various positions passing time on their flight, from knitting to reading magazines to eating meals. The model is from the late 1950s and has areas of loss and surface dirt. Several missing parts, such as feet, hands, and armrests, were cast in plaster from intact components and are carefully being shaped and finished, and will eventually be painted to completely blend in. .\n.\n.\n#conservation #museum #museummonday #aircraftmodel #airplanemodel #plastercast #behindthescenes",
  "taken_at": "2019-08-12T18:41:37+00:00",
  "location": "",
  "path": "photos/201908/15b00a7463f4cad8f9d38b21401b61f7.jpg"
}
...and so on
]
```

The default query mode is to ensure that all queries match but you can also specify that only one or more queries need to match by passing the `-query-mode ANY` flag:

```
 $> ./bin/emit \
 	-json \
	-format-json \
	-query 'caption=behindthescenes' \
	-query 'caption=conservation' \
	-query-mode ANY \
	-media-uri file:///usr/local/instagram/sfomuseum_20201008_part_2/media.json

[
{
  "caption": "There's always something going on here, no matter what time of the year it is. This new year we'll be presenting a exhibitions on hand-shaped wooden surfboards. Our photographer constructed this custom set up to photograph these oversized objects. Stay tuned for more updates about \"Reflections in Wood — Surfboards and Shapers!\"\n.\n.\n.\n#SurfboardsAndShapers #behindthescenes #photography#surfboards #surfer #shapersurfboards #shapers",
  "taken_at": "2019-01-02T18:43:14+00:00",
  "location": "SFO Museum",
  "path": "photos/201901/b3093ec7359023b0f8d5c69d817e7bc9.jpg"
}
... and so on
]
```

## See also

* https://github.com/aaronland/go-json-query

