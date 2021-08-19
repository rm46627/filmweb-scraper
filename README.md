# want
Want is a data scraper for https://www.filmweb.pl/. Visits movies pages from "want to see" page of given user,  
creates txt file with every movie that is available on any vod and lists those vod sites.



#### how to use:
```
go run main.go tosee `username`
```
or
```
go install /path/to/want
want tosee `username`
```

#### example output:
```
vod counter:

netflix: 3
chili.com: 1
apple: 2

address: https://www.filmweb.pl/film/Believe+Me%3A+The+Abduction+of+Lisa+McVey-2018-813938/vod
vod:
	[netflix]
===========

address: https://www.filmweb.pl/film/Le+Mans+%2766-2019-705791/vod
vod:
	[netflix apple]
===========

address: https://www.filmweb.pl/film/Kosmiczny+mecz-1996-582/vod
vod:
	[netflix apple chili.com]
===========
```
