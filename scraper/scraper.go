package scraper

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

// FilterStringsContains return slice of strings which contain given filters.
func FilterStringsContains(strSlice []string, canContain ...string) []string {
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

// FilterStringsDoesNotContain return slice of strings which does not contain given filters.
func FilterStringsDoesNotContain(strSlice []string, cannotContain string) []string {
	var out []string
	for _, str := range strSlice {
		if strings.Contains(str, cannotContain) {
			continue
		}
		out = append(out, str)
	}
	return out
}

// AddSuffixToSliceOfString adds suffix to every string from given slice.
func AddSuffixToSliceOfString(slice *[]string, suffix string) {
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

// GettingThatVod search for vod addresses in given slice of links, write data to []Movie and counts occurrence of each site.
func GettingThatVod(links []string, movieSlice *[]Movie, vodCounter *map[string]int) {
	done := len(links)
	for progress, linkToVodPage := range links {
		address := baseURL + linkToVodPage
		linksOnVodPage, err := GetLinksFromHTML(address)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot collect links from HTML: %v\n", err)
		}
		filteredLinks := FilterStringsContains(linksOnVodPage,
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

// GetLinksFromHTML find links in the given address.
// Returns slice of string and error.
func GetLinksFromHTML(address string) ([]string, error) {
	resp, err := http.Get(address)
	if err != nil {
		return nil, fmt.Errorf("there were too many redirects or if there was an HTTP protocol error")
	}
	var links []string
	z := html.NewTokenizer(resp.Body)
	var addressSuffix bool
	addressSuffix = address[len(address)-3:] == "vod"

	for {
		// for tokenization error
		tokenType := z.Next()

		switch tokenType {
		case html.ErrorToken:
			return removeDuplicateValues(links), nil
		case html.StartTagToken, html.EndTagToken:
			// if address ends with "vod"
			// need to collect links that dont have the "title" attribute
			if !addressSuffix {
				token := z.Token()
				if token.Data == "a" {
					var value string
					var titleAttr bool
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

// WriteToFile parses data from slice of Movie struct to txt file
func WriteToFile(movieSlice *[]Movie, vodCounter *map[string]int) error {

	filename := os.Args[1] + ".txt"
	// check for old file
	if _, exist := os.Stat(filename); exist == nil {
		os.Remove(filename)
	}

	if len(*movieSlice) == 0 {
		return fmt.Errorf("didnt catch any links, check provided username")
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening the %s file: %v", filename, err)
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

// Movie struct contains address to the filmweb page of the movie and
// slice of the names of the websites where movie can be viewed.
type Movie struct {
	Address string
	Vod     []string
}

// MovieSlice func return slice of Movie struct
func MovieSlice() []Movie {
	var slice []Movie
	return slice
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
