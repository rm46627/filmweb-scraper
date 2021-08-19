package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "want",
	Long: `Want is a data scraper for https://www.filmweb.pl/. 
	Visits movies pages from \"want to see\" page of given user, 
	creates txt file with every movie that is available on any vod and 
	lists those vod sites.
	Usage:
	want tosee [username]`,
}
