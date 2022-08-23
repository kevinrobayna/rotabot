package main

import "fmt"

var (
	AppName = "rotabot"
	Version = "unknown"
	Date    = "unknown"
)

func main() {
	fmt.Println(fmt.Sprintf("Hello world from %s running %s built on %s\n", AppName, Version, Date))
}
