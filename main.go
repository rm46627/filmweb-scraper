package main

import (
	"github.com/rm46627/want/cmd"
)

const (
	s = `Want is a data scraper for https://www.filmweb.pl/. 
Visits movies pages from \"want to see\" page of given user, 
creates txt file with every movie that is available on any vod and 
lists those vod sites.

Usage:
want tosee [username]
`
)

func main() {
	cmd.RootCmd.SetHelpTemplate(s)
	cmd.RootCmd.Execute()
}
