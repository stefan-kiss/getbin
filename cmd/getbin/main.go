package main

import (
	"flag"
	"github.com/cavaliercoder/grab"
	"github.com/k0kubun/pp"
	"log"
)

func main() {

	url := flag.String("url", "http://www.golang-book.com/public/pdf/gobook.pdf", "url of object to download")
	resp, err := grab.Get(".", *url)
	flag.Parse()

	if err != nil {
		log.Fatal(err)
	}

	pp.Print(resp)

}
