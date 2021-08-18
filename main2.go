package main

import (
	"fmt"
	"os"

	"github.com/rm46627/fws/scraper"
)

func main() {
	// check for args
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error reading username\n")
		os.Exit(1)
	}

	address := "https://www.filmweb.pl/user/" + os.Args[1] + "/wantToSee?page=1"
	movieSlice := scraper.DataSlice()
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
}
