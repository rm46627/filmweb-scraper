# want
Want is a CLI data scraper for https://www.filmweb.pl/. Visits movies from "want to see" page of given user,  
creates txt file and lists movies that are available on any vod site.

#### how to use:
```
$ go run main.go tosee `username`
```
or
```
$ go install /path/to/want
$ want tosee `username`
```

#### example output:
```
vod counter:

apple: 4
canalplus: 1
chili.com: 4
kinopodbaranami.pl: 2
mojeekino.pl: 2
netflix: 2
nowehoryzonty.pl: 2
player.pl: 1
vod.tvp.pl: 1

address: https://www.filmweb.pl/film/Gdyby+ulica+Beale+umia%C5%82a+m%C3%B3wi%C4%87-2018-805105/
vod: 
	apple
	chili.com
-----------------------------
address: https://www.filmweb.pl/film/To+my-2019-816982/
vod: 
	chili.com

----------------------------------------
address: https://www.filmweb.pl/film/Pewnego+dnia-2018-807899/
vod: 
	kinopodbaranami.pl

----------------------------------------
address: https://www.filmweb.pl/film/Tootsie-1982-1043/
vod: 
	apple

----------------------------------------
address: https://www.filmweb.pl/film/Toy+Story+4-2019-632733/
vod: 
	canalplus
	player.pl
	apple
	chili.com

----------------------------------------
address: https://www.filmweb.pl/film/Styks-2018-804145/
vod: 
	vod.tvp.pl

----------------------------------------
address: https://www.filmweb.pl/film/Przyn%C4%99ta-2019-834132/
vod: 
	nowehoryzonty.pl

----------------------------------------
address: https://www.filmweb.pl/film/Pami%C4%85tka-2019-818322/
vod: 
	netflix

----------------------------------------
address: https://www.filmweb.pl/film/Kwiat+szcz%C4%99%C5%9Bcia-2019-822592/
vod: 
	kinopodbaranami.pl
	mojeekino.pl

----------------------------------------
address: https://www.filmweb.pl/film/Opowie%C5%9B%C4%87+o+trzech+siostrach-2018-817881/
vod: 
	mojeekino.pl

----------------------------------------
address: https://www.filmweb.pl/film/Metamorfoza+ptak%C3%B3w-2020-849024/
vod: 
	nowehoryzonty.pl

----------------------------------------
```
