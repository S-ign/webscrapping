package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func removeDuplicateValues(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func removeSubstringFromSlice(stringSlice []string, subString string) ([]string, error) {
	newSlice := []string{}
	for _, s := range stringSlice {
		if !strings.Contains(s, subString) {
			newSlice = append(newSlice, s)
		}
	}
	return newSlice, nil
}

func getLinks() ([]string, error) {
	c := colly.NewCollector(
		// define the domain name you would like to scrap
		colly.AllowedDomains("jayfeldmanwellness.com"),
	)

	articleLinks := []string{}
	links := []string{}

	for i := 1; i < 4; i++ {
		// define the classes or id you would like search by
		c.OnHTML(".blog_holder", func(e *colly.HTMLElement) {
			// define what you would like to return
			links = e.ChildAttrs("a", "href")
			//fmt.Println(links)
		})
		for _, l := range links {
			if len(l) > 20 {
				articleLinks = append(articleLinks, l)
			}
		}

		page := "https://jayfeldmanwellness.com/articles/page/" + strconv.Itoa(i)
		if i == 1 {
			page = "https://jayfeldmanwellness.com/articles"
		}
		c.Visit(page)
	}

	articleLinks = removeDuplicateValues(articleLinks)

	articleLinks, err := removeSubstringFromSlice(articleLinks, "category")
	check(err)

	articleLinks, err = removeSubstringFromSlice(articleLinks, "articles")
	check(err)

	articleLinks, err = removeSubstringFromSlice(articleLinks, "comments")
	check(err)

	articleLinks, err = removeSubstringFromSlice(articleLinks, "respond")
	check(err)

	return articleLinks, nil
}

func getTags(links []string, selector string) ([]string, error) {
	c := colly.NewCollector(
		// define the domain name you would like to scrap
		colly.AllowedDomains("jayfeldmanwellness.com"),
	)

	tags := []string{}

	for _, l := range links {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			tag, err := e.DOM.Html()
			check(err)

			tags = append(tags, tag)
		})

		c.Visit(l)
	}

	return tags, nil
}

func main() {
	links, err := getLinks()
	check(err)

	tags, err := getTags(links, ".blog_single")
	check(err)

	f, err := os.Create("articles.html")
	check(err)

	defer f.Close()

	f.WriteString(fmt.Sprintf(`
	<html>
		<head></head>
		<title></title>
		<body>
			%s
		</body>
	</html>
	`, strings.Join(tags, "\n")))

	for _, l := range links {
		fmt.Println(l)
	}

}
