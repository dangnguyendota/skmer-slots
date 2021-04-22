package main

import (
	"./chinese_carnival"
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
		classic.Start()
	} else if *cc {
		chinese_carnival.Start()
	} else if *fb {
		football.Start()
	}

}
