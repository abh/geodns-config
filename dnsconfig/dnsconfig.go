package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/abh/geodns-config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var VERSION string = "2.2.1"
var buildTime string
var gitVersion string

var (
	zonesFile       = flag.String("config", "config/zones.json", "zones.json configuration file")
	outputDir       = flag.String("output", "dns", "output directory")
	showVersionFlag = flag.Bool("version", false, "Show dnsconfig version")
	httpPort        = flag.Int("httpport", 0, "Enable HTTP interface on port")
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

	if *httpPort > 0 {
		setupDaemon(*httpPort)
	} else {
		runOnce()
	}
}

func runOnce() {
	zones := new(dnsconfig.Zones)

	err := zones.LoadZonesConfig(*zonesFile)
	if err != nil {
		log.Printf("Could not open '%s': %s", *zonesFile, err)
		os.Exit(2)
	}

	BuildAll(zones)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templateFile, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return

	}
	w.Write(templateFile)
}

func RestTest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]int{"foo": 123, "bar": 456}); err != nil {
		log.Printf("encoding response: %s", err)
	}
}

func setupDaemon(port int) {

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api/test", RestTest).Methods("GET")
	router.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CombinedLoggingHandler(os.Stdout, router))

	if err != nil {
		log.Fatalln(err)
	}

}
