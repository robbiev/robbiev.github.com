package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

type ByTimeDesc []indexEntry

func (a ByTimeDesc) Len() int           { return len(a) }
func (a ByTimeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimeDesc) Less(i, j int) bool { return a[i].time > a[j].time }

type indexEntry struct {
	html []*html.Node
	time int64
}

func createIndexEntry(title string, date string, path string) []*html.Node {
	b, err := ioutil.ReadFile("index-entry-template.html")
	exitOnErr(err)
	entryHTML, err := html.ParseFragment(bytes.NewReader(b), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	exitOnErr(err)

	for _, entry := range entryHTML {
		titleNode := queryHTML(entry, hasType(atom.A))
		if titleNode != nil {
			titleNode.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: title,
			})

			var found bool
			for i, attr := range titleNode.Attr {
				if attr.Key == "href" {
					fmt.Println(attr)
					attr.Val = path
					titleNode.Attr[i] = attr
					found = true
					break
				}
			}
			if !found {
				titleNode.Attr = append(titleNode.Attr, html.Attribute{
					Key: "href",
					Val: path,
				})
			}
		}
		dateNode := queryHTML(entry, hasClass("date"))
		if dateNode != nil {
			dateNode.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: date,
			})
		}
	}
	return entryHTML
}

func indexHTML(entries []indexEntry) *html.Node {
	b, err := ioutil.ReadFile("index-template.html")
	exitOnErr(err)
	entryHTML, err := html.Parse(bytes.NewReader(b))
	exitOnErr(err)

	home := queryHTML(entryHTML, hasClass("home"))
	for _, v := range entries {
		for _, html := range v.html {
			home.AppendChild(html)
		}
	}

	return entryHTML
}

func main() {
	files, err := ioutil.ReadDir(blogEntryLocation)
	exitOnErr(err)

	var indexEntries []indexEntry

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
			for _, eh := range entryHTML {
				entrye.AppendChild(eh)
			}

			t, err := time.Parse("January 2, 2006", dateText)
			exitOnErr(err)

			targetDirStart := t.Format("2006/01/02/")
			name := f.Name()[0 : len(f.Name())-len(filepath.Ext(f.Name()))]
			targetDirLoc := filepath.Join(targetDirStart, name)
			targetDir := filepath.Join(baseLocation, targetDirLoc)

			exitOnErr(os.RemoveAll(targetDir))
			exitOnErr(os.MkdirAll(targetDir, 0755))

			targetFile, err := os.Create(filepath.Join(targetDir, "index.html"))
			exitOnErr(err)
			exitOnErr(html.Render(targetFile, template))
			targetFile.Close()

			indexEntries = append(indexEntries, indexEntry{
				html: createIndexEntry(titleText, dateText, targetDirLoc+"/"),
				time: t.Unix(),
			})

			// render the resulting blog entry page
			//exitOnErr(html.Render(os.Stdin, template))
		}
	}

	sort.Sort(ByTimeDesc(indexEntries))

	targetFile, err := os.Create(filepath.Join(baseLocation, "index.html"))
	exitOnErr(err)
	exitOnErr(html.Render(targetFile, indexHTML(indexEntries)))
	targetFile.Close()
}
