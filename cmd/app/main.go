package main

import "fmt"

var (
	AppName = "rotabot"
	Version = "unknown"
	Date    = "unknown"
)

func main() {
	fmt.Printf("Hello world from %s running %s built on %s\n", AppName, Version, Date)
}
