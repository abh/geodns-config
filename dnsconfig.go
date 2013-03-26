package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var VERSION string = "2.0.0"
var gitVersion string

var (
	zonesFile = flag.String("config", "config/zones.json", "zones.json configuration file")
)

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}

	log.SetPrefix("geodns ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func main() {

	zones := new(Zones)

	err := zones.LoadZonesConfig(*zonesFile)
	if err != nil {
		log.Printf("Could not open '%s': %s", *zonesFile, err)
		os.Exit(2)
	}

	for _, zone := range zones.All() {

		log.Printf("Building %s\n", zone.Name)
		err = zone.LoadConfig()
		if err != nil {
			log.Printf("Could not load configuration for '%s': %s", zone.Name, err)
			continue
		}

		js, err := zone.BuildJSON()
		if err != nil {
			log.Printf("Could not build DNS data for '%s': %s", zone.Name, err)
			continue
		}

		fmt.Print(js)

	}

}
