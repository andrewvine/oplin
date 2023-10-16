package main

import (
	"oplin/internal/env"
	"oplin/internal/lineage/wiring"
	"flag"
	"fmt"
	"log"
	"strconv"
)

var webPort int;

func init() {
	flag.IntVar(&webPort, "web_port", 8080, "the port the webserver listens on")
}

func main() {
	flag.Parse()
	env.Setup()
	r := wiring.NewGinEngine()
	err := wiring.SetupLineage(r)
	if err != nil {
		log.Fatal(err)
	}

	portString := strconv.Itoa(webPort)
	r.Run(fmt.Sprintf(":%s",portString))
}
