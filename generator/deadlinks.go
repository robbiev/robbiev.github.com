package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var seen map[string]struct{}
var client http.Client

func findLinks(base string, location string, parent string, isHTTP bool) {
	if _, ok := seen[location]; ok {
		return
	}
	seen[location] = struct{}{}

	var r io.ReadCloser
	if isHTTP {
		resp, err := client.Get(location)
		if err != nil {
			fmt.Print("\033[2K\r")
			fmt.Printf("FAILED: %s\nfrom: %s\nto: %s\n", err, parent, location)
			// os.Exit(1)
			return
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Print("\033[2K\r")
			fmt.Printf("FAILED: %d\nfrom: %s\nto: %s\n", resp.StatusCode, parent, location)
			// os.Exit(1)
			return
		}
		return
	} else {
		var err error
		r, err = os.Open(filepath.Join(base, location))
		exitOnErr(err)
	}

	doc, err := html.Parse(r)
	r.Close()
	exitOnErr(err)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string
			for _, v := range n.Attr {
				if v.Key == "href" {
					href = v.Val
				}
			}

			isHTTP := strings.HasPrefix(href, "http")

			if href == "" {
				href = "index.html"
			} else if !isHTTP && href[len(href)-1:] == "/" {
				href = href + "index.html"
			}
			fmt.Print("\033[2K\r")
			fmt.Printf("%s: %s", location, href)
			findLinks(base, href, location, isHTTP)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getBlogRoot() string {
	dir, err := os.Getwd()
	exitOnErr(err)
	fmt.Printf("working directory: %s\n", dir)
	for {
		if _, err := os.Stat(filepath.Join(dir, "blog_entries")); err == nil {
			break
		}
		dir = filepath.Dir(dir)
		if dir == "/" || dir == "." {
			fmt.Fprintln(os.Stderr, "can't find blog root")
			os.Exit(1)
		}
	}
	return dir
}

func main() {
	blogRoot := getBlogRoot()
	fmt.Println(blogRoot)
	seen = make(map[string]struct{})
	client = http.Client{
		Timeout: 5 * time.Second,
	}
	findLinks(blogRoot, "index.html", "/", false)
	fmt.Println()
}
