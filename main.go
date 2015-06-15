package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	//"os"
	"github.com/BurntSushi/toml"
)

type Preview struct {
	Width  int
	Height int
	Method string
}

type Config struct {
	Previews       map[string]Preview
	Asset_Bucket   string
	Asset_Prefix   string
	Preview_Bucket string
	Preview_Prefix string
	StorageDomain  string
	ListenPort     int
	ListenPortSSL  int
}

var configuration Config

func init() {
}

func main() {
	if _, err := toml.DecodeFile("config.toml", &configuration); err != nil {
		// handle error
		panic(err)
	}

	r := mux.NewRouter()

	http.Handle("/", r)
	registerHandlers(r)
	//fmt.Println(configuration.Previews)
	//	fmt.Println(configuration.Previews["small"])
	fmt.Printf("I'm listening...\n")
	http.ListenAndServe(fmt.Sprintf(":%d", configuration.ListenPort), nil)
}
