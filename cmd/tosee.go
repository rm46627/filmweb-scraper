package cmd

import (
	"fmt"
	"os"

	"github.com/rm46627/want/scraper"
	"github.com/spf13/cobra"
)

var tosee = &cobra.Command{
	Use:   "tosee",
	Short: "Command to provide the username.",
	Run: func(cmd *cobra.Command, args []string) {
		// check for args
		if len(args) > 1 || len(args) == 0 {
			fmt.Fprintf(os.Stderr, "error reading username\n")
			os.Exit(1)
		}

		address := "https://www.filmweb.pl/user/" + args[0] + "/wantToSee?page=1"
		movieSlice := scraper.MovieSlice()
		vodCounter := make(map[string]int)

		// getting links from want to see page of specific user
		links, err := scraper.GetLinksFromHTML(address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		// filters links that contain "/film" and dont contain "/ranking/"
		filteredLinks := scraper.FilterStringsContains(links, "/film/")
		filteredLinks = scraper.FilterStringsDoesNotContain(filteredLinks, "/ranking/")
		scraper.AddSuffixToSliceOfString(&filteredLinks, "/vod")

		// writing data to slice of Movie struct and vodCounter
		scraper.GettingThatVod(filteredLinks, &movieSlice, &vodCounter)

		// writing everything to file
		scraper.WriteToFile(&movieSlice, &vodCounter)

		fmt.Println("Success!")

	},
}

func init() {
	RootCmd.AddCommand(tosee)
}
