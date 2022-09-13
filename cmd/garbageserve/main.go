package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

func main() {
	baseLocation := getBlogRoot()
	fmt.Printf("blog root: %s\n", baseLocation)
	go func() {
		time.Sleep(time.Second)
		cmd := exec.Command("open", "http://localhost:8080")
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(baseLocation))))
}
