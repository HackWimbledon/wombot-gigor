package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var brains *Brains
var brainmanager *BrainManager

func main() {
	brains = new(Brains)
	err := brains.Initialise()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	brainmanager = new(BrainManager)
	err = brainmanager.Initialise(brains)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router := httprouter.New()
	router.GET("/brains", getBrains)
	router.POST("/brain/:brainid", getBrain)
	router.NotFound = http.FileServer(http.Dir("static"))
	log.Fatal(http.ListenAndServe(":8080", router))

}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func getBrains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(brains)
}

func getBrain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	brainid := ps.ByName("brainid")
	fmt.Fprintf(w, "BRAIN! %q", html.EscapeString(brains.Brains[brainid].Start))
}
