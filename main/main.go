package main

import (
    "flag"
    "log"
    "os"

    "github.com/mycrawler/crawler"
)

var url string
var selector string
var sel string
var dir string
var saveHTML bool
var sleep int

func init() {
    flag.StringVar(&url, "u", "", "target url")
    flag.StringVar(&selector, "selector", "", "selector")
    flag.StringVar(&sel, "sel", "document.querySelector(\"body\")", `sel.`)
    flag.StringVar(&dir, "d", "./", "dir to place data.")
    flag.BoolVar(&saveHTML, "s", false, "save html to file.")
    flag.IntVar(&sleep, "sleep", 3, "time wait before save html content.")
}

func main() {

    flag.Parse()

    // url = `https://www.veer.com/search-image/fengjing/`
    // dir = `remote`
    // sleep := 3

    checkParaValid()

    crawler.CrawlPicFromUrl(url, selector, dir, sel, sleep, saveHTML)
}

func checkParaValid() {
    if url == "" {
        log.Fatalf("url is empty.")
    }
    checkDirectoryValid(dir)
}

func checkDirectoryValid(dir string) {
    _, err := os.Stat(dir)
    if err != nil {
        err = os.Mkdir(dir, os.ModePerm)
        if err != nil {
            log.Fatalf("failed to create dir %s,exit", dir)
        }
    }
}
