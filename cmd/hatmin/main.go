package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/amlweems/atmin"
)

func main() {
	reqFlag := flag.String("request", "", "path to file containing initial request")
	addrFlag := flag.String("addr", "example.org:443", "server to send network requests to")
	tlsFlag := flag.Bool("tls", true, "enable or disable TLS")
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

	if *dryRunFlag {
		ex := atmin.HTTPExecutor{Addr: *addrFlag, TLS: *tlsFlag}
		out := ex.Execute(in)
		os.Stdout.Write(out)
		return
	}

	m := atmin.NewMinimizer(in).ExecuteHTTP(*addrFlag, *tlsFlag).ValidateString(*needleFlag)
	min := m.Minimize()

	os.Stdout.Write(min)
	os.Stdout.Write([]byte("\n"))
}
