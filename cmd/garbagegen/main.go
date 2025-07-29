package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "embed"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

//go:embed post-template.html
var postTemplate []byte

//go:embed index-entry-template.html
var indexEntryTemplate []byte

//go:embed index-template.html
var indexTemplate []byte

var baseLocation string

type queryFunc func(*html.Node) *html.Node

type indexEntry struct {
	html []*html.Node
	time int64
}

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

func fakeBodyNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	}
}

func createIndexEntry(title string, date string, path string) []*html.Node {
	entryHTML, err := html.ParseFragment(bytes.NewReader(indexEntryTemplate), fakeBodyNode())
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
	entryHTML, err := html.Parse(bytes.NewReader(indexTemplate))
	exitOnErr(err)

	home := queryHTML(entryHTML, hasClass("home"))
	for _, v := range entries {
		for _, html := range v.html {
			home.AppendChild(html)
		}
	}

	return entryHTML
}

func generateEntries(location string, indexEntries []indexEntry, postProc func(bytes.Buffer) bytes.Buffer) []indexEntry {
	files, err := os.ReadDir(location)
	exitOnErr(err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		// read the blog entry
		var titleText, dateText string
		var bodyText bytes.Buffer
		includeInIndex := true
		{
			p := filepath.Join(location, f.Name())
			srcf, err := os.Open(p)
			exitOnErr(err)

			scan := bufio.NewScanner(srcf)

			// get the title
			scan.Scan()
			exitOnErr(scan.Err())
			includeInIndex = !strings.HasPrefix(scan.Text(), "-")
			titleText = strings.TrimPrefix(scan.Text(), "-")

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
		template, err := html.Parse(bytes.NewReader(postTemplate))
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

		if !includeInIndex {
			continue
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

type autoLinkClassAdder struct{}

func (r *autoLinkClassAdder) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindAutoLink, r.renderLink)
}

func (r *autoLinkClassAdder) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		w.WriteString("</a>")
		return ast.WalkContinue, nil
	}
	n := node.(*ast.AutoLink)
	if n.AutoLinkType == ast.AutoLinkEmail {
		return ast.WalkContinue, nil
	}
	u := n.URL(source)
	w.WriteString(`<a class="url" href="`)
	w.Write(util.EscapeHTML(u))
	w.WriteString(`">`)
	w.Write(util.EscapeHTML(u))
	return ast.WalkContinue, nil
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
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithRendererOptions(
				goldmarkhtml.WithUnsafe(),
				renderer.WithNodeRenderers(
					util.Prioritized(&autoLinkClassAdder{}, 100),
				),
			),
		)
		var buf bytes.Buffer
		if err := md.Convert(b.Bytes(), &buf); err != nil {
			panic(err)
		}
		return buf
	})

	sort.Slice(indexEntries, func(i, j int) bool {
		return indexEntries[i].time > indexEntries[j].time
	})

	targetFile := filepath.Join(baseLocation, "index.html")
	target, err := os.Create(targetFile)
	fmt.Printf("generating: %s\n", targetFile)
	exitOnErr(err)
	exitOnErr(html.Render(target, createIndexHTML(indexEntries)))
	target.Close()
}
