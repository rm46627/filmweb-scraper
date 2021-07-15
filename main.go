package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

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
		address := baseURL + linkToVodPage
		linksOnVodPage, err := getLinksFromHTML(address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot collect links from HTML: %v\n", err)
		}
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
		movie := Movie{Address: address[:len(address)-3], Vod: vodNames}
		(*movieSlice) = append((*movieSlice), movie)

		fmt.Printf("progress: %d/%d\n", progress, done)
	}
}

// find links in the given address
func getLinksFromHTML(address string) ([]string, error) {
	resp, err := http.Get(address)
	if err != nil {
		return nil, fmt.Errorf("there were too many redirects or if there was an HTTP protocol error")
	}
	var links []string
	z := html.NewTokenizer(resp.Body)
	var addressSuffix bool = address[len(address)-3:] == "vod"

	for {
		// for tokenization error
		tokenType := z.Next()

		switch tokenType {
		case html.ErrorToken:
			return removeDuplicateValues(links), nil
		case html.StartTagToken, html.EndTagToken:
			// if address ends with "vod"
			// need to collect linsk that dont have the "title" attribute
			if !addressSuffix {
				token := z.Token()
				if token.Data == "a" {
					var value string
					var titleAttr bool = false
					for _, attr := range token.Attr {
						if attr.Key == "href" {
							value = attr.Val
						} else if attr.Key == "title" {
							titleAttr = true
						}
					}
					if !titleAttr {
						links = append(links, value)
					}
				}
				// collects all the links
			} else {
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

func writeToFile(movieSlice *[]Movie, vodCounter *map[string]int) error {

	filename := os.Args[1] + ".txt"
	// check for old file
	if _, exist := os.Stat(filename); exist == nil {
		os.Remove(filename)
	}

	if len(*movieSlice) == 0 {
		return fmt.Errorf("didnt catch any links, chceck provided username")
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("openning the %s file: %v", filename, err)
	}

	datawriter := bufio.NewWriter(file)

	t, err := template.New("T1").Parse(template1)
	if err != nil {
		return fmt.Errorf("parsing T1 template body: %v", err)
	}

	err = t.Execute(datawriter, *vodCounter)
	if err != nil {
		return fmt.Errorf("applying T1 template to data obj: %v", err)
	}

	for _, arg := range *movieSlice {
		removeDuplicateValues(arg.Vod)
		t, err := template.New("T2").Parse(template2)
		if err != nil {
			return fmt.Errorf("parsing T2 template body for %s: %v", arg.Address, err)
		}

		err = t.Execute(datawriter, arg)
		if err != nil {
			return fmt.Errorf("applying T2 template to data obj for %s: %v", arg.Address, err)
		}
	}

	datawriter.Flush()
	file.Close()
	return nil
}

type Movie struct {
	Address string
	Vod     []string
}

const (
	baseURL = "https://www.filmweb.pl"

	template1 = `
	vod counter:
	
	{{range $Name, $Value := . -}}
	{{$Name}}: {{$Value}}
	{{end}}`

	template2 = `
	address: {{.Address}}
	vod: 
	{{range $Name := .Vod}}	{{$Name}}
	{{end}}
	----------------------------------------`
)

func main() {
	// check for args
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error reading username\n")
		os.Exit(1)
	}

	address := "https://www.filmweb.pl/user/" + os.Args[1] + "/wantToSee?page=1"
	var movieSlice []Movie
	vodCounter := make(map[string]int)

	// getting links from want to see page of specyfic user
	links, err := getLinksFromHTML(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// filters links that contain "/film" and dont contain "/ranking/"
	filteredLinks := filterStrings_Contains(links, "/film/")
	filteredLinks = filterStrings_DoesNotContain(filteredLinks, "/ranking/")
	addSuffixToSliceOfString(&filteredLinks, "/vod")

	// writing data to slice of Movie struct and vodCounter
	gettingThatVod(filteredLinks, &movieSlice, &vodCounter)

	// writing everything to file
	writeToFile(&movieSlice, &vodCounter)

	fmt.Println("Success!")
}
