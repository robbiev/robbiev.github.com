package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"io/ioutil"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	baseLocation      = "../"
	blogEntryLocation = baseLocation + "blog_entries/"
)

type queryFunc func(*html.Node) *html.Node

func hasClass(class string) queryFunc {
	return func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode {
			if s, ok := findAttr(n, "class"); ok && s == class {
				return n
			}
		}
		return nil
	}
}

func hasType(a atom.Atom) queryFunc {
	return func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.DataAtom == a {
			return n
		}
		return nil
	}
}

func findAttr(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func queryHTML(n *html.Node, f queryFunc) *html.Node {
	if n := f(n); n != nil {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := queryHTML(c, f)
		if result != nil {
			return result
		}
	}

	return nil
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadTemplate() (*html.Node, error) {
	f, err := os.Open("post-template.html")
	exitOnErr(err)
	defer f.Close()
	return html.Parse(f)
}

func main() {
	files, err := ioutil.ReadDir(blogEntryLocation)
	exitOnErr(err)

	for _, f := range files {
		if !f.IsDir() {
			fmt.Println()
			fmt.Println(f.Name())

			// read the blog entry
			p := filepath.Join(blogEntryLocation, f.Name())
			srcf, err := os.Open(p)
			exitOnErr(err)

			scan := bufio.NewScanner(srcf)

			// get the title
			scan.Scan()
			exitOnErr(scan.Err())
			titleText := scan.Text()

			// get the date
			scan.Scan()
			exitOnErr(scan.Err())
			dateText := scan.Text()

			// read the rest of the body
			var buf bytes.Buffer
			for scan.Scan() {
				buf.Write(scan.Bytes())
				buf.WriteByte('\n')
			}

			exitOnErr(scan.Err())
			srcf.Close()

			// read the blog entry body as HTML
			entryHTML, err := html.ParseFragment(&buf, &html.Node{
				Type:     html.ElementNode,
				Data:     "body",
				DataAtom: atom.Body,
			})
			exitOnErr(err)

			// get the blog entry page template
			template, err := loadTemplate()
			exitOnErr(err)

			title := queryHTML(template, hasType(atom.Title))
			heading := queryHTML(template, hasType(atom.H1))
			date := queryHTML(template, hasClass("date"))
			entrye := queryHTML(template, hasClass("entry"))

			// set the blog entry data in the blog entry page template
			title.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: titleText,
			})
			heading.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: titleText,
			})
			date.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: dateText,
			})

			t, err := time.Parse("January 2, 2006", dateText)
			exitOnErr(err)
			fmt.Println(t.Format("02"))
			fmt.Println(t.Format("01"))
			fmt.Println(t.Year())

			for _, eh := range entryHTML {
				entrye.AppendChild(eh)
			}

			// render the resulting blog entry page
			//exitOnErr(html.Render(os.Stdin, template))
		}
	}
}
