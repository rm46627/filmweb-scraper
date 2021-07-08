package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func filterStrings_Contains(strSlice []string, canContain ...string) []string {
	var out []string
	for _, str := range strSlice {
		for _, cmp := range canContain {
			if strings.Contains(str, cmp) {
				out = append(out, str)
			}
		}
	}
	return out
}

func filterStrings_DoesNotContain(strSlice []string, cannotContain string) []string {
	var out []string
	for _, str := range strSlice {
		if strings.Contains(str, cannotContain) {
			continue
		}
		out = append(out, str)
	}
	return out
}

func addSuffixToSliceOfString(slice *[]string, suffix string) {
	for i, arg := range *slice {
		if !strings.HasSuffix(arg, suffix) {
			(*slice)[i] = (*slice)[i] + suffix
		}
	}
}

func removeDuplicateValues(strSlice []string) []string {
	strings := make(map[string]bool)
	var out []string
	for _, arg := range strSlice {
		if !strings[arg] {
			out = append(out, arg)
			strings[arg] = true
		}
	}
	return out
}

func gettingThatVod(links []string, movieSlice *[]Movie, vodCounter *map[string]int) {
	done := len(links)
	for progress, linkToVodPage := range links {
		adres := baseURL + linkToVodPage
		linksOnVodPage := getLinksFromHTML(adres)
		filteredLinks := filterStrings_Contains(linksOnVodPage,
			"itunes.apple.com",
			"netflix.com",
			"hbogo.pl",
			"www.primevideo.com",
			"canalplus.com",
			"player.pl",
			"playnow.pl",
			"chili.com",
			"vod.tvp.pl",
			"nowehoryzonty.pl",
			"vod.mdag.pl",
			"piecsmakow.pl",
			"mojeekino.pl",
			"ninateka.pl",
			"vod.pl",
			"dafilms.pl",
			"cineman.pl",
			"outfilm.pl",
			"kinopodbaranami.pl")
		vodNames := linksToVodName(vodCounter, filteredLinks,
			"apple",
			"netflix",
			"hbo",
			"primevideo",
			"canalplus",
			"player.pl",
			"playnow.pl",
			"chili.com",
			"vod.tvp.pl",
			"nowehoryzonty.pl",
			"vod.mdag.pl",
			"piecsmakow.pl",
			"mojeekino.pl",
			"ninateka.pl",
			"vod.pl",
			"dafilms.pl",
			"cineman.pl",
			"outfilm.pl",
			"kinopodbaranami.pl")

		if len(vodNames) < 1 {
			continue
		}
		movie := Movie{Adres: adres, Vod: vodNames}
		(*movieSlice) = append((*movieSlice), movie)

		fmt.Printf("progress: %d/%d\n", progress, done)
	}
}

func getLinksFromHTML(adres string) []string {
	resp, err := http.Get(adres)
	if err != nil {
		log.Fatal(err)
	}

	var links []string
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return removeDuplicateValues(links)
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}

		}
	}
}

func linksToVodName(vodCounter *map[string]int, filteredLinks []string, names ...string) []string {
	var out []string
	for _, link := range filteredLinks {
		for _, name := range names {
			if strings.Contains(link, name) {
				(*vodCounter)[name]++
				out = append(out, name)
			}
		}
	}
	return out
}

type Movie struct {
	Adres string
	Vod   []string
}

const (
	baseURL = "https://www.filmweb.pl"
)

func main() {

	adres := "https://www.filmweb.pl/user/" + os.Args[1] + "/wantToSee?page=1"
	var movieSlice []Movie
	vodCounter := make(map[string]int)

	// getting links from want to see page of specyfic user
	links := getLinksFromHTML(adres)

	// filters links that contain "/film" and dont contain "/ranking/"
	filteredLinks := filterStrings_Contains(links, "/film/")
	filteredLinks = filterStrings_DoesNotContain(filteredLinks, "/ranking/")
	addSuffixToSliceOfString(&filteredLinks, "/vod")

	// writing data to slice of Movie struct and vodCounter
	gettingThatVod(filteredLinks, &movieSlice, &vodCounter)

	// writing everything to data []string
	var data []string
	for _, obj := range movieSlice {
		dataInfo := "\nadres: " + obj.Adres + "\n" + "vod:"
		data = append(data, dataInfo)
		for _, vod := range obj.Vod {
			str := "\t" + vod
			data = append(data, str)
		}
		data = append(data, "===========")
	}

	data = append(data, "\nvod counter: \n")

	for key, value := range vodCounter {
		str := key + ": " + strconv.Itoa(value)
		data = append(data, str)
	}

	// writting data []string to file
	filename := os.Args[1] + ".txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range data {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()

	fmt.Println("writting done")
}
