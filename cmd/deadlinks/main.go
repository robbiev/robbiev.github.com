package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	numWorkers = 5
)

var (
	workerCh = make(chan job, 5*numWorkers)
	seen     = make(map[string]struct{})
	client   = http.Client{
		Timeout: 10 * time.Second,
	}
	wg         sync.WaitGroup
	printMutex sync.Mutex
)

type job struct {
	sourceLocation string
	targetLocation string
}

func clearLine() {
	fmt.Print("\033[2K\r")
}

func checkLink(j job) {
	if _, ok := seen[j.targetLocation]; ok {
		return
	}
	seen[j.targetLocation] = struct{}{}
	wg.Add(1)
	workerCh <- j
}

// local links need to be interpreted to exist on the file system
func indexify(href string) string {
	if href == "" {
		return "index.html"
	}

	if href[len(href)-1:] == "/" {
		return href + "index.html"
	}

	return href
}

func findLinks(base string, j job, n *html.Node) {
	print := func(nextJob job) job {
		printMutex.Lock()
		clearLine()
		fmt.Printf("%s: %s", nextJob.sourceLocation, nextJob.targetLocation)
		printMutex.Unlock()
		return nextJob
	}

	if n.Type == html.ElementNode && n.Data == atom.A.String() {
		var href string
		for _, v := range n.Attr {
			if v.Key == atom.Href.String() {
				href = v.Val
			}
		}

		if strings.HasPrefix(href, "http") {
			checkLink(print(job{
				targetLocation: href,
				sourceLocation: j.targetLocation,
			}))
		} else {
			findLinksInFile(base, print(job{
				targetLocation: indexify(href),
				sourceLocation: j.targetLocation,
			}))
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findLinks(base, j, c)
	}
}

func findLinksInFile(base string, j job) {
	if _, ok := seen[j.targetLocation]; ok {
		return
	}
	seen[j.targetLocation] = struct{}{}

	r, err := os.Open(filepath.Join(base, j.targetLocation))
	exitOnErr(err)

	doc, err := html.Parse(r)
	r.Close()
	exitOnErr(err)

	findLinks(base, j, doc)
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
	fmt.Println(strings.Repeat("=", 10))
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

func linkChecker() {
	print := func(f func()) {
		printMutex.Lock()
		clearLine()
		f()
		fmt.Println(strings.Repeat("=", 10))
		printMutex.Unlock()
	}

	for job := range workerCh {
		resp, err := client.Get(job.targetLocation)

		if err != nil {
			print(func() {
				fmt.Printf("FAILED: %s\nfrom: %s\nto: %s\n", err, job.sourceLocation, job.targetLocation)
			})
		} else {
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				print(func() {
					fmt.Printf("FAILED: %d\nfrom: %s\nto: %s\n", resp.StatusCode, job.sourceLocation, job.targetLocation)
				})
			}
		}

		wg.Done()
	}
}

func main() {
	// toggle line wrapping so clearing lines is easier
	fmt.Print("\033[?7l")
	defer fmt.Print("\033[?7h")

	for i := 0; i < numWorkers; i++ {
		go linkChecker()
	}

	findLinksInFile(getBlogRoot(), job{
		targetLocation: "index.html",
		sourceLocation: "/", // whatever
	})

	wg.Wait()
	close(workerCh)

	fmt.Println()
}
