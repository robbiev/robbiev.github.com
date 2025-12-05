package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

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

func slugify(title string) string {
	lower := strings.ToLower(title)
	lowerHyphenated := strings.ReplaceAll(lower, " ", "-")
	// allow list letters, digits, and hyphens
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		return -1
	}, lowerHyphenated)
}

func main() {
	baseLocation := getBlogRoot()
	fmt.Printf("blog root: %s\n", baseLocation)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: garbagenew <title>")
		os.Exit(1)
	}
	title := strings.Join(os.Args[1:], " ")

	// Create a new file in blog_entries/title-slug.md
	file, err := os.Create(filepath.Join(baseLocation, "blog_entries", fmt.Sprintf("%s.md", slugify(title))))
	exitOnErr(err)
	defer file.Close()

	// Every file should start with something like this:
	// How I Set Up Xv6 On macOS
	// August 1, 2024
	// <blank line>
	file.WriteString(fmt.Sprintf("%s\n", title))
	file.WriteString(time.Now().Format("January 2, 2006\n"))
	file.WriteString("\n")
}
