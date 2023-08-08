package main

import (
	"fmt"
	"time"
)

var crawler WebCrawler

// main will start the script
func main() {
	crawler = NewWebCrawler()
	fmt.Println("Hello Parser!")
	fmt.Printf("Started at: %s\n", time.Now())
	RunCrawler()
	fmt.Printf("Finished at: %s\n", time.Now())
}
