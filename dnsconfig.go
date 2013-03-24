package main

import (
	"log"
)

var VERSION string = "2.0.0"
var gitVersion string

func main() {

}

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}

	log.SetPrefix("geodns ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}
