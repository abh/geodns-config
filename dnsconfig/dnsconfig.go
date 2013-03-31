package main

import (
	"flag"
	"fmt"
	"github.com/devel/dnsconfig"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var VERSION string = "2.1.0"
var buildTime string
var gitVersion string

var (
	zonesFile       = flag.String("config", "config/zones.json", "zones.json configuration file")
	outputDir       = flag.String("output", "dns", "output directory")
	showVersionFlag = flag.Bool("version", false, "Show dnsconfig version")
	Verbose         = flag.Bool("verbose", false, "verbose output")
)

func init() {
	if len(gitVersion) > 0 {
		VERSION = VERSION + "/" + gitVersion
	}

	log.SetPrefix("geodns ")
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

func BuildAll(zones *dnsconfig.Zones) {

	for _, zone := range zones.All() {

		log.Printf("Building %s\n", zone.Name)
		err := zone.LoadConfig()
		zone.Verbose = *Verbose
		if err != nil {
			log.Printf("Could not load configuration for '%s': %s", zone.Name, err)
			continue
		}

		js, err := zone.BuildJSON()
		if err != nil {
			log.Printf("Could not build DNS data for '%s': %s", zone.Name, err)
			continue
		}

		fileName := filepath.Join(*outputDir, zone.Name+".json")
		err = ioutil.WriteFile(fileName, []byte(js), 0644)
		if err != nil {
			log.Printf("Could not write '%s' to '%s': %s", zone.Name, fileName, err)
			continue
		}
	}
}

func main() {

	flag.Parse()

	if *showVersionFlag {
		fmt.Println("dnsconfig", VERSION, buildTime)
		os.Exit(0)
	}

	zones := new(dnsconfig.Zones)

	err := zones.LoadZonesConfig(*zonesFile)
	if err != nil {
		log.Printf("Could not open '%s': %s", *zonesFile, err)
		os.Exit(2)
	}

	BuildAll(zones)

}
