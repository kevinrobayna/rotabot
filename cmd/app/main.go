package main

import (
	"fmt"
)

var (
	AppName = "rotabot"
	Version = "unknown"
	Date    = "unknown"
)

func hello() string {
	return fmt.Sprintf("Hello world from %s running %s built on %s", AppName, Version, Date)
}

func main() {
	fmt.Println(hello())
}
