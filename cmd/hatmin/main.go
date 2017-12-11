package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/amlweems/atmin"
)

func main() {
	reqFlag := flag.String("request", "", "path to file containing initial request")
	urlFlag := flag.String("url", "https://example.org", "http server to send requests to")
	needleFlag := flag.String("needle", "", "string to search for in responses which indicates a valid response")
	dryRunFlag := flag.Bool("dry-run", false, "make a single request and return the response")
	flag.Parse()

	if *reqFlag == "" || (!*dryRunFlag && *needleFlag == "") {
		log.Printf("-request and -needle required")
	}

	in, err := ioutil.ReadFile(*reqFlag)
	if err != nil {
		log.Fatal(err)
	}

	url, err := url.Parse(*urlFlag)
	if *dryRunFlag {
		ex := atmin.HTTPExecutor{URL: url}
		out := ex.Execute(in)
		os.Stdout.Write(out)
		return
	}

	m := atmin.NewMinimizer(in).ExecuteHTTP(url).ValidateString(*needleFlag)
	min := m.Minimize()

	os.Stdout.Write(min)
	os.Stdout.Write([]byte("\n"))
}
