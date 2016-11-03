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

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var baseLocation string

type queryFunc func(*html.Node) *html.Node

type postProcFunc func(bytes.Buffer) bytes.Buffer

type indexEntry struct {
	html []*html.Node
	time int64
}

type ByTimeDesc []indexEntry

func (a ByTimeDesc) Len() int           { return len(a) }
func (a ByTimeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimeDesc) Less(i, j int) bool { return a[i].time > a[j].time }

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
	f, err := os.Open(filepath.Join(baseLocation, "generator/post-template.html"))
	exitOnErr(err)
	defer f.Close()
	return html.Parse(f)
}

func fakeBodyNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	}
}

func createIndexEntry(title string, date string, path string) []*html.Node {
	b, err := ioutil.ReadFile(filepath.Join(baseLocation, "generator/index-entry-template.html"))
	exitOnErr(err)
	entryHTML, err := html.ParseFragment(bytes.NewReader(b), fakeBodyNode())
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

func createIndexHTML(entries []indexEntry) *html.Node {
	b, err := ioutil.ReadFile(filepath.Join(baseLocation, "generator/index-template.html"))
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

func generateEntries(location string, indexEntries []indexEntry, postProc postProcFunc) []indexEntry {
	files, err := ioutil.ReadDir(location)
	exitOnErr(err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		// read the blog entry
		var titleText, dateText string
		var bodyText bytes.Buffer
		{
			p := filepath.Join(location, f.Name())
			srcf, err := os.Open(p)
			exitOnErr(err)

			scan := bufio.NewScanner(srcf)

			// get the title
			scan.Scan()
			exitOnErr(scan.Err())
			titleText = scan.Text()

			// get the date
			scan.Scan()
			exitOnErr(scan.Err())
			dateText = scan.Text()

			// read the rest of the body
			for scan.Scan() {
				bodyText.Write(scan.Bytes())
				bodyText.WriteByte('\n')
			}

			exitOnErr(scan.Err())
			srcf.Close()

			bodyText = postProc(bodyText)
		}

		// get the blog entry page template
		template, err := loadTemplate()
		exitOnErr(err)

		// set the blog entry data in the blog entry page template
		{
			title := queryHTML(template, hasType(atom.Title))
			heading := queryHTML(template, hasType(atom.H1))
			date := queryHTML(template, hasClass("date"))
			entrye := queryHTML(template, hasClass("entry"))

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

			// read the blog entry body as HTML
			entryHTML, err := html.ParseFragment(&bodyText, fakeBodyNode())
			exitOnErr(err)

			for _, eh := range entryHTML {
				entrye.AppendChild(eh)
			}
		}

		t, err := time.Parse("January 2, 2006", dateText)
		exitOnErr(err)

		var targetPath, targetDir string
		{
			targetPathStart := t.Format("2006/01/02/")
			fileNameWithoutExt := f.Name()[0 : len(f.Name())-len(filepath.Ext(f.Name()))]
			targetPath = filepath.Join(targetPathStart, fileNameWithoutExt)
			targetDir = filepath.Join(baseLocation, targetPath)
		}

		exitOnErr(os.RemoveAll(targetDir))
		exitOnErr(os.MkdirAll(targetDir, 0755))

		// write to file
		{
			targetFile := filepath.Join(targetDir, "index.html")
			fmt.Printf("generating: %s\n", targetFile)
			target, err := os.Create(targetFile)
			exitOnErr(err)
			exitOnErr(html.Render(target, template))
			target.Close()
		}

		indexEntries = append(indexEntries, indexEntry{
			html: createIndexEntry(titleText, dateText, targetPath+"/"),
			time: t.Unix(),
		})
	}
	return indexEntries
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
	baseLocation = getBlogRoot()
	fmt.Printf("blog root: %s\n", baseLocation)

	blogEntryLocation := filepath.Join(baseLocation, "blog_entries")
	blogEntryLocationMD := filepath.Join(baseLocation, "blog_entries_md")

	var indexEntries []indexEntry

	indexEntries = generateEntries(blogEntryLocation, indexEntries, func(b bytes.Buffer) bytes.Buffer {
		return b
	})

	indexEntries = generateEntries(blogEntryLocationMD, indexEntries, func(b bytes.Buffer) bytes.Buffer {
		commonHtmlFlags := 0
		commonExtensions := 0 |
			blackfriday.EXTENSION_TABLES |
			blackfriday.EXTENSION_FENCED_CODE |
			blackfriday.EXTENSION_AUTOLINK |
			blackfriday.EXTENSION_STRIKETHROUGH |
			blackfriday.EXTENSION_DEFINITION_LISTS
		renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
		out := blackfriday.MarkdownOptions(b.Bytes(), renderer, blackfriday.Options{
			Extensions: commonExtensions,
		})
		return *bytes.NewBuffer(out)
	})

	sort.Sort(ByTimeDesc(indexEntries))

	targetFile := filepath.Join(baseLocation, "index.html")
	target, err := os.Create(targetFile)
	fmt.Printf("generating: %s\n", targetFile)
	exitOnErr(err)
	exitOnErr(html.Render(target, createIndexHTML(indexEntries)))
	target.Close()
}
