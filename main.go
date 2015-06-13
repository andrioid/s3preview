package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	//"os"
	"github.com/BurntSushi/toml"
)

// S3 creds (delete before github push)
// AKIAIKE57LTDKBK5W6FA
// 6sFwzVUXFLmzhMewi0jQQXBL0OJsTUmKA3TSTRkz

var asset_bucket = "andridk-assets"
var preview_bucket = "andridk-assets"
var preview_prefix = "s3preview"

type PreviewMethod int

type Preview struct {
	Width  int
	Height int
	Method string
}

type Config struct {
	Previews       map[string]Preview
	Asset_Bucket   string
	Preview_Bucket string
	Preview_Prefix string
}

var configuraton Config

func init() {
}

func main() {
	if _, err := toml.DecodeFile("config.toml", &configuraton); err != nil {
		// handle error
		panic(err)
	}
	configuraton.Asset_Bucket = "andridk-assets"
	configuraton.Preview_Bucket = "andridk-assets"
	configuraton.Preview_Prefix = "s3preview"

	r := mux.NewRouter()

	http.Handle("/", r)
	registerHandlers(r)
	fmt.Println(configuraton.Previews)
	//	fmt.Println(configuraton.Previews["small"])
	fmt.Printf("I'm listening...\n")
	http.ListenAndServe(":8097", nil)
}
