package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
)

var script = flag.String("s", "na", "path to the script you'd like to run")
var port = flag.String("p", "na", "port to run on")
var gate = flag.String("g", "na", "gate to run on")

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func recieveInfo(rw http.ResponseWriter, req *http.Request) {
	_, err := exec.Command("/bin/sh", *script).Output()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *script == "na" {
		log.Fatal("Bad script path")
	}

	if *port == "na" {
		log.Fatal("Bad port")
	}

	if *gate == "na" {
		log.Fatal("Bad gate")
	}

	http.HandleFunc("/"+*gate, recieveInfo)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
