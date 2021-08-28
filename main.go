package main

import (
	"./carnival"
	"./classic"
	"./football"
	"flag"
)

var cs = flag.Bool("classic", false, "")
var cc = flag.Bool("chinese", false, "")
var fb = flag.Bool("football", false, "")

func main() {
	flag.Parse()
	if *cs {
		classic.Gen()
	} else if *cc {
		carnival.Gen()
	} else if *fb {
		football.Gen()
	} else {
		carnival.Gen()
	}
}
